# Docker Security - Безопасность контейнеров

## Описание
В этом разделе рассматриваются best practices по безопасности Docker контейнеров и образов.

## Упражнение 1: Сканирование образов на уязвимости

### Задание
Научитесь сканировать Docker образы на наличие уязвимостей.

### Использование Docker Scout (встроено в Docker Desktop)

```bash
# Сканирование локального образа
docker scout cves nginx:latest

# Сканирование с подробным выводом
docker scout cves --format json nginx:latest > scan-report.json

# Сравнение двух образов
docker scout compare --to nginx:latest nginx:alpine

# Рекомендации по исправлению
docker scout recommendations nginx:latest

# Просмотр CVE в образе
docker scout cves --only-severity critical,high nginx:latest
```

### Использование Trivy

```bash
# Установка Trivy (macOS)
brew install trivy

# Сканирование образа
trivy image nginx:latest

# Сканирование с фильтрацией по критичности
trivy image --severity HIGH,CRITICAL nginx:latest

# Сканирование локального Dockerfile
trivy config Dockerfile

# Вывод в JSON
trivy image --format json nginx:latest > trivy-report.json

# Игнорирование неисправленных уязвимостей
trivy image --ignore-unfixed nginx:latest
```

### Создание безопасного образа

```dockerfile
# Dockerfile.secure
# Используем конкретную версию (не latest)
FROM python:3.11.6-slim-bookworm

# Обновляем пакеты и устанавливаем security обновления
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Создаем непривилегированного пользователя
RUN groupadd -r appuser -g 1000 && \
    useradd -r -u 1000 -g appuser appuser

WORKDIR /app

# Копируем с правильными правами
COPY --chown=appuser:appuser requirements.txt .

# Устанавливаем зависимости
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

COPY --chown=appuser:appuser . .

# Переключаемся на непривилегированного пользователя
USER appuser

# Ограничиваем возможности
# Использовать при docker run: --cap-drop=ALL --cap-add=NET_BIND_SERVICE

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8000/health || exit 1

CMD ["python", "app.py"]
```

### Вопросы для самопроверки
1. Почему важно использовать конкретные версии образов?
2. Что такое CVE (Common Vulnerabilities and Exposures)?
3. Как часто нужно сканировать образы?

---

## Упражнение 2: Запуск контейнеров от непривилегированного пользователя

### Задание
Создайте контейнеры, работающие от имени обычного пользователя.

### Неправильный подход (запуск от root)

```dockerfile
# Dockerfile.root
FROM nginx:alpine

# ❌ ПЛОХО: По умолчанию nginx работает от root
COPY index.html /usr/share/nginx/html/

CMD ["nginx", "-g", "daemon off;"]
```

### Правильный подход (запуск от пользователя)

```dockerfile
# Dockerfile.nonroot
FROM nginx:alpine

# Создание непривилегированного пользователя
RUN addgroup -g 1001 -S nginxuser && \
    adduser -S nginxuser -u 1001 -G nginxuser

# Изменение владельца файлов
RUN chown -R nginxuser:nginxuser /var/cache/nginx && \
    chown -R nginxuser:nginxuser /var/log/nginx && \
    chown -R nginxuser:nginxuser /etc/nginx/conf.d && \
    touch /var/run/nginx.pid && \
    chown -R nginxuser:nginxuser /var/run/nginx.pid

# Копирование конфигурации для работы без root
COPY nginx-nonroot.conf /etc/nginx/nginx.conf
COPY index.html /usr/share/nginx/html/

USER nginxuser

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
```

```nginx
# nginx-nonroot.conf
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    # PID файл в доступной директории
    pid /tmp/nginx.pid;
    
    server {
        listen 8080;  # Непривилегированный порт
        server_name localhost;
        root /usr/share/nginx/html;
        index index.html;
        
        location / {
            try_files $uri $uri/ =404;
        }
    }
}
```

### Команды для выполнения

```bash
# Создать файлы

# Собрать образы
docker build -t nginx:root -f Dockerfile.root .
docker build -t nginx:nonroot -f Dockerfile.nonroot .

# Запустить и проверить пользователя
docker run -d --name nginx-root nginx:root
docker run -d --name nginx-nonroot nginx:nonroot

# Проверить от какого пользователя работает процесс
docker exec nginx-root ps aux
docker exec nginx-nonroot ps aux

# Попробовать записать файл (должно быть запрещено для nonroot)
docker exec nginx-root touch /test.txt && echo "✓ Root может писать"
docker exec nginx-nonroot touch /test.txt 2>&1 | grep -q "Permission denied" && echo "✓ User не может писать в корень"

# Очистка
docker stop nginx-root nginx-nonroot
docker rm nginx-root nginx-nonroot
```

### Вопросы для самопроверки
1. Почему запуск от root опасен?
2. Какие директории должны быть доступны для записи?
3. Можно ли использовать порты < 1024 без root?

---

## Упражнение 3: Ограничение возможностей контейнера

### Задание
Используйте Linux capabilities и security опции.

### Команды для выполнения

```bash
# 1. Запуск с минимальными capabilities
docker run -d \
  --name secure-app \
  --cap-drop=ALL \
  --cap-add=NET_BIND_SERVICE \
  --read-only \
  --tmpfs /tmp \
  --security-opt=no-new-privileges:true \
  nginx:alpine

# 2. Просмотр capabilities контейнера
docker exec secure-app sh -c "apk add --no-cache libcap && capsh --print"

# 3. Запуск с ограничениями ресурсов
docker run -d \
  --name limited-app \
  --memory="256m" \
  --memory-swap="256m" \
  --cpus="0.5" \
  --pids-limit 100 \
  --ulimit nofile=1024:1024 \
  nginx:alpine

# 4. Проверка ограничений
docker stats limited-app --no-stream

# 5. Запуск с read-only файловой системой
docker run -d \
  --name readonly-app \
  --read-only \
  --tmpfs /tmp:rw,noexec,nosuid,size=64m \
  --tmpfs /run:rw,noexec,nosuid,size=64m \
  nginx:alpine

# 6. Проверка, что ФС read-only
docker exec readonly-app touch /test.txt 2>&1 | grep -q "Read-only file system" && echo "✓ ФС защищена"

# 7. Использование AppArmor/SELinux
docker run -d \
  --name selinux-app \
  --security-opt label=level:s0:c100,c200 \
  nginx:alpine

# 8. Запрет повышения привилегий
docker run -d \
  --name no-new-privs \
  --security-opt=no-new-privileges:true \
  nginx:alpine

# Очистка
docker stop secure-app limited-app readonly-app selinux-app no-new-privs
docker rm secure-app limited-app readonly-app selinux-app no-new-privs
```

### Linux Capabilities

```bash
# Полезные capabilities:
# NET_BIND_SERVICE - привязка к портам < 1024
# SETUID, SETGID - изменение UID/GID
# NET_ADMIN - сетевое администрирование
# SYS_TIME - изменение системного времени
# CHOWN - изменение владельца файлов

# ✅ Минимальный набор для веб-приложения
--cap-drop=ALL --cap-add=NET_BIND_SERVICE

# ✅ Для базы данных
--cap-drop=ALL --cap-add=SETUID --cap-add=SETGID --cap-add=DAC_OVERRIDE

# ❌ Опасные capabilities (избегать!)
--cap-add=SYS_ADMIN    # Почти как root
--cap-add=SYS_MODULE   # Загрузка kernel модулей
--privileged           # Все capabilities
```

### Вопросы для самопроверки
1. Что такое Linux capabilities?
2. Зачем использовать `--read-only`?
3. Что делает `--security-opt=no-new-privileges:true`?

---

## Упражнение 4: Секреты и чувствительные данные

### Задание
Научитесь безопасно работать с секретами.

### ❌ НЕПРАВИЛЬНЫЕ способы

```dockerfile
# Dockerfile.bad-secrets

# ❌ ПЛОХО: Пароли в ENV
ENV DATABASE_PASSWORD=supersecret123

# ❌ ПЛОХО: Секреты в слоях
RUN echo "api_key=12345" > /app/config.txt

# ❌ ПЛОХО: Копирование .env файлов
COPY .env /app/
```

### ✅ ПРАВИЛЬНЫЕ способы

#### 1. Docker Secrets (для Docker Swarm)

```bash
# Создать секрет
echo "supersecret123" | docker secret create db_password -

# Использовать в сервисе
docker service create \
  --name myapp \
  --secret db_password \
  myapp:latest

# В контейнере секрет доступен в /run/secrets/db_password
```

#### 2. Передача через переменные окружения (для docker-compose)

```yaml
# docker-compose.secrets.yml
version: '3.8'

services:
  app:
    image: myapp:latest
    environment:
      - DB_PASSWORD=${DB_PASSWORD}
    # Или из файла:
    env_file:
      - .env.production
```

```bash
# .env.production (добавить в .gitignore!)
DB_PASSWORD=secret_from_file
API_KEY=key_from_file
```

#### 3. Build secrets (для Dockerfile)

```dockerfile
# Dockerfile.build-secrets
FROM alpine:3.18

WORKDIR /app

# Использование build secret
RUN --mount=type=secret,id=npmrc,target=/root/.npmrc \
    npm install

# Секрет не попадает в финальный образ!
```

```bash
# Сборка с секретом
docker build \
  --secret id=npmrc,src=$HOME/.npmrc \
  -t app:latest \
  -f Dockerfile.build-secrets .
```

#### 4. Использование внешних secret managers

```python
# app.py
import boto3
import os

def get_secret():
    client = boto3.client('secretsmanager')
    response = client.get_secret_value(SecretId='prod/db/password')
    return response['SecretString']

# Или использовать vault, AWS Secrets Manager, Azure Key Vault и т.д.
```

### Команды для проверки секретов в образах

```bash
# Проверить историю образа на наличие секретов
docker history myapp:latest

# Проверить переменные окружения
docker inspect myapp:latest --format='{{.Config.Env}}'

# Извлечь и просканировать файловую систему
docker save myapp:latest -o myapp.tar
tar -xf myapp.tar
grep -r "password\|secret\|key" .
```

### Best practices

```bash
# ✅ Используйте .dockerignore
cat > .dockerignore << 'EOF'
.env
.env.*
*.pem
*.key
*.crt
.git
secrets/
EOF

# ✅ Никогда не коммитьте секреты в Git
git secrets --install
git secrets --register-aws

# ✅ Используйте переменные окружения
docker run -e DB_PASSWORD="$DB_PASSWORD" myapp

# ✅ Ротация секретов
# Регулярно меняйте секреты и обновляйте контейнеры
```

---

## Упражнение 5: Сканирование конфигурации

### Задание
Проверьте Docker конфигурацию на безопасность.

### Docker Bench Security

```bash
# Запуск Docker Bench Security
docker run -it --rm \
  --net host \
  --pid host \
  --userns host \
  --cap-add audit_control \
  -v /etc:/etc:ro \
  -v /usr/bin/containerd:/usr/bin/containerd:ro \
  -v /usr/bin/runc:/usr/bin/runc:ro \
  -v /usr/lib/systemd:/usr/lib/systemd:ro \
  -v /var/lib:/var/lib:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  docker/docker-bench-security

# Анализ результатов
```

### Hadolint для Dockerfile

```bash
# Установка
brew install hadolint

# Проверка Dockerfile
hadolint Dockerfile

# Проверка с игнорированием определенных правил
hadolint --ignore DL3006 --ignore DL3018 Dockerfile

# Вывод в JSON
hadolint --format json Dockerfile > hadolint-report.json
```

### Пример правильного Dockerfile

```dockerfile
# Dockerfile.best-practices

# DL3006: Always tag the version of an image explicitly
FROM python:3.11.6-slim-bookworm

# Метаданные
LABEL maintainer="security@example.com"
LABEL version="1.0"
LABEL security.scan="2024-01-15"

# DL3008: Pin versions in apt-get install
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates=20230311 \
        curl=7.88.1-10+deb12u5 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# DL3013: Pin versions in pip install
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# DL3002: Last USER should not be root
RUN useradd -m -u 1000 appuser

WORKDIR /app
COPY --chown=appuser:appuser . .

USER appuser

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
    CMD curl -f http://localhost:8000/health || exit 1

EXPOSE 8000

CMD ["python", "app.py"]
```

---

## Упражнение 6: Сетевая безопасность

### Задание
Настройте безопасную сетевую изоляцию.

### docker-compose с сетевой сегментацией

```yaml
# docker-compose.network-security.yml
version: '3.8'

services:
  # Frontend - публичный доступ
  frontend:
    image: nginx:alpine
    networks:
      - public
    ports:
      - "80:80"
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M

  # Backend - закрытый доступ
  backend:
    image: myapp:latest
    networks:
      - internal
      - public
    environment:
      - DB_HOST=database
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M

  # Database - полностью изолированная
  database:
    image: postgres:15-alpine
    networks:
      - internal
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db_password
    secrets:
      - db_password
    volumes:
      - db_data:/var/lib/postgresql/data
    # Нет проброса портов = нет внешнего доступа

networks:
  # Публичная сеть
  public:
    driver: bridge
    ipam:
      config:
        - subnet: 172.25.0.0/24

  # Внутренняя сеть (нет доступа к интернету)
  internal:
    driver: bridge
    internal: true
    ipam:
      config:
        - subnet: 172.26.0.0/24

secrets:
  db_password:
    file: ./secrets/db_password.txt

volumes:
  db_data:
```

### Использование firewall правил

```bash
# Ограничить доступ к Docker API
# Только локальный доступ
sudo ufw allow from 127.0.0.1 to any port 2375

# Использовать TLS для Docker daemon
# /etc/docker/daemon.json
{
  "tls": true,
  "tlsverify": true,
  "tlscacert": "/etc/docker/ca.pem",
  "tlscert": "/etc/docker/server-cert.pem",
  "tlskey": "/etc/docker/server-key.pem"
}
```

---

## Итоговое задание: Security Checklist

Создайте безопасное Docker приложение, проверив все пункты:

### Dockerfile Security Checklist

- [ ] Используются конкретные версии базовых образов (не latest)
- [ ] Минимальный базовый образ (alpine, distroless)
- [ ] Обновлены все пакеты
- [ ] Удалены build зависимости в финальном образе
- [ ] Используется multi-stage build
- [ ] Контейнер работает от непривилегированного пользователя
- [ ] Установлены минимальные права на файлы
- [ ] Нет секретов в образе
- [ ] Используется .dockerignore
- [ ] Добавлен HEALTHCHECK

### Runtime Security Checklist

- [ ] Удалены ненужные capabilities (`--cap-drop=ALL`)
- [ ] Read-only файловая система где возможно
- [ ] Ограничены ресурсы (memory, cpu)
- [ ] Используется `--security-opt=no-new-privileges`
- [ ] Секреты через переменные окружения или secrets API
- [ ] Сетевая сегментация настроена
- [ ] Логирование настроено
- [ ] Регулярное обновление образов

### Мониторинг Security Checklist

- [ ] Регулярное сканирование образов на уязвимости
- [ ] Мониторинг логов на подозрительную активность
- [ ] Аудит контейнеров (кто, когда, что запустил)
- [ ] Проверка конфигурации (Docker Bench Security)
- [ ] Валидация Dockerfile (Hadolint)

## Полезные инструменты

```bash
# Сканеры уязвимостей
- Trivy
- Clair
- Anchore
- Snyk

# Анализ Dockerfile
- Hadolint
- Dockerfile Linter

# Мониторинг
- Docker Bench Security
- Falco
- Sysdig

# Secret management
- HashiCorp Vault
- AWS Secrets Manager
- Azure Key Vault
- Docker Secrets
```

## Чек-лист освоенных навыков

- [ ] Сканирование образов на уязвимости
- [ ] Создание непривилегированных контейнеров
- [ ] Ограничение Linux capabilities
- [ ] Безопасная работа с секретами
- [ ] Использование read-only файловых систем
- [ ] Ограничение ресурсов контейнеров
- [ ] Сетевая изоляция
- [ ] Проверка Dockerfile на best practices
- [ ] Регулярное обновление образов
- [ ] Мониторинг безопасности

