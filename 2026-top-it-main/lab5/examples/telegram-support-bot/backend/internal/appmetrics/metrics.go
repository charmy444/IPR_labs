package appmetrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)

	HTTPRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"method", "route", "status"},
	)

	// SupportTicketsCreated counts user messages persisted as support tickets (Telegram bot path).
	SupportTicketsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "support_tickets_created_total",
		Help: "Total number of user support messages stored",
	})

	// SupportResponsesSent counts successful support responses sent via API.
	SupportResponsesSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "support_responses_sent_total",
		Help: "Total number of support responses created via HTTP API",
	})
)

// GinHTTPMiddleware records request counts and latency using route template (low cardinality).
func GinHTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		method := c.Request.Method

		HTTPRequestsTotal.WithLabelValues(method, route, status).Inc()
		HTTPRequestDurationSeconds.WithLabelValues(method, route, status).Observe(time.Since(start).Seconds())
	}
}
