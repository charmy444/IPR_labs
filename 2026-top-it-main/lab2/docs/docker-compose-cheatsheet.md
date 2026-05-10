# Docker Compose Cheat Sheet

## Основные команды

### Управление сервисами

```bash
# Запуск всех сервисов
docker-compose up
docker-compose up -d                 # Фоновый режим
docker-compose up --build            # С пересборкой образов
docker-compose up --force-recreate   # Пересоздать контейнеры
docker-compose up --no-deps web      # Без зависимостей

# Остановка сервисов
docker-compose stop                  # Остановить сервисы
docker-compose down                  # Остановить и удалить
docker-compose down -v               # + удалить volumes
docker-compose down --rmi all        # + удалить образы
docker-compose down --remove-orphans # + orphan контейнеры

# Пауза/возобновление
docker-compose pause
docker-compose unpause

# Перезапуск
docker-compose restart
docker-compose restart web           # Конкретный сервис
```

### Сборка

```bash
# Сборка образов
docker-compose build
docker-compose build --no-cache      # Без кэша
docker-compose build --pull          # Обновить базовые образы
docker-compose build web             # Конкретный сервис
```

### Просмотр и логи

```bash
# Статус сервисов
docker-compose ps
docker-compose ps -a                 # Все сервисы

# Логи
docker-compose logs
docker-compose logs -f               # Follow mode
docker-compose logs --tail=100 web   # Последние 100 строк
docker-compose logs --since 30m      # За последние 30 минут

# События
docker-compose events
docker-compose events --json         # В JSON формате
```

### Выполнение команд

```bash
# Exec в запущенном контейнере
docker-compose exec web bash
docker-compose exec web ls -la
docker-compose exec -T web command   # Без TTY

# Run одноразовой команды
docker-compose run web python manage.py migrate
docker-compose run --rm web npm test # С автоудалением
docker-compose run -e DEBUG=1 web    # С env переменной
```

### Масштабирование

```bash
# Масштабировать сервис
docker-compose up -d --scale web=3
docker-compose up -d --scale worker=5

# Нельзя масштабировать сервисы с container_name!
```

### Конфигурация

```bash
# Проверка синтаксиса
docker-compose config

# Список сервисов
docker-compose config --services

# Разрешенная конфигурация
docker-compose config --resolve-image-digests

# Использование нескольких файлов
docker-compose -f docker-compose.yml -f docker-compose.prod.yml config
```

### Порты

```bash
# Показать порты сервиса
docker-compose port web 80
docker-compose port --index=2 web 80  # Для масштабированных сервисов
```

### Топология

```bash
# Процессы в контейнерах
docker-compose top
docker-compose top web
```

---

## docker-compose.yml структура

### Полный пример

```yaml
version: '3.8'

# Сервисы
services:
  # Веб-приложение
  web:
    # Образ из Docker Hub
    image: nginx:alpine
    
    # Или сборка из Dockerfile
    build:
      context: ./web
      dockerfile: Dockerfile
      args:
        NODE_ENV: production
      target: production
      cache_from:
        - myapp:latest
    
    # Имя контейнера
    container_name: my-web-app
    
    # Hostname
    hostname: web-server
    
    # Порты
    ports:
      - "8080:80"              # host:container
      - "443:443"
      - "127.0.0.1:8081:80"    # Только localhost
    
    # Expose (только внутри Docker)
    expose:
      - "3000"
    
    # Переменные окружения
    environment:
      - NODE_ENV=production
      - DEBUG=false
      - API_URL=http://api:3000
    
    # Или из файла
    env_file:
      - .env
      - .env.production
    
    # Volumes
    volumes:
      - ./src:/app/src:ro            # Bind mount (read-only)
      - app-data:/data               # Named volume
      - /tmp                         # Anonymous volume
      - type: bind                   # Long syntax
        source: ./config
        target: /etc/config
        read_only: true
    
    # Сети
    networks:
      - frontend
      - backend
    
    # Зависимости
    depends_on:
      db:
        condition: service_healthy
      cache:
        condition: service_started
    
    # Команда
    command: nginx -g 'daemon off;'
    
    # Entrypoint
    entrypoint: /docker-entrypoint.sh
    
    # Рабочая директория
    working_dir: /app
    
    # Пользователь
    user: "1000:1000"
    
    # Политика перезапуска
    restart: unless-stopped
    # no, always, on-failure, unless-stopped
    
    # Healthcheck
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 40s
    
    # Ограничения ресурсов
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 256M
      replicas: 3
    
    # Метки
    labels:
      - "com.example.description=Web service"
      - "com.example.version=1.0"
    
    # Логирование
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
    
    # DNS
    dns:
      - 8.8.8.8
      - 8.8.4.4
    dns_search:
      - dc1.example.com
    
    # Extra hosts
    extra_hosts:
      - "host.docker.internal:host-gateway"
      - "api.example.com:192.168.1.100"
    
    # Tmpfs
    tmpfs:
      - /tmp
      - /run:size=100M
    
    # Security
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
    
    # Devices
    devices:
      - "/dev/ttyUSB0:/dev/ttyUSB0"
    
    # Profiles (опциональные сервисы)
    profiles:
      - debug
      - monitoring

  # База данных
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    volumes:
      - db-data:/var/lib/postgresql/data
      - ./init-db:/docker-entrypoint-initdb.d:ro
    networks:
      - backend
    secrets:
      - db_password
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Cache
  cache:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data
    networks:
      - backend

# Networks
networks:
  frontend:
    driver: bridge
    ipam:
      driver: default
      config:
        - subnet: 172.20.0.0/16
          gateway: 172.20.0.1
  
  backend:
    driver: bridge
    internal: true  # Нет доступа к внешней сети
    ipam:
      config:
        - subnet: 172.21.0.0/16

# Volumes
volumes:
  app-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/app
  
  db-data:
    driver: local
    labels:
      - "com.example.backup=true"
  
  redis-data:
    external: true  # Уже существующий volume

# Secrets
secrets:
  db_password:
    file: ./secrets/db_password.txt
  
  api_key:
    external: true  # Уже существующий secret

# Configs
configs:
  nginx_config:
    file: ./nginx.conf
```

---

## Паттерны и примеры

### 1. Разделение окружений

```yaml
# docker-compose.yml (базовая конфигурация)
version: '3.8'
services:
  web:
    build: .
    ports:
      - "8080:80"
```

```yaml
# docker-compose.override.yml (автоматически применяется)
version: '3.8'
services:
  web:
    volumes:
      - ./src:/app/src
    environment:
      - DEBUG=true
```

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  web:
    restart: always
    environment:
      - DEBUG=false
    deploy:
      replicas: 3
```

```bash
# Разработка (использует override автоматически)
docker-compose up

# Production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 2. Множественные базы данных

```yaml
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres-data:/var/lib/postgresql/data
  
  mysql:
    image: mysql:8
    environment:
      MYSQL_ROOT_PASSWORD: secret
    volumes:
      - mysql-data:/var/lib/mysql
  
  mongodb:
    image: mongo:6
    volumes:
      - mongo-data:/data/db

volumes:
  postgres-data:
  mysql-data:
  mongo-data:
```

### 3. Microservices с API Gateway

```yaml
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - auth
      - users
      - products
    networks:
      - public

  auth:
    build: ./services/auth
    networks:
      - public
      - private
    depends_on:
      - db
  
  users:
    build: ./services/users
    networks:
      - public
      - private
    depends_on:
      - db
  
  products:
    build: ./services/products
    networks:
      - public
      - private
    depends_on:
      - db

  db:
    image: postgres:15
    networks:
      - private
    volumes:
      - db-data:/var/lib/postgresql/data

networks:
  public:
  private:
    internal: true

volumes:
  db-data:
```

### 4. Мониторинг стек

```yaml
services:
  app:
    build: .
    labels:
      - "prometheus.scrape=true"
      - "prometheus.port=9090"
  
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    ports:
      - "9090:9090"
    profiles:
      - monitoring
  
  grafana:
    image: grafana/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana
    ports:
      - "3000:3000"
    profiles:
      - monitoring

volumes:
  prometheus-data:
  grafana-data:
```

```bash
# Запуск с мониторингом
docker-compose --profile monitoring up -d
```

### 5. Development с hot-reload

```yaml
services:
  frontend:
    build:
      context: ./frontend
      target: development
    volumes:
      - ./frontend/src:/app/src
      - /app/node_modules  # Не перезаписывать node_modules
    environment:
      - CHOKIDAR_USEPOLLING=true  # Для hot-reload на Windows/Mac
    command: npm run dev
    ports:
      - "3000:3000"

  backend:
    build:
      context: ./backend
      target: development
    volumes:
      - ./backend:/app
      - /app/__pycache__
      - /app/.pytest_cache
    environment:
      - FLASK_ENV=development
      - FLASK_DEBUG=1
    command: flask run --host=0.0.0.0
```

---

## Переменные в docker-compose

### .env файл

```bash
# .env
TAG=latest
DB_PASSWORD=secret123
API_PORT=8080
```

### Использование в docker-compose.yml

```yaml
services:
  web:
    image: myapp:${TAG:-latest}      # Значение по умолчанию
    ports:
      - "${API_PORT}:80"
    environment:
      - DB_PASSWORD=${DB_PASSWORD}
```

### Интерполяция переменных

```yaml
services:
  web:
    # Простая подстановка
    image: "myapp:${TAG}"
    
    # Значение по умолчанию
    image: "myapp:${TAG:-latest}"
    
    # Сообщение об ошибке если не установлена
    image: "myapp:${TAG:?TAG must be set}"
    
    # Замена если пустая
    image: "myapp:${TAG:+production}"
```

---

## Полезные команды

### Отладка

```bash
# Проверить конфигурацию
docker-compose config

# Проверить с подстановкой переменных
docker-compose config --resolve-image-digests

# Какие volumes будут использованы
docker-compose config --volumes

# Показать имена контейнеров
docker-compose ps -q

# Получить IP адрес контейнера
docker-compose exec web hostname -i
```

### Очистка

```bash
# Удалить остановленные контейнеры
docker-compose rm

# Удалить orphan контейнеры
docker-compose down --remove-orphans

# Полная очистка
docker-compose down -v --rmi all --remove-orphans
```

### Работа с несколькими файлами

```bash
# Использовать несколько compose файлов
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Переменная окружения для порядка файлов
export COMPOSE_FILE=docker-compose.yml:docker-compose.override.yml:docker-compose.local.yml
docker-compose up -d
```

---

## Best Practices

### 1. Используйте version 3.8+

```yaml
version: '3.8'  # Последняя версия формата
```

### 2. Явные зависимости

```yaml
services:
  app:
    depends_on:
      db:
        condition: service_healthy  # Дождаться healthcheck
```

### 3. Healthchecks

```yaml
services:
  db:
    healthcheck:
      test: ["CMD-SHELL", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
```

### 4. Переменные окружения

```yaml
# Используйте .env файлы вместо hardcode
env_file:
  - .env
  - .env.local
```

### 5. Именованные volumes

```yaml
# ✅ Хорошо: Именованные
volumes:
  - db-data:/var/lib/postgresql/data

# ❌ Плохо: Анонимные
volumes:
  - /var/lib/postgresql/data
```

### 6. Сетевая изоляция

```yaml
networks:
  frontend:  # Публичная сеть
  backend:   # Приватная сеть
    internal: true
```

---

## Алиасы для удобства

```bash
# Добавьте в ~/.bashrc или ~/.zshrc

alias dc='docker-compose'
alias dcup='docker-compose up -d'
alias dcdown='docker-compose down'
alias dclogs='docker-compose logs -f'
alias dcps='docker-compose ps'
alias dcexec='docker-compose exec'
alias dcbuild='docker-compose build --no-cache'
alias dcrestart='docker-compose restart'
```

---

Этот cheat sheet охватывает все основные операции с Docker Compose!

