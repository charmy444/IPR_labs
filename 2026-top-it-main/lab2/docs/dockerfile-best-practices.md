# Dockerfile Best Practices

## Содержание
- [Общие принципы](#общие-принципы)
- [Выбор базового образа](#выбор-базового-образа)
- [Оптимизация слоев](#оптимизация-слоев)
- [Безопасность](#безопасность)
- [Производительность](#производительность)
- [Примеры](#примеры)

---

## Общие принципы

### 1. Используйте .dockerignore

```bash
# .dockerignore
node_modules
npm-debug.log
.git
.gitignore
README.md
.env
.env.*
*.md
.DS_Store
.vscode
.idea
__pycache__
*.pyc
*.pyo
*.pyd
.pytest_cache
.coverage
*.log
tmp/
temp/
```

**Зачем:**
- Уменьшает размер контекста сборки
- Ускоряет сборку
- Предотвращает попадание ненужных файлов в образ

### 2. Используйте конкретные версии

❌ **Плохо:**
```dockerfile
FROM python:latest
FROM node
```

✅ **Хорошо:**
```dockerfile
FROM python:3.11.6-slim-bookworm
FROM node:18.17.0-alpine
```

**Преимущества:**
- Предсказуемые сборки
- Избежание breaking changes
- Лучшая воспроизводимость

### 3. Минимизируйте количество слоев

❌ **Плохо:**
```dockerfile
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y git
RUN apt-get clean
```

✅ **Хорошо:**
```dockerfile
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
```

---

## Выбор базового образа

### Рекомендуемые базовые образы

| Язык/Платформа | Минимальный | Оптимальный | Размер |
|----------------|-------------|-------------|--------|
| Python | `python:3.11-alpine` | `python:3.11-slim` | 50-180MB |
| Node.js | `node:18-alpine` | `node:18-slim` | 40-200MB |
| Go | `scratch` | `alpine:3.18` | 0-7MB |
| Java | `eclipse-temurin:17-jre-alpine` | `eclipse-temurin:17-jre` | 150-250MB |
| .NET | `mcr.microsoft.com/dotnet/runtime:7.0-alpine` | `mcr.microsoft.com/dotnet/aspnet:7.0` | 100-200MB |

### Alpine vs Slim vs Full

**Alpine** (минимальный размер)
```dockerfile
FROM python:3.11-alpine
# Размер: ~50MB
# Плюсы: Очень маленький
# Минусы: musl libc (проблемы совместимости), медленная сборка нативных пакетов
```

**Slim** (оптимальный)
```dockerfile
FROM python:3.11-slim
# Размер: ~120MB
# Плюсы: Небольшой, glibc, хорошая совместимость
# Минусы: Чуть больше Alpine
```

**Full** (для разработки)
```dockerfile
FROM python:3.11
# Размер: ~900MB
# Плюсы: Все инструменты включены
# Минусы: Очень большой, не для production
```

### Distroless образы (Google)

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o app

FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/app /
CMD ["/app"]
```

**Преимущества:**
- Минимальный размер
- Нет shell, package manager
- Максимальная безопасность
- Только runtime dependencies

---

## Оптимизация слоев

### 1. Правильный порядок COPY

❌ **Плохо:**
```dockerfile
# Копирует все сразу, ломает кэш при любом изменении
COPY . .
RUN npm install
```

✅ **Хорошо:**
```dockerfile
# Копирует сначала зависимости, затем код
COPY package*.json ./
RUN npm ci --only=production
COPY . .
```

### 2. Multi-stage builds

❌ **Плохо:**
```dockerfile
FROM node:18
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build
CMD ["node", "dist/server.js"]
# Размер: ~1GB
```

✅ **Хорошо:**
```dockerfile
# Этап сборки
FROM node:18 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production
FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
CMD ["node", "dist/server.js"]
# Размер: ~150MB
```

### 3. Использование build cache

```dockerfile
# Кэшируется если requirements.txt не изменился
COPY requirements.txt .
RUN pip install -r requirements.txt

# Изменения в коде не влияют на слой выше
COPY . .
```

### 4. Минимизация слоев с RUN

❌ **Плохо:**
```dockerfile
RUN apt-get update
RUN apt-get install -y curl
RUN curl -o file https://example.com/file
RUN tar -xzf file
RUN rm file
# 5 слоев, file остается в промежуточных слоях
```

✅ **Хорошо:**
```dockerfile
RUN apt-get update && \
    apt-get install -y curl && \
    curl -o file https://example.com/file && \
    tar -xzf file && \
    rm file && \
    apt-get purge -y curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
# 1 слой, file удален в том же слое
```

---

## Безопасность

### 1. Не используйте root

❌ **Плохо:**
```dockerfile
FROM node:18
COPY . /app
WORKDIR /app
CMD ["node", "server.js"]
# Запускается от root
```

✅ **Хорошо:**
```dockerfile
FROM node:18
RUN groupadd -r nodejs && useradd -r -g nodejs nodejs
WORKDIR /app
COPY --chown=nodejs:nodejs . .
USER nodejs
CMD ["node", "server.js"]
```

### 2. Не храните секреты в образе

❌ **Плохо:**
```dockerfile
ENV API_KEY=supersecret123
COPY .env /app/
RUN echo "password=secret" > config.txt
```

✅ **Хорошо:**
```dockerfile
# Использовать переменные окружения при запуске
# docker run -e API_KEY=$API_KEY myapp

# Или Docker secrets
# docker secret create api_key -
```

### 3. Используйте HEALTHCHECK

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8000/health || exit 1
```

### 4. Сканируйте на уязвимости

```bash
# Перед публикацией образа
trivy image myapp:latest
docker scout cves myapp:latest
```

### 5. Используйте read-only где возможно

```dockerfile
# Минимальные права на файлы
COPY --chown=user:user --chmod=555 app.py /app/

# Read-only корневая ФС
# docker run --read-only --tmpfs /tmp myapp
```

### 6. Обновляйте базовый образ

```dockerfile
FROM python:3.11-slim

# Обновление пакетов безопасности
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
```

---

## Производительность

### 1. Используйте BuildKit

```bash
export DOCKER_BUILDKIT=1
docker build -t myapp .
```

**Преимущества:**
- Параллельная сборка
- Улучшенный кэш
- Secrets для build
- SSH forwarding

### 2. Используйте cache mounts

```dockerfile
# Кэш pip пакетов между сборками
RUN --mount=type=cache,target=/root/.cache/pip \
    pip install -r requirements.txt

# Кэш npm пакетов
RUN --mount=type=cache,target=/root/.npm \
    npm ci --only=production
```

### 3. Удаляйте ненужное в том же слое

```dockerfile
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential && \
    pip install -r requirements.txt && \
    apt-get purge -y build-essential && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
```

### 4. Используйте .dockerignore

```bash
# Исключить тяжелые папки
node_modules/
venv/
.git/
*.log
```

---

## Примеры

### Python приложение

```dockerfile
# Используем multi-stage build
FROM python:3.11-slim AS builder

# Build зависимости
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gcc \
        libc6-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Кэшируем зависимости
COPY requirements.txt .
RUN pip wheel --no-cache-dir --wheel-dir /wheels -r requirements.txt

# Production stage
FROM python:3.11-slim

# Runtime зависимости
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl && \
    rm -rf /var/lib/apt/lists/*

# Непривилегированный пользователь
RUN useradd -m -u 1000 appuser

WORKDIR /app

# Копируем wheels из builder
COPY --from=builder /wheels /wheels
RUN pip install --no-cache-dir --no-index --find-links=/wheels /wheels/* && \
    rm -rf /wheels

# Копируем приложение
COPY --chown=appuser:appuser . .

USER appuser

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=3s \
    CMD curl -f http://localhost:8000/health || exit 1

CMD ["python", "app.py"]
```

### Node.js приложение

```dockerfile
# Development dependencies
FROM node:18-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci

# Builder
FROM node:18-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# Production
FROM node:18-alpine AS production

# Безопасность
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

WORKDIR /app

# Копируем только production зависимости
COPY package*.json ./
RUN npm ci --only=production && \
    npm cache clean --force

# Копируем собранное приложение
COPY --from=builder --chown=nodejs:nodejs /app/dist ./dist

USER nodejs

EXPOSE 3000

ENV NODE_ENV=production

HEALTHCHECK --interval=30s --timeout=3s \
    CMD node healthcheck.js || exit 1

CMD ["node", "dist/server.js"]
```

### Go приложение

```dockerfile
# Builder
FROM golang:1.21-alpine AS builder

# Build зависимости
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Кэшируем модули
COPY go.mod go.sum ./
RUN go mod download

# Сборка
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o app .

# Production (минимальный образ)
FROM scratch

# Копируем CA сертификаты для HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Копируем timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Копируем бинарник
COPY --from=builder /build/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
```

---

## Чек-лист Best Practices

### Размер и производительность
- [ ] Используется базовый образ минимального размера (alpine/slim)
- [ ] Применяется multi-stage build
- [ ] Правильный порядок COPY (зависимости → код)
- [ ] Минимизировано количество слоев
- [ ] Удалены build зависимости
- [ ] Очищен package manager cache
- [ ] Используется .dockerignore

### Безопасность
- [ ] Используется конкретная версия базового образа
- [ ] Контейнер запускается от непривилегированного пользователя
- [ ] Нет секретов в образе или истории
- [ ] Добавлен HEALTHCHECK
- [ ] Обновлены пакеты безопасности
- [ ] Образ сканируется на уязвимости

### Качество кода
- [ ] Добавлены LABEL с метаданными
- [ ] Используется WORKDIR вместо cd
- [ ] Используется COPY вместо ADD (где возможно)
- [ ] CMD и ENTRYPOINT в exec форме
- [ ] Явно указан EXPOSE
- [ ] Инструкции отсортированы логически

### Документация
- [ ] README.md с инструкциями по сборке
- [ ] Комментарии в Dockerfile
- [ ] Примеры использования
- [ ] Переменные окружения документированы

---

## Инструменты для проверки

### Hadolint (Dockerfile linter)

```bash
# Установка
brew install hadolint

# Проверка
hadolint Dockerfile

# С игнорированием правил
hadolint --ignore DL3006 --ignore DL3018 Dockerfile
```

### Dive (анализ слоев)

```bash
# Установка
brew install dive

# Анализ
dive myapp:latest
```

### Trivy (сканер уязвимостей)

```bash
# Установка
brew install trivy

# Сканирование
trivy image myapp:latest
trivy image --severity HIGH,CRITICAL myapp:latest
```

---

## Дополнительные ресурсы

- [Dockerfile Reference](https://docs.docker.com/engine/reference/builder/)
- [Best practices for writing Dockerfiles](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Docker Security](https://docs.docker.com/engine/security/)
- [BuildKit](https://docs.docker.com/build/buildkit/)

---

Следуйте этим best practices, и ваши Docker образы будут быстрыми, безопасными и оптимизированными! 🐳

