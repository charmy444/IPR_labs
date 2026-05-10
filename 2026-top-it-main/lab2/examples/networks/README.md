# Docker Networks - Сетевое взаимодействие контейнеров

## Описание
Docker networks позволяют контейнерам взаимодействовать друг с другом и с внешним миром.

## Упражнение 1: Типы сетей Docker

### Задание
Изучите основные типы сетей в Docker.

### Команды для выполнения

```bash
# 1. Просмотр существующих сетей
docker network ls

# 2. Подробная информация о сети
docker network inspect bridge

# 3. Создание bridge сети
docker network create my-bridge-network

# 4. Создание сетей разных типов
docker network create --driver bridge app-network
docker network create --driver overlay swarm-network 2>/dev/null || echo "Требуется Docker Swarm"
docker network create --driver host host-network 2>/dev/null || echo "host драйвер не поддерживается так"

# 5. Просмотр созданных сетей
docker network ls --filter driver=bridge

# 6. Удаление сети
docker network rm my-bridge-network

# 7. Очистка неиспользуемых сетей
docker network prune -f
```

### Типы сетей Docker

| Тип | Описание | Использование |
|-----|----------|---------------|
| **bridge** | Частная сеть внутри хоста | По умолчанию для контейнеров |
| **host** | Использует сеть хоста | Высокая производительность |
| **none** | Нет сетевого интерфейса | Изолированные контейнеры |
| **overlay** | Мульти-хост сеть | Docker Swarm кластеры |
| **macvlan** | MAC адрес контейнера | Legacy приложения |

### Вопросы для самопроверки
1. В чем разница между bridge и host сетью?
2. Когда использовать overlay сеть?
3. Что такое default bridge сеть?

---

## Упражнение 2: Контейнеры в одной сети

### Задание
Создайте два контейнера в одной сети и настройте их взаимодействие.

### Команды для выполнения

```bash
# 1. Создать пользовательскую bridge сеть
docker network create app-net

# 2. Запустить первый контейнер (backend)
docker run -d \
  --name backend \
  --network app-net \
  alpine sleep 3600

# 3. Запустить второй контейнер (frontend)
docker run -d \
  --name frontend \
  --network app-net \
  alpine sleep 3600

# 4. Проверить сетевые настройки
docker network inspect app-net

# 5. Проверить связь между контейнерами по имени
docker exec frontend ping -c 3 backend

# 6. Проверить обратную связь
docker exec backend ping -c 3 frontend

# 7. Запустить веб-сервер в backend
docker exec -d backend sh -c "
  echo 'Hello from backend!' > /tmp/index.html;
  cd /tmp && nohup busybox httpd -f -p 8000 &
"

# 8. Получить ответ от backend через frontend
docker exec frontend wget -q -O- http://backend:8000/index.html

# 9. Попробовать подключиться к контейнеру вне сети
docker run --rm alpine ping -c 1 backend
# Должна быть ошибка: bad address

# 10. Очистка
docker stop backend frontend
docker rm backend frontend
docker network rm app-net
```

### Вопросы для самопроверки
1. Почему контейнеры могут общаться по именам?
2. Работает ли DNS в default bridge сети?
3. Как контейнер получает IP адрес?

---

## Упражнение 3: Множественные сети

### Задание
Подключите контейнер к нескольким сетям одновременно.

### Команды для выполнения

```bash
# 1. Создать две сети
docker network create frontend-net
docker network create backend-net

# 2. Запустить базу данных только в backend сети
docker run -d \
  --name database \
  --network backend-net \
  -e POSTGRES_PASSWORD=secret \
  postgres:15-alpine

# 3. Запустить API в обеих сетях
docker run -d \
  --name api \
  --network backend-net \
  alpine sleep 3600

docker network connect frontend-net api

# 4. Запустить frontend только в frontend сети
docker run -d \
  --name web \
  --network frontend-net \
  nginx:alpine

# 5. Проверить подключения
docker network inspect frontend-net --format '{{range .Containers}}{{.Name}} {{end}}'
docker network inspect backend-net --format '{{range .Containers}}{{.Name}} {{end}}'

# 6. API может общаться с базой данных
docker exec api ping -c 2 database

# 7. Web может общаться с API
docker exec web ping -c 2 api

# 8. Web НЕ может общаться с базой данных
docker exec web ping -c 1 database 2>&1 | grep -q "bad address" && echo "✓ Доступ запрещен (как и должно быть)"

# 9. Отключить API от frontend сети
docker network disconnect frontend-net api

# 10. Теперь web не может общаться с API
docker exec web ping -c 1 api 2>&1 | grep -q "bad address" && echo "✓ Доступ запрещен"

# 11. Очистка
docker stop database api web
docker rm database api web
docker network rm frontend-net backend-net
```

### Архитектура

```
┌─────────────────┐
│   frontend-net  │
│                 │
│  ┌─────┐ ┌────┐│
│  │ web │─│api ││─┐
│  └─────┘ └────┘│ │
└─────────────────┘ │
                    │
┌───────────────────│─┐
│   backend-net     │ │
│                   │ │
│  ┌────┐ ┌────────┴─┴┐
│  │ db │─│    api    │
│  └────┘ └───────────┘
└─────────────────────┘
```

### Вопросы для самопроверки
1. Сколько сетей может иметь один контейнер?
2. Как контейнер узнает, в каких сетях он находится?
3. Зачем разделять контейнеры по разным сетям?

---

## Упражнение 4: Проброс портов

### Задание
Изучите различные способы публикации портов контейнера.

### Команды для выполнения

```bash
# 1. Проброс на определенный порт хоста
docker run -d --name web1 -p 8080:80 nginx:alpine

# 2. Проброс на случайный порт хоста
docker run -d --name web2 -p 80 nginx:alpine

# 3. Проброс на определенный IP и порт
docker run -d --name web3 -p 127.0.0.1:8081:80 nginx:alpine

# 4. Проброс нескольких портов
docker run -d --name web4 \
  -p 8082:80 \
  -p 8443:443 \
  nginx:alpine

# 5. Проброс UDP порта
docker run -d --name dns -p 53:53/udp alpine sleep 3600

# 6. Просмотр проброшенных портов
docker ps --format "table {{.Names}}\t{{.Ports}}"

# 7. Получить конкретный порт
docker port web2 80

# 8. Проверить доступность
curl http://localhost:8080
curl http://localhost:8081

# 9. Узнать порт web2
WEB2_PORT=$(docker port web2 80 | cut -d: -f2)
curl http://localhost:$WEB2_PORT

# 10. Очистка
docker stop web1 web2 web3 web4 dns
docker rm web1 web2 web3 web4 dns
```

### Форматы проброса портов

```bash
# Синтаксис: -p [HOST_IP:]HOST_PORT:CONTAINER_PORT[/PROTOCOL]

# Примеры:
-p 8080:80                    # Порт 8080 хоста → порт 80 контейнера
-p 127.0.0.1:8080:80         # Только локальный доступ
-p 8080:80/tcp               # Явно указан TCP (по умолчанию)
-p 53:53/udp                 # UDP порт
-p 8080-8090:8080-8090       # Диапазон портов
```

### Вопросы для самопроверки
1. Можно ли пробросить один порт хоста на несколько контейнеров?
2. Как узнать, какой случайный порт был назначен?
3. В чем разница между `-p` и `--expose`?

---

## Упражнение 5: DNS и Service Discovery

### Задание
Изучите механизм DNS в Docker сетях.

### Команды для выполнения

```bash
# 1. Создать сеть
docker network create --driver bridge \
  --subnet 172.25.0.0/16 \
  --ip-range 172.25.5.0/24 \
  --gateway 172.25.5.254 \
  custom-network

# 2. Запустить контейнеры с заданными IP
docker run -d \
  --name service1 \
  --network custom-network \
  --ip 172.25.5.10 \
  --network-alias api \
  --network-alias backend \
  alpine sleep 3600

docker run -d \
  --name service2 \
  --network custom-network \
  --ip 172.25.5.20 \
  --network-alias api \
  alpine sleep 3600

docker run -d \
  --name client \
  --network custom-network \
  alpine sleep 3600

# 3. Проверить DNS разрешение по имени контейнера
docker exec client nslookup service1

# 4. Проверить DNS разрешение по алиасу
docker exec client nslookup api

# 5. Получить все IP для алиаса (round-robin)
docker exec client sh -c "
  for i in 1 2 3 4 5; do
    nslookup api | grep 'Address:' | tail -1
  done
"

# 6. Проверить /etc/hosts в контейнере
docker exec client cat /etc/hosts

# 7. Проверить /etc/resolv.conf
docker exec client cat /etc/resolv.conf

# 8. Ping по алиасу
docker exec client ping -c 2 api

# 9. Очистка
docker stop service1 service2 client
docker rm service1 service2 client
docker network rm custom-network
```

### Вопросы для самопроверки
1. Как работает DNS в пользовательских bridge сетях?
2. Что такое network alias?
3. Можно ли несколько контейнеров иметь один алиас?

---

## Упражнение 6: Изоляция сетей

### Задание
Создайте изолированные сети для различных окружений.

### docker-compose.yml

```yaml
# docker-compose-networks.yml
version: '3.8'

services:
  # Производственная БД
  prod-db:
    image: postgres:15-alpine
    container_name: prod-database
    environment:
      POSTGRES_PASSWORD: prod-secret
      POSTGRES_DB: proddb
    networks:
      - production
    volumes:
      - prod-db-data:/var/lib/postgresql/data

  # Тестовая БД
  test-db:
    image: postgres:15-alpine
    container_name: test-database
    environment:
      POSTGRES_PASSWORD: test-secret
      POSTGRES_DB: testdb
    networks:
      - testing
    volumes:
      - test-db-data:/var/lib/postgresql/data

  # Production API
  prod-api:
    image: alpine
    container_name: prod-api
    command: sh -c "apk add --no-cache postgresql-client && sleep 3600"
    networks:
      - production
      - frontend
    depends_on:
      - prod-db

  # Test API
  test-api:
    image: alpine
    container_name: test-api
    command: sh -c "apk add --no-cache postgresql-client && sleep 3600"
    networks:
      - testing
    depends_on:
      - test-db

  # Frontend (доступ только к prod-api)
  frontend:
    image: nginx:alpine
    container_name: frontend-web
    networks:
      - frontend
    ports:
      - "8080:80"
    depends_on:
      - prod-api

networks:
  production:
    driver: bridge
    internal: false
    ipam:
      config:
        - subnet: 172.20.0.0/16
  
  testing:
    driver: bridge
    internal: true  # Нет доступа к внешней сети
    ipam:
      config:
        - subnet: 172.21.0.0/16
  
  frontend:
    driver: bridge
    ipam:
      config:
        - subnet: 172.22.0.0/16

volumes:
  prod-db-data:
  test-db-data:
```

### Команды для выполнения

```bash
# Создать файл docker-compose-networks.yml (см. выше)

# Запустить
docker-compose -f docker-compose-networks.yml up -d

# Проверить сети
docker network ls | grep docker-compose-networks

# prod-api может подключиться к prod-db
docker exec prod-api psql -h prod-db -U postgres -d proddb -c "SELECT version();"

# test-api может подключиться к test-db
docker exec test-api psql -h test-db -U postgres -d testdb -c "SELECT version();"

# prod-api НЕ может подключиться к test-db
docker exec prod-api ping -c 1 test-db 2>&1 | grep -q "bad address" && echo "✓ Изоляция работает"

# test-api НЕ может подключиться к интернету (internal network)
docker exec test-api ping -c 1 8.8.8.8 2>&1 | grep -q "Network is unreachable" && echo "✓ Internal сеть работает"

# frontend может подключиться к prod-api
docker exec frontend ping -c 2 prod-api

# Просмотреть топологию сетей
docker network inspect docker-compose-networks_production --format '{{range .Containers}}{{.Name}} {{end}}'
docker network inspect docker-compose-networks_testing --format '{{range .Containers}}{{.Name}} {{end}}'
docker network inspect docker-compose-networks_frontend --format '{{range .Containers}}{{.Name}} {{end}}'

# Остановить
docker-compose -f docker-compose-networks.yml down -v
```

### Вопросы для самопроверки
1. Что означает `internal: true` для сети?
2. Зачем разделять production и testing сети?
3. Как ограничить доступ между сетями?

---

## Упражнение 7: Мониторинг сетевого трафика

### Задание
Изучите сетевую активность контейнеров.

### Команды для выполнения

```bash
# 1. Создать сеть и запустить контейнеры
docker network create monitor-net

docker run -d \
  --name server \
  --network monitor-net \
  nginx:alpine

docker run -d \
  --name client \
  --network monitor-net \
  alpine sleep 3600

# 2. Установить инструменты мониторинга
docker exec server apk add --no-cache tcpdump iftop
docker exec client apk add --no-cache curl

# 3. Запустить tcpdump в одном терминале
# docker exec server tcpdump -i eth0 -n

# 4. Генерировать трафик
docker exec client sh -c "
  for i in {1..10}; do
    curl -s http://server > /dev/null;
    sleep 1;
  done
"

# 5. Проверить сетевую статистику контейнера
docker stats server --no-stream

# 6. Просмотреть сетевые интерфейсы
docker exec server ifconfig
docker exec server ip addr show

# 7. Проверить сетевые подключения
docker exec server netstat -an 2>/dev/null || docker exec server ss -an

# 8. Просмотреть таблицу маршрутизации
docker exec server route -n
docker exec server ip route

# 9. Проверить открытые порты
docker exec server netstat -tlnp 2>/dev/null || docker exec server ss -tlnp

# 10. Очистка
docker stop server client
docker rm server client
docker network rm monitor-net
```

### Полезные команды для отладки сети

```bash
# Проверить связность
docker exec CONTAINER ping -c 3 TARGET

# Проверить открытые порты
docker exec CONTAINER nc -zv HOST PORT

# Трассировка маршрута
docker exec CONTAINER traceroute TARGET

# DNS lookup
docker exec CONTAINER nslookup HOSTNAME
docker exec CONTAINER dig HOSTNAME

# Проверить HTTP подключение
docker exec CONTAINER curl -v http://TARGET

# Просмотр сетевых подключений
docker exec CONTAINER netstat -tupn
```

---

## Итоговое задание

Создайте микросервисную архитектуру с правильной сетевой сегментацией:

1. **Public Network**: Nginx reverse proxy
2. **Application Network**: Backend API сервисы
3. **Database Network**: PostgreSQL, Redis
4. **Monitoring Network**: Prometheus, Grafana
5. **Admin Network**: Административные инструменты

Требования:
- Frontend имеет доступ только к reverse proxy
- API сервисы изолированы от прямого доступа
- База данных доступна только для API
- Monitoring может получать метрики от всех сервисов
- Используйте custom DNS aliases
- Настройте правильные subnet для каждой сети
- Добавьте healthcheck endpoints

## Чек-лист освоенных навыков

- [ ] Понимание типов сетей Docker
- [ ] Создание пользовательских bridge сетей
- [ ] Подключение контейнеров к нескольким сетям
- [ ] Настройка проброса портов
- [ ] Работа с DNS и network aliases
- [ ] Создание изолированных (internal) сетей
- [ ] Настройка подсетей и IP адресов
- [ ] Мониторинг сетевого трафика
- [ ] Отладка сетевых проблем
- [ ] Создание сетевой сегментации в docker-compose

## Полезные команды

```bash
# Информация о всех сетях
docker network ls --format "table {{.Name}}\t{{.Driver}}\t{{.Scope}}"

# Найти сети контейнера
docker inspect CONTAINER --format '{{range $k, $v := .NetworkSettings.Networks}}{{$k}} {{end}}'

# Получить IP адрес контейнера
docker inspect CONTAINER --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'

# Удалить все пользовательские сети
docker network prune -f

# Создать сеть с IPv6
docker network create --ipv6 --subnet=2001:db8:1::/64 my-ipv6-net
```

