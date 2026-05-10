# Docker Compose - Оркестрация контейнеров

## Описание
Docker Compose позволяет определять и запускать многоконтейнерные приложения с помощью YAML файлов.

## Упражнение 1: Простой docker-compose.yml

### Задание
Создайте простую конфигурацию для веб-приложения.

### Создайте docker-compose.yml

```yaml
# docker-compose-simple.yml
version: '3.8'

services:
  web:
    image: nginx:alpine
    ports:
      - "8080:80"
    container_name: simple-web
```

### Команды для выполнения

```bash
# Создать файл
cat > docker-compose-simple.yml << 'EOF'
version: '3.8'

services:
  web:
    image: nginx:alpine
    ports:
      - "8080:80"
    container_name: simple-web
EOF

# Запустить сервисы
docker-compose -f docker-compose-simple.yml up -d

# Проверить статус
docker-compose -f docker-compose-simple.yml ps

# Просмотреть логи
docker-compose -f docker-compose-simple.yml logs

# Остановить сервисы
docker-compose -f docker-compose-simple.yml down
```

### Вопросы для самопроверки
1. Что означает флаг `-d` в команде `docker-compose up`?
2. Зачем указывать `container_name`?
3. Как просмотреть логи конкретного сервиса?

---

## Упражнение 2: Веб-приложение с базой данных

### Задание
Создайте конфигурацию для Flask приложения с PostgreSQL.

### Создайте приложение

```python
# app.py
from flask import Flask, jsonify
import psycopg2
import os
import time

app = Flask(__name__)

def get_db_connection():
    max_retries = 5
    for i in range(max_retries):
        try:
            conn = psycopg2.connect(
                host=os.environ.get('DB_HOST', 'db'),
                database=os.environ.get('DB_NAME', 'testdb'),
                user=os.environ.get('DB_USER', 'postgres'),
                password=os.environ.get('DB_PASSWORD', 'secret')
            )
            return conn
        except psycopg2.OperationalError:
            if i < max_retries - 1:
                time.sleep(2)
            else:
                raise

@app.route('/')
def home():
    return jsonify({'message': 'Flask + PostgreSQL in Docker!'})

@app.route('/db-test')
def db_test():
    try:
        conn = get_db_connection()
        cur = conn.cursor()
        cur.execute('SELECT version();')
        db_version = cur.fetchone()
        cur.close()
        conn.close()
        return jsonify({
            'status': 'success',
            'database': db_version[0]
        })
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        }), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

### Создайте Dockerfile

```dockerfile
# Dockerfile
FROM python:3.11-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y postgresql-client && \
    apt-get clean

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

CMD ["python", "app.py"]
```

### Создайте requirements.txt

```txt
Flask==3.0.0
psycopg2-binary==2.9.9
```

### Создайте docker-compose.yml

```yaml
# docker-compose-db.yml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    container_name: postgres-db
    environment:
      POSTGRES_DB: testdb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  web:
    build: .
    container_name: flask-app
    environment:
      DB_HOST: db
      DB_NAME: testdb
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

### Команды для выполнения

```bash
# Создать все файлы (используйте содержимое выше)

# Собрать и запустить
docker-compose -f docker-compose-db.yml up -d --build

# Проверить статус
docker-compose -f docker-compose-db.yml ps

# Проверить логи
docker-compose -f docker-compose-db.yml logs -f web

# Проверить приложение
curl http://localhost:5000/
curl http://localhost:5000/db-test

# Остановить и удалить
docker-compose -f docker-compose-db.yml down -v
```

### Вопросы для самопроверки
1. Что делает `depends_on` с условием `service_healthy`?
2. Зачем нужен `volumes` для базы данных?
3. Как работает `healthcheck` в PostgreSQL сервисе?

---

## Упражнение 3: Полный стек приложения (Frontend + Backend + Database)

### Задание
Создайте полноценное приложение с frontend, backend и базой данных.

### Backend (Node.js + Express)

```javascript
// backend/server.js
const express = require('express');
const { Pool } = require('pg');
const cors = require('cors');

const app = express();
app.use(cors());
app.use(express.json());

const pool = new Pool({
  host: process.env.DB_HOST || 'db',
  port: 5432,
  database: process.env.DB_NAME || 'appdb',
  user: process.env.DB_USER || 'postgres',
  password: process.env.DB_PASSWORD || 'secret',
});

// Создание таблицы при старте
pool.query(`
  CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  )
`).catch(err => console.error('Error creating table:', err));

app.get('/api/messages', async (req, res) => {
  try {
    const result = await pool.query('SELECT * FROM messages ORDER BY created_at DESC');
    res.json(result.rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

app.post('/api/messages', async (req, res) => {
  try {
    const { content } = req.body;
    const result = await pool.query(
      'INSERT INTO messages (content) VALUES ($1) RETURNING *',
      [content]
    );
    res.json(result.rows[0]);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, '0.0.0.0', () => {
  console.log(`Backend server running on port ${PORT}`);
});
```

```json
// backend/package.json
{
  "name": "backend",
  "version": "1.0.0",
  "main": "server.js",
  "dependencies": {
    "express": "^4.18.2",
    "pg": "^8.11.3",
    "cors": "^2.8.5"
  }
}
```

```dockerfile
# backend/Dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install --production
COPY server.js .
EXPOSE 3000
CMD ["node", "server.js"]
```

### Frontend (HTML + JavaScript)

```html
<!-- frontend/index.html -->
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Docker Full Stack App</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            padding: 30px;
        }
        h1 {
            color: #667eea;
            margin-bottom: 30px;
            text-align: center;
        }
        .input-group {
            display: flex;
            gap: 10px;
            margin-bottom: 30px;
        }
        input {
            flex: 1;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 6px;
            font-size: 16px;
        }
        button {
            padding: 12px 30px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 16px;
        }
        button:hover { background: #5568d3; }
        .messages {
            display: flex;
            flex-direction: column;
            gap: 15px;
        }
        .message {
            background: #f5f5f5;
            padding: 15px;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        .message-content {
            font-size: 16px;
            margin-bottom: 8px;
        }
        .message-time {
            font-size: 12px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🐳 Docker Full Stack App</h1>
        <div class="input-group">
            <input type="text" id="messageInput" placeholder="Введите сообщение...">
            <button onclick="addMessage()">Отправить</button>
        </div>
        <div id="messages" class="messages"></div>
    </div>

    <script>
        const API_URL = 'http://localhost:3000/api';

        async function loadMessages() {
            try {
                const response = await fetch(`${API_URL}/messages`);
                const messages = await response.json();
                displayMessages(messages);
            } catch (error) {
                console.error('Error loading messages:', error);
            }
        }

        async function addMessage() {
            const input = document.getElementById('messageInput');
            const content = input.value.trim();
            
            if (!content) return;

            try {
                await fetch(`${API_URL}/messages`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ content })
                });
                input.value = '';
                loadMessages();
            } catch (error) {
                console.error('Error adding message:', error);
            }
        }

        function displayMessages(messages) {
            const container = document.getElementById('messages');
            container.innerHTML = messages.map(msg => `
                <div class="message">
                    <div class="message-content">${msg.content}</div>
                    <div class="message-time">${new Date(msg.created_at).toLocaleString('ru-RU')}</div>
                </div>
            `).join('');
        }

        // Загрузка сообщений при старте
        loadMessages();
        
        // Обновление каждые 3 секунды
        setInterval(loadMessages, 3000);

        // Enter для отправки
        document.getElementById('messageInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') addMessage();
        });
    </script>
</body>
</html>
```

```dockerfile
# frontend/Dockerfile
FROM nginx:alpine
COPY index.html /usr/share/nginx/html/
EXPOSE 80
```

### Docker Compose

```yaml
# docker-compose-fullstack.yml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    container_name: fullstack-db
    environment:
      POSTGRES_DB: appdb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    volumes:
      - db_data:/var/lib/postgresql/data
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 3s
      retries: 5

  backend:
    build: ./backend
    container_name: fullstack-backend
    environment:
      DB_HOST: db
      DB_NAME: appdb
      DB_USER: postgres
      DB_PASSWORD: secret
      PORT: 3000
    ports:
      - "3000:3000"
    depends_on:
      db:
        condition: service_healthy
    networks:
      - app-network
    restart: unless-stopped

  frontend:
    build: ./frontend
    container_name: fullstack-frontend
    ports:
      - "8080:80"
    depends_on:
      - backend
    networks:
      - app-network

volumes:
  db_data:

networks:
  app-network:
    driver: bridge
```

### Команды для выполнения

```bash
# Создать структуру
mkdir -p fullstack/{backend,frontend}
cd fullstack

# Создать файлы (используйте содержимое выше)

# Запустить весь стек
docker-compose -f docker-compose-fullstack.yml up -d --build

# Проверить статус всех сервисов
docker-compose -f docker-compose-fullstack.yml ps

# Просмотреть логи
docker-compose -f docker-compose-fullstack.yml logs -f

# Открыть в браузере http://localhost:8080

# Остановить и удалить все
docker-compose -f docker-compose-fullstack.yml down -v
```

---

## Упражнение 4: Расширенные возможности Docker Compose

### Задание
Изучите расширенные функции Docker Compose.

### Создайте продвинутый docker-compose.yml

```yaml
# docker-compose-advanced.yml
version: '3.8'

services:
  # Redis для кэширования
  cache:
    image: redis:7-alpine
    container_name: app-cache
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    networks:
      - backend-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  # PostgreSQL база данных
  database:
    image: postgres:15-alpine
    container_name: app-database
    environment:
      POSTGRES_DB: ${DB_NAME:-appdb}
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-secret}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init-db:/docker-entrypoint-initdb.d
    networks:
      - backend-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-postgres}"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Backend API
  api:
    build:
      context: ./api
      dockerfile: Dockerfile
      args:
        NODE_ENV: production
    container_name: app-api
    environment:
      NODE_ENV: production
      DB_HOST: database
      DB_PORT: 5432
      REDIS_HOST: cache
      REDIS_PORT: 6379
    ports:
      - "${API_PORT:-4000}:4000"
    depends_on:
      database:
        condition: service_healthy
      cache:
        condition: service_healthy
    networks:
      - backend-network
      - frontend-network
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M

  # Nginx как reverse proxy
  nginx:
    image: nginx:alpine
    container_name: app-nginx
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./static:/usr/share/nginx/html:ro
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - api
    networks:
      - frontend-network
    restart: always

  # Мониторинг с Prometheus (опционально)
  prometheus:
    image: prom/prometheus:latest
    container_name: app-prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus_data:/prometheus
    ports:
      - "9090:9090"
    networks:
      - monitoring-network
    profiles:
      - monitoring

  # Grafana для визуализации (опционально)
  grafana:
    image: grafana/grafana:latest
    container_name: app-grafana
    environment:
      GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_PASSWORD:-admin}
    volumes:
      - grafana_data:/var/lib/grafana
    ports:
      - "3001:3000"
    depends_on:
      - prometheus
    networks:
      - monitoring-network
    profiles:
      - monitoring

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local
  prometheus_data:
    driver: local
  grafana_data:
    driver: local

networks:
  backend-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
  frontend-network:
    driver: bridge
  monitoring-network:
    driver: bridge
```

### Создайте .env файл

```env
# .env
DB_NAME=appdb
DB_USER=postgres
DB_PASSWORD=supersecret123
API_PORT=4000
GRAFANA_PASSWORD=secure_password
```

### Команды для выполнения

```bash
# Проверить конфигурацию
docker-compose -f docker-compose-advanced.yml config

# Запустить только основные сервисы
docker-compose -f docker-compose-advanced.yml up -d

# Запустить с мониторингом
docker-compose -f docker-compose-advanced.yml --profile monitoring up -d

# Масштабировать сервис
docker-compose -f docker-compose-advanced.yml up -d --scale api=3

# Просмотреть использование ресурсов
docker stats

# Выполнить команду в сервисе
docker-compose -f docker-compose-advanced.yml exec database psql -U postgres -d appdb

# Просмотреть логи конкретного сервиса
docker-compose -f docker-compose-advanced.yml logs -f api

# Перезапустить сервис
docker-compose -f docker-compose-advanced.yml restart api

# Остановить все
docker-compose -f docker-compose-advanced.yml down -v
```

---

## Полезные команды Docker Compose

```bash
# Сборка без кэша
docker-compose build --no-cache

# Пересоздать контейнеры
docker-compose up -d --force-recreate

# Удалить orphan контейнеры
docker-compose down --remove-orphans

# Просмотр переменных окружения сервиса
docker-compose config | grep -A 10 environment

# Приостановить/возобновить сервисы
docker-compose pause
docker-compose unpause

# Выполнить одноразовую команду
docker-compose run --rm web python manage.py migrate

# Просмотр портов
docker-compose port web 80
```

## Итоговое задание

Создайте полноценное многоконтейнерное приложение, которое включает:
- Frontend (React/Vue/Angular или статический HTML)
- Backend API (Node.js/Python/Go)
- База данных (PostgreSQL/MySQL/MongoDB)
- Кэш (Redis)
- Reverse proxy (Nginx)

Требования:
- Все сервисы должны быть в отдельных сетях по назначению
- База данных должна сохранять данные в volume
- Healthchecks для критичных сервисов
- Использование .env для конфигурации
- Ограничения ресурсов для сервисов

## Чек-лист освоенных навыков

- [ ] Создание базового docker-compose.yml
- [ ] Работа с несколькими сервисами
- [ ] Использование volumes для постоянного хранения
- [ ] Настройка сетей между контейнерами
- [ ] Использование depends_on с условиями
- [ ] Добавление healthchecks
- [ ] Работа с переменными окружения и .env файлами
- [ ] Использование profiles для опциональных сервисов
- [ ] Ограничение ресурсов контейнеров
- [ ] Масштабирование сервисов

