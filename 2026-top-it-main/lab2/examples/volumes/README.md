# Docker Volumes - Постоянное хранение данных

## Описание
Docker volumes используются для сохранения данных между перезапусками контейнеров и для обмена данными между контейнерами.

## Упражнение 1: Создание и использование именованных томов

### Задание
Научитесь создавать и использовать Docker volumes.

### Команды для выполнения

```bash
# 1. Просмотр существующих томов
docker volume ls

# 2. Создание нового тома
docker volume create my-data

# 3. Просмотр информации о томе
docker volume inspect my-data

# 4. Запуск контейнера с томом
docker run -d \
  --name db-container \
  -v my-data:/var/lib/postgresql/data \
  -e POSTGRES_PASSWORD=secret \
  postgres:15-alpine

# 5. Проверка, что данные сохраняются
docker exec db-container psql -U postgres -c "CREATE TABLE test (id INT, name TEXT);"
docker exec db-container psql -U postgres -c "INSERT INTO test VALUES (1, 'Docker Volumes');"

# 6. Удаление контейнера
docker stop db-container
docker rm db-container

# 7. Запуск нового контейнера с тем же томом
docker run -d \
  --name db-container-new \
  -v my-data:/var/lib/postgresql/data \
  -e POSTGRES_PASSWORD=secret \
  postgres:15-alpine

# 8. Проверка, что данные сохранились
docker exec db-container-new psql -U postgres -c "SELECT * FROM test;"

# 9. Очистка
docker stop db-container-new
docker rm db-container-new
docker volume rm my-data
```

### Вопросы для самопроверки
1. Где физически хранятся данные Docker volumes?
2. Что произойдет с данными при удалении контейнера?
3. Можно ли использовать один том для нескольких контейнеров одновременно?

---

## Упражнение 2: Bind mounts

### Задание
Изучите работу с bind mounts для монтирования директорий хоста.

### Создайте тестовые файлы

```bash
# Создать директорию
mkdir -p ~/docker-test/html
cd ~/docker-test

# Создать HTML файл
cat > html/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>Docker Bind Mount</title>
    <style>
        body {
            font-family: Arial;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            background: linear-gradient(45deg, #3498db, #8e44ad);
        }
        .container {
            background: white;
            padding: 40px;
            border-radius: 10px;
            text-align: center;
            box-shadow: 0 10px 40px rgba(0,0,0,0.3);
        }
        h1 { color: #3498db; }
        p { color: #555; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🐳 Docker Bind Mount</h1>
        <p>Этот файл монтирован с хоста!</p>
        <p>Попробуйте изменить его и обновить страницу.</p>
    </div>
</body>
</html>
EOF

# Запустить Nginx с bind mount
docker run -d \
  --name web-bind \
  -p 8080:80 \
  -v $(pwd)/html:/usr/share/nginx/html:ro \
  nginx:alpine

# Проверить в браузере: http://localhost:8080

# Изменить файл на хосте
cat > html/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>Обновлено!</title>
    <style>
        body {
            font-family: Arial;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            background: linear-gradient(45deg, #e74c3c, #f39c12);
        }
        .container {
            background: white;
            padding: 40px;
            border-radius: 10px;
            text-align: center;
            box-shadow: 0 10px 40px rgba(0,0,0,0.3);
        }
        h1 { color: #e74c3c; }
    </style>
</head>
<body>
    <div class="container">
        <h1>✨ Файл обновлен!</h1>
        <p>Изменения применились мгновенно!</p>
    </div>
</body>
</html>
EOF

# Обновить страницу в браузере - увидите изменения

# Очистка
docker stop web-bind
docker rm web-bind
cd ~
rm -rf ~/docker-test
```

### Вопросы для самопроверки
1. В чем разница между volume и bind mount?
2. Что означает флаг `:ro` в bind mount?
3. Когда лучше использовать bind mount вместо volume?

---

## Упражнение 3: Работа с несколькими томами

### Задание
Создайте приложение, использующее несколько томов для разных типов данных.

### Команды для выполнения

```bash
# Создать тома
docker volume create app-data
docker volume create app-logs
docker volume create app-config

# Просмотреть созданные тома
docker volume ls | grep app-

# Запустить контейнер с несколькими томами
docker run -d \
  --name multi-volume-app \
  -v app-data:/data \
  -v app-logs:/var/log/app \
  -v app-config:/etc/app/config \
  busybox sh -c "
    echo 'Application started' > /var/log/app/app.log;
    echo 'user=admin' > /etc/app/config/settings.conf;
    echo 'test data' > /data/test.txt;
    tail -f /var/log/app/app.log
  "

# Проверить содержимое томов через другой контейнер
docker run --rm -v app-data:/data alpine cat /data/test.txt
docker run --rm -v app-logs:/logs alpine cat /logs/app.log
docker run --rm -v app-config:/config alpine cat /config/settings.conf

# Создать резервную копию тома
docker run --rm \
  -v app-data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/app-data-backup.tar.gz -C /data .

# Восстановить из резервной копии в новый том
docker volume create app-data-restored
docker run --rm \
  -v app-data-restored:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/app-data-backup.tar.gz -C /data

# Проверить восстановленные данные
docker run --rm -v app-data-restored:/data alpine cat /data/test.txt

# Очистка
docker stop multi-volume-app
docker rm multi-volume-app
docker volume rm app-data app-logs app-config app-data-restored
rm app-data-backup.tar.gz
```

### Вопросы для самопроверки
1. Как создать резервную копию Docker volume?
2. Можно ли скопировать данные между томами?
3. Как восстановить том из резервной копии?

---

## Упражнение 4: tmpfs mounts

### Задание
Изучите временные файловые системы в памяти (tmpfs).

### Команды для выполнения

```bash
# Запустить контейнер с tmpfs
docker run -d \
  --name tmpfs-test \
  --tmpfs /tmp:rw,noexec,nosuid,size=100m \
  alpine sh -c "
    echo 'Data in memory' > /tmp/test.txt;
    df -h /tmp;
    sleep 3600
  "

# Проверить содержимое tmpfs
docker exec tmpfs-test cat /tmp/test.txt

# Проверить размер
docker exec tmpfs-test df -h /tmp

# Попробовать записать много данных
docker exec tmpfs-test sh -c "dd if=/dev/zero of=/tmp/large bs=1M count=200"
# Должна быть ошибка: No space left on device

# Перезапустить контейнер
docker restart tmpfs-test

# Проверить, что данные исчезли
docker exec tmpfs-test ls /tmp

# Очистка
docker stop tmpfs-test
docker rm tmpfs-test
```

### Использование в docker-compose

```yaml
# docker-compose-tmpfs.yml
version: '3.8'

services:
  app:
    image: alpine
    command: sh -c "echo 'test' > /tmp/data && sleep 3600"
    tmpfs:
      - /tmp:size=100M,mode=1777
      - /run:size=50M
```

### Вопросы для самопроверки
1. Когда следует использовать tmpfs вместо volume?
2. Что происходит с данными в tmpfs при перезапуске контейнера?
3. Как ограничить размер tmpfs?

---

## Упражнение 5: Обмен данными между контейнерами

### Задание
Создайте несколько контейнеров, которые обмениваются данными через общий том.

### Команды для выполнения

```bash
# Создать общий том
docker volume create shared-data

# Контейнер 1: Писатель (Writer)
docker run -d \
  --name writer \
  -v shared-data:/data \
  alpine sh -c "
    while true; do
      date >> /data/log.txt;
      echo 'Message from writer' >> /data/log.txt;
      sleep 5;
    done
  "

# Контейнер 2: Читатель (Reader)
docker run -d \
  --name reader \
  -v shared-data:/data:ro \
  alpine sh -c "
    while true; do
      echo '=== Latest logs ===';
      tail -n 5 /data/log.txt;
      sleep 10;
    done
  "

# Просмотреть логи читателя
docker logs -f reader

# Остановить через 30 секунд (Ctrl+C)

# Контейнер 3: Обработчик (Processor)
docker run -d \
  --name processor \
  -v shared-data:/data \
  alpine sh -c "
    while true; do
      if [ -f /data/log.txt ]; then
        lines=\$(wc -l < /data/log.txt);
        echo \"Total lines: \$lines\" > /data/stats.txt;
        echo \"Last update: \$(date)\" >> /data/stats.txt;
      fi
      sleep 15;
    done
  "

# Проверить статистику
sleep 20
docker exec processor cat /data/stats.txt

# Очистка
docker stop writer reader processor
docker rm writer reader processor
docker volume rm shared-data
```

### Практический пример: Приложение для обработки логов

```yaml
# docker-compose-shared-volume.yml
version: '3.8'

services:
  # Генератор логов
  log-generator:
    image: alpine
    volumes:
      - logs:/var/log
    command: >
      sh -c "while true; do
        echo \"[$$(date)] Log entry $$RANDOM\" >> /var/log/app.log;
        sleep 2;
      done"

  # Анализатор логов
  log-analyzer:
    image: alpine
    volumes:
      - logs:/var/log:ro
      - analytics:/var/analytics
    command: >
      sh -c "while true; do
        if [ -f /var/log/app.log ]; then
          grep -c 'Log entry' /var/log/app.log > /var/analytics/count.txt;
          tail -n 10 /var/log/app.log > /var/analytics/recent.txt;
        fi
        sleep 10;
      done"

  # Веб-сервер для просмотра
  web-viewer:
    image: nginx:alpine
    volumes:
      - analytics:/usr/share/nginx/html:ro
    ports:
      - "8080:80"

volumes:
  logs:
  analytics:
```

```bash
# Запустить
docker-compose -f docker-compose-shared-volume.yml up -d

# Просмотреть файлы
curl http://localhost:8080/count.txt
curl http://localhost:8080/recent.txt

# Остановить
docker-compose -f docker-compose-shared-volume.yml down -v
```

---

## Упражнение 6: Volume драйверы

### Задание
Изучите различные драйверы для томов.

### Local драйвер (по умолчанию)

```bash
# Создать том с явным указанием драйвера
docker volume create \
  --driver local \
  --opt type=none \
  --opt device=/tmp/my-volume \
  --opt o=bind \
  custom-local-volume

# Использовать том
docker run --rm -v custom-local-volume:/data alpine ls -la /data

# Удалить
docker volume rm custom-local-volume
```

### NFS volume (сетевое хранилище)

```bash
# Создать NFS том (требуется NFS сервер)
# docker volume create \
#   --driver local \
#   --opt type=nfs \
#   --opt o=addr=nfs-server.example.com,rw \
#   --opt device=:/path/to/share \
#   nfs-volume
```

### Примеры в docker-compose

```yaml
# docker-compose-volume-drivers.yml
version: '3.8'

services:
  app:
    image: alpine
    command: sh -c "echo 'test' > /data/test.txt && sleep 3600"
    volumes:
      - local-volume:/data
      - bind-volume:/bind-data

volumes:
  # Обычный именованный том
  local-volume:
    driver: local

  # Том с bind mount
  bind-volume:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /tmp/docker-bind

  # Том в памяти (tmpfs)
  # tmpfs-volume:
  #   driver: local
  #   driver_opts:
  #     type: tmpfs
  #     device: tmpfs
  #     o: size=100m

  # NFS том
  # nfs-volume:
  #   driver: local
  #   driver_opts:
  #     type: nfs
  #     o: addr=nfs-server,rw
  #     device: ":/exported/path"
```

---

## Упражнение 7: Управление жизненным циклом томов

### Задание
Изучите best practices для управления томами.

### Команды для выполнения

```bash
# 1. Просмотр всех томов с подробностями
docker volume ls
docker system df -v

# 2. Фильтрация томов
docker volume ls --filter "dangling=true"
docker volume ls --filter "name=app"

# 3. Удаление неиспользуемых томов
docker volume prune -f

# 4. Создание тома с метками
docker volume create \
  --label project=myapp \
  --label environment=production \
  myapp-prod-data

# 5. Фильтрация по меткам
docker volume ls --filter "label=project=myapp"

# 6. Инспекция тома
docker volume inspect myapp-prod-data

# 7. Очистка
docker volume rm myapp-prod-data
```

### Best Practices

```bash
# ❌ ПЛОХО: Анонимные тома
docker run -v /data nginx

# ✅ ХОРОШО: Именованные тома
docker run -v my-data:/data nginx

# ✅ ХОРОШО: Использование меток
docker volume create --label backup=daily app-data

# ✅ ХОРОШО: Регулярное резервное копирование
docker run --rm \
  -v app-data:/data:ro \
  -v $(pwd):/backup \
  alpine tar czf /backup/backup-$(date +%Y%m%d).tar.gz /data

# ✅ ХОРОШО: Очистка старых томов
docker volume prune --filter "until=720h" -f
```

---

## Итоговое задание

Создайте приложение со следующей архитектурой томов:

1. **База данных**: Используйте именованный том для PostgreSQL
2. **Логи приложения**: Используйте bind mount для доступа к логам с хоста
3. **Кэш**: Используйте tmpfs для временных данных
4. **Конфигурация**: Используйте read-only bind mount для конфигурационных файлов
5. **Общие данные**: Создайте том для обмена данными между микросервисами

Требования:
- Реализовать через docker-compose
- Добавить метки для всех томов
- Создать скрипт для резервного копирования
- Настроить автоматическую очистку логов

## Чек-лист освоенных навыков

- [ ] Создание именованных томов
- [ ] Использование bind mounts
- [ ] Работа с tmpfs
- [ ] Обмен данными между контейнерами
- [ ] Создание резервных копий томов
- [ ] Восстановление данных из резервных копий
- [ ] Использование различных volume драйверов
- [ ] Фильтрация и управление томами
- [ ] Применение меток к томам
- [ ] Очистка неиспользуемых томов
- [ ] Настройка read-only доступа
- [ ] Ограничение размера tmpfs

## Полезные команды

```bash
# Просмотр использования места
docker system df -v | grep -A 20 "Local Volumes"

# Экспорт тома
docker run --rm -v SOURCE:/from -v $(pwd):/to alpine cp -av /from /to

# Клонирование тома
docker volume create new-volume
docker run --rm -v old-volume:/from -v new-volume:/to alpine cp -av /from/. /to/

# Очистка всех томов (ОСТОРОЖНО!)
docker volume rm $(docker volume ls -q)

# Инспекция точки монтирования
docker volume inspect VOLUME_NAME --format '{{ .Mountpoint }}'
```

