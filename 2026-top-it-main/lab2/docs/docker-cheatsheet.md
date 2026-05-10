# Docker Cheat Sheet

## Содержание
- [Работа с образами](#работа-с-образами)
- [Работа с контейнерами](#работа-с-контейнерами)
- [Dockerfile](#dockerfile)
- [Docker Compose](#docker-compose)
- [Volumes](#volumes)
- [Networks](#networks)
- [Система и очистка](#система-и-очистка)
- [Логи и мониторинг](#логи-и-мониторинг)

---

## Работа с образами

### Основные команды

```bash
# Поиск образа на Docker Hub
docker search nginx

# Загрузка образа
docker pull nginx:latest
docker pull nginx:1.25-alpine

# Просмотр локальных образов
docker images
docker image ls

# Фильтрация образов
docker images --filter "dangling=true"
docker images --filter "before=nginx:latest"

# Просмотр информации об образе
docker image inspect nginx:latest

# История образа (слои)
docker image history nginx:latest

# Тегирование образа
docker tag nginx:latest myregistry.com/nginx:v1.0

# Удаление образа
docker rmi nginx:latest
docker image rm nginx:latest

# Удаление нескольких образов
docker rmi $(docker images -q)

# Сохранение образа в файл
docker save nginx:latest > nginx.tar
docker save nginx:latest -o nginx.tar

# Загрузка образа из файла
docker load < nginx.tar
docker load -i nginx.tar

# Экспорт контейнера в образ
docker export container_id > container.tar
docker import container.tar myimage:tag
```

---

## Работа с контейнерами

### Запуск контейнеров

```bash
# Простой запуск
docker run hello-world

# Запуск в фоновом режиме
docker run -d nginx:latest

# Запуск с именем
docker run -d --name my-nginx nginx:latest

# Интерактивный запуск
docker run -it ubuntu:22.04 /bin/bash

# Запуск с автоудалением
docker run --rm ubuntu:22.04 echo "Hello"

# Проброс портов
docker run -d -p 8080:80 nginx:latest
docker run -d -p 127.0.0.1:8080:80 nginx:latest

# Переменные окружения
docker run -d -e MYSQL_ROOT_PASSWORD=secret mysql:8

# Монтирование volumes
docker run -d -v /host/path:/container/path nginx
docker run -d -v myvolume:/data nginx

# Ограничение ресурсов
docker run -d --memory="512m" --cpus="1.0" nginx
docker run -d --memory-reservation="256m" nginx

# Сеть
docker run -d --network my-network nginx

# Перезапуск
docker run -d --restart=always nginx
docker run -d --restart=on-failure:5 nginx
docker run -d --restart=unless-stopped nginx
```

### Управление контейнерами

```bash
# Список запущенных контейнеров
docker ps

# Список всех контейнеров
docker ps -a

# Форматированный вывод
docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Status}}"

# Фильтрация
docker ps --filter "status=running"
docker ps --filter "ancestor=nginx"

# Остановка контейнера
docker stop container_id
docker stop container_name

# Остановка с таймаутом
docker stop -t 30 container_id

# Принудительная остановка
docker kill container_id

# Запуск остановленного контейнера
docker start container_id

# Перезапуск
docker restart container_id

# Пауза/возобновление
docker pause container_id
docker unpause container_id

# Удаление контейнера
docker rm container_id

# Принудительное удаление
docker rm -f container_id

# Удаление всех остановленных
docker container prune
```

### Взаимодействие с контейнерами

```bash
# Выполнение команды
docker exec container_id ls -la

# Интерактивный shell
docker exec -it container_id /bin/bash
docker exec -it container_id /bin/sh

# Выполнение от другого пользователя
docker exec -u root container_id whoami

# Просмотр логов
docker logs container_id
docker logs -f container_id          #Follow
docker logs --tail 100 container_id  # Последние 100 строк
docker logs --since 30m container_id # За последние 30 минут

# Копирование файлов
docker cp file.txt container_id:/path/
docker cp container_id:/path/file.txt ./

# Просмотр изменений в ФС
docker diff container_id

# Информация о контейнере
docker inspect container_id

# Статистика
docker stats                    # Все контейнеры
docker stats container_id       # Конкретный контейнер
docker stats --no-stream        # Одноразовый вывод

# Процессы в контейнере
docker top container_id

# Порты контейнера
docker port container_id

# Attach к контейнеру
docker attach container_id
```

### Обновление контейнера

```bash
# Обновить ресурсы
docker update --memory="1g" --cpus="2" container_id

# Обновить политику перезапуска
docker update --restart=always container_id
```

---

## Dockerfile

### Основные инструкции

```dockerfile
# Базовый образ
FROM python:3.11-slim

# Метаданные
LABEL maintainer="email@example.com"
LABEL version="1.0"
LABEL description="My application"

# Аргументы сборки
ARG NODE_VERSION=18
ARG APP_DIR=/app

# Переменные окружения
ENV APP_HOME=/app
ENV PYTHONUNBUFFERED=1

# Рабочая директория
WORKDIR /app

# Копирование файлов
COPY file.txt /app/
COPY --chown=user:user . /app/
ADD https://example.com/file.tar.gz /app/

# Выполнение команд
RUN apt-get update && apt-get install -y package
RUN pip install -r requirements.txt

# Открытие порта
EXPOSE 8000

# Том
VOLUME ["/data"]

# Пользователь
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8000/ || exit 1

# Точка входа
ENTRYPOINT ["python"]
CMD ["app.py"]
```

### Сборка образов

```bash
# Простая сборка
docker build -t myapp:latest .

# С указанием Dockerfile
docker build -t myapp:latest -f Dockerfile.prod .

# С аргументами
docker build --build-arg NODE_VERSION=18 -t myapp .

# Без кэша
docker build --no-cache -t myapp .

# Целевой этап (multi-stage)
docker build --target production -t myapp:prod .

# С pull новых базовых образов
docker build --pull -t myapp .

# Использование BuildKit
DOCKER_BUILDKIT=1 docker build -t myapp .

# С секретами
docker build --secret id=npmrc,src=$HOME/.npmrc -t myapp .
```

---

## Docker Compose

### Основные команды

```bash
# Запуск сервисов
docker-compose up
docker-compose up -d                    # Фоновый режим
docker-compose up --build               # С пересборкой
docker-compose up --force-recreate      # Пересоздать контейнеры

# Масштабирование
docker-compose up -d --scale web=3

# Остановка сервисов
docker-compose stop
docker-compose down                     # Остановка и удаление
docker-compose down -v                  # + удаление volumes
docker-compose down --remove-orphans    # + orphan контейнеры

# Просмотр сервисов
docker-compose ps
docker-compose ps -a

# Логи
docker-compose logs
docker-compose logs -f                  # Follow
docker-compose logs web                 # Конкретный сервис
docker-compose logs --tail=100 web      # Последние 100 строк

# Выполнение команд
docker-compose exec web ls -la
docker-compose exec web /bin/bash

# Запуск одноразовой команды
docker-compose run web python manage.py migrate
docker-compose run --rm web npm test

# Сборка
docker-compose build
docker-compose build --no-cache
docker-compose build --pull

# Конфигурация
docker-compose config                   # Проверка синтаксиса
docker-compose config --services        # Список сервисов

# Пауза/возобновление
docker-compose pause
docker-compose unpause

# Перезапуск
docker-compose restart
docker-compose restart web

# События
docker-compose events

# Порты
docker-compose port web 80

# Топология
docker-compose top
```

### docker-compose.yml примеры

```yaml
version: '3.8'

services:
  web:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        NODE_ENV: production
    image: myapp:latest
    container_name: web-app
    ports:
      - "8080:80"
      - "443:443"
    environment:
      - NODE_ENV=production
      - DB_HOST=db
    env_file:
      - .env
    volumes:
      - ./src:/app/src:ro
      - app-data:/data
    networks:
      - frontend
      - backend
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 256M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 30s
      timeout: 3s
      retries: 3
    labels:
      - "traefik.enable=true"

  db:
    image: postgres:15-alpine
    volumes:
      - db-data:/var/lib/postgresql/data
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password
    networks:
      - backend

networks:
  frontend:
    driver: bridge
  backend:
    driver: bridge
    internal: true

volumes:
  app-data:
  db-data:

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

---

## Volumes

### Управление volumes

```bash
# Создание volume
docker volume create myvolume

# Список volumes
docker volume ls

# Информация о volume
docker volume inspect myvolume

# Удаление volume
docker volume rm myvolume

# Удаление неиспользуемых
docker volume prune

# Использование volume
docker run -v myvolume:/data nginx
docker run --mount source=myvolume,target=/data nginx
```

### Типы монтирования

```bash
# Named volume
docker run -v myvolume:/data nginx

# Bind mount
docker run -v /host/path:/container/path nginx
docker run -v $(pwd):/app nginx

# tmpfs (в памяти)
docker run --tmpfs /tmp:rw,noexec,nosuid,size=100m nginx

# Read-only
docker run -v myvolume:/data:ro nginx
```

### Резервное копирование

```bash
# Backup
docker run --rm \
  -v myvolume:/data:ro \
  -v $(pwd):/backup \
  alpine tar czf /backup/backup.tar.gz -C /data .

# Restore
docker run --rm \
  -v myvolume:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/backup.tar.gz -C /data
```

---

## Networks

### Управление сетями

```bash
# Создание сети
docker network create mynet
docker network create --driver bridge mynet
docker network create --subnet=172.20.0.0/16 mynet

# Список сетей
docker network ls

# Информация о сети
docker network inspect mynet

# Подключение контейнера
docker network connect mynet container_id

# Отключение контейнера
docker network disconnect mynet container_id

# Удаление сети
docker network rm mynet

# Удаление неиспользуемых
docker network prune
```

### Типы сетей

```bash
# Bridge (по умолчанию)
docker network create --driver bridge mynet

# Host (использует сеть хоста)
docker run --network host nginx

# None (нет сети)
docker run --network none alpine

# Overlay (для Swarm)
docker network create --driver overlay mynet

# Macvlan
docker network create -d macvlan \
  --subnet=192.168.1.0/24 \
  --gateway=192.168.1.1 \
  -o parent=eth0 macnet
```

---

## Система и очистка

### Информация о системе

```bash
# Информация о Docker
docker info
docker version

# Использование диска
docker system df
docker system df -v

# События
docker events
docker events --since '2024-01-01T00:00:00'
```

### Очистка

```bash
# Удаление остановленных контейнеров
docker container prune

# Удаление неиспользуемых образов
docker image prune
docker image prune -a            # Все неиспользуемые

# Удаление неиспользуемых volumes
docker volume prune

# Удаление неиспользуемых сетей
docker network prune

# Удаление build cache
docker builder prune

# Полная очистка
docker system prune
docker system prune -a           # + неиспользуемые образы
docker system prune -a --volumes # + volumes
```

---

## Логи и мониторинг

### Логи

```bash
# Просмотр логов
docker logs container_id

# Follow режим
docker logs -f container_id

# Последние N строк
docker logs --tail 100 container_id

# С временными метками
docker logs -t container_id

# За период
docker logs --since 30m container_id
docker logs --until 2024-01-01T00:00:00 container_id

# Настройка логирования при запуске
docker run --log-driver json-file \
  --log-opt max-size=10m \
  --log-opt max-file=3 \
  nginx
```

### Мониторинг

```bash
# Статистика ресурсов
docker stats
docker stats --no-stream
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"

# Процессы
docker top container_id

# События
docker events
docker events --filter 'type=container'
docker events --filter 'event=start'
```

---

## Полезные алиасы

```bash
# Добавьте в ~/.bashrc или ~/.zshrc

# Удалить все остановленные контейнеры
alias dprune='docker container prune -f'

# Удалить все неиспользуемые образы
alias diprune='docker image prune -a -f'

# Полная очистка
alias dclean='docker system prune -a --volumes -f'

# Остановить все контейнеры
alias dstop='docker stop $(docker ps -q)'

# Удалить все контейнеры
alias drm='docker rm $(docker ps -aq)'

# Удалить все образы
alias drmi='docker rmi $(docker images -q)'

# Войти в контейнер
alias dexec='docker exec -it'

# Логи
alias dlogs='docker logs -f'

# Статистика
alias dstats='docker stats --no-stream'

# Docker Compose
alias dc='docker-compose'
alias dcup='docker-compose up -d'
alias dcdown='docker-compose down'
alias dclogs='docker-compose logs -f'
alias dcps='docker-compose ps'
```

---

## Переменные окружения

```bash
# Docker Host
export DOCKER_HOST=tcp://192.168.1.100:2376

# TLS настройки
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=/path/to/certs

# BuildKit
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

# Proxy настройки
export HTTP_PROXY=http://proxy:8080
export HTTPS_PROXY=http://proxy:8080
export NO_PROXY=localhost,127.0.0.1
```

---

## Быстрые команды

```bash
# Узнать IP контейнера
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' container_id

# Узнать все переменные окружения
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' container_id

# Найти контейнеры по образу
docker ps --filter ancestor=nginx

# Остановить все контейнеры
docker stop $(docker ps -q)

# Удалить все volumes
docker volume rm $(docker volume ls -q)

# Размер образа
docker images --format "{{.Repository}}:{{.Tag}} {{.Size}}"

# Количество слоев
docker history --no-trunc nginx:latest | wc -l

# Найти большие файлы в образе
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  wagoodman/dive:latest nginx:latest
```

---

## Отладка

```bash
# Проверить подключение
docker run --rm busybox ping -c 3 google.com

# DNS проверка
docker run --rm busybox nslookup google.com

# Проверить порт
docker run --rm busybox nc -zv host 80

# Проверить curl
docker run --rm appropriate/curl curl -I https://google.com

# Debug mode
docker run -it --entrypoint /bin/sh nginx:alpine

# Инспекция сети
docker run --rm --net container:web nicolaka/netshoot
```

---

Этот cheat sheet покрывает большинство повседневных операций с Docker. Сохраните его для быстрого доступа!

