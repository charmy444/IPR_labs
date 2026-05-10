# Telegram Support Bot

Телеграм-бот для техподдержки с веб-дашбордом для операторов.

## Что это

- **Telegram-бот** (Go) — приём сообщений, регистрация пользователей, команды `/start`, `/help`.
- **REST API** (Go, Gin) — список сообщений, ответы, статистика.
- **Веб-дашборд** (Next.js) — просмотр диалогов, ответы пользователям, статистика.
- **PostgreSQL** — хранение пользователей, сообщений и ответов поддержки.

## Быстрый старт

1. Скопировать `.env.example` в `.env` и указать `BOT_TOKEN` (получить у [@BotFather](https://t.me/BotFather)).
2. Запуск: `docker-compose up -d`
3. Дашборд: http://localhost:3000  
   API: http://localhost:8080  
   Проверка: http://localhost:8080/health

## Архитектура

```
Telegram Bot (Go)  →  PostgreSQL  ←  REST API (Go)
                              ↑
                    Next.js дашборд
```

## Переменные окружения

| Переменная | Описание |
|------------|----------|
| `BOT_TOKEN` | Токен бота от @BotFather |
| `DATABASE_URL` | Строка подключения к PostgreSQL (по умолчанию: `postgres://postgres:postgres@localhost:5432/support_bot?sslmode=disable`) |
| `SERVER_PORT` | Порт API (по умолчанию `8080`) |
| `SERVER_HOST` | Хост API (по умолчанию `0.0.0.0`) |
| `NEXT_PUBLIC_API_URL` | URL API для фронтенда (в Docker: `http://backend:8080`) |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | (Опционально) OTLP HTTP endpoint для трейсов, например `http://tempo:4318` в Docker с [observability](observability/README.md) |
| `OTEL_SERVICE_NAME` | Имя сервиса в трейсах (по умолчанию `telegram-support-backend`) |

## Observability (лабораторная №7)

Backend отдаёт метрики Prometheus на **`GET /metrics`** (счётчики HTTP, гистограмма задержек, бизнес-метрики `support_tickets_created_total` и `support_responses_sent_total`). При заданном `OTEL_EXPORTER_OTLP_ENDPOINT` включается трейсинг (OpenTelemetry → OTLP).

Локальный стек **Prometheus + Grafana + Tempo** (порт Grafana **3001**, чтобы не конфликтовать с frontend на 3000):

```bash
docker compose -f docker-compose.yml -f docker-compose.observability.yml up -d
```

Подробности и проверки — в [observability/README.md](observability/README.md). Развёртывание мониторинга в Kubernetes — в [lab7/examples/telegram-support-observability](../../../lab7/examples/telegram-support-observability/README.md).

## API (кратко)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/messages` | Список сообщений (`?limit=50&offset=0`) |
| GET | `/api/messages/:id` | Сообщение по ID с ответами |
| GET | `/api/messages/unread` | Непрочитанные сообщения |
| POST | `/api/messages/:id/read` | Отметить как прочитанное |
| POST | `/api/responses` | Отправить ответ пользователю (JSON: `message_id`, `response_text`) |
| GET | `/api/responses/:messageId` | Ответы по сообщению |
| GET | `/api/stats` | Статистика |
| GET | `/api/users` | Список пользователей с числом сообщений |
| GET | `/api/users/:id/messages` | Сообщения пользователя |
| GET | `/metrics` | Метрики Prometheus (лаб. 7) |

## Структура проекта

```
telegram-support-bot/
├── backend/                    # Go-сервис (бот + API)
│   ├── cmd/telegram-support-bot/main.go
│   ├── internal/
│   │   ├── appmetrics/         # Метрики Prometheus (лаб. 7)
│   │   ├── bot/                # Логика Telegram-бота
│   │   ├── config/
│   │   ├── models/
│   │   ├── repository/         # Работа с БД
│   │   ├── server/             # HTTP API
│   │   └── telemetry/          # OpenTelemetry / OTLP (лаб. 7)
│   └── migrations/             # SQL-миграции
├── frontend/                   # Next.js (React)
│   └── app/
│       ├── page.tsx            # Главная
│       └── dashboard/         # Дашборд поддержки
├── observability/              # Prometheus, Grafana, Tempo (лаб. 7)
├── docker-compose.yml
├── docker-compose.observability.yml
├── .env.example
├── k8s/ # Kubernetes: только приложение (см. k8s/README.md)
└── README.md
```

## Kubernetes

Манифесты в каталоге [`k8s/`](k8s/README.md) описывают **frontend и backend**, без PostgreSQL. База развёртывается отдельно: пример инфраструктуры [telegram-support-infra](../../../lab6/examples/telegram-support-infra/README.md). Порядок работы, Kustomize и Helm — в [лабораторной работе №6](../../../lab6/README.md).

## Команды бота

- `/start` — регистрация пользователя.
- `/help` — справка по командам.
- Любой текст — создаётся обращение в поддержку; сообщение сохраняется в БД и отображается в дашборде.

## Стек

- **Backend:** Go, Gin, go-telegram-bot-api, PostgreSQL (lib/pq).
- **Frontend:** Next.js 14, React, TypeScript, Axios, Bootstrap 5.
- **БД:** PostgreSQL 15, миграции из папки `migrations/`.
