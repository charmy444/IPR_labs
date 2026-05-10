# Лабораторная работа №3: Docker Compose и Multi-stage Builds

## Цель работы

Освоить создание многоконтейнерных приложений с помощью Docker Compose, научиться оптимизировать Docker образы с использованием multi-stage builds, а также изучить продвинутые возможности оркестрации контейнеров.

## Требования

- Docker установлен и работает
- Docker Compose установлен (версия 2.x)
- Успешно выполнена Лабораторная работа №2
- Базовые знания работы с Docker (Dockerfile, образы, контейнеры)

---

## Содержание

- [Часть 1: Введение в Docker Compose](#часть-1-введение-в-docker-compose)
- [Часть 2: Структура docker-compose.yml](#часть-2-структура-docker-composeyml)
- [Часть 3: Работа с несколькими сервисами](#часть-3-работа-с-несколькими-сервисами)
- [Часть 4: Сети и Volumes в Docker Compose](#часть-4-сети-и-volumes-в-docker-compose)
- [Часть 5: Переменные окружения и .env файлы](#часть-5-переменные-окружения-и-env-файлы)
- [Часть 6: Healthchecks и зависимости](#часть-6-healthchecks-и-зависимости)
- [Часть 7: Multi-stage Builds](#часть-7-multi-stage-builds)
- [Часть 8: Профили и расширенные возможности](#часть-8-профили-и-расширенные-возможности)
- [Часть 9: Практический проект - Fullstack приложение](#часть-9-практический-проект---fullstack-приложение)
- [Задание для самостоятельной работы](#задание-для-самостоятельной-работы)
- [Контрольные вопросы](#контрольные-вопросы)

---

## Часть 1: Введение в Docker Compose

### Что такое Docker Compose?

**Docker Compose** — это инструмент для определения и запуска многоконтейнерных Docker-приложений. С помощью YAML-файла вы описываете все сервисы вашего приложения, а затем одной командой создаёте и запускаете их.

### Почему Docker Compose?

| Без Compose | С Compose |
|-------------|-----------|
| Множество команд `docker run` | Одна команда `docker-compose up` |
| Ручное создание сетей | Автоматические сети |
| Сложное управление зависимостями | `depends_on` с условиями |
| Разрозненные конфигурации | Всё в одном файле |
| Сложное масштабирование | `--scale service=N` |

### Основные команды Docker Compose

| Команда | Описание |
|---------|----------|
| `docker-compose up` | Создать и запустить контейнеры |
| `docker-compose up -d` | Запустить в фоновом режиме |
| `docker-compose down` | Остановить и удалить контейнеры |
| `docker-compose ps` | Показать запущенные сервисы |
| `docker-compose logs` | Показать логи |
| `docker-compose logs -f` | Следить за логами |
| `docker-compose build` | Собрать образы |
| `docker-compose pull` | Скачать образы |
| `docker-compose exec` | Выполнить команду в контейнере |
| `docker-compose restart` | Перезапустить сервисы |
| `docker-compose stop` | Остановить сервисы |
| `docker-compose start` | Запустить остановленные сервисы |

---

### Практическое задание 1.1: Первый docker-compose.yml

**Создайте простую конфигурацию:**

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  web:
    image: nginx:alpine
    ports:
      - "8080:80"
    container_name: my-nginx
```

**Выполните команды:**
```bash
# Создать директорию проекта
mkdir compose-intro
cd compose-intro

# Создать docker-compose.yml (содержимое выше)

# Запустить
docker-compose up -d

# Проверить статус
docker-compose ps

# Посмотреть логи
docker-compose logs

# Открыть в браузере: http://localhost:8080

# Остановить
docker-compose down
```

**Что произошло:**
1. Docker Compose прочитал файл `docker-compose.yml`
2. Скачал образ `nginx:alpine` (если его нет локально)
3. Создал сеть `compose-intro_default`
4. Создал и запустил контейнер `my-nginx`
5. Пробросил порт 8080 хоста на порт 80 контейнера

---

## Часть 2: Структура docker-compose.yml

### Основные секции

```yaml
version: '3.8'          # Версия формата Compose

services:               # Определение сервисов (контейнеров)
  service_name:
    image: image:tag    # Образ для использования
    build: ./path       # ИЛИ путь к Dockerfile
    ports:
      - "host:container"
    volumes:
      - ./local:/container
    environment:
      - VAR=value
    depends_on:
      - other_service

volumes:                # Определение volumes
  volume_name:

networks:               # Определение сетей
  network_name:
```

### Параметры сервиса

| Параметр | Описание | Пример |
|----------|----------|--------|
| `image` | Образ для использования | `image: nginx:alpine` |
| `build` | Путь к Dockerfile | `build: ./app` |
| `ports` | Проброс портов | `ports: ["8080:80"]` |
| `volumes` | Монтирование директорий | `volumes: ["./data:/app/data"]` |
| `environment` | Переменные окружения | `environment: [DEBUG=true]` |
| `env_file` | Файл с переменными | `env_file: .env` |
| `depends_on` | Зависимости | `depends_on: [db]` |
| `container_name` | Имя контейнера | `container_name: my-app` |
| `restart` | Политика перезапуска | `restart: unless-stopped` |
| `networks` | Сети | `networks: [frontend, backend]` |
| `command` | Переопределение CMD | `command: python app.py` |
| `entrypoint` | Переопределение ENTRYPOINT | `entrypoint: /start.sh` |

---

### Практическое задание 2.1: Сборка из Dockerfile

**Структура проекта:**
```
flask-compose/
├── docker-compose.yml
├── Dockerfile
├── app.py
└── requirements.txt
```

**app.py:**
```python
from flask import Flask, jsonify
import os
import socket

app = Flask(__name__)

@app.route('/')
def home():
    return jsonify({
        'message': 'Hello from Docker Compose!',
        'hostname': socket.gethostname(),
        'environment': os.getenv('FLASK_ENV', 'production')
    })

@app.route('/health')
def health():
    return jsonify({'status': 'healthy'}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
```

**requirements.txt:**
```
Flask==3.0.0
```

**Dockerfile:**
```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 5000

CMD ["python", "app.py"]
```

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  web:
    build: .
    ports:
      - "5000:5000"
    environment:
      - FLASK_ENV=development
    volumes:
      - ./app.py:/app/app.py
    container_name: flask-app
```

**Выполните команды:**
```bash
# Собрать и запустить
docker-compose up -d --build

# Проверить
curl http://localhost:5000
curl http://localhost:5000/health

# Посмотреть логи
docker-compose logs -f web

# Остановить
docker-compose down
```

**Преимущества:**
- При изменении `app.py` код обновляется автоматически (bind mount)
- Флаг `--build` пересобирает образ при изменениях

---

## Часть 3: Работа с несколькими сервисами

### Практическое задание 3.1: Веб-приложение с базой данных

**Структура проекта:**
```
webapp-db/
├── docker-compose.yml
├── app/
│   ├── Dockerfile
│   ├── app.py
│   └── requirements.txt
└── init-db/
    └── init.sql
```

**app/app.py:**
```python
from flask import Flask, jsonify, request
import psycopg2
import os
import time

app = Flask(__name__)

def get_db_connection():
    """Получение соединения с БД с повторными попытками"""
    max_retries = 10
    for i in range(max_retries):
        try:
            conn = psycopg2.connect(
                host=os.environ.get('DB_HOST', 'db'),
                database=os.environ.get('DB_NAME', 'appdb'),
                user=os.environ.get('DB_USER', 'postgres'),
                password=os.environ.get('DB_PASSWORD', 'secret')
            )
            return conn
        except psycopg2.OperationalError as e:
            if i < max_retries - 1:
                print(f"Database not ready, retrying... ({i+1}/{max_retries})")
                time.sleep(2)
            else:
                raise e

@app.route('/')
def home():
    return jsonify({
        'message': 'Flask + PostgreSQL with Docker Compose!',
        'endpoints': ['/users', '/users/<id>', '/health', '/db-status']
    })

@app.route('/health')
def health():
    return jsonify({'status': 'healthy'}), 200

@app.route('/db-status')
def db_status():
    try:
        conn = get_db_connection()
        cur = conn.cursor()
        cur.execute('SELECT version();')
        version = cur.fetchone()[0]
        cur.close()
        conn.close()
        return jsonify({
            'status': 'connected',
            'database': version
        })
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        }), 500

@app.route('/users', methods=['GET', 'POST'])
def users():
    conn = get_db_connection()
    cur = conn.cursor()
    
    if request.method == 'POST':
        data = request.get_json()
        cur.execute(
            'INSERT INTO users (name, email) VALUES (%s, %s) RETURNING id, name, email, created_at',
            (data['name'], data['email'])
        )
        user = cur.fetchone()
        conn.commit()
        cur.close()
        conn.close()
        return jsonify({
            'id': user[0],
            'name': user[1],
            'email': user[2],
            'created_at': str(user[3])
        }), 201
    
    cur.execute('SELECT id, name, email, created_at FROM users ORDER BY id')
    users = cur.fetchall()
    cur.close()
    conn.close()
    
    return jsonify([{
        'id': u[0],
        'name': u[1],
        'email': u[2],
        'created_at': str(u[3])
    } for u in users])

@app.route('/users/<int:user_id>')
def get_user(user_id):
    conn = get_db_connection()
    cur = conn.cursor()
    cur.execute('SELECT id, name, email, created_at FROM users WHERE id = %s', (user_id,))
    user = cur.fetchone()
    cur.close()
    conn.close()
    
    if user:
        return jsonify({
            'id': user[0],
            'name': user[1],
            'email': user[2],
            'created_at': str(user[3])
        })
    return jsonify({'error': 'User not found'}), 404

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
```

**app/requirements.txt:**
```
Flask==3.0.0
psycopg2-binary==2.9.9
```

**app/Dockerfile:**
```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 5000

CMD ["python", "app.py"]
```

**init-db/init.sql:**
```sql
-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Добавление тестовых данных
INSERT INTO users (name, email) VALUES 
    ('Иван Петров', 'ivan@example.com'),
    ('Мария Сидорова', 'maria@example.com'),
    ('Алексей Козлов', 'alexey@example.com');
```

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    container_name: postgres-db
    environment:
      POSTGRES_DB: appdb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init-db:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  web:
    build: ./app
    container_name: flask-app
    environment:
      DB_HOST: db
      DB_NAME: appdb
      DB_USER: postgres
      DB_PASSWORD: secret
    ports:
      - "5000:5000"
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

volumes:
  postgres_data:
```

**Выполните команды:**
```bash
# Создать директории
mkdir -p webapp-db/{app,init-db}
cd webapp-db

# Создать все файлы (содержимое выше)

# Запустить
docker-compose up -d --build

# Проверить статус
docker-compose ps

# Проверить API
curl http://localhost:5000/
curl http://localhost:5000/db-status
curl http://localhost:5000/users

# Добавить пользователя
curl -X POST http://localhost:5000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Новый Пользователь", "email": "new@example.com"}'

# Получить пользователя
curl http://localhost:5000/users/1

# Посмотреть логи всех сервисов
docker-compose logs -f

# Выполнить SQL в базе данных
docker-compose exec db psql -U postgres -d appdb -c "SELECT * FROM users;"

# Остановить и удалить (с данными)
docker-compose down -v
```

**Что изучили:**
- Связь сервисов через сеть
- Использование `depends_on` с `condition: service_healthy`
- Инициализация БД через `docker-entrypoint-initdb.d`
- Именованные volumes для persistence

---

## Часть 4: Сети и Volumes в Docker Compose

### Сети в Docker Compose

По умолчанию Docker Compose создаёт одну сеть для всех сервисов в проекте. Но можно создавать изолированные сети:

```yaml
version: '3.8'

services:
  frontend:
    image: nginx:alpine
    networks:
      - frontend-network
      
  backend:
    image: python:3.11-slim
    networks:
      - frontend-network
      - backend-network
      
  database:
    image: postgres:15-alpine
    networks:
      - backend-network

networks:
  frontend-network:
    driver: bridge
  backend-network:
    driver: bridge
    internal: true  # Нет доступа к интернету
```

### Типы networks

| Тип | Описание | Когда использовать |
|-----|----------|-------------------|
| `bridge` | По умолчанию, изолированная сеть | Большинство случаев |
| `host` | Сеть хоста | Максимальная производительность |
| `none` | Без сети | Полная изоляция |
| `overlay` | Между хостами | Docker Swarm |

### Volumes в Docker Compose

| Тип | Синтаксис | Описание |
|-----|-----------|----------|
| Именованный | `volume_name:/path` | Управляется Docker |
| Bind mount | `./local:/container` | Директория хоста |
| tmpfs | `type: tmpfs` | В памяти |

```yaml
version: '3.8'

services:
  app:
    image: myapp
    volumes:
      # Именованный volume
      - app_data:/app/data
      # Bind mount
      - ./config:/app/config:ro
      # tmpfs
      - type: tmpfs
        target: /app/temp

volumes:
  app_data:
    driver: local
```

---

### Практическое задание 4.1: Изолированные сети

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  # Nginx - доступен извне, связан с backend
  nginx:
    image: nginx:alpine
    container_name: proxy
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf:ro
    networks:
      - frontend
    depends_on:
      - api

  # API - связан с nginx и базой данных
  api:
    build: ./api
    container_name: api-server
    environment:
      DATABASE_URL: postgresql://postgres:secret@db:5432/appdb
    networks:
      - frontend
      - backend
    depends_on:
      db:
        condition: service_healthy

  # База данных - только в backend сети
  db:
    image: postgres:15-alpine
    container_name: database
    environment:
      POSTGRES_DB: appdb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    volumes:
      - db_data:/var/lib/postgresql/data
    networks:
      - backend
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Redis - только в backend сети
  cache:
    image: redis:7-alpine
    container_name: redis-cache
    networks:
      - backend
    volumes:
      - redis_data:/data

volumes:
  db_data:
  redis_data:

networks:
  frontend:
    driver: bridge
  backend:
    driver: bridge
    internal: true  # Нет доступа к интернету
```

**Проверка изоляции:**
```bash
# Запустить
docker-compose up -d

# Проверить сети
docker network ls | grep compose

# Проверить, что nginx не может достучаться до db
docker-compose exec nginx ping -c 2 db
# Должна быть ошибка!

# API может достучаться до db
docker-compose exec api ping -c 2 db
# Должно работать

# Остановить
docker-compose down -v
```

---

## Часть 5: Переменные окружения и .env файлы

### Способы задания переменных

**1. Непосредственно в docker-compose.yml:**
```yaml
services:
  app:
    environment:
      - DEBUG=true
      - DATABASE_URL=postgresql://localhost/db
```

**2. Через файл .env:**
```yaml
services:
  app:
    env_file:
      - .env
```

**3. Подстановка переменных из .env:**
```yaml
services:
  db:
    image: postgres:${POSTGRES_VERSION:-15}-alpine
    environment:
      POSTGRES_PASSWORD: ${DB_PASSWORD}
```

### Практическое задание 5.1: Работа с .env

**Структура проекта:**
```
env-demo/
├── docker-compose.yml
├── .env
├── .env.example
└── app/
    ├── Dockerfile
    └── app.py
```

**.env:**
```env
# Версии образов
POSTGRES_VERSION=15
REDIS_VERSION=7

# База данных
DB_NAME=production_db
DB_USER=admin
DB_PASSWORD=super_secret_password_123
DB_HOST=db
DB_PORT=5432

# Приложение
APP_PORT=8000
APP_DEBUG=false
APP_SECRET_KEY=your-secret-key-here

# Redis
REDIS_HOST=cache
REDIS_PORT=6379
```

**.env.example:**
```env
# Скопируйте в .env и заполните реальными значениями

# Версии образов
POSTGRES_VERSION=15
REDIS_VERSION=7

# База данных
DB_NAME=
DB_USER=
DB_PASSWORD=
DB_HOST=db
DB_PORT=5432

# Приложение
APP_PORT=8000
APP_DEBUG=true
APP_SECRET_KEY=

# Redis
REDIS_HOST=cache
REDIS_PORT=6379
```

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  db:
    image: postgres:${POSTGRES_VERSION:-15}-alpine
    container_name: ${COMPOSE_PROJECT_NAME:-app}-db
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  cache:
    image: redis:${REDIS_VERSION:-7}-alpine
    container_name: ${COMPOSE_PROJECT_NAME:-app}-cache
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data

  app:
    build: ./app
    container_name: ${COMPOSE_PROJECT_NAME:-app}-api
    ports:
      - "${APP_PORT:-8000}:8000"
    environment:
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - DEBUG=${APP_DEBUG}
      - SECRET_KEY=${APP_SECRET_KEY}
    depends_on:
      db:
        condition: service_healthy
      cache:
        condition: service_started

volumes:
  postgres_data:
  redis_data:
```

**Проверка конфигурации:**
```bash
# Показать результирующую конфигурацию
docker-compose config

# Запустить
docker-compose up -d

# Проверить переменные внутри контейнера
docker-compose exec app env | grep DB

# Остановить
docker-compose down -v
```

**Важно:** Никогда не коммитьте `.env` файл с реальными паролями! Добавьте его в `.gitignore`.

---

## Часть 6: Healthchecks и зависимости

### Что такое Healthcheck?

Healthcheck позволяет Docker'у проверять, работает ли контейнер правильно.

**Параметры healthcheck:**

| Параметр | Описание | Значение по умолчанию |
|----------|----------|----------------------|
| `test` | Команда проверки | - |
| `interval` | Интервал между проверками | 30s |
| `timeout` | Таймаут проверки | 30s |
| `retries` | Количество повторов | 3 |
| `start_period` | Начальная задержка | 0s |

### Типы depends_on

**Простой (старый способ):**
```yaml
services:
  app:
    depends_on:
      - db
```
⚠️ Не ждёт готовности сервиса!

**С условием (рекомендуется):**
```yaml
services:
  app:
    depends_on:
      db:
        condition: service_healthy
      cache:
        condition: service_started
```

### Условия зависимостей

| Условие | Описание |
|---------|----------|
| `service_started` | Контейнер запущен |
| `service_healthy` | Контейнер прошёл healthcheck |
| `service_completed_successfully` | Контейнер завершился успешно |

---

### Практическое задание 6.1: Healthchecks

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  # PostgreSQL с healthcheck
  postgres:
    image: postgres:15-alpine
    container_name: postgres-health
    environment:
      POSTGRES_DB: testdb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d testdb"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
    volumes:
      - postgres_data:/var/lib/postgresql/data

  # Redis с healthcheck
  redis:
    image: redis:7-alpine
    container_name: redis-health
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    volumes:
      - redis_data:/data

  # MySQL с healthcheck
  mysql:
    image: mysql:8
    container_name: mysql-health
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: testdb
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-uroot", "-psecret"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s
    volumes:
      - mysql_data:/var/lib/mysql

  # MongoDB с healthcheck
  mongo:
    image: mongo:6
    container_name: mongo-health
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: secret
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 20s
    volumes:
      - mongo_data:/data/db

  # Приложение, ожидающее все БД
  app:
    image: alpine:3.18
    container_name: app-health
    command: sh -c "echo 'All databases are ready!' && sleep infinity"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      mysql:
        condition: service_healthy
      mongo:
        condition: service_healthy

volumes:
  postgres_data:
  redis_data:
  mysql_data:
  mongo_data:
```

**Команды:**
```bash
# Запустить и следить за статусом
docker-compose up -d
docker-compose ps

# Следить за здоровьем
watch -n 1 'docker-compose ps'

# Проверить логи healthcheck
docker inspect --format='{{json .State.Health}}' postgres-health | jq

# Остановить
docker-compose down -v
```

---

## Часть 7: Multi-stage Builds

### Что такое Multi-stage Build?

Multi-stage build позволяет использовать несколько этапов (FROM) в одном Dockerfile:
- **Этап сборки:** устанавливаем все инструменты, компилируем код
- **Финальный этап:** копируем только готовые артефакты

**Преимущества:**
- 🔥 Значительно меньший размер образа
- 🔒 Меньше уязвимостей (нет лишних пакетов)
- ⚡ Быстрее деплой
- 🛠️ Разделение сборки и runtime

### Практическое задание 7.1: Go приложение

**main.go:**
```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "runtime"
)

type Response struct {
    Message     string `json:"message"`
    Hostname    string `json:"hostname"`
    GoVersion   string `json:"go_version"`
    Platform    string `json:"platform"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    hostname, _ := os.Hostname()
    
    response := Response{
        Message:   "Hello from Go Multi-stage Build!",
        Hostname:  hostname,
        GoVersion: runtime.Version(),
        Platform:  runtime.GOOS + "/" + runtime.GOARCH,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/health", healthHandler)
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    fmt.Printf("Server starting on port %s\n", port)
    http.ListenAndServe(":"+port, nil)
}
```

**Dockerfile БЕЗ multi-stage (для сравнения):**
```dockerfile
# Dockerfile.single
FROM golang:1.21

WORKDIR /app
COPY main.go .

RUN go build -o server main.go

EXPOSE 8080
CMD ["./server"]
```

**Dockerfile С multi-stage:**
```dockerfile
# Dockerfile.multi

# ========================================
# Этап 1: Сборка
# ========================================
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем исходный код
COPY main.go .

# Компилируем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o server main.go

# ========================================
# Этап 2: Финальный образ
# ========================================
FROM alpine:3.18

# Устанавливаем CA сертификаты (для HTTPS запросов)
RUN apk --no-cache add ca-certificates

# Создаём непривилегированного пользователя
RUN adduser -D -u 1000 appuser

WORKDIR /app

# Копируем ТОЛЬКО бинарник из этапа сборки
COPY --from=builder /app/server .

# Меняем владельца
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
```

**Сравнение размеров:**
```bash
# Собрать оба варианта
docker build -t go-app:single -f Dockerfile.single .
docker build -t go-app:multi -f Dockerfile.multi .

# Сравнить размеры
docker images | grep go-app

# Результат:
# go-app:single  ~850MB
# go-app:multi   ~15MB  (в ~57 раз меньше!)

# Запустить и проверить
docker run -d -p 8080:8080 --name go-multi go-app:multi
curl http://localhost:8080

# Очистка
docker stop go-multi && docker rm go-multi
```

---

### Практическое задание 7.2: Node.js/React приложение

**Структура проекта:**
```
react-multi/
├── Dockerfile
├── nginx.conf
├── package.json
├── public/
│   └── index.html
└── src/
    ├── App.js
    └── index.js
```

**public/index.html:**
```html
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>React Multi-stage</title>
</head>
<body>
    <div id="root"></div>
</body>
</html>
```

**src/index.js:**
```javascript
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<React.StrictMode><App /></React.StrictMode>);
```

**src/App.js:**
```javascript
import React, { useState, useEffect } from 'react';

function App() {
    const [time, setTime] = useState(new Date().toLocaleTimeString());

    useEffect(() => {
        const timer = setInterval(() => {
            setTime(new Date().toLocaleTimeString());
        }, 1000);
        return () => clearInterval(timer);
    }, []);

    return (
        <div style={{
            minHeight: '100vh',
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            fontFamily: 'system-ui, -apple-system, sans-serif'
        }}>
            <div style={{
                background: 'white',
                padding: '3rem',
                borderRadius: '1rem',
                boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
                textAlign: 'center'
            }}>
                <h1 style={{ color: '#667eea', marginBottom: '1rem' }}>
                    🐳 Multi-stage Build
                </h1>
                <p style={{ fontSize: '1.25rem', color: '#4a5568' }}>
                    React + Nginx = Production Ready
                </p>
                <p style={{ 
                    fontSize: '2rem', 
                    fontWeight: 'bold',
                    color: '#764ba2',
                    marginTop: '1.5rem'
                }}>
                    {time}
                </p>
            </div>
        </div>
    );
}

export default App;
```

**package.json:**
```json
{
  "name": "react-multi-stage",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-scripts": "5.0.1"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test"
  },
  "browserslist": {
    "production": [">0.2%", "not dead", "not op_mini all"],
    "development": ["last 1 chrome version"]
  }
}
```

**nginx.conf:**
```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    # Gzip сжатие
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # SPA routing - все запросы на index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Кеширование статических файлов
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Безопасность
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
}
```

**Dockerfile:**
```dockerfile
# ========================================
# Этап 1: Установка зависимостей
# ========================================
FROM node:18-alpine AS deps

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

# ========================================
# Этап 2: Сборка приложения
# ========================================
FROM node:18-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY public ./public
COPY src ./src

# Сборка production версии
RUN npm run build

# ========================================
# Этап 3: Production образ
# ========================================
FROM nginx:alpine AS production

# Удаляем дефолтную конфигурацию
RUN rm /etc/nginx/conf.d/default.conf

# Копируем нашу конфигурацию
COPY nginx.conf /etc/nginx/conf.d/

# Копируем собранное приложение
COPY --from=builder /app/build /usr/share/nginx/html

# Nginx запускается от root, но worker процессы от nginx
# Устанавливаем правильные права
RUN chown -R nginx:nginx /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
```

**Команды:**
```bash
# Собрать образ
docker build -t react-app:multi .

# Проверить размер
docker images react-app:multi
# Результат: ~40MB (вместо ~1GB для node:18)

# Запустить
docker run -d -p 3000:80 --name react-multi react-app:multi

# Открыть в браузере: http://localhost:3000

# Остановить
docker stop react-multi && docker rm react-multi
```

---

### Практическое задание 7.3: Python приложение

**Dockerfile.python-multi:**
```dockerfile
# ========================================
# Этап 1: Сборка зависимостей
# ========================================
FROM python:3.11-slim AS builder

WORKDIR /app

# Установка build зависимостей
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gcc \
        libc6-dev \
        libpq-dev && \
    rm -rf /var/lib/apt/lists/*

# Создаём виртуальное окружение
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Устанавливаем зависимости
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# ========================================
# Этап 2: Production образ
# ========================================
FROM python:3.11-slim AS production

# Установка только runtime зависимостей
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        libpq5 && \
    rm -rf /var/lib/apt/lists/*

# Создаём пользователя
RUN useradd -m -u 1000 appuser

WORKDIR /app

# Копируем виртуальное окружение
COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Копируем приложение
COPY --chown=appuser:appuser . .

USER appuser

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/health')" || exit 1

CMD ["gunicorn", "--bind", "0.0.0.0:8000", "--workers", "4", "app:app"]
```

---

### Таблица сравнения размеров образов

| Язык/Технология | Без multi-stage | С multi-stage | Экономия |
|-----------------|-----------------|---------------|----------|
| Go | ~850 MB | ~15 MB | 98% |
| Node.js/React | ~1.1 GB | ~40 MB | 96% |
| Python/Flask | ~900 MB | ~150 MB | 83% |
| Java/Spring | ~500 MB | ~100 MB | 80% |
| Rust | ~1.5 GB | ~10 MB | 99% |

---

## Часть 8: Профили и расширенные возможности

### Профили (Profiles)

Профили позволяют запускать определённые сервисы только когда нужно:

```yaml
version: '3.8'

services:
  app:
    image: myapp:latest
    ports:
      - "8080:8080"

  # Только для разработки
  debug:
    image: debug-tools
    profiles:
      - dev
      - debug

  # Только для production
  monitoring:
    image: prometheus
    profiles:
      - prod
      - monitoring
```

**Использование:**
```bash
# Только основной сервис
docker-compose up -d

# С профилем dev
docker-compose --profile dev up -d

# С несколькими профилями
docker-compose --profile dev --profile monitoring up -d
```

---

### Практическое задание 8.1: Профили для разных окружений

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  # Основное приложение (всегда запускается)
  app:
    build: ./app
    container_name: main-app
    ports:
      - "8080:8080"
    environment:
      - ENV=${ENV:-production}
    depends_on:
      db:
        condition: service_healthy

  # База данных (всегда запускается)
  db:
    image: postgres:15-alpine
    container_name: app-db
    environment:
      POSTGRES_DB: appdb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  # ===== Профиль: dev =====
  # Adminer для работы с БД
  adminer:
    image: adminer
    container_name: adminer
    ports:
      - "8081:8080"
    depends_on:
      - db
    profiles:
      - dev
      - tools

  # Hot reload watcher
  watcher:
    build:
      context: ./app
      target: development
    container_name: dev-watcher
    volumes:
      - ./app:/app
    command: npm run watch
    profiles:
      - dev

  # ===== Профиль: monitoring =====
  # Prometheus
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus_data:/prometheus
    ports:
      - "9090:9090"
    profiles:
      - monitoring
      - prod

  # Grafana
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    volumes:
      - grafana_data:/var/lib/grafana
    ports:
      - "3000:3000"
    depends_on:
      - prometheus
    profiles:
      - monitoring
      - prod

  # ===== Профиль: testing =====
  # Тестовая БД
  test-db:
    image: postgres:15-alpine
    container_name: test-db
    environment:
      POSTGRES_DB: testdb
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
    profiles:
      - testing

  # Test runner
  test-runner:
    build:
      context: ./app
      target: testing
    container_name: test-runner
    depends_on:
      - test-db
    command: npm run test
    profiles:
      - testing

volumes:
  postgres_data:
  prometheus_data:
  grafana_data:
```

**Команды:**
```bash
# Только основные сервисы (app + db)
docker-compose up -d

# Разработка (+ adminer, watcher)
docker-compose --profile dev up -d

# Production с мониторингом
docker-compose --profile prod up -d

# Тестирование
docker-compose --profile testing up -d

# Остановить все
docker-compose --profile dev --profile monitoring --profile testing down
```

---

### Расширенные опции

**Ограничение ресурсов:**
```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

**Logging:**
```yaml
services:
  app:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

**Restart политики:**
```yaml
services:
  app:
    restart: unless-stopped  # или: no, always, on-failure
```

**Extra hosts:**
```yaml
services:
  app:
    extra_hosts:
      - "host.docker.internal:host-gateway"
      - "myhost:192.168.1.100"
```

---

## Часть 9: Практический проект - Fullstack приложение

### Практическое задание 9.1: Полноценное веб-приложение

Создадим полноценное приложение со всеми компонентами:
- Frontend (React/Nginx)
- Backend (Node.js/Express)
- База данных (PostgreSQL)
- Кэш (Redis)
- Reverse Proxy (Nginx)

**Структура проекта:**
```
fullstack-app/
├── docker-compose.yml
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── .env
├── .env.example
├── nginx/
│   └── nginx.conf
├── backend/
│   ├── Dockerfile
│   ├── package.json
│   └── src/
│       └── server.js
└── frontend/
    ├── Dockerfile
    ├── nginx.conf
    ├── package.json
    ├── public/
    │   └── index.html
    └── src/
        ├── App.js
        └── index.js
```

**backend/src/server.js:**
```javascript
const express = require('express');
const { Pool } = require('pg');
const redis = require('redis');
const cors = require('cors');

const app = express();
app.use(cors());
app.use(express.json());

// PostgreSQL connection
const pool = new Pool({
    host: process.env.DB_HOST || 'db',
    port: process.env.DB_PORT || 5432,
    database: process.env.DB_NAME || 'appdb',
    user: process.env.DB_USER || 'postgres',
    password: process.env.DB_PASSWORD || 'secret',
});

// Redis connection
const redisClient = redis.createClient({
    url: `redis://${process.env.REDIS_HOST || 'cache'}:${process.env.REDIS_PORT || 6379}`
});
redisClient.connect().catch(console.error);

// Initialize database
const initDB = async () => {
    try {
        await pool.query(`
            CREATE TABLE IF NOT EXISTS tasks (
                id SERIAL PRIMARY KEY,
                title VARCHAR(255) NOT NULL,
                completed BOOLEAN DEFAULT FALSE,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        `);
        console.log('Database initialized');
    } catch (err) {
        console.error('Database init error:', err);
    }
};

// Routes
app.get('/api/health', (req, res) => {
    res.json({ status: 'healthy', timestamp: new Date().toISOString() });
});

app.get('/api/tasks', async (req, res) => {
    try {
        // Check cache
        const cached = await redisClient.get('tasks');
        if (cached) {
            return res.json({ tasks: JSON.parse(cached), source: 'cache' });
        }

        // Query database
        const result = await pool.query('SELECT * FROM tasks ORDER BY created_at DESC');
        
        // Cache for 30 seconds
        await redisClient.setEx('tasks', 30, JSON.stringify(result.rows));
        
        res.json({ tasks: result.rows, source: 'database' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

app.post('/api/tasks', async (req, res) => {
    try {
        const { title } = req.body;
        const result = await pool.query(
            'INSERT INTO tasks (title) VALUES ($1) RETURNING *',
            [title]
        );
        
        // Invalidate cache
        await redisClient.del('tasks');
        
        res.status(201).json(result.rows[0]);
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

app.put('/api/tasks/:id', async (req, res) => {
    try {
        const { id } = req.params;
        const { completed } = req.body;
        const result = await pool.query(
            'UPDATE tasks SET completed = $1 WHERE id = $2 RETURNING *',
            [completed, id]
        );
        
        // Invalidate cache
        await redisClient.del('tasks');
        
        if (result.rows.length === 0) {
            return res.status(404).json({ error: 'Task not found' });
        }
        res.json(result.rows[0]);
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

app.delete('/api/tasks/:id', async (req, res) => {
    try {
        const { id } = req.params;
        await pool.query('DELETE FROM tasks WHERE id = $1', [id]);
        
        // Invalidate cache
        await redisClient.del('tasks');
        
        res.status(204).send();
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

const PORT = process.env.PORT || 3000;

initDB().then(() => {
    app.listen(PORT, '0.0.0.0', () => {
        console.log(`Backend server running on port ${PORT}`);
    });
});
```

**backend/package.json:**
```json
{
  "name": "fullstack-backend",
  "version": "1.0.0",
  "main": "src/server.js",
  "scripts": {
    "start": "node src/server.js",
    "dev": "nodemon src/server.js"
  },
  "dependencies": {
    "express": "^4.18.2",
    "pg": "^8.11.3",
    "redis": "^4.6.10",
    "cors": "^2.8.5"
  },
  "devDependencies": {
    "nodemon": "^3.0.1"
  }
}
```

**backend/Dockerfile:**
```dockerfile
# Multi-stage build для backend

# ========================================
# Этап 1: Development
# ========================================
FROM node:18-alpine AS development

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 3000

CMD ["npm", "run", "dev"]

# ========================================
# Этап 2: Production dependencies
# ========================================
FROM node:18-alpine AS deps

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

# ========================================
# Этап 3: Production
# ========================================
FROM node:18-alpine AS production

RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY --chown=nodejs:nodejs . .

USER nodejs

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

CMD ["node", "src/server.js"]
```

**frontend/src/App.js:**
```javascript
import React, { useState, useEffect } from 'react';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000/api';

function App() {
    const [tasks, setTasks] = useState([]);
    const [newTask, setNewTask] = useState('');
    const [loading, setLoading] = useState(true);
    const [source, setSource] = useState('');

    const fetchTasks = async () => {
        try {
            const response = await fetch(`${API_URL}/tasks`);
            const data = await response.json();
            setTasks(data.tasks);
            setSource(data.source);
            setLoading(false);
        } catch (error) {
            console.error('Error fetching tasks:', error);
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTasks();
    }, []);

    const addTask = async (e) => {
        e.preventDefault();
        if (!newTask.trim()) return;

        try {
            await fetch(`${API_URL}/tasks`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ title: newTask })
            });
            setNewTask('');
            fetchTasks();
        } catch (error) {
            console.error('Error adding task:', error);
        }
    };

    const toggleTask = async (id, completed) => {
        try {
            await fetch(`${API_URL}/tasks/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ completed: !completed })
            });
            fetchTasks();
        } catch (error) {
            console.error('Error updating task:', error);
        }
    };

    const deleteTask = async (id) => {
        try {
            await fetch(`${API_URL}/tasks/${id}`, { method: 'DELETE' });
            fetchTasks();
        } catch (error) {
            console.error('Error deleting task:', error);
        }
    };

    return (
        <div style={{
            minHeight: '100vh',
            background: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%)',
            padding: '2rem',
            fontFamily: 'system-ui, -apple-system, sans-serif'
        }}>
            <div style={{
                maxWidth: '600px',
                margin: '0 auto',
                background: 'rgba(255, 255, 255, 0.1)',
                backdropFilter: 'blur(10px)',
                borderRadius: '1rem',
                padding: '2rem',
                boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.5)'
            }}>
                <h1 style={{
                    color: '#e94560',
                    textAlign: 'center',
                    marginBottom: '1.5rem',
                    fontSize: '2rem'
                }}>
                    🐳 Docker Fullstack Todo
                </h1>

                <div style={{
                    textAlign: 'center',
                    marginBottom: '1rem',
                    padding: '0.5rem',
                    background: source === 'cache' ? '#10b981' : '#3b82f6',
                    borderRadius: '0.5rem',
                    color: 'white',
                    fontSize: '0.875rem'
                }}>
                    Data source: {source || 'loading...'}
                </div>

                <form onSubmit={addTask} style={{
                    display: 'flex',
                    gap: '0.5rem',
                    marginBottom: '1.5rem'
                }}>
                    <input
                        type="text"
                        value={newTask}
                        onChange={(e) => setNewTask(e.target.value)}
                        placeholder="Добавить задачу..."
                        style={{
                            flex: 1,
                            padding: '0.75rem 1rem',
                            border: 'none',
                            borderRadius: '0.5rem',
                            fontSize: '1rem',
                            background: 'rgba(255, 255, 255, 0.9)'
                        }}
                    />
                    <button type="submit" style={{
                        padding: '0.75rem 1.5rem',
                        background: '#e94560',
                        color: 'white',
                        border: 'none',
                        borderRadius: '0.5rem',
                        cursor: 'pointer',
                        fontSize: '1rem',
                        fontWeight: 'bold'
                    }}>
                        Добавить
                    </button>
                </form>

                {loading ? (
                    <p style={{ color: 'white', textAlign: 'center' }}>Загрузка...</p>
                ) : (
                    <ul style={{ listStyle: 'none', padding: 0 }}>
                        {tasks.map(task => (
                            <li key={task.id} style={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: '0.75rem',
                                padding: '1rem',
                                background: 'rgba(255, 255, 255, 0.1)',
                                borderRadius: '0.5rem',
                                marginBottom: '0.5rem'
                            }}>
                                <input
                                    type="checkbox"
                                    checked={task.completed}
                                    onChange={() => toggleTask(task.id, task.completed)}
                                    style={{
                                        width: '1.25rem',
                                        height: '1.25rem',
                                        cursor: 'pointer'
                                    }}
                                />
                                <span style={{
                                    flex: 1,
                                    color: task.completed ? '#9ca3af' : 'white',
                                    textDecoration: task.completed ? 'line-through' : 'none'
                                }}>
                                    {task.title}
                                </span>
                                <button
                                    onClick={() => deleteTask(task.id)}
                                    style={{
                                        background: '#ef4444',
                                        color: 'white',
                                        border: 'none',
                                        borderRadius: '0.25rem',
                                        padding: '0.25rem 0.75rem',
                                        cursor: 'pointer'
                                    }}
                                >
                                    ✕
                                </button>
                            </li>
                        ))}
                    </ul>
                )}

                {tasks.length === 0 && !loading && (
                    <p style={{ color: '#9ca3af', textAlign: 'center' }}>
                        Нет задач. Добавьте первую!
                    </p>
                )}
            </div>
        </div>
    );
}

export default App;
```

**frontend/Dockerfile:**
```dockerfile
# Multi-stage build для frontend

# ========================================
# Этап 1: Dependencies
# ========================================
FROM node:18-alpine AS deps

WORKDIR /app

COPY package*.json ./
RUN npm ci

# ========================================
# Этап 2: Build
# ========================================
FROM node:18-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG REACT_APP_API_URL
ENV REACT_APP_API_URL=$REACT_APP_API_URL

RUN npm run build

# ========================================
# Этап 3: Production
# ========================================
FROM nginx:alpine AS production

RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/

COPY --from=builder /app/build /usr/share/nginx/html

RUN chown -R nginx:nginx /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
```

**docker-compose.yml (базовый):**
```yaml
version: '3.8'

services:
  # PostgreSQL Database
  db:
    image: postgres:15-alpine
    container_name: fullstack-db
    environment:
      POSTGRES_DB: ${DB_NAME:-appdb}
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-secret}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - backend
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-postgres}"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Redis Cache
  cache:
    image: redis:7-alpine
    container_name: fullstack-cache
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    networks:
      - backend
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  # Backend API
  backend:
    build:
      context: ./backend
      target: production
    container_name: fullstack-backend
    environment:
      DB_HOST: db
      DB_PORT: 5432
      DB_NAME: ${DB_NAME:-appdb}
      DB_USER: ${DB_USER:-postgres}
      DB_PASSWORD: ${DB_PASSWORD:-secret}
      REDIS_HOST: cache
      REDIS_PORT: 6379
      PORT: 3000
    networks:
      - backend
      - frontend
    depends_on:
      db:
        condition: service_healthy
      cache:
        condition: service_healthy
    restart: unless-stopped

  # Frontend
  frontend:
    build:
      context: ./frontend
      target: production
      args:
        REACT_APP_API_URL: ${API_URL:-http://localhost:3000/api}
    container_name: fullstack-frontend
    ports:
      - "80:80"
    networks:
      - frontend
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  backend:
    driver: bridge
  frontend:
    driver: bridge
```

**docker-compose.dev.yml (для разработки):**
```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      target: development
    volumes:
      - ./backend:/app
      - /app/node_modules
    ports:
      - "3000:3000"
    environment:
      NODE_ENV: development

  frontend:
    build:
      context: ./frontend
      target: deps
    command: npm start
    volumes:
      - ./frontend:/app
      - /app/node_modules
    ports:
      - "3001:3000"
    environment:
      REACT_APP_API_URL: http://localhost:3000/api

  # Adminer для работы с БД
  adminer:
    image: adminer
    container_name: fullstack-adminer
    ports:
      - "8080:8080"
    networks:
      - backend
    depends_on:
      - db
```

**.env:**
```env
# Database
DB_NAME=appdb
DB_USER=postgres
DB_PASSWORD=supersecret123

# API
API_URL=http://localhost:3000/api

# Compose
COMPOSE_PROJECT_NAME=fullstack
```

**Команды для работы:**
```bash
# Разработка
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Production
docker-compose up -d --build

# Просмотр логов
docker-compose logs -f

# Статус
docker-compose ps

# Остановка
docker-compose down

# Остановка с удалением данных
docker-compose down -v
```

---

## Задание для самостоятельной работы

### Создание многоконтейнерного приложения

**Задачи:**
1. Создайте многоконтейнерное приложение с использованием Docker Compose
2. Используйте multi-stage builds для оптимизации образов
3. Выложите проект в репозиторий gitlab.mai.ru с документацией

**Требования:**

#### 1. Docker Compose (30 баллов)

| Критерий | Баллы |
|----------|-------|
| Минимум 3 сервиса (frontend/backend/db) | 10 |
| Использование networks для изоляции | 5 |
| Использование volumes для persistence | 5 |
| Healthchecks для критичных сервисов | 5 |
| Правильное использование depends_on | 5 |

#### 2. Multi-stage Builds (20 баллов)

| Критерий | Баллы |
|----------|-------|
| Минимум один multi-stage Dockerfile | 10 |
| Размер финального образа < 200MB | 5 |
| Использование непривилегированного пользователя | 5 |

#### 3. Конфигурация (15 баллов)

| Критерий | Баллы |
|----------|-------|
| Использование .env файлов | 5 |
| .env.example для документации | 5 |
| Профили для разных окружений (dev/prod) | 5 |

#### 4. Работающее приложение (25 баллов)

| Критерий | Баллы |
|----------|-------|
| Приложение запускается одной командой | 10 |
| Все сервисы взаимодействуют корректно | 10 |
| Данные сохраняются после перезапуска | 5 |

#### 5. Документация (10 баллов)

| Критерий | Баллы |
|----------|-------|
| README с инструкциями по запуску | 5 |
| Описание архитектуры приложения | 5 |

### Критерии оценки

| Критерий | Баллы |
|----------|-------|
| Docker Compose конфигурация | 30 |
| Multi-stage builds | 20 |
| Конфигурация и переменные окружения | 15 |
| Работающее приложение | 25 |
| Документация | 10 |
| **ИТОГО** | **100** |

### Примеры проектов

1. **Todo List приложение**
   - React frontend
   - Node.js/Express backend
   - PostgreSQL или MongoDB

2. **Блог платформа**
   - Vue.js frontend
   - Python/FastAPI backend
   - PostgreSQL + Redis

3. **Чат приложение**
   - React frontend
   - Node.js + Socket.io backend
   - Redis для сообщений

4. **API для интернет-магазина**
   - Angular frontend
   - Go backend
   - PostgreSQL + Redis

---

## Контрольные вопросы

### Базовые вопросы

1. **Что такое Docker Compose и зачем он нужен?**
   - Какие проблемы он решает?

2. **Какая структура у файла docker-compose.yml?**
   - Назовите основные секции.

3. **Чем отличается `docker-compose up` от `docker-compose start`?**
   - Когда использовать каждую команду?

4. **Что делает параметр `depends_on`?**
   - В чём разница между `depends_on` и `depends_on` с условием?

5. **Как работают healthchecks в Docker Compose?**
   - Какие параметры можно настроить?

6. **Чем отличается bind mount от named volume?**
   - Когда использовать каждый тип?

### Продвинутые вопросы

7. **Что такое multi-stage build и зачем он нужен?**
   - Какие преимущества даёт?

8. **Как передать переменные окружения в docker-compose?**
   - Назовите минимум 3 способа.

9. **Зачем использовать отдельные сети в Docker Compose?**
   - Приведите пример изоляции.

10. **Как масштабировать сервисы в Docker Compose?**
    - Какие ограничения существуют?

11. **Что такое профили (profiles) в Docker Compose?**
    - Приведите пример использования.

12. **Как оптимизировать кэширование слоёв в Dockerfile?**
    - Какой порядок инструкций правильный?

13. **Чем отличается `docker-compose down` от `docker-compose down -v`?**
    - Когда использовать флаг `-v`?

14. **Как выполнить команду внутри запущенного сервиса?**
    - Какие команды для этого используются?

15. **Что такое service_healthy условие?**
    - Как его настроить?

---

## Шпаргалка команд Docker Compose

```bash
# ============ УПРАВЛЕНИЕ ============
docker-compose up                    # Запустить все сервисы
docker-compose up -d                 # Запустить в фоне
docker-compose up -d --build         # Собрать и запустить
docker-compose down                  # Остановить и удалить
docker-compose down -v               # + удалить volumes
docker-compose down --rmi all        # + удалить образы

# ============ СТАТУС ============
docker-compose ps                    # Статус сервисов
docker-compose ps -a                 # Все контейнеры
docker-compose top                   # Процессы в контейнерах

# ============ ЛОГИ ============
docker-compose logs                  # Логи всех сервисов
docker-compose logs -f               # Follow логи
docker-compose logs -f service_name  # Логи одного сервиса
docker-compose logs --tail 100       # Последние 100 строк

# ============ УПРАВЛЕНИЕ СЕРВИСАМИ ============
docker-compose start                 # Запустить остановленные
docker-compose stop                  # Остановить
docker-compose restart               # Перезапустить
docker-compose pause                 # Приостановить
docker-compose unpause               # Возобновить

# ============ СБОРКА ============
docker-compose build                 # Собрать образы
docker-compose build --no-cache      # Без кэша
docker-compose build service_name    # Собрать один сервис

# ============ ВЫПОЛНЕНИЕ КОМАНД ============
docker-compose exec service_name sh  # Войти в контейнер
docker-compose exec db psql -U user  # Команда в контейнере
docker-compose run --rm service cmd  # Одноразовый контейнер

# ============ КОНФИГУРАЦИЯ ============
docker-compose config                # Проверить конфигурацию
docker-compose config --services     # Список сервисов

# ============ ПРОФИЛИ ============
docker-compose --profile dev up -d   # С профилем
docker-compose --profile dev --profile tools up -d

# ============ МАСШТАБИРОВАНИЕ ============
docker-compose up -d --scale web=3   # 3 инстанса сервиса

# ============ ОЧИСТКА ============
docker-compose rm -f                 # Удалить контейнеры
docker system prune                  # Очистить систему
```

---

## Дополнительные материалы

### Официальная документация

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Compose File Reference](https://docs.docker.com/compose/compose-file/)
- [Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Networking in Compose](https://docs.docker.com/compose/networking/)

### Обучающие ресурсы

- [Docker Compose Tutorial](https://docs.docker.com/compose/gettingstarted/)
- [Play with Docker](https://labs.play-with-docker.com/) - Онлайн песочница
- [Awesome Docker Compose](https://github.com/docker/awesome-compose) - Примеры проектов

### Инструменты

- [Docker Desktop](https://www.docker.com/products/docker-desktop) - GUI для Docker
- [Dive](https://github.com/wagoodman/dive) - Анализ слоёв образа
- [Hadolint](https://github.com/hadolint/hadolint) - Линтер для Dockerfile
- [ctop](https://github.com/bcicen/ctop) - Top для контейнеров

### Наши материалы

- `docker-compose/` - Примеры Docker Compose конфигураций
- `examples/ml-ollama/` - **ML окружение с Ollama + Open WebUI** (локальный ChatGPT)
- Примеры multi-stage builds — см. `lab2/examples/multi-stage/`

---

## Следующие шаги

После успешного выполнения этой лабораторной работы вы будете готовы к:

**Дальнейшее изучение:**
- Kubernetes для оркестрации контейнеров
- Docker Swarm для кластеризации
- CI/CD pipelines с Docker
- Container security best practices

---

## Часто задаваемые вопросы

**В: Можно ли использовать docker compose (без дефиса)?**  
О: Да, начиная с Docker Compose v2, команда `docker compose` (без дефиса) является рекомендуемой. `docker-compose` работает для совместимости.

**В: Как посмотреть, почему контейнер не стартует?**  
О: Используйте `docker-compose logs service_name` и `docker-compose ps -a` для диагностики.

**В: Данные пропадают после `docker-compose down`?**  
О: Данные в named volumes сохраняются. Используйте флаг `-v` только если хотите удалить volumes.

**В: Как обновить только один сервис?**  
О: `docker-compose up -d --build service_name` пересоберёт и перезапустит только указанный сервис.

**В: Почему multi-stage build так важен?**  
О: Он уменьшает размер образа в 10-50 раз, что критично для безопасности и скорости деплоя.

---

**Успехов в освоении Docker Compose и Multi-stage Builds!** 🐳

---

© 2024 МАИ - Московский Авиационный Институт
