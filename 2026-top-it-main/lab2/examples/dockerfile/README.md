# Создание Docker образов с помощью Dockerfile

## Описание
В этом разделе вы научитесь создавать собственные Docker образы с использованием Dockerfile.

## Упражнение 1: Простейший Dockerfile

### Задание
Создайте простой Docker образ на базе Ubuntu с установленным curl.

### Создайте Dockerfile

```dockerfile
# Dockerfile.simple
FROM ubuntu:22.04

# Обновление пакетов и установка curl
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Установка рабочей директории
WORKDIR /app

# Команда по умолчанию
CMD ["curl", "--version"]
```

### Команды для выполнения

```bash
# Создать файл
cat > Dockerfile.simple << 'EOF'
FROM ubuntu:22.04
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
CMD ["curl", "--version"]
EOF

# Собрать образ
docker build -t my-ubuntu:v1 -f Dockerfile.simple .

# Запустить контейнер
docker run --rm my-ubuntu:v1

# Запустить с другой командой
docker run --rm my-ubuntu:v1 curl -I https://google.com
```

### Вопросы для самопроверки
1. Что делает инструкция `FROM`?
2. Зачем объединять команды через `&&` в `RUN`?
3. В чем разница между `CMD` и `RUN`?

---

## Упражнение 2: Веб-приложение на Python

### Задание
Создайте Docker образ для простого Python веб-приложения с использованием Flask.

### Создайте файлы приложения

```python
# app.py
from flask import Flask, jsonify
import os
import socket

app = Flask(__name__)

@app.route('/')
def home():
    return jsonify({
        'message': 'Hello from Docker!',
        'hostname': socket.gethostname(),
        'environment': os.environ.get('APP_ENV', 'development')
    })

@app.route('/health')
def health():
    return jsonify({'status': 'healthy'}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
```

```txt
# requirements.txt
Flask==3.0.0
Werkzeug==3.0.1
```

### Создайте Dockerfile

```dockerfile
# Dockerfile.python
FROM python:3.11-slim

# Установка рабочей директории
WORKDIR /app

# Копирование файла зависимостей
COPY requirements.txt .

# Установка зависимостей
RUN pip install --no-cache-dir -r requirements.txt

# Копирование исходного кода
COPY app.py .

# Открытие порта
EXPOSE 5000

# Переменная окружения
ENV APP_ENV=production

# Запуск приложения
CMD ["python", "app.py"]
```

### Команды для выполнения

```bash
# Создать файлы
cat > app.py << 'EOF'
from flask import Flask, jsonify
import os
import socket

app = Flask(__name__)

@app.route('/')
def home():
    return jsonify({
        'message': 'Hello from Docker!',
        'hostname': socket.gethostname(),
        'environment': os.environ.get('APP_ENV', 'development')
    })

@app.route('/health')
def health():
    return jsonify({'status': 'healthy'}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
EOF

cat > requirements.txt << 'EOF'
Flask==3.0.0
Werkzeug==3.0.1
EOF

cat > Dockerfile.python << 'EOF'
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY app.py .
EXPOSE 5000
ENV APP_ENV=production
CMD ["python", "app.py"]
EOF

# Собрать образ
docker build -t flask-app:v1 -f Dockerfile.python .

# Запустить контейнер
docker run -d -p 5000:5000 --name my-flask-app flask-app:v1

# Проверить работу
curl http://localhost:5000/
curl http://localhost:5000/health

# Остановить и удалить
docker stop my-flask-app
docker rm my-flask-app
```

### Вопросы для самопроверки
1. Зачем использовать `WORKDIR`?
2. Почему `requirements.txt` копируется отдельно от `app.py`?
3. Что делает инструкция `EXPOSE`?

---

## Упражнение 3: Node.js приложение

### Задание
Создайте Docker образ для Node.js приложения с Express.

### Создайте файлы приложения

```javascript
// server.js
const express = require('express');
const app = express();
const port = process.env.PORT || 3000;

app.get('/', (req, res) => {
  res.json({
    message: 'Node.js app running in Docker',
    version: process.version,
    timestamp: new Date().toISOString()
  });
});

app.get('/api/info', (req, res) => {
  res.json({
    uptime: process.uptime(),
    memory: process.memoryUsage(),
    pid: process.pid
  });
});

app.listen(port, '0.0.0.0', () => {
  console.log(`Server running on port ${port}`);
});
```

```json
// package.json
{
  "name": "docker-node-app",
  "version": "1.0.0",
  "description": "Simple Node.js app for Docker",
  "main": "server.js",
  "scripts": {
    "start": "node server.js"
  },
  "dependencies": {
    "express": "^4.18.2"
  }
}
```

### Создайте Dockerfile

```dockerfile
# Dockerfile.node
FROM node:18-alpine

# Создание директории приложения
WORKDIR /usr/src/app

# Копирование package.json и package-lock.json
COPY package*.json ./

# Установка зависимостей
RUN npm ci --only=production

# Копирование исходного кода
COPY server.js .

# Приложение работает на порту 3000
EXPOSE 3000

# Создание пользователя без прав root
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Переключение на непривилегированного пользователя
USER nodejs

# Запуск приложения
CMD ["node", "server.js"]
```

### Создайте .dockerignore

```
# .dockerignore
node_modules
npm-debug.log
.git
.gitignore
README.md
.env
.DS_Store
```

### Команды для выполнения

```bash
# Создать файлы
cat > server.js << 'EOF'
const express = require('express');
const app = express();
const port = process.env.PORT || 3000;

app.get('/', (req, res) => {
  res.json({
    message: 'Node.js app running in Docker',
    version: process.version,
    timestamp: new Date().toISOString()
  });
});

app.get('/api/info', (req, res) => {
  res.json({
    uptime: process.uptime(),
    memory: process.memoryUsage(),
    pid: process.pid
  });
});

app.listen(port, '0.0.0.0', () => {
  console.log(`Server running on port ${port}`);
});
EOF

cat > package.json << 'EOF'
{
  "name": "docker-node-app",
  "version": "1.0.0",
  "description": "Simple Node.js app for Docker",
  "main": "server.js",
  "scripts": {
    "start": "node server.js"
  },
  "dependencies": {
    "express": "^4.18.2"
  }
}
EOF

cat > .dockerignore << 'EOF'
node_modules
npm-debug.log
.git
.gitignore
README.md
.env
.DS_Store
EOF

cat > Dockerfile.node << 'EOF'
FROM node:18-alpine
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm ci --only=production
COPY server.js .
EXPOSE 3000
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001
USER nodejs
CMD ["node", "server.js"]
EOF

# Собрать образ
docker build -t node-app:v1 -f Dockerfile.node .

# Запустить контейнер
docker run -d -p 3000:3000 --name my-node-app node-app:v1

# Проверить работу
curl http://localhost:3000/
curl http://localhost:3000/api/info

# Остановить и удалить
docker stop my-node-app
docker rm my-node-app
```

### Вопросы для самопроверки
1. Зачем нужен файл `.dockerignore`?
2. Почему используется `npm ci` вместо `npm install`?
3. Зачем создавать непривилегированного пользователя?

---

## Упражнение 4: Использование ARG и ENV

### Задание
Создайте Dockerfile с использованием аргументов сборки и переменных окружения.

### Создайте Dockerfile

```dockerfile
# Dockerfile.args
# Аргументы сборки
ARG PYTHON_VERSION=3.11
ARG APP_DIR=/app

# Базовый образ с использованием аргумента
FROM python:${PYTHON_VERSION}-slim

# Информация о образе
LABEL maintainer="student@mai.ru"
LABEL version="1.0"
LABEL description="Python app with build arguments"

# Переменные окружения
ENV APP_HOME=${APP_DIR} \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

WORKDIR ${APP_HOME}

# Аргумент для версии приложения
ARG APP_VERSION=1.0.0
ENV VERSION=${APP_VERSION}

# Создание файла с версией
RUN echo "Version: ${VERSION}" > version.txt && \
    echo "Python: ${PYTHON_VERSION}" >> version.txt

# Команда по умолчанию
CMD ["sh", "-c", "cat version.txt && python --version"]
```

### Команды для выполнения

```bash
# Создать Dockerfile
cat > Dockerfile.args << 'EOF'
ARG PYTHON_VERSION=3.11
ARG APP_DIR=/app

FROM python:${PYTHON_VERSION}-slim

LABEL maintainer="student@mai.ru"
LABEL version="1.0"
LABEL description="Python app with build arguments"

ENV APP_HOME=${APP_DIR} \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

WORKDIR ${APP_HOME}

ARG APP_VERSION=1.0.0
ENV VERSION=${APP_VERSION}

RUN echo "Version: ${VERSION}" > version.txt && \
    echo "Python: ${PYTHON_VERSION}" >> version.txt

CMD ["sh", "-c", "cat version.txt && python --version"]
EOF

# Собрать образ с аргументами по умолчанию
docker build -t python-arg-app:default -f Dockerfile.args .

# Собрать образ с кастомными аргументами
docker build \
  --build-arg PYTHON_VERSION=3.10 \
  --build-arg APP_VERSION=2.0.0 \
  --build-arg APP_DIR=/myapp \
  -t python-arg-app:custom \
  -f Dockerfile.args .

# Запустить контейнеры
docker run --rm python-arg-app:default
docker run --rm python-arg-app:custom

# Переопределить переменную окружения при запуске
docker run --rm -e VERSION=3.0.0 python-arg-app:custom sh -c "echo Version: \$VERSION"
```

### Вопросы для самопроверки
1. В чем разница между `ARG` и `ENV`?
2. Можно ли изменить `ARG` при запуске контейнера?
3. Какие переменные окружения можно переопределить при `docker run`?

---

## Упражнение 5: ENTRYPOINT vs CMD

### Задание
Поймите разницу между ENTRYPOINT и CMD.

### Пример 1: Только CMD

```dockerfile
# Dockerfile.cmd
FROM alpine:3.18
CMD ["echo", "Hello from CMD"]
```

### Пример 2: Только ENTRYPOINT

```dockerfile
# Dockerfile.entrypoint
FROM alpine:3.18
ENTRYPOINT ["echo", "Hello from ENTRYPOINT"]
```

### Пример 3: ENTRYPOINT + CMD

```dockerfile
# Dockerfile.both
FROM alpine:3.18
ENTRYPOINT ["echo"]
CMD ["Hello from both!"]
```

### Команды для выполнения

```bash
# Создать Dockerfiles
cat > Dockerfile.cmd << 'EOF'
FROM alpine:3.18
CMD ["echo", "Hello from CMD"]
EOF

cat > Dockerfile.entrypoint << 'EOF'
FROM alpine:3.18
ENTRYPOINT ["echo", "Hello from ENTRYPOINT"]
EOF

cat > Dockerfile.both << 'EOF'
FROM alpine:3.18
ENTRYPOINT ["echo"]
CMD ["Hello from both!"]
EOF

# Собрать образы
docker build -t test-cmd -f Dockerfile.cmd .
docker build -t test-entrypoint -f Dockerfile.entrypoint .
docker build -t test-both -f Dockerfile.both .

# Тестирование CMD
docker run --rm test-cmd
docker run --rm test-cmd echo "Override CMD"  # Полностью заменяет CMD

# Тестирование ENTRYPOINT
docker run --rm test-entrypoint
docker run --rm test-entrypoint "Additional argument"  # Добавляет аргумент

# Тестирование ENTRYPOINT + CMD
docker run --rm test-both
docker run --rm test-both "Custom message"  # Заменяет только CMD
```

### Вопросы для самопроверки
1. Что происходит с `CMD` при передаче аргументов в `docker run`?
2. Можно ли переопределить `ENTRYPOINT` при запуске?
3. Когда лучше использовать `ENTRYPOINT`, а когда `CMD`?

---

## Упражнение 6: Многострочный Dockerfile для веб-приложения

### Задание
Создайте полноценный Dockerfile с использованием всех изученных инструкций.

### Создайте комплексное приложение

```python
# main.py
from flask import Flask, render_template_string
import os

app = Flask(__name__)

HTML_TEMPLATE = """
<!DOCTYPE html>
<html>
<head>
    <title>Docker Demo</title>
    <style>
        body { font-family: Arial; margin: 40px; background: #f0f0f0; }
        .container { background: white; padding: 20px; border-radius: 8px; }
        h1 { color: #2196F3; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🐳 Docker Demo Application</h1>
        <p><strong>Version:</strong> {{ version }}</p>
        <p><strong>Environment:</strong> {{ env }}</p>
        <p><strong>Port:</strong> {{ port }}</p>
    </div>
</body>
</html>
"""

@app.route('/')
def home():
    return render_template_string(
        HTML_TEMPLATE,
        version=os.environ.get('APP_VERSION', 'unknown'),
        env=os.environ.get('APP_ENV', 'development'),
        port=os.environ.get('PORT', '8000')
    )

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 8000))
    app.run(host='0.0.0.0', port=port)
```

### Создайте Dockerfile

```dockerfile
# Dockerfile.complete
FROM python:3.11-slim AS base

# Метаданные
LABEL maintainer="student@mai.ru" \
      version="1.0.0" \
      description="Complete Flask application"

# Аргументы сборки
ARG APP_USER=appuser
ARG APP_UID=1000
ARG APP_GID=1000

# Переменные окружения
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# Установка системных зависимостей
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Создание непривилегированного пользователя
RUN groupadd -g ${APP_GID} ${APP_USER} && \
    useradd -m -u ${APP_UID} -g ${APP_GID} -s /bin/bash ${APP_USER}

# Установка рабочей директории
WORKDIR /app

# Копирование и установка зависимостей
COPY --chown=${APP_USER}:${APP_USER} requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Копирование приложения
COPY --chown=${APP_USER}:${APP_USER} main.py .

# Переключение на непривилегированного пользователя
USER ${APP_USER}

# Открытие порта
EXPOSE 8000

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8000/ || exit 1

# Переменные окружения приложения
ENV APP_VERSION=1.0.0 \
    APP_ENV=production \
    PORT=8000

# Точка входа и команда по умолчанию
ENTRYPOINT ["python"]
CMD ["main.py"]
```

### Команды для выполнения

```bash
# Создать файлы приложения
# (используйте файлы выше)

# Собрать образ
docker build -t complete-app:v1 -f Dockerfile.complete .

# Запустить с различными конфигурациями
docker run -d -p 8000:8000 --name app1 complete-app:v1

docker run -d -p 8001:8001 \
  -e PORT=8001 \
  -e APP_ENV=staging \
  --name app2 \
  complete-app:v1

# Проверить healthcheck
docker inspect app1 | grep -A 10 Health

# Просмотреть приложение
curl http://localhost:8000/
curl http://localhost:8001/

# Очистка
docker stop app1 app2
docker rm app1 app2
```

---

## Итоговое задание

Создайте Dockerfile для приложения на вашем любимом языке программирования, которое:

1. Использует многоэтапную сборку (если применимо)
2. Устанавливает зависимости из файла
3. Копирует исходный код
4. Работает от имени непривилегированного пользователя
5. Имеет healthcheck
6. Использует переменные окружения
7. Имеет метаданные (LABEL)
8. Оптимизировано по размеру

## Чек-лист освоенных навыков

- [ ] Создание базового Dockerfile
- [ ] Использование инструкций FROM, RUN, COPY, WORKDIR
- [ ] Работа с CMD и ENTRYPOINT
- [ ] Использование ARG и ENV
- [ ] Добавление LABEL для метаданных
- [ ] Оптимизация слоев образа
- [ ] Создание .dockerignore файла
- [ ] Работа с непривилегированными пользователями
- [ ] Добавление HEALTHCHECK
- [ ] Использование EXPOSE для документирования портов

