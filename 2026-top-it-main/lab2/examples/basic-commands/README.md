# Базовые команды Docker

## Описание
В этом разделе вы изучите основные команды Docker для работы с контейнерами и образами.

## Упражнение 1: Работа с образами

### Задание
Освойте команды для работы с Docker образами.

### Команды для выполнения

```bash
# 1. Проверить версию Docker
docker --version
docker version

# 2. Получить информацию о системе Docker
docker info

# 3. Загрузить образ из Docker Hub
docker pull nginx:latest
docker pull python:3.11-slim
docker pull node:18-alpine

# 4. Просмотреть список локальных образов
docker images
docker image ls

# 5. Просмотреть подробную информацию об образе
docker image inspect nginx:latest

# 6. Просмотреть историю образа (слои)
docker image history nginx:latest

# 7. Удалить образ
docker rmi nginx:latest
# или
docker image rm nginx:latest
```

### Вопросы для самопроверки
1. Чем отличается `docker images` от `docker image ls`?
2. Что показывает команда `docker image history`?
3. Почему нельзя удалить образ, если он используется контейнером?

---

## Упражнение 2: Запуск и управление контейнерами

### Задание
Научитесь запускать контейнеры и управлять ими.

### Команды для выполнения

```bash
# 1. Запустить контейнер в фоновом режиме
docker run -d --name my-nginx nginx:latest

# 2. Запустить контейнер с проброшенными портами
docker run -d -p 8080:80 --name web-server nginx:latest

# 3. Запустить контейнер с переменными окружения
docker run -d -e MYSQL_ROOT_PASSWORD=secret --name my-db mysql:8

# 4. Запустить контейнер в интерактивном режиме
docker run -it ubuntu:22.04 /bin/bash
# Внутри контейнера выполните: apt update && apt install -y curl
# Выйдите: exit

# 5. Просмотреть список запущенных контейнеров
docker ps

# 6. Просмотреть все контейнеры (включая остановленные)
docker ps -a

# 7. Остановить контейнер
docker stop my-nginx

# 8. Запустить остановленный контейнер
docker start my-nginx

# 9. Перезапустить контейнер
docker restart my-nginx

# 10. Удалить контейнер
docker stop my-nginx
docker rm my-nginx

# 11. Удалить контейнер принудительно
docker rm -f web-server
```

### Вопросы для самопроверки
1. В чем разница между `docker run` и `docker start`?
2. Что означает флаг `-d` в команде `docker run`?
3. Как пробросить порт 3000 контейнера на порт 8080 хоста?

---

## Упражнение 3: Взаимодействие с контейнерами

### Задание
Научитесь взаимодействовать с запущенными контейнерами.

### Команды для выполнения

```bash
# 1. Запустить контейнер для экспериментов
docker run -d --name test-container ubuntu:22.04 sleep infinity

# 2. Выполнить команду внутри контейнера
docker exec test-container ls -la /

# 3. Открыть интерактивную оболочку в контейнере
docker exec -it test-container /bin/bash
# Выйдите: exit

# 4. Просмотреть логи контейнера
docker logs test-container

# 5. Следить за логами в реальном времени
docker logs -f test-container
# Остановите: Ctrl+C

# 6. Просмотреть последние 50 строк логов
docker logs --tail 50 test-container

# 7. Просмотреть статистику использования ресурсов
docker stats test-container
# Остановите: Ctrl+C

# 8. Просмотреть информацию о контейнере
docker inspect test-container

# 9. Просмотреть процессы внутри контейнера
docker top test-container

# 10. Скопировать файл из контейнера на хост
docker exec test-container sh -c "echo 'Hello Docker' > /tmp/test.txt"
docker cp test-container:/tmp/test.txt ./test.txt
cat test.txt

# 11. Скопировать файл с хоста в контейнер
echo "Hello from host" > host-file.txt
docker cp host-file.txt test-container:/tmp/
docker exec test-container cat /tmp/host-file.txt

# 12. Очистка
docker rm -f test-container
rm test.txt host-file.txt
```

### Вопросы для самопроверки
1. В чем разница между `docker exec` и `docker run`?
2. Как просмотреть логи контейнера за последние 5 минут?
3. Что покажет команда `docker inspect`?

---

## Упражнение 4: Очистка системы

### Задание
Научитесь очищать Docker систему от неиспользуемых ресурсов.

### Команды для выполнения

```bash
# 1. Создать несколько тестовых контейнеров
docker run -d --name test1 nginx:alpine
docker run -d --name test2 nginx:alpine
docker run -d --name test3 nginx:alpine

# 2. Остановить контейнеры
docker stop test1 test2 test3

# 3. Удалить все остановленные контейнеры
docker container prune -f

# 4. Удалить неиспользуемые образы
docker image prune -f

# 5. Удалить все неиспользуемые образы (не только dangling)
docker image prune -a -f

# 6. Удалить неиспользуемые тома
docker volume prune -f

# 7. Удалить неиспользуемые сети
docker network prune -f

# 8. Удалить все неиспользуемые ресурсы одной командой
docker system prune -f

# 9. Удалить все (включая тома)
docker system prune -a --volumes -f

# 10. Просмотреть использование дискового пространства Docker
docker system df

# 11. Подробный отчет об использовании места
docker system df -v
```

### Вопросы для самопроверки
1. Что такое "dangling" образы?
2. В чем разница между `docker image prune` и `docker image prune -a`?
3. Безопасно ли выполнять `docker system prune -a --volumes`?

---

## Упражнение 5: Поиск и работа с Docker Hub

### Задание
Научитесь искать образы и работать с Docker Hub.

### Команды для выполнения

```bash
# 1. Поиск образов на Docker Hub
docker search nginx

# 2. Поиск с ограничением результатов
docker search --limit 5 python

# 3. Поиск только официальных образов
docker search --filter "is-official=true" redis

# 4. Загрузить конкретную версию образа
docker pull redis:7.2-alpine

# 5. Загрузить все теги образа (не рекомендуется)
# docker pull -a redis  # Осторожно: загрузит все версии!

# 6. Просмотреть доступные теги на Docker Hub (через браузер)
# https://hub.docker.com/_/redis?tab=tags

# 7. Создать свой тег для образа
docker tag redis:7.2-alpine my-redis:latest

# 8. Просмотреть образы
docker images | grep redis

# 9. Удалить теги
docker rmi my-redis:latest
docker rmi redis:7.2-alpine
```

### Практическое задание
Найдите на Docker Hub образ PostgreSQL версии 15, загрузите его и запустите контейнер с базой данных, используя переменные окружения для установки пароля.

### Вопросы для самопроверки
1. Как узнать, какие теги доступны для образа?
2. В чем разница между официальными и неофициальными образами?
3. Что означает тег `alpine` в названии образа?

---

## Дополнительные команды для продвинутых пользователей

```bash
# Экспорт и импорт контейнеров
docker export my-container > container.tar
docker import container.tar my-image:latest

# Сохранение и загрузка образов
docker save nginx:latest > nginx.tar
docker load < nginx.tar

# Просмотр изменений в файловой системе контейнера
docker diff container-name

# Создание образа из контейнера
docker commit container-name new-image:tag

# Ограничение ресурсов при запуске
docker run -d --memory="512m" --cpus="1.0" nginx:alpine

# Запуск с автоматическим перезапуском
docker run -d --restart=always nginx:alpine

# Просмотр событий Docker в реальном времени
docker events

# Пауза и возобновление контейнера
docker pause container-name
docker unpause container-name
```

## Итоговое задание

Выполните следующую последовательность действий:

1. Загрузите образ `python:3.11-slim`
2. Запустите контейнер с этим образом в интерактивном режиме
3. Установите внутри контейнера пакет `requests` через pip
4. Создайте из этого контейнера новый образ с именем `my-python:latest`
5. Запустите новый контейнер из созданного образа
6. Проверьте, что пакет `requests` установлен
7. Очистите все созданные ресурсы

### Решение (попробуйте сначала самостоятельно!)

```bash
# Шаг 1
docker pull python:3.11-slim

# Шаг 2
docker run -it --name python-custom python:3.11-slim /bin/bash

# Шаг 3 (внутри контейнера)
pip install requests
exit

# Шаг 4
docker commit python-custom my-python:latest

# Шаг 5
docker run -it --rm my-python:latest python -c "import requests; print(requests.__version__)"

# Шаг 6 (должна вывестись версия requests)

# Шаг 7
docker rm python-custom
docker rmi my-python:latest
docker rmi python:3.11-slim
```

## Чек-лист освоенных навыков

- [ ] Установка и проверка Docker
- [ ] Загрузка образов из Docker Hub
- [ ] Просмотр и удаление образов
- [ ] Запуск контейнеров в различных режимах
- [ ] Управление жизненным циклом контейнеров
- [ ] Выполнение команд внутри контейнеров
- [ ] Просмотр логов и статистики
- [ ] Копирование файлов между хостом и контейнером
- [ ] Очистка системы от неиспользуемых ресурсов
- [ ] Поиск образов на Docker Hub
- [ ] Работа с тегами образов

