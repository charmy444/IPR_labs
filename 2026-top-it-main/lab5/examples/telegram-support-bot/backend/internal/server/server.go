package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"telegram-support-bot/internal/appmetrics"
	"telegram-support-bot/internal/bot"
	"telegram-support-bot/internal/config"
	"telegram-support-bot/internal/models"
	"telegram-support-bot/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Server struct {
	config *config.Config
	repo   *repository.Repository
	bot    *bot.Bot
	router *gin.Engine
}

func NewServer(cfg *config.Config, repo *repository.Repository, botInstance *bot.Bot) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Use(otelgin.Middleware(cfg.OTELServiceName))
	router.Use(appmetrics.GinHTTPMiddleware())

	s := &Server{
		config: cfg,
		repo:   repo,
		bot:    botInstance,
		router: router,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// API маршруты
	api := s.router.Group("/api")
	{
		// Сообщения
		api.GET("/messages", s.getMessages)
		api.GET("/messages/:id", s.getMessage)
		api.GET("/messages/unread", s.getUnreadMessages)
		api.POST("/messages/:id/read", s.markAsRead)

		// Ответы поддержки
		api.POST("/responses", s.createResponse)
		api.GET("/responses/:messageId", s.getResponses)

		// Статистика
		api.GET("/stats", s.getStats)

		// Пользователи
		api.GET("/users", s.getUsers)
		api.GET("/users/:id/messages", s.getUserMessages)
	}

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// getMessages возвращает список сообщений
func (s *Server) getMessages(c *gin.Context) {
	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	messages, err := s.repo.GetMessagesWithResponses(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"limit":    limit,
		"offset":   offset,
	})
}

// getMessage возвращает сообщение по ID
func (s *Server) getMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	message, err := s.repo.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if message == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	responses, err := s.repo.GetResponsesByMessageID(c.Request.Context(), message.Message.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   message,
		"responses": responses,
	})
}

// getUnreadMessages возвращает непрочитанные сообщения
func (s *Server) getUnreadMessages(c *gin.Context) {
	messages, err := s.repo.GetUnreadMessages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

// markAsRead помечает сообщение как прочитанное
func (s *Server) markAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	err = s.repo.MarkMessagesAsRead(c.Request.Context(), []int64{id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// createResponse создает ответ поддержки
func (s *Server) createResponse(c *gin.Context) {
	var req struct {
		MessageID        int64  `json:"message_id" binding:"required"`
		ResponseText     string `json:"response_text" binding:"required"`
		SupportAgentName string `json:"support_agent_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем сообщение
	message, err := s.repo.GetMessageByID(c.Request.Context(), req.MessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if message == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Получаем пользователя
	user, err := s.repo.GetUserByMessageID(c.Request.Context(), req.MessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Отправляем ответ через бота (просто сообщение, без reply)
	_, err = s.bot.SendMessageToUser(user.ID, req.ResponseText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message: " + err.Error()})
		return
	}

	// Создаем запись ответа в базе данных
	response := &models.SupportResponse{
		MessageID: message.ID,
		StaffID:   user.ID, // Используем ID пользователя (в будущем нужно создать отдельного агента поддержки)
		Content:   req.ResponseText,
	}

	err = s.repo.CreateSupportResponse(c.Request.Context(), response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	appmetrics.SupportResponsesSent.Inc()
	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"response": response,
	})
}

// getResponses возвращает ответы на сообщение
func (s *Server) getResponses(c *gin.Context) {
	messageIDStr := c.Param("messageId")
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	responses, err := s.repo.GetResponsesByMessageID(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"responses": responses,
	})
}

// getStats возвращает статистику
func (s *Server) getStats(c *gin.Context) {
	stats, err := s.repo.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// UserWithMessageCount — пользователь с количеством сообщений
type UserWithMessageCount struct {
	models.User
	MessageCount int `json:"message_count"`
}

// getUsers возвращает список пользователей с количеством сообщений
func (s *Server) getUsers(c *gin.Context) {
	messages, err := s.repo.GetMessagesWithResponses(c.Request.Context(), 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	usersMap := make(map[int64]*models.User)
	countMap := make(map[int64]int)
	for _, msg := range messages {
		usersMap[msg.User.ID] = &msg.User
		countMap[msg.User.ID]++
	}

	users := make([]UserWithMessageCount, 0, len(usersMap))
	for _, user := range usersMap {
		users = append(users, UserWithMessageCount{
			User:         *user,
			MessageCount: countMap[user.ID],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// getUserMessages возвращает сообщения пользователя
func (s *Server) getUserMessages(c *gin.Context) {
	userID := c.Param("id")

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	messages, err := s.repo.GetMessagesByUserIDWithResponses(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

// Start запускает сервер
func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         ":" + s.config.ServerPort,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv.ListenAndServe()
}

// Shutdown останавливает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	// Gin engine doesn't have Close method, so we just return nil
	return nil
}

// Вспомогательные функции
// parseUUID was removed as we're using int64 IDs instead of UUIDs
