# Часто задаваемые вопросы (FAQ)

## Содержание
- [Общие вопросы](#общие-вопросы)
- [Установка и настройка](#установка-и-настройка)
- [Работа с образами](#работа-с-образами)
- [Работа с контейнерами](#работа-с-контейнерами)
- [Docker Compose](#docker-compose)
- [Сети и volumes](#сети-и-volumes)
- [Производительность](#производительность)
- [Безопасность](#безопасность)
- [Troubleshooting](#troubleshooting)

---

## Общие вопросы

### В чем разница между Docker и виртуальной машиной?

**Docker (контейнеры):**
- Использует ядро хост-системы
- Легковесный (мегабайты)
- Быстрый старт (секунды)
- Меньше изоляции
- Эффективное использование ресурсов

**Виртуальная машина:**
- Полная виртуализация с гостевой ОС
- Тяжеловесная (гигабайты)
- Медленный старт (минуты)
- Полная изоляция
- Больше накладных расходов

### Когда использовать Docker, а когда VM?

**Используйте Docker:**
- Микросервисы
- Изоляция приложений
- CI/CD pipelines
- Разработка и тестирование
- Быстрое масштабирование

**Используйте VM:**
- Нужна другая ОС
- Максимальная изоляция
- Legacy приложения
- Долгоживущие stateful сервисы

### Что такое Docker Hub?

Docker Hub - это публичный registry для Docker образов. Это как GitHub для Docker образов, где можно:
- Искать и загружать готовые образы
- Публиковать свои образы
- Автоматически собирать образы из GitHub
- Хранить private образы

---

## Установка и настройка

### Как установить Docker на разных ОС?

**macOS:**
```bash
# Через Homebrew
brew install --cask docker

# Или скачать Docker Desktop с docker.com
```

**Windows:**
```bash
# Скачать Docker Desktop с docker.com
# Требуется WSL 2
```

**Linux (Ubuntu/Debian):**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### Как проверить, что Docker работает?

```bash
docker --version
docker run hello-world
docker info
```

### Нужно ли sudo для Docker команд?

На Linux по умолчанию да, но можно добавить пользователя в группу docker:
```bash
sudo usermod -aG docker $USER
# Выйти и войти снова
```

На macOS и Windows (Docker Desktop) sudo не требуется.

### Как настроить Docker за proxy?

Создайте файл `/etc/systemd/system/docker.service.d/http-proxy.conf`:
```ini
[Service]
Environment="HTTP_PROXY=http://proxy:8080"
Environment="HTTPS_PROXY=http://proxy:8080"
Environment="NO_PROXY=localhost,127.0.0.1"
```

Перезапустите Docker:
```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

---

## Работа с образами

### Как узнать, какие образы доступны на Docker Hub?

```bash
docker search nginx
docker search --filter "is-official=true" python
```

Или через веб-интерфейс: https://hub.docker.com/

### Чем отличается образ с тегом `alpine` от обычного?

**Alpine** - это минималистичный Linux дистрибутив:
- Очень маленький размер (~5MB базовый образ)
- Использует musl libc вместо glibc
- Может быть несовместимость с некоторыми приложениями
- Отлично для production

**Обычный** (на базе Debian/Ubuntu):
- Больший размер (~100-900MB)
- Использует glibc
- Лучшая совместимость
- Больше инструментов из коробки

### Как уменьшить размер Docker образа?

1. Используйте alpine или slim базовые образы
2. Multi-stage builds
3. Объединяйте RUN команды
4. Удаляйте кэш пакетных менеджеров
5. Используйте .dockerignore
6. Копируйте только нужные файлы

### Что такое dangling образы?

Dangling образы - это образы без тега (показываются как `<none>`), обычно остаются после пересборки образов.

Удаление:
```bash
docker image prune
```

### Как сохранить образ в файл?

```bash
# Сохранить
docker save nginx:latest > nginx.tar
docker save nginx:latest -o nginx.tar

# Загрузить
docker load < nginx.tar
docker load -i nginx.tar
```

---

## Работа с контейнерами

### В чем разница между `docker run` и `docker start`?

- `docker run` - создает НОВЫЙ контейнер из образа
- `docker start` - запускает СУЩЕСТВУЮЩИЙ остановленный контейнер

### Как попасть внутрь запущенного контейнера?

```bash
# Если есть bash
docker exec -it container_name /bin/bash

# Если только sh (alpine)
docker exec -it container_name /bin/sh

# Или attach (но это подключает к главному процессу)
docker attach container_name
```

### Почему контейнер сразу останавливается?

Контейнер работает, пока работает главный процесс (CMD/ENTRYPOINT). Если процесс завершился - контейнер останавливается.

Для отладки:
```bash
docker logs container_name
docker run -it image_name /bin/sh
```

### Как скопировать файлы между хостом и контейнером?

```bash
# С хоста в контейнер
docker cp file.txt container_name:/path/

# Из контейнера на хост
docker cp container_name:/path/file.txt ./
```

### Что делает флаг `--rm`?

Автоматически удаляет контейнер после остановки:
```bash
docker run --rm ubuntu echo "Hello"
```

Полезно для одноразовых команд.

---

## Docker Compose

### В чем разница между `docker-compose.yml` и `docker-compose.override.yml`?

- `docker-compose.yml` - базовая конфигурация
- `docker-compose.override.yml` - автоматически применяется поверх базовой (для локальной разработки)

### Как использовать несколько compose файлов?

```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

Или через переменную:
```bash
export COMPOSE_FILE=docker-compose.yml:docker-compose.prod.yml
docker-compose up -d
```

### Как масштабировать сервис?

```bash
docker-compose up -d --scale web=3
```

Примечание: нельзя масштабировать сервисы с `container_name` или фиксированными портами.

### Почему `depends_on` не ждет готовности сервиса?

`depends_on` только ждет запуска контейнера, не готовности приложения.

Решения:
1. Используйте healthcheck:
```yaml
depends_on:
  db:
    condition: service_healthy
```

2. Используйте wait-for-it скрипт
3. Добавьте retry логику в приложение

### Как передать переменные окружения?

```yaml
services:
  app:
    environment:
      - DEBUG=true
    # или
    env_file:
      - .env
```

Через командную строку:
```bash
DEBUG=true docker-compose up
```

---

## Сети и volumes

### Как контейнеры общаются между собой?

В пользовательских bridge сетях (и в Docker Compose) контейнеры видят друг друга по имени через встроенный DNS.

```bash
# Создать сеть
docker network create mynet

# Запустить контейнеры
docker run -d --network mynet --name app1 nginx
docker run -d --network mynet --name app2 alpine

# app2 может пинговать app1 по имени
docker exec app2 ping app1
```

### Что лучше: named volume или bind mount?

**Named volume:**
- Управляется Docker
- Лучшая производительность
- Кроссплатформенность
- Для production

**Bind mount:**
- Управляется пользователем
- Прямой доступ к файлам
- Для разработки
- Проблемы с производительностью на Mac/Windows

### Как сохранить данные из volume?

```bash
docker run --rm \
  -v myvolume:/data:ro \
  -v $(pwd):/backup \
  alpine tar czf /backup/backup.tar.gz -C /data .
```

### Что такое tmpfs?

Временная файловая система в памяти:
```bash
docker run --tmpfs /tmp:rw,size=100m alpine
```

Данные исчезают при остановке контейнера. Полезно для:
- Временных файлов
- Кэшей
- Секретов

---

## Производительность

### Docker медленный на macOS/Windows?

Да, особенно с bind mounts. Причины:
- Виртуализация (Docker работает в VM)
- Файловая система (перевод между хостом и VM)

Решения:
1. Используйте named volumes вместо bind mounts
2. Используйте cached или delegated режимы:
   ```yaml
   volumes:
     - ./src:/app/src:cached
   ```
3. Используйте docker-sync

### Как ускорить сборку образов?

1. Правильный порядок COPY (зависимости перед кодом)
2. Используйте .dockerignore
3. BuildKit: `export DOCKER_BUILDKIT=1`
4. Кэш слоев: `--cache-from`
5. Multi-stage builds

### Как ограничить ресурсы контейнера?

```bash
docker run -d \
  --memory="512m" \
  --cpus="1.0" \
  --pids-limit=100 \
  nginx
```

В docker-compose:
```yaml
deploy:
  resources:
    limits:
      cpus: '1.0'
      memory: 512M
```

### Почему Docker занимает много места?

```bash
# Проверить использование
docker system df

# Очистить
docker system prune -a --volumes
```

---

## Безопасность

### Безопасно ли запускать контейнеры от root?

Нет! Всегда создавайте непривилегированного пользователя:

```dockerfile
RUN useradd -m -u 1000 appuser
USER appuser
```

### Как сканировать образы на уязвимости?

```bash
# Trivy
trivy image myapp:latest

# Docker Scout
docker scout cves myapp:latest

# Snyk
snyk container test myapp:latest
```

### Как безопасно хранить секреты?

❌ **Плохо:**
```dockerfile
ENV API_KEY=secret123
```

✅ **Хорошо:**
```bash
# Переменные окружения
docker run -e API_KEY=$API_KEY myapp

# Docker secrets (Swarm)
echo "secret" | docker secret create api_key -

# External secret manager (Vault, AWS Secrets Manager)
```

### Что такое rootless Docker?

Docker, работающий от имени обычного пользователя (не root). Повышает безопасность, но имеет ограничения.

---

## Troubleshooting

### "Cannot connect to the Docker daemon"

```bash
# Проверить, запущен ли Docker
sudo systemctl status docker  # Linux
open -a Docker  # macOS

# Проверить DOCKER_HOST
echo $DOCKER_HOST
unset DOCKER_HOST
```

### "No space left on device"

```bash
# Проверить место
docker system df
df -h

# Очистить
docker system prune -a --volumes
```

### "Port is already allocated"

Порт занят другим процессом:
```bash
# Найти процесс
sudo lsof -i :8080

# Или использовать другой порт
docker run -p 8081:80 nginx
```

### "driver failed programming external connectivity"

```bash
# Перезапустить Docker
sudo systemctl restart docker  # Linux

# Или перезагрузить систему
```

### Контейнер не может подключиться к интернету

```bash
# Проверить DNS
docker run --rm alpine ping -c 3 8.8.8.8
docker run --rm alpine nslookup google.com

# Исправить DNS в /etc/docker/daemon.json
{
  "dns": ["8.8.8.8", "8.8.4.4"]
}
```

---

## Полезные ресурсы

- [Docker Documentation](https://docs.docker.com/)
- [Docker Hub](https://hub.docker.com/)
- [Play with Docker](https://labs.play-with-docker.com/)
- [Docker Forum](https://forums.docker.com/)
- [Stack Overflow - Docker](https://stackoverflow.com/questions/tagged/docker)

---

Если ваш вопрос не нашелся в FAQ, обратитесь к [документации](https://docs.docker.com/) или [troubleshooting guide](./examples/troubleshooting/README.md).

