# Docker Quick Start Guide

Быстрое руководство для начала работы с Docker.

## За 5 минут

### 1. Установка

**macOS:**
```bash
brew install --cask docker
# Или скачать Docker Desktop с docker.com
```

**Windows:**
- Скачать Docker Desktop с docker.com
- Требуется WSL 2

**Linux:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### 2. Проверка

```bash
docker --version
docker run hello-world
```

### 3. Ваш первый контейнер

```bash
# Запустить веб-сервер
docker run -d -p 8080:80 nginx

# Открыть в браузере
# http://localhost:8080

# Остановить
docker stop $(docker ps -q)
```

## За 15 минут

### Создание простого приложения

**1. Создайте файлы:**

```python
# app.py
from flask import Flask, jsonify

app = Flask(__name__)

@app.route('/')
def home():
    return jsonify({'message': 'Hello Docker!', 'status': 'running'})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

```txt
# requirements.txt
Flask==3.0.0
```

```dockerfile
# Dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 5000

CMD ["python", "app.py"]
```

**2. Соберите и запустите:**

```bash
# Собрать образ
docker build -t myapp .

# Запустить контейнер
docker run -d -p 5000:5000 --name myapp myapp

# Проверить
curl http://localhost:5000

# Остановить
docker stop myapp
docker rm myapp
```

## За 30 минут

### Многоконтейнерное приложение

**1. Структура проекта:**
```
myproject/
├── docker-compose.yml
├── backend/
│   ├── Dockerfile
│   ├── app.py
│   └── requirements.txt
└── frontend/
    ├── Dockerfile
    └── index.html
```

**2. Backend (backend/app.py):**

```python
from flask import Flask, jsonify
from flask_cors import CORS
import redis

app = Flask(__name__)
CORS(app)

cache = redis.Redis(host='cache', port=6379)

@app.route('/')
def home():
    count = cache.incr('hits')
    return jsonify({
        'message': 'Hello from Docker!',
        'visits': count
    })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

**3. Backend Dockerfile:**

```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 5000

CMD ["python", "app.py"]
```

**4. Frontend (frontend/index.html):**

```html
<!DOCTYPE html>
<html>
<head>
    <title>Docker App</title>
    <style>
        body {
            font-family: Arial;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .container {
            background: white;
            padding: 40px;
            border-radius: 10px;
            text-align: center;
            box-shadow: 0 10px 40px rgba(0,0,0,0.3);
        }
        button {
            padding: 10px 20px;
            font-size: 16px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        }
        button:hover { background: #5568d3; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🐳 Docker Multi-Container App</h1>
        <p id="message">Loading...</p>
        <p>Visits: <span id="visits">0</span></p>
        <button onclick="loadData()">Refresh</button>
    </div>

    <script>
        async function loadData() {
            const response = await fetch('http://localhost:5000/');
            const data = await response.json();
            document.getElementById('message').textContent = data.message;
            document.getElementById('visits').textContent = data.visits;
        }
        loadData();
    </script>
</body>
</html>
```

**5. Frontend Dockerfile:**

```dockerfile
FROM nginx:alpine
COPY index.html /usr/share/nginx/html/
```

**6. docker-compose.yml:**

```yaml
version: '3.8'

services:
  frontend:
    build: ./frontend
    ports:
      - "8080:80"
    depends_on:
      - backend

  backend:
    build: ./backend
    ports:
      - "5000:5000"
    depends_on:
      - cache
    environment:
      - REDIS_HOST=cache

  cache:
    image: redis:alpine
    volumes:
      - cache-data:/data

volumes:
  cache-data:
```

**7. Backend requirements.txt:**

```txt
Flask==3.0.0
flask-cors==4.0.0
redis==5.0.1
```

**8. Запуск:**

```bash
# Запустить все сервисы
docker-compose up -d

# Открыть в браузере
# http://localhost:8080

# Проверить логи
docker-compose logs -f

# Остановить
docker-compose down
```

## Основные команды

### Образы

```bash
docker pull nginx                    # Загрузить образ
docker images                        # Список образов
docker build -t myapp .             # Собрать образ
docker rmi myapp                     # Удалить образ
```

### Контейнеры

```bash
docker run -d nginx                  # Запустить контейнер
docker ps                            # Список контейнеров
docker stop container_id             # Остановить
docker rm container_id               # Удалить
docker logs container_id             # Логи
docker exec -it container_id bash    # Войти внутрь
```

### Docker Compose

```bash
docker-compose up -d                 # Запустить
docker-compose ps                    # Статус
docker-compose logs -f               # Логи
docker-compose down                  # Остановить
```

### Очистка

```bash
docker system prune                  # Очистить все
docker image prune -a                # Удалить образы
docker volume prune                  # Удалить volumes
```

## Что дальше?

1. **Изучите основы:**
   - [Docker Cheat Sheet](./docker-cheatsheet.md)
   - [Dockerfile Best Practices](./dockerfile-best-practices.md)
   - [Docker Compose Cheat Sheet](./docker-compose-cheatsheet.md)

2. **Практикуйтесь:**
   - [Базовые команды](../examples/basic-commands/README.md)
   - [Создание Dockerfile](../examples/dockerfile/README.md)
   - [Docker Compose](../examples/docker-compose/README.md)

3. **Изучите продвинутые темы:**
   - [Multi-Stage Builds](../examples/multi-stage/README.md)
   - [Volumes](../examples/volumes/README.md)
   - [Networks](../examples/networks/README.md)
   - [Security](../examples/security/README.md)

4. **Решайте проблемы:**
   - [Troubleshooting](../examples/troubleshooting/README.md)
   - [FAQ](./faq.md)

## Полезные ресурсы

- [Docker Documentation](https://docs.docker.com/)
- [Docker Hub](https://hub.docker.com/)
- [Play with Docker](https://labs.play-with-docker.com/) - Онлайн песочница
- [Docker Forum](https://forums.docker.com/)

---

**Поздравляем!** Вы сделали первые шаги с Docker 🎉

Переходите к [основной лабораторной работе](../README.md) для более глубокого изучения.

