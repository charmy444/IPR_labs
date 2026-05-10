# Docker Troubleshooting - Решение проблем

## Описание
Руководство по диагностике и решению типичных проблем Docker.

## Категория 1: Проблемы с контейнерами

### Проблема 1: Контейнер сразу останавливается

```bash
# Симптомы
docker run -d myapp:latest
docker ps  # Контейнер не запущен

# Диагностика
docker ps -a  # Посмотреть все контейнеры
docker logs <container_id>  # Проверить логи
docker inspect <container_id>  # Подробная информация

# Возможные причины и решения:

# 1. Ошибка в приложении
docker logs <container_id>
# Решение: Исправить код приложения

# 2. Некорректная команда CMD/ENTRYPOINT
docker inspect <container_id> --format='{{.Config.Cmd}}'
# Решение: Исправить Dockerfile

# 3. Приложение завершается сразу
# Решение: Убедиться, что приложение работает постоянно
# Например, для тестирования:
docker run -d myapp:latest sleep infinity
docker exec -it <container_id> /bin/sh

# 4. Отсутствуют зависимости
docker run -it myapp:latest /bin/sh
# Попробовать запустить приложение вручную
```

### Проблема 2: Контейнер не отвечает

```bash
# Диагностика
docker ps  # Проверить статус
docker top <container_id>  # Процессы в контейнере
docker stats <container_id>  # Использование ресурсов

# Проверить healthcheck
docker inspect <container_id> --format='{{json .State.Health}}'

# Проверить логи
docker logs --tail 100 -f <container_id>

# Решения:

# 1. Контейнер перегружен
docker stats <container_id>
# Увеличить лимиты:
docker update --memory="1g" --cpus="2" <container_id>

# 2. Deadlock или зависание
docker exec <container_id> ps aux
# Перезапустить:
docker restart <container_id>

# 3. Ожидает ввода
docker attach <container_id>
# Или запустить в интерактивном режиме

# 4. Проблемы с сетью
docker exec <container_id> ping 8.8.8.8
docker exec <container_id> curl http://example.com
```

### Проблема 3: "Permission denied" ошибки

```bash
# Симптомы
Error: permission denied

# Диагностика
docker exec <container_id> whoami
docker exec <container_id> ls -la /target/directory

# Решения:

# 1. Контейнер работает от непривилегированного пользователя
# Изменить владельца в Dockerfile:
COPY --chown=user:user . /app

# 2. Проблемы с volume permissions
docker run -v /host/path:/container/path \
  --user $(id -u):$(id -g) \
  myapp:latest

# 3. SELinux проблемы (на RHEL/CentOS/Fedora)
# Добавить :z или :Z к volume
docker run -v /host/path:/container/path:z myapp:latest

# 4. Нужны дополнительные capabilities
docker run --cap-add=DAC_OVERRIDE myapp:latest
```

---

## Категория 2: Проблемы с образами

### Проблема 4: Ошибки при сборке образа

```bash
# ERROR: failed to solve with frontend dockerfile.v0

# Причина 1: Синтаксическая ошибка в Dockerfile
# Решение: Проверить Dockerfile на ошибки
hadolint Dockerfile

# Причина 2: Недоступен базовый образ
docker pull nginx:latest  # Проверить доступность

# Причина 3: Кэш поврежден
docker builder prune -af  # Очистить кэш сборки
docker system prune -af   # Полная очистка

# Причина 4: Недостаточно места
df -h  # Проверить место на диске
docker system df  # Использование Docker
docker system prune -a --volumes -f  # Очистить

# Причина 5: Проблемы с сетью при COPY/ADD
# Использовать .dockerignore
# Или проверить proxy настройки
```

### Проблема 5: Образ слишком большой

```bash
# Диагностика
docker images myapp
docker history myapp:latest

# Решения:

# 1. Использовать multi-stage build
FROM golang:1.21 AS builder
# сборка
FROM alpine:3.18
COPY --from=builder /app/binary /app/

# 2. Очищать кэш в том же RUN
RUN apt-get update && \
    apt-get install -y package && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 3. Использовать меньший базовый образ
FROM alpine:3.18  # вместо ubuntu
FROM python:3.11-slim  # вместо python:3.11

# 4. Анализ слоев
docker history myapp:latest
# Или с помощью dive:
# docker run --rm -it \
#   -v /var/run/docker.sock:/var/run/docker.sock \
#   wagoodman/dive:latest myapp:latest

# 5. Удалить build зависимости
RUN apk add --no-cache --virtual .build-deps gcc musl-dev && \
    pip install package && \
    apk del .build-deps
```

### Проблема 6: "No space left on device"

```bash
# Диагностика
df -h
docker system df
docker system df -v

# Решения:

# 1. Очистить неиспользуемые ресурсы
docker system prune -a --volumes -f

# 2. Удалить старые образы
docker image prune -a -f

# 3. Удалить dangling volumes
docker volume prune -f

# 4. Удалить build cache
docker builder prune -af

# 5. Найти большие образы/контейнеры
docker images --format "{{.Repository}}:{{.Tag}} {{.Size}}" | sort -k2 -h
docker ps -as

# 6. Изменить Docker data directory
# /etc/docker/daemon.json
{
  "data-root": "/mnt/docker-data"
}
# Перезапустить Docker
sudo systemctl restart docker
```

---

## Категория 3: Сетевые проблемы

### Проблема 7: Контейнер не доступен по порту

```bash
# Симптомы
curl http://localhost:8080
# curl: (7) Failed to connect to localhost port 8080

# Диагностика
docker ps  # Проверить PORTS
docker port <container_id>  # Какие порты пробрасываются

# Решения:

# 1. Порт не пробрасывается
# Остановить и запустить с -p:
docker stop <container_id>
docker rm <container_id>
docker run -d -p 8080:80 myapp:latest

# 2. Приложение слушает неправильный интерфейс
# В приложении использовать 0.0.0.0, не 127.0.0.1:
# app.listen(8080, '0.0.0.0')  # ✅
# app.listen(8080, '127.0.0.1')  # ❌

# 3. Порт уже занят на хосте
sudo lsof -i :8080  # Проверить кто использует порт
# Использовать другой порт:
docker run -d -p 8081:80 myapp:latest

# 4. Firewall блокирует порт
sudo ufw status
sudo ufw allow 8080

# 5. Приложение не запустилось
docker logs <container_id>
docker exec <container_id> netstat -tlnp
```

### Проблема 8: Контейнеры не видят друг друга

```bash
# Симптомы
docker exec frontend ping backend
# ping: bad address 'backend'

# Диагностика
docker network ls
docker network inspect bridge

# Решения:

# 1. Контейнеры в разных сетях
# Создать общую сеть:
docker network create mynet
docker network connect mynet container1
docker network connect mynet container2

# 2. Используется default bridge (нет DNS)
# Создать пользовательскую bridge сеть:
docker network create --driver bridge mynet
docker run -d --network mynet --name backend nginx
docker run -d --network mynet --name frontend alpine

# 3. Неправильное имя хоста
docker network inspect mynet
# Использовать container_name или network alias

# 4. Порт не открыт в контейнере
docker exec backend netstat -tlnp
# Убедиться что приложение слушает порт

# Docker Compose пример:
version: '3.8'
services:
  backend:
    image: nginx
    networks:
      - mynet
  frontend:
    image: alpine
    command: ping backend
    networks:
      - mynet
networks:
  mynet:
```

### Проблема 9: Контейнер не имеет доступа в интернет

```bash
# Симптомы
docker exec <container_id> ping 8.8.8.8
# Network is unreachable

# Диагностика
docker exec <container_id> cat /etc/resolv.conf
docker exec <container_id> ip route

# Решения:

# 1. DNS проблемы
docker run --dns 8.8.8.8 --dns 8.8.4.4 myapp:latest

# Или в /etc/docker/daemon.json:
{
  "dns": ["8.8.8.8", "8.8.4.4"]
}

# 2. Сеть помечена как internal
docker network inspect mynet | grep internal
# Пересоздать без internal:
docker network rm mynet
docker network create mynet

# 3. Проблемы с iptables
sudo iptables -L DOCKER-USER
# Проверить правила

# 4. Proxy настройки
docker run -e HTTP_PROXY=http://proxy:8080 myapp:latest

# 5. IPv4 forwarding отключен
sysctl net.ipv4.ip_forward
# Включить:
sudo sysctl -w net.ipv4.ip_forward=1
```

---

## Категория 4: Проблемы с volumes

### Проблема 10: Данные не сохраняются

```bash
# Симптомы
# После перезапуска контейнера данные исчезли

# Диагностика
docker inspect <container_id> | grep -A 10 Mounts

# Решения:

# 1. Используется анонимный volume
# Использовать именованный volume:
docker volume create mydata
docker run -v mydata:/data myapp:latest

# 2. Не указан volume
docker run -v /var/lib/postgresql/data postgres:15

# 3. Volume удален вместе с контейнером
# Не использовать -v при docker rm:
docker rm container_name  # ✅ volume сохранен
docker rm -v container_name  # ❌ volume удален

# Docker Compose:
volumes:
  mydata:  # Явно объявить volume
```

### Проблема 11: Медленная работа с volume на macOS/Windows

```bash
# Проблема: Bind mounts очень медленные

# Решения:

# 1. Использовать cached mode (macOS)
-v /host/path:/container/path:cached

# 2. Использовать delegated mode
-v /host/path:/container/path:delegated

# 3. Использовать именованные volumes вместо bind mounts
docker volume create mydata
docker run -v mydata:/app/data myapp:latest

# 4. Использовать docker-sync (для больших проектов)
# gem install docker-sync
# docker-sync start

# Docker Compose:
volumes:
  - mydata:/app/data  # ✅ Быстро
  # - ./src:/app/src  # ❌ Медленно на Mac/Windows
```

---

## Категория 5: Docker Compose проблемы

### Проблема 12: "Service failed to build"

```bash
# Диагностика
docker-compose build --no-cache service_name

# Решения:

# 1. Проблема с контекстом сборки
docker-compose config  # Проверить конфигурацию

# 2. Неправильный путь к Dockerfile
services:
  app:
    build:
      context: ./app
      dockerfile: Dockerfile  # Проверить путь

# 3. Build args не переданы
services:
  app:
    build:
      args:
        - NODE_VERSION=18

# 4. Кэш проблемы
docker-compose build --no-cache --pull
```

### Проблема 13: "depends_on не работает"

```bash
# Проблема: База данных еще не готова когда app стартует

# Решение 1: Использовать condition
services:
  db:
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 5s
  app:
    depends_on:
      db:
        condition: service_healthy

# Решение 2: Wait-for-it скрипт
# В Dockerfile:
COPY wait-for-it.sh /usr/local/bin/
# В docker-compose:
command: ["./wait-for-it.sh", "db:5432", "--", "python", "app.py"]

# Решение 3: Retry логика в приложении
import time
def connect_db():
    for i in range(30):
        try:
            conn = connect()
            return conn
        except:
            time.sleep(1)
```

---

## Категория 6: Производительность

### Проблема 14: Контейнер работает медленно

```bash
# Диагностика
docker stats <container_id>

# Проверить ограничения
docker inspect <container_id> | grep -i memory
docker inspect <container_id> | grep -i cpu

# Решения:

# 1. Увеличить лимиты ресурсов
docker run -m 2g --cpus="2" myapp:latest

# 2. Проверить I/O
docker stats  # Смотреть BLOCK I/O

# 3. Оптимизировать образ
# Использовать multi-stage build
# Уменьшить количество слоев
# Использовать alpine образы

# 4. Проверить логирование
docker inspect <container_id> | grep LogConfig
# Ограничить размер логов:
docker run --log-opt max-size=10m --log-opt max-file=3 myapp

# 5. Использовать tmpfs для временных данных
docker run --tmpfs /tmp:rw,noexec,nosuid,size=1g myapp
```

### Проблема 15: Сборка образа занимает много времени

```bash
# Решения:

# 1. Правильный порядок COPY в Dockerfile
# ❌ Плохо:
COPY . .
RUN npm install

# ✅ Хорошо:
COPY package*.json ./
RUN npm ci
COPY . .

# 2. Использовать .dockerignore
cat > .dockerignore << 'EOF'
node_modules
.git
*.md
.env
EOF

# 3. Использовать BuildKit
export DOCKER_BUILDKIT=1
docker build .

# 4. Кэширование слоев
# Использовать --cache-from:
docker build --cache-from myapp:latest -t myapp:new .

# 5. Параллельная сборка stage-ов
# BuildKit автоматически параллелит независимые этапы
```

---

## Полезные команды для диагностики

```bash
# === Логи и мониторинг ===
docker logs -f --tail 100 <container_id>
docker logs --since 30m <container_id>
docker logs --until 2023-01-01T00:00:00 <container_id>
docker stats
docker events

# === Информация о контейнере ===
docker inspect <container_id>
docker top <container_id>
docker port <container_id>
docker diff <container_id>

# === Сеть ===
docker network ls
docker network inspect bridge
docker exec <container_id> netstat -tlnp
docker exec <container_id> ping 8.8.8.8
docker exec <container_id> nslookup example.com
docker exec <container_id> curl -v http://target

# === Файловая система ===
docker exec <container_id> df -h
docker exec <container_id> du -sh /app
docker exec <container_id> ls -lah /
docker cp <container_id>:/path/to/file ./local/

# === Процессы ===
docker exec <container_id> ps aux
docker exec <container_id> top
docker exec <container_id> free -m

# === Диски и volumes ===
docker system df
docker volume ls
docker volume inspect <volume_name>

# === Debug ===
docker run -it --rm --entrypoint /bin/sh myapp:latest
docker exec -it <container_id> /bin/sh
docker attach <container_id>
```

## Чек-лист диагностики

Когда что-то не работает, проверьте:

### 1. Контейнер
- [ ] Контейнер запущен? `docker ps`
- [ ] Что в логах? `docker logs <id>`
- [ ] Какой статус? `docker inspect <id>`
- [ ] Какие процессы? `docker top <id>`
- [ ] Использование ресурсов? `docker stats <id>`

### 2. Сеть
- [ ] Порты пробрасываются? `docker port <id>`
- [ ] Контейнеры в одной сети? `docker network inspect`
- [ ] DNS работает? `docker exec <id> nslookup hostname`
- [ ] Ping проходит? `docker exec <id> ping target`
- [ ] Порт открыт? `docker exec <id> netstat -tlnp`

### 3. Volumes
- [ ] Volume примонтирован? `docker inspect <id> | grep Mounts`
- [ ] Права доступа? `docker exec <id> ls -la /path`
- [ ] Volume существует? `docker volume ls`
- [ ] Место на диске? `docker system df`

### 4. Образ
- [ ] Образ существует? `docker images`
- [ ] Правильная версия? `docker inspect image`
- [ ] Нет уязвимостей? `docker scan image`
- [ ] Размер нормальный? `docker history image`

### 5. Система
- [ ] Место на диске? `df -h`
- [ ] Docker работает? `docker info`
- [ ] Версия Docker? `docker version`
- [ ] Нет проблем с памятью? `free -m`

## Итоговое задание

Создайте скрипт для автоматической диагностики Docker проблем:

```bash
#!/bin/bash
# docker-debug.sh

CONTAINER_ID=$1

if [ -z "$CONTAINER_ID" ]; then
    echo "Usage: $0 <container_id>"
    exit 1
fi

echo "=== Container Status ==="
docker ps -a --filter id=$CONTAINER_ID

echo -e "\n=== Container Logs (last 50 lines) ==="
docker logs --tail 50 $CONTAINER_ID

echo -e "\n=== Container Inspect ==="
docker inspect $CONTAINER_ID | jq '.[]| {State, RestartCount, NetworkSettings}'

echo -e "\n=== Resource Usage ==="
docker stats --no-stream $CONTAINER_ID

echo -e "\n=== Processes ==="
docker top $CONTAINER_ID

echo -e "\n=== Network ==="
docker exec $CONTAINER_ID netstat -tlnp 2>/dev/null || echo "netstat not available"

echo -e "\n=== Disk Usage ==="
docker exec $CONTAINER_ID df -h 2>/dev/null || echo "df not available"

echo -e "\n=== Environment ==="
docker exec $CONTAINER_ID env | sort

echo "=== Done ==="
```

## Чек-лист освоенных навыков

- [ ] Диагностика проблем с запуском контейнеров
- [ ] Анализ логов Docker
- [ ] Решение проблем с правами доступа
- [ ] Устранение ошибок сборки образов
- [ ] Оптимизация размера образов
- [ ] Диагностика сетевых проблем
- [ ] Решение проблем с DNS
- [ ] Работа с volumes и данными
- [ ] Отладка Docker Compose
- [ ] Оптимизация производительности
- [ ] Использование инструментов мониторинга
- [ ] Очистка Docker системы

## Полезные ресурсы

- [Docker Documentation](https://docs.docker.com/)
- [Docker Forum](https://forums.docker.com/)
- [Stack Overflow - Docker](https://stackoverflow.com/questions/tagged/docker)
- [Docker GitHub Issues](https://github.com/moby/moby/issues)

