# Лабораторная работа №2: Docker - Контейнеризация приложений

## Цель работы

Научиться создавать Docker образы для собственных приложений, работать с Docker Hub и контейнеризировать существующие проекты.

## Требования

- Docker установлен и работает
- Аккаунт на [Docker Hub](https://hub.docker.com/)
- Базовые знания работы с контейнерами (из предыдущих занятий)

---

## Содержание

- [Часть 1: Основные команды Docker](#часть-1-основные-команды-docker)
- [Часть 2: Создание Dockerfile](#часть-2-создание-dockerfile)
- [Часть 3: ENTRYPOINT vs CMD](#часть-3-entrypoint-vs-cmd)
- [Часть 4: Работа с несколькими базовыми образами](#часть-4-работа-с-несколькими-базовыми-образами)
- [Часть 5: Дебаг контейнеров](#часть-5-дебаг-контейнеров)
- [Часть 6: Docker Hub](#часть-6-docker-hub)
- [Часть 7: Оптимизация образов](#часть-7-оптимизация-образов)
- [Часть 8: Bind mounts для разработки](#часть-8-bind-mounts-для-разработки)
- [Часть 9: Сборка from scratch (опционально)](#часть-9-сборка-from-scratch-опционально)
- [Задание для самостоятельной работы](#задание-для-самостоятельной-работы)
- [Контрольные вопросы](#контрольные-вопросы)

---

## Часть 1: Основные команды Docker

### Команда: `docker run`

**Что делает:** Создает и запускает новый контейнер из образа.

**Синтаксис:**
```bash
docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
```

**Основные опции:**

| Опция | Что делает | Пример |
|-------|-----------|--------|
| `-d` | Запуск в фоновом режиме (detached) | `docker run -d nginx` |
| `-it` | Интерактивный режим с терминалом | `docker run -it ubuntu bash` |
| `--name` | Задать имя контейнеру | `docker run --name web nginx` |
| `-p` | Проброс порта (host:container) | `docker run -p 8080:80 nginx` |
| `-v` | Монтирование volume | `docker run -v /data:/app/data nginx` |
| `-e` | Установить переменную окружения | `docker run -e DEBUG=true app` |
| `--rm` | Автоудаление контейнера после остановки | `docker run --rm ubuntu echo hi` |

**Примеры:**

```bash
# 1. Простой запуск
docker run hello-world
# Что происходит: Скачивает образ (если нет) → Создает контейнер → Запускает → Выводит текст → Останавливается

# 2. Интерактивный режим
docker run -it ubuntu:22.04 bash
# Что происходит: Запускает Ubuntu и открывает bash терминал
# Вы внутри контейнера! Можете выполнять команды: ls, pwd, apt update
# Выход: exit

# 3. Фоновый режим (daemon)
docker run -d --name my-nginx nginx:latest
# Что происходит: Запускает Nginx в фоне, возвращает ID контейнера
# Контейнер продолжает работать

# 4. С пробросом порта
docker run -d -p 8080:80 nginx
# Что происходит: Nginx внутри слушает порт 80
#                  Снаружи доступен на порту 8080
# Откройте браузер: http://localhost:8080

# 5. С переменными окружения
docker run -d -e MYSQL_ROOT_PASSWORD=secret mysql:8
# Что происходит: Запускает MySQL с паролем root = "secret"

# 6. Автоудаление после выполнения
docker run --rm ubuntu echo "Hello Docker"
# Что происходит: Выводит "Hello Docker" и сразу удаляет контейнер
```

**Когда использовать:**
- `docker run` - когда нужно создать НОВЫЙ контейнер
- Каждый `docker run` создает новый контейнер!

---

### Команда: `docker ps`

**Что делает:** Показывает список контейнеров.

**Синтаксис:**
```bash
docker ps [OPTIONS]
```

**Опции:**

| Опция | Что делает |
|-------|-----------|
| *(без опций)* | Показать только запущенные контейнеры |
| `-a` | Показать ВСЕ контейнеры (и остановленные) |
| `-q` | Показать только ID контейнеров |
| `--filter` | Фильтровать контейнеры |
| `--format` | Форматировать вывод |

**Примеры:**

```bash
# Запущенные контейнеры
docker ps

# Все контейнеры
docker ps -a

# Только ID запущенных
docker ps -q

# Фильтр по имени
docker ps --filter "name=web"

# Фильтр по статусу
docker ps --filter "status=exited"

# Кастомный формат
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

**Что показывает вывод:**
- **CONTAINER ID** - уникальный ID контейнера (первые 12 символов)
- **IMAGE** - из какого образа создан
- **COMMAND** - какая команда выполняется
- **CREATED** - когда создан
- **STATUS** - статус (Up = запущен, Exited = остановлен)
- **PORTS** - пробрасываемые порты
- **NAMES** - имя контейнера

---

### Команда: `docker stop`

**Что делает:** Останавливает запущенный контейнер.

**Как работает:**
1. Отправляет сигнал SIGTERM (мягкая остановка)
2. Ждет 10 секунд
3. Если не остановился - отправляет SIGKILL (жесткая остановка)

**Примеры:**

```bash
# Остановить по имени
docker stop my-nginx

# Остановить по ID
docker stop abc123def456

# Остановить с таймаутом
docker stop -t 30 my-nginx
# Ждет 30 секунд перед SIGKILL

# Остановить все запущенные
docker stop $(docker ps -q)
```

**Когда использовать:**
- Для корректной остановки контейнера
- Данные в volumes сохраняются
- Контейнер можно запустить снова через `docker start`

---

### Команда: `docker start`

**Что делает:** Запускает остановленный контейнер.

**Отличие от `docker run`:**
- `docker run` - создает НОВЫЙ контейнер
- `docker start` - запускает СУЩЕСТВУЮЩИЙ контейнер

**Примеры:**

```bash
# Запустить остановленный контейнер
docker start my-nginx

# Запустить и присоединиться к выводу
docker start -a my-nginx

# Запустить интерактивно
docker start -ai my-ubuntu
```

---

### Команда: `docker exec`

**Что делает:** Выполняет команду ВНУТРИ запущенного контейнера.

**Синтаксис:**
```bash
docker exec [OPTIONS] CONTAINER COMMAND [ARG...]
```

**Опции:**

| Опция | Что делает |
|-------|-----------|
| `-it` | Интерактивный режим |
| `-d` | Выполнить в фоне |
| `-u` | От имени какого пользователя |
| `-e` | Установить переменную окружения |
| `-w` | Рабочая директория |

**Примеры:**

```bash
# Выполнить команду
docker exec my-nginx nginx -v

# Открыть shell (bash)
docker exec -it my-nginx /bin/bash

# Открыть shell (sh) - для alpine
docker exec -it my-alpine /bin/sh

# Посмотреть файлы
docker exec my-nginx ls -la /usr/share/nginx/html

# Выполнить от root
docker exec -it -u root my-nginx bash

# Установить переменную окружения
docker exec -it -e DEBUG=true my-app bash
```

**Когда использовать:**
- Для отладки запущенного контейнера
- Для выполнения команд внутри контейнера
- Контейнер ДОЛЖЕН быть запущен!

---

### Команда: `docker logs`

**Что делает:** Показывает логи контейнера (stdout и stderr).

**Опции:**

| Опция | Что делает |
|-------|-----------|
| `-f` | Follow - следить в реальном времени |
| `--tail N` | Показать последние N строк |
| `--since` | Логи с определенного времени |
| `--until` | Логи до определенного времени |
| `-t` | Показать timestamp |

**Примеры:**

```bash
# Все логи
docker logs my-nginx

# Последние 100 строк
docker logs --tail 100 my-nginx

# Следить в реальном времени
docker logs -f my-nginx
# Остановить: Ctrl+C

# С временными метками
docker logs -t my-nginx

# Логи за последние 30 минут
docker logs --since 30m my-nginx

# Логи с определенной даты
docker logs --since 2024-10-17T10:00:00 my-nginx
```

**Когда использовать:**
- Для отладки проблем
- Для мониторинга работы приложения
- Работает даже если контейнер остановлен!

---

### Команда: `docker rm`

**Что делает:** Удаляет контейнер.

**Важно:** Контейнер должен быть остановлен!

**Примеры:**

```bash
# Удалить остановленный контейнер
docker rm my-nginx

# Принудительно удалить (даже запущенный)
docker rm -f my-nginx

# Удалить несколько контейнеров
docker rm container1 container2 container3

# Удалить все остановленные контейнеры
docker rm $(docker ps -aq)

# Или более безопасно
docker container prune
```

---

### Команда: `docker images`

**Что делает:** Показывает список локальных образов.

**Примеры:**

```bash
# Все образы
docker images

# Фильтр по имени
docker images nginx

# Только ID
docker images -q

# Показать dangling образы (без тега)
docker images -f "dangling=true"
```

---

### Команда: `docker pull`

**Что делает:** Загружает образ из Docker Hub (или другого registry).

**Примеры:**

```bash
# Загрузить latest версию
docker pull nginx

# Загрузить конкретную версию
docker pull nginx:1.25-alpine

# Загрузить с другого registry
docker pull gcr.io/google-containers/busybox
```

**Когда использовать:**
- Чтобы загрузить образ до запуска контейнера
- `docker run` автоматически делает `pull`, если образа нет локально

---

### Команда: `docker build`

**Что делает:** Собирает образ из Dockerfile.

**Синтаксис:**
```bash
docker build [OPTIONS] PATH
```

**Опции:**

| Опция | Что делает |
|-------|-----------|
| `-t` | Задать имя и тег образу |
| `-f` | Указать другой Dockerfile |
| `--no-cache` | Не использовать кэш |
| `--build-arg` | Передать аргументы сборки |
| `--target` | Собрать конкретный stage (multi-stage) |

**Примеры:**

```bash
# Собрать образ из текущей директории
docker build -t myapp:latest .

# С другим Dockerfile
docker build -t myapp -f Dockerfile.prod .

# Без кэша
docker build --no-cache -t myapp .

# С аргументами
docker build --build-arg VERSION=1.0 -t myapp .
```

---

### Команда: `docker attach`

**Что делает:** Присоединяется к STDIN/STDOUT/STDERR контейнера.

**Отличие от `docker exec`:**
- `attach` - подключается к главному процессу контейнера
- `exec` - запускает НОВЫЙ процесс в контейнере

**Примеры:**

```bash
# Присоединиться к контейнеру
docker attach my-container

# Отключиться БЕЗ остановки контейнера:
# Нажать: Ctrl+P, затем Ctrl+Q
```

**Когда использовать:**
- Для просмотра вывода приложения
- Редко используется, обычно используют `docker logs` или `docker exec`

---

### Команда: `docker inspect`

**Что делает:** Показывает подробную информацию о контейнере или образе в JSON формате.

**Примеры:**

```bash
# Вся информация
docker inspect my-nginx

# Конкретное поле
docker inspect --format='{{.State.Status}}' my-nginx

# IP адрес контейнера
docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' my-nginx

# Переменные окружения
docker inspect --format='{{range .Config.Env}}{{println .}}{{end}}' my-nginx
```

---

## Практическое задание 1.1: Работа с командами

Выполните следующую последовательность команд и запишите результаты:

```bash
# 1. Загрузить образ
docker pull nginx:alpine

# 2. Посмотреть образы
docker images

# 3. Запустить контейнер
docker run -d --name test-nginx -p 8080:80 nginx:alpine

# 4. Проверить запущенные контейнеры
docker ps

# 5. Посмотреть логи
docker logs test-nginx

# 6. Выполнить команду внутри
docker exec test-nginx ls /usr/share/nginx/html

# 7. Открыть shell
docker exec -it test-nginx /bin/sh
# Внутри выполните: ls, pwd, ps aux, exit

# 8. Посмотреть детальную информацию
docker inspect test-nginx

# 9. Остановить
docker stop test-nginx

# 10. Посмотреть все контейнеры
docker ps -a

# 11. Запустить снова
docker start test-nginx

# 12. Остановить и удалить
docker stop test-nginx
docker rm test-nginx

# 13. Удалить образ
docker rmi nginx:alpine
```

**Вопросы для самопроверки:**
1. Чем отличается `docker run` от `docker start`?
2. Что показывает `docker ps` без опций?
3. Как выполнить команду внутри контейнера?
4. Как посмотреть логи контейнера в реальном времени?
5. Можно ли удалить запущенный контейнер?

---

## Часть 2: Создание Dockerfile

### Основные инструкции Dockerfile

#### FROM - Базовый образ

**Что делает:** Указывает базовый образ для сборки.

**Синтаксис:**
```dockerfile
FROM image:tag
```

**Примеры:**
```dockerfile
# Python 3.11 (полная версия ~900MB)
FROM python:3.11

# Python 3.11 slim (~120MB)
FROM python:3.11-slim

# Python 3.11 alpine (~50MB)
FROM python:3.11-alpine

# Ubuntu
FROM ubuntu:22.04

# Node.js
FROM node:18-alpine

# Nginx
FROM nginx:alpine
```

**Когда использовать какой образ:**
- `latest` - ❌ не используйте! (версия может измениться)
- `alpine` - ✅ минимальный размер (но могут быть проблемы совместимости)
- `slim` - ✅ баланс размера и совместимости
- `полная версия` - для разработки, если нужны все инструменты

---

#### WORKDIR - Рабочая директория

**Что делает:** Устанавливает рабочую директорию. Если не существует - создает.

**Синтаксис:**
```dockerfile
WORKDIR /path/to/directory
```

**Пример:**
```dockerfile
FROM python:3.11-slim

# Устанавливаем рабочую директорию
WORKDIR /app

# Все последующие команды выполняются в /app
COPY app.py .           # Копирует в /app/app.py
RUN ls -la              # Показывает содержимое /app
```

**Зачем нужно:**
- Избежать беспорядка в корневой директории
- Упрощает команды (не нужно писать полные пути)
- Можно использовать несколько раз в Dockerfile

---

#### COPY - Копирование файлов

**Что делает:** Копирует файлы с хоста в образ.

**Синтаксис:**
```dockerfile
COPY <source> <destination>
```

**Примеры:**
```dockerfile
# Скопировать один файл
COPY app.py /app/app.py

# Скопировать в текущую WORKDIR
COPY app.py .

# Скопировать директорию
COPY ./src /app/src

# Скопировать все файлы
COPY . /app

# Скопировать несколько файлов
COPY app.py requirements.txt /app/
```

**Важно:**
- Source - относительно директории сборки (где Dockerfile)
- Destination - внутри образа

---

#### RUN - Выполнение команд

**Что делает:** Выполняет команду во время СБОРКИ образа.

**Два формата:**

1. **Shell форма:**
```dockerfile
RUN command param1 param2
```

2. **Exec форма:**
```dockerfile
RUN ["executable", "param1", "param2"]
```

**Примеры:**
```dockerfile
# Установка пакетов
RUN apt-get update && apt-get install -y curl

# Установка Python зависимостей
RUN pip install flask

# Создание директории
RUN mkdir /data

# Множественные команды
RUN apt-get update && \
    apt-get install -y curl wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
```

**Важно:**
- Каждый `RUN` создает новый слой в образе
- Лучше объединять команды через `&&` для уменьшения слоев

---

#### CMD - Команда по умолчанию

**Что делает:** Указывает команду, которая выполнится при ЗАПУСКЕ контейнера.

**Три формата:**

1. **Exec форма (рекомендуемая):**
```dockerfile
CMD ["executable", "param1", "param2"]
```

2. **Shell форма:**
```dockerfile
CMD command param1 param2
```

3. **Параметры для ENTRYPOINT:**
```dockerfile
CMD ["param1", "param2"]
```

**Примеры:**
```dockerfile
# Exec форма
CMD ["python", "app.py"]
CMD ["nginx", "-g", "daemon off;"]
CMD ["node", "server.js"]

# Shell форма
CMD python app.py
CMD echo "Hello Docker"
```

**Важно:**
- Можно переопределить при `docker run`
- Только один `CMD` в Dockerfile (если несколько - работает последний)

---

#### ENTRYPOINT - Точка входа

**Что делает:** Указывает исполняемый файл, который ВСЕГДА выполнится.

**Два формата:**

1. **Exec форма (рекомендуемая):**
```dockerfile
ENTRYPOINT ["executable", "param1"]
```

2. **Shell форма:**
```dockerfile
ENTRYPOINT command param1
```

**Примеры:**
```dockerfile
# Python приложение
ENTRYPOINT ["python"]
CMD ["app.py"]

# Nginx
ENTRYPOINT ["nginx"]
CMD ["-g", "daemon off;"]

# Скрипт инициализации
ENTRYPOINT ["/docker-entrypoint.sh"]
```

---

#### EXPOSE - Документирование портов

**Что делает:** Документирует, какой порт использует приложение.

**Синтаксис:**
```dockerfile
EXPOSE port [port/protocol...]
```

**Примеры:**
```dockerfile
# Один порт
EXPOSE 8000

# Несколько портов
EXPOSE 80 443

# С протоколом
EXPOSE 53/udp
EXPOSE 53/tcp
```

**Важно:**
- Это только ДОКУМЕНТАЦИЯ!
- Порт все равно нужно пробрасывать через `-p` при `docker run`

---

#### ENV - Переменные окружения

**Что делает:** Устанавливает переменные окружения.

**Синтаксис:**
```dockerfile
ENV key=value
ENV key1=value1 key2=value2
```

**Примеры:**
```dockerfile
# Одна переменная
ENV APP_VERSION=1.0.0

# Несколько переменных
ENV DEBUG=false \
    PORT=8000 \
    DATABASE_URL=postgresql://localhost/db
```

**Использование:**
```dockerfile
ENV APP_HOME=/app
WORKDIR $APP_HOME
```

---

### Практическое задание 2.1: Простой Dockerfile

Создайте Python приложение:

**app.py:**
```python
print("Hello from Docker!")
print("My first containerized app")
```

**Dockerfile:**
```dockerfile
# Базовый образ Python 3.11 slim
FROM python:3.11-slim

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем наш файл
COPY app.py .

# Команда по умолчанию
CMD ["python", "app.py"]
```

**Сборка и запуск:**
```bash
# Собрать образ
docker build -t my-first-app .

# Запустить
docker run my-first-app

# Должно вывести:
# Hello from Docker!
# My first containerized app
```

---

### Практическое задание 2.2: Веб-приложение

**app.py:**
```python
from flask import Flask

app = Flask(__name__)

@app.route('/')
def home():
    return 'Hello from Docker!'

@app.route('/info')
def info():
    import sys
    return f'Python version: {sys.version}'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

**requirements.txt:**
```
Flask==3.0.0
```

**Dockerfile:**
```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Копируем requirements первым (для кэширования)
COPY requirements.txt .

# Устанавливаем зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Копируем приложение
COPY app.py .

# Документируем порт
EXPOSE 5000

# Запускаем приложение
CMD ["python", "app.py"]
```

**Сборка и запуск:**
```bash
docker build -t flask-app .
docker run -d -p 5000:5000 --name myapp flask-app

# Проверка
curl http://localhost:5000
curl http://localhost:5000/info

# Просмотр логов
docker logs -f myapp

# Остановка
docker stop myapp
docker rm myapp
```

---

## Часть 3: ENTRYPOINT vs CMD

### Понимание разницы

**CMD** - команда по умолчанию, легко переопределяется  
**ENTRYPOINT** - основная команда, к ней добавляются аргументы

### Вариант 1: Только CMD

**Dockerfile:**
```dockerfile
FROM alpine:3.18
CMD ["echo", "Hello from CMD"]
```

**Использование:**
```bash
# Выполнит команду по умолчанию
docker run test-cmd
# Вывод: Hello from CMD

# CMD ПОЛНОСТЬЮ заменяется
docker run test-cmd echo "Goodbye"
# Вывод: Goodbye
```

**Когда использовать:** Когда нужна гибкость и команду часто меняют.

---

### Вариант 2: Только ENTRYPOINT

**Dockerfile:**
```dockerfile
FROM alpine:3.18
ENTRYPOINT ["echo"]
```

**Использование:**
```bash
# Ошибка! Нет аргументов для echo
docker run test-entrypoint

# Аргументы ДОБАВЛЯЮТСЯ к ENTRYPOINT
docker run test-entrypoint "Hello World"
# Выполнит: echo "Hello World"
# Вывод: Hello World

# Можно передать любые аргументы
docker run test-entrypoint "Line 1" "Line 2"
# Выполнит: echo "Line 1" "Line 2"
```

**Когда использовать:** Когда образ делает одну задачу с разными параметрами.

---

### Вариант 3: ENTRYPOINT + CMD (комбинация)

**Dockerfile:**
```dockerfile
FROM alpine:3.18
ENTRYPOINT ["echo"]
CMD ["Hello from Docker!"]
```

**Использование:**
```bash
# Выполнит: echo "Hello from Docker!"
docker run test-both
# Вывод: Hello from Docker!

# CMD заменяется, ENTRYPOINT остается
docker run test-both "Custom message"
# Выполнит: echo "Custom message"
# Вывод: Custom message
```

**Когда использовать:** Когда есть команда по умолчанию, но можно передать другие аргументы.

---

### Практический пример: Утилита grep

**Dockerfile:**
```dockerfile
FROM alpine:3.18

# Установим текстовый файл для примера
RUN echo "apple" > /data.txt && \
    echo "banana" >> /data.txt && \
    echo "orange" >> /data.txt

ENTRYPOINT ["grep"]
CMD ["apple", "/data.txt"]
```

**Использование:**
```bash
# По умолчанию ищет "apple"
docker run grep-util
# Вывод: apple

# Поиск другого слова
docker run grep-util "banana" /data.txt
# Вывод: banana

# С опциями grep
docker run grep-util -i "ORANGE" /data.txt
# Вывод: orange
```

---

### Переопределение ENTRYPOINT

Если всё-таки нужно переопределить ENTRYPOINT:

```bash
docker run --entrypoint /bin/sh test-both
```

---

### Строковый vs JSON формат

**Строковый формат** (Shell форма):
```dockerfile
CMD python app.py
ENTRYPOINT python app.py
```

**Что происходит:** Команда выполняется через `/bin/sh -c`

**JSON формат** (Exec форма) - РЕКОМЕНДУЕТСЯ:
```dockerfile
CMD ["python", "app.py"]
ENTRYPOINT ["python", "app.py"]
```

**Что происходит:** Команда выполняется напрямую, без shell

**Различия:**

| Аспект | Shell форма | Exec форма |
|--------|-------------|-----------|
| Синтаксис | `CMD python app.py` | `CMD ["python", "app.py"]` |
| Shell переменные | ✅ Работают (`$HOME`) | ❌ Не работают |
| Сигналы (SIGTERM) | ❌ Не доходят до приложения | ✅ Доходят напрямую |
| PID 1 | `/bin/sh` | Ваше приложение |
| Рекомендуется | Редко | ✅ Всегда |

**Пример с переменными:**
```dockerfile
# Shell форма - переменные работают
CMD echo "Home directory: $HOME"

# Exec форма - переменные НЕ работают
CMD ["echo", "Home directory: $HOME"]
# Выведет буквально: Home directory: $HOME

# Решение для Exec формы:
CMD ["/bin/sh", "-c", "echo Home directory: $HOME"]
```

---

### Практическое задание 3.1

Создайте три Dockerfile и протестируйте разницу:

**Dockerfile.cmd:**
```dockerfile
FROM alpine:3.18
CMD ["echo", "Only CMD"]
```

**Dockerfile.entrypoint:**
```dockerfile
FROM alpine:3.18
ENTRYPOINT ["echo"]
CMD ["Default message"]
```

**Dockerfile.both:**
```dockerfile
FROM alpine:3.18
ENTRYPOINT ["echo", "Prefix:"]
CMD ["Default message"]
```

**Тестирование:**
```bash
# Собрать
docker build -t test-cmd -f Dockerfile.cmd .
docker build -t test-entrypoint -f Dockerfile.entrypoint .
docker build -t test-both -f Dockerfile.both .

# Тест 1: Без аргументов
docker run test-cmd
docker run test-entrypoint
docker run test-both

# Тест 2: С аргументами
docker run test-cmd "New message"
docker run test-entrypoint "New message"
docker run test-both "New message"
```

**Вопросы:**
1. Что вывела каждая команда?
2. Как изменился вывод с аргументами?
3. Какой вариант для какой задачи подходит?

---

## Часть 4: Работа с несколькими базовыми образами

### Когда использовать разные базовые образы

| Язык | Базовый образ | Размер | Когда использовать |
|------|---------------|--------|-------------------|
| Python | `python:3.11` | ~900MB | Разработка, нужны все инструменты |
| Python | `python:3.11-slim` | ~120MB | Production, веб-приложения |
| Python | `python:3.11-alpine` | ~50MB | Минимальный размер |
| Node.js | `node:18` | ~900MB | Разработка |
| Node.js | `node:18-alpine` | ~170MB | Production |
| Go | `golang:1.21` | ~800MB | Только для сборки |
| Go | `alpine:3.18` | ~7MB | Runtime для Go бинарников |
| Java | `openjdk:17` | ~470MB | Java приложения |
| Nginx | `nginx:alpine` | ~40MB | Веб-сервер |
| Ubuntu | `ubuntu:22.04` | ~80MB | Универсальный |
| Alpine | `alpine:3.18` | ~7MB | Минимальная база |

---

### Задание 4.1: Python приложение с зависимостями

**app.py:**
```python
from flask import Flask, jsonify
import requests

app = Flask(__name__)

@app.route('/')
def home():
    return jsonify({
        'message': 'Hello Docker!',
        'status': 'running'
    })

@app.route('/external')
def external():
    # Проверяем, что requests работает
    response = requests.get('https://api.github.com')
    return jsonify({
        'github_status': response.status_code
    })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

**requirements.txt:**
```
Flask==3.0.0
requests==2.31.0
```

**Dockerfile:**
```dockerfile
# Используем slim образ
FROM python:3.11-slim

# Рабочая директория
WORKDIR /app

# Копируем requirements.txt ПЕРВЫМ
# Это важно для кэширования слоев!
COPY requirements.txt .

# Устанавливаем зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Теперь копируем код
COPY app.py .

# Порт приложения
EXPOSE 5000

# Команда запуска
CMD ["python", "app.py"]
```

**Почему этот порядок важен?**

1. requirements.txt меняется редко → слой кэшируется
2. app.py меняется часто → только этот слой пересобирается
3. Сборка образа происходит БЫСТРЕЕ!

**Запуск:**
```bash
docker build -t flask-with-deps .
docker run -d -p 5000:5000 flask-with-deps

curl http://localhost:5000/
curl http://localhost:5000/external
```

---

### Задание 4.2: Node.js приложение

**server.js:**
```javascript
const express = require('express');
const app = express();
const port = 3000;

app.get('/', (req, res) => {
  res.json({
    message: 'Node.js in Docker!',
    version: process.version,
    timestamp: new Date().toISOString()
  });
});

app.get('/info', (req, res) => {
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

**package.json:**
```json
{
  "name": "docker-node-app",
  "version": "1.0.0",
  "main": "server.js",
  "scripts": {
    "start": "node server.js"
  },
  "dependencies": {
    "express": "^4.18.2"
  }
}
```

**Dockerfile:**
```dockerfile
FROM node:18-alpine

WORKDIR /app

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm ci --only=production

# Копируем код
COPY server.js .

EXPOSE 3000

CMD ["node", "server.js"]
```

**Почему `npm ci` а не `npm install`?**
- `npm ci` - быстрее, для production
- `npm ci` - использует package-lock.json точно
- `npm install` - для разработки

**Запуск:**
```bash
docker build -t node-app .
docker run -d -p 3000:3000 node-app

curl http://localhost:3000/
curl http://localhost:3000/info
```

---

### Задание 4.3: Сравнение размеров образов

Создайте три варианта Python образа:

**Dockerfile.full:**
```dockerfile
FROM python:3.11
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY app.py .
CMD ["python", "app.py"]
```

**Dockerfile.slim:**
```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY app.py .
CMD ["python", "app.py"]
```

**Dockerfile.alpine:**
```dockerfile
FROM python:3.11-alpine
WORKDIR /app
COPY requirements.txt .
# Alpine может требовать дополнительные пакеты
RUN apk add --no-cache gcc musl-dev linux-headers && \
    pip install --no-cache-dir -r requirements.txt && \
    apk del gcc musl-dev linux-headers
COPY app.py .
CMD ["python", "app.py"]
```

**Сборка и сравнение:**
```bash
docker build -t app:full -f Dockerfile.full .
docker build -t app:slim -f Dockerfile.slim .
docker build -t app:alpine -f Dockerfile.alpine .

# Сравнить размеры
docker images | grep app:

# Результат примерно:
# app:full    ~950MB
# app:slim    ~180MB
# app:alpine  ~80MB
```

**Выводы:**
- `full` - для разработки, все инструменты включены
- `slim` - ✅ лучший выбор для production
- `alpine` - минимальный размер, но могут быть проблемы с зависимостями

---

### Задание 4.4: Статический сайт с Nginx

Создайте простой статический веб-сайт и упакуйте его в Docker образ с Nginx.

**Структура проекта:**
```
my-website/
├── Dockerfile
├── index.html
├── styles.css
└── app.js
```

**index.html:**
```html
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Docker Website</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
    <div class="container">
        <h1>Привет из Docker! 🐳</h1>
        <p>Это статический сайт, запущенный в Nginx контейнере</p>
        <button id="clickBtn">Нажми меня!</button>
        <p id="counter">Нажатий: 0</p>
    </div>
    <script src="app.js"></script>
</body>
</html>
```

**styles.css:**
```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: Arial, sans-serif;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    min-height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
}

.container {
    background: white;
    padding: 3rem;
    border-radius: 15px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    text-align: center;
    max-width: 500px;
}

h1 {
    color: #333;
    margin-bottom: 1rem;
    font-size: 2.5rem;
}

p {
    color: #666;
    margin-bottom: 1.5rem;
    font-size: 1.2rem;
}

button {
    background: #667eea;
    color: white;
    border: none;
    padding: 15px 40px;
    font-size: 1.1rem;
    border-radius: 50px;
    cursor: pointer;
    transition: all 0.3s ease;
}

button:hover {
    background: #764ba2;
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

#counter {
    margin-top: 1.5rem;
    font-size: 1.5rem;
    font-weight: bold;
    color: #667eea;
}
```

**app.js:**
```javascript
let count = 0;
const btn = document.getElementById('clickBtn');
const counter = document.getElementById('counter');

btn.addEventListener('click', () => {
    count++;
    counter.textContent = `Нажатий: ${count}`;
    
    // Небольшая анимация
    counter.style.transform = 'scale(1.2)';
    setTimeout(() => {
        counter.style.transform = 'scale(1)';
    }, 200);
});

// Добавим transition для плавности
counter.style.transition = 'transform 0.2s ease';
```

**Dockerfile:**
```dockerfile
# Используем официальный образ Nginx на Alpine
FROM nginx:alpine

# Удаляем дефолтную страницу Nginx
RUN rm -rf /usr/share/nginx/html/*

# Копируем наши файлы в директорию Nginx
COPY index.html /usr/share/nginx/html/
COPY styles.css /usr/share/nginx/html/
COPY app.js /usr/share/nginx/html/

# Nginx слушает порт 80
EXPOSE 80

# CMD уже определена в базовом образе nginx
# Запускает: nginx -g 'daemon off;'
```

**Сборка и запуск:**

```bash
# Создать директорию проекта
mkdir my-website
cd my-website

# Создать файлы (скопируйте содержимое выше)

# Собрать образ
docker build -t my-website:1.0 .

# Запустить контейнер
docker run -d -p 8080:80 --name website my-website:1.0

# Открыть в браузере
# http://localhost:8080
```

**Проверка:**
```bash
# Посмотреть логи
docker logs website

# Посмотреть файлы внутри
docker exec website ls -la /usr/share/nginx/html/

# Остановить
docker stop website
docker rm website
```

**Улучшенная версия с конфигурацией Nginx:**

Создайте файл `nginx.conf`:
```nginx
server {
    listen 80;
    server_name localhost;
    
    root /usr/share/nginx/html;
    index index.html;
    
    # Включить gzip сжатие
    gzip on;
    gzip_types text/css application/javascript text/html;
    
    # Кэширование статики
    location ~* \.(css|js|jpg|jpeg|png|gif|ico)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # Основной location
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

**Dockerfile с кастомной конфигурацией:**
```dockerfile
FROM nginx:alpine

# Удаляем дефолтную конфигурацию
RUN rm /etc/nginx/conf.d/default.conf

# Копируем нашу конфигурацию
COPY nginx.conf /etc/nginx/conf.d/

# Копируем сайт
COPY index.html /usr/share/nginx/html/
COPY styles.css /usr/share/nginx/html/
COPY app.js /usr/share/nginx/html/

EXPOSE 80
```

**Сборка:**
```bash
docker build -t my-website:2.0 .
docker run -d -p 8080:80 my-website:2.0

# Проверить размер образа
docker images my-website
# Размер: ~40-45MB (благодаря nginx:alpine)!
```

**Разработка с bind mount:**
```bash
# Запустить с монтированием локальных файлов
docker run -d \
  -p 8080:80 \
  -v $(pwd):/usr/share/nginx/html \
  --name dev-website \
  nginx:alpine

# Теперь можно редактировать HTML/CSS/JS
# и сразу видеть изменения в браузере!
```

**Преимущества этого подхода:**
- ✅ Очень маленький размер образа (~40MB)
- ✅ Nginx быстро отдает статику
- ✅ Production-ready решение
- ✅ Легко масштабировать
- ✅ Удобная разработка с bind mount

---

## Часть 5: Дебаг контейнеров

### Основные техники отладки

#### Техника 1: Просмотр логов

```bash
# Все логи
docker logs container_name

# Последние 100 строк
docker logs --tail 100 container_name

# Следить в реальном времени
docker logs -f container_name

# С временными метками
docker logs -t container_name

# За последний час
docker logs --since 1h container_name
```

**Когда использовать:** Первое, что нужно проверить при проблемах!

---

#### Техника 2: Вход в контейнер

```bash
# Bash (если есть)
docker exec -it container_name /bin/bash

# Sh (alpine образы)
docker exec -it container_name /bin/sh

# Внутри контейнера можно:
ls -la              # Посмотреть файлы
ps aux              # Посмотреть процессы
env                 # Посмотреть переменные
cat /etc/hosts      # Посмотреть hosts
netstat -tulpn      # Посмотреть порты
```

---

#### Техника 3: Выполнение команд

```bash
# Проверить, что файл существует
docker exec container_name ls -la /app/config.json

# Прочитать файл
docker exec container_name cat /app/config.json

# Проверить процессы
docker exec container_name ps aux

# Проверить порты
docker exec container_name netstat -tulpn

# Проверить переменные окружения
docker exec container_name env
```

---

#### Техника 4: Docker inspect

```bash
# Вся информация о контейнере
docker inspect container_name

# Статус контейнера
docker inspect --format='{{.State.Status}}' container_name

# IP адрес
docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' container_name

# Переменные окружения
docker inspect --format='{{range .Config.Env}}{{println .}}{{end}}' container_name

# Монтированные volumes
docker inspect --format='{{range .Mounts}}{{.Source}} -> {{.Destination}}{{println}}{{end}}' container_name
```

---

#### Техника 5: Docker stats

```bash
# Использование ресурсов
docker stats container_name

# Все контейнеры
docker stats

# Одноразовый вывод
docker stats --no-stream
```

**Что показывает:**
- CPU % - использование процессора
- MEM USAGE / LIMIT - память
- MEM % - процент от лимита
- NET I/O - сетевой трафик
- BLOCK I/O - дисковый ввод/вывод

---

### Отладка упавшего контейнера

**Проблема:** Контейнер сразу падает при запуске.

**Шаг 1: Посмотреть логи**
```bash
# Найти ID последнего контейнера
docker ps -a

# Посмотреть логи
docker logs <container_id>
```

**Шаг 2: Запустить с другой командой**
```bash
# Переопределить CMD
docker run -it my-image /bin/sh

# Внутри проверить:
ls -la /app
cat /app/config.json
python app.py  # Запустить вручную
```

**Шаг 3: Переопределить ENTRYPOINT**
```bash
docker run -it --entrypoint /bin/sh my-image
```

**Шаг 4: Проверить, что файлы скопировались**
```bash
docker run --rm my-image ls -la /app
```

---

### Практическое задание 5.1: Отладка проблемного приложения

**Создайте проблемное приложение:**

**app.py:**
```python
import os

# Ошибка: файл не существует
config_file = os.environ.get('CONFIG_FILE', '/app/config.json')
with open(config_file, 'r') as f:
    config = f.read()
    print(f"Config: {config}")
```

**Dockerfile:**
```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY app.py .
CMD ["python", "app.py"]
```

**Задание:**
1. Соберите образ
2. Запустите контейнер
3. Контейнер упадет - найдите причину
4. Исправьте проблему

**Отладка:**
```bash
# Собрать
docker build -t debug-app .

# Запустить - упадет
docker run debug-app

# Посмотреть логи
docker ps -a
docker logs <container_id>
# Увидите: FileNotFoundError: /app/config.json

# Решение 1: Создать файл
echo '{"debug": true}' > config.json

# Dockerfile (исправленный):
FROM python:3.11-slim
WORKDIR /app
COPY config.json .
COPY app.py .
CMD ["python", "app.py"]

# Решение 2: Использовать переменную
docker run -e CONFIG_FILE=/dev/null debug-app
```

## Часть 6: Docker Hub

### Что такое Docker Hub?

Docker Hub - это реестр (registry) Docker образов:
- Публичные образы (nginx, python, mysql)
- Ваши собственные образы
- Бесплатные и платные аккаунты

**Официальный сайт:** https://hub.docker.com/

---

### Задание 6.1: Регистрация и login

**Шаг 1: Создание аккаунта**

1. Перейдите на https://hub.docker.com/signup
2. Заполните форму регистрации
3. Подтвердите email
4. Запомните ваш username

**Шаг 2: Вход через терминал**

```bash
# Войти в Docker Hub
docker login

# Введите username и password
# Login Succeeded

# Проверить, что вошли
docker info | grep Username
```

**Выход:**
```bash
docker logout
```

---

### Задание 6.2: Теги и версии

**Что такое теги?**

Тег - это метка версии образа: `image:tag`

**Примеры:**
- `nginx:latest` - последняя версия
- `nginx:1.25` - конкретная версия
- `python:3.11-slim` - версия + вариант
- `myapp:v1.0.0` - ваша версия

**Правила именования образов для Docker Hub:**

```
username/repository:tag
```

**Примеры:**
- `johndoe/myapp:latest`
- `johndoe/myapp:1.0.0`
- `johndoe/myapp:dev`

**Создание тегов:**

```bash
# Собрать с именем для Docker Hub
docker build -t username/myapp:1.0.0 .

# Или создать тег для существующего образа
docker tag myapp:latest username/myapp:1.0.0

# Создать несколько тегов
docker tag username/myapp:1.0.0 username/myapp:latest
docker tag username/myapp:1.0.0 username/myapp:stable

# Посмотреть все образы
docker images username/myapp
```

---

### Задание 6.3: Push в Docker Hub

**Загрузка образа:**

```bash
# 1. Убедитесь, что вы вошли
docker login

# 2. Соберите образ с правильным именем
docker build -t username/myapp:1.0.0 .

# 3. Загрузите на Docker Hub
docker push username/myapp:1.0.0

# Вывод:
# The push refers to repository [docker.io/username/myapp]
# 1.0.0: digest: sha256:abc123... size: 1234
```

**Загрузка нескольких тегов:**

```bash
docker push username/myapp:1.0.0
docker push username/myapp:latest
```

**Проверка:**

1. Откройте https://hub.docker.com/
2. Войдите в аккаунт
3. Перейдите в Repositories
4. Увидите ваш образ!

---

### Задание 6.4: Pull образа

**Скачивание образа:**

```bash
# Скачать latest версию
docker pull username/myapp

# Скачать конкретную версию
docker pull username/myapp:1.0.0

# Скачать с другого registry
docker pull gcr.io/google-containers/busybox
```

**Проверка:**

```bash
# Посмотреть скачанные образы
docker images

# Запустить
docker run username/myapp:1.0.0
```

---

### Задание 6.5: Практика - полный цикл

**Создайте и опубликуйте образ:**

**app.py:**
```python
print("My Docker Hub app!")
print("Version 1.0.0")
```

**Dockerfile:**
```dockerfile
FROM python:3.11-alpine
WORKDIR /app
COPY app.py .
CMD ["python", "app.py"]
```

**Шаги:**

```bash
# 1. Войти
docker login

# 2. Собрать (замените username на ваш!)
docker build -t username/hello-docker:1.0.0 .

# 3. Создать теги
docker tag username/hello-docker:1.0.0 username/hello-docker:latest

# 4. Загрузить
docker push username/hello-docker:1.0.0
docker push username/hello-docker:latest

# 5. Удалить локальные образы
docker rmi username/hello-docker:1.0.0
docker rmi username/hello-docker:latest

# 6. Скачать заново
docker pull username/hello-docker:latest

# 7. Запустить
docker run username/hello-docker:latest
```

**Поздравляю!** Теперь любой человек может скачать ваш образ:
```bash
docker run username/hello-docker:latest
```

---

### Приватные репозитории

**Docker Hub предоставляет:**
- Бесплатно: 1 приватный репозиторий
- Платно: неограниченное количество

**Создать приватный репозиторий:**
1. Docker Hub → Create Repository
2. Выбрать "Private"
3. Загружать так же через `docker push`

---

## Часть 7: Оптимизация образов

### Почему размер важен?

- ⚡ Быстрее скачивание
- ⚡ Быстрее сборка
- 💾 Меньше места на диске
- 🔒 Меньше поверхность атаки
- 💰 Дешевле хранение и трафик

---

### Техника 1: Выбор правильного базового образа

**Плохо:**
```dockerfile
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y python3
# Размер: ~200MB
```

**Хорошо:**
```dockerfile
FROM python:3.11-slim
# Размер: ~120MB
```

**Отлично:**
```dockerfile
FROM python:3.11-alpine
# Размер: ~50MB
```

---

### Техника 2: Объединение команд RUN

**Плохо (создает много слоев):**
```dockerfile
FROM ubuntu:22.04
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y wget
RUN apt-get install -y vim
RUN apt-get clean
# 5 слоев, каждый добавляет размер
```

**Хорошо (один слой):**
```dockerfile
FROM ubuntu:22.04
RUN apt-get update && \
    apt-get install -y \
        curl \
        wget \
        vim && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
# 1 слой, очистка в том же слое
```

---

### Техника 3: Использование --no-cache-dir

**Плохо:**
```dockerfile
RUN pip install flask
# pip сохраняет кэш ~10-20MB
```

**Хорошо:**
```dockerfile
RUN pip install --no-cache-dir flask
# Кэш не сохраняется
```

Для apt:
```dockerfile
RUN apt-get update && \
    apt-get install -y package && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
```

---

### Техника 4: .dockerignore

Создайте `.dockerignore` чтобы не копировать ненужное:

**.dockerignore:**
```
# Git
.git
.gitignore

# Python
__pycache__
*.pyc
*.pyo
*.pyd
.Python
env/
venv/
*.egg-info/

# Node.js
node_modules/
npm-debug.log

# Тесты и документация
tests/
docs/
*.md
README.md

# IDE
.vscode/
.idea/
*.swp

# Logs
*.log
logs/

# OS
.DS_Store
Thumbs.db

# Временные файлы
*.tmp
temp/
```

**Зачем:**
- Уменьшает размер контекста сборки
- Ускоряет `COPY . .`
- Не копирует секреты (.env файлы)

---

### Техника 5: Порядок инструкций (кэширование слоев)

**Плохо:**
```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY . .                    # Копируем всё сразу
RUN pip install -r requirements.txt
CMD ["python", "app.py"]
```

**Проблема:** При изменении app.py, pip install выполнится заново!

**Хорошо:**
```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .     # Сначала requirements
RUN pip install --no-cache-dir -r requirements.txt
COPY . .                    # Потом остальное
CMD ["python", "app.py"]
```

**Преимущество:** Изменение app.py не вызывает переустановку зависимостей!

---

### Практическое задание 7.1: Оптимизация образа

**Создайте неоптимизированный Dockerfile:**

**Dockerfile.bad:**
```dockerfile
FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install -y python3
RUN apt-get install -y python3-pip
RUN apt-get install -y curl
RUN apt-get install -y vim

WORKDIR /app

COPY . .

RUN pip3 install flask

CMD ["python3", "app.py"]
```

**Оптимизированный Dockerfile:**

**Dockerfile.good:**
```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Только необходимые пакеты
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Зависимости первыми
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Код последним
COPY app.py .

CMD ["python", "app.py"]
```

**.dockerignore:**
```
.git
__pycache__
*.pyc
*.log
.DS_Store
```

**Сравнение:**
```bash
docker build -t app:bad -f Dockerfile.bad .
docker build -t app:good -f Dockerfile.good .

docker images | grep app:

# Результат:
# app:bad   ~400MB
# app:good  ~130MB
```

**Экономия: ~270MB (67%)!**

---

## Часть 8: Bind mounts для разработки

### Что такое bind mount?

**Bind mount** - монтирование директории хоста в контейнер.

**Зачем:**
- Разработка без пересборки образа
- Изменения кода сразу видны в контейнере
- Сохранение данных на хосте

---

### Синтаксис bind mount

```bash
docker run -v /host/path:/container/path image

# Или современный синтаксис
docker run --mount type=bind,source=/host/path,target=/container/path image
```

---

### Задание 8.1: Разработка Python приложения

**app.py:**
```python
print("Version 1")
```

**Dockerfile:**
```dockerfile
FROM python:3.11-slim
WORKDIR /app
CMD ["python", "app.py"]
```

**Собрать образ:**
```bash
docker build -t dev-app .
```

**Запуск с bind mount:**

```bash
# Linux/Mac
docker run -v $(pwd):/app dev-app

# Windows PowerShell
docker run -v ${PWD}:/app dev-app

# Вывод: Version 1
```

**Изменить код БЕЗ пересборки:**

```bash
# Изменить файл
echo 'print("Version 2")' > app.py

# Запустить снова
docker run -v $(pwd):/app dev-app

# Вывод: Version 2
```

**Магия!** Код обновился без пересборки образа!

---

### Задание 8.2: Flask с live reload

**app.py:**
```python
from flask import Flask

app = Flask(__name__)

@app.route('/')
def home():
    return 'Hello Docker - Version 1'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
```

**requirements.txt:**
```
Flask==3.0.0
```

**Dockerfile:**
```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
CMD ["python", "app.py"]
```

**Собрать:**
```bash
docker build -t flask-dev .
```

**Запустить с bind mount:**
```bash
docker run -d \
  -p 5000:5000 \
  -v $(pwd):/app \
  --name flask-dev \
  flask-dev

# Проверить
curl http://localhost:5000
# Вывод: Hello Docker - Version 1
```

**Изменить код:**
```bash
# Поменять app.py
sed -i 's/Version 1/Version 2/' app.py

# Или отредактировать в редакторе
```

**Проверить:**
```bash
curl http://localhost:5000
# Вывод: Hello Docker - Version 2
```

**Flask автоматически перезагрузился!** (режим debug=True)

---

### Задание 8.3: Монтирование для конфигурации

```bash
# Создать конфиг на хосте
echo "DEBUG=true" > config.env

# Запустить с монтированием конфига
docker run -v $(pwd)/config.env:/app/config.env myapp
```

---

### Задание 8.4: Только для чтения (read-only)

```bash
# Монтировать только для чтения
docker run -v $(pwd):/app:ro myapp

# Контейнер не сможет изменить файлы на хосте
```

---

### Различия: Bind mount vs Volume

| Аспект | Bind Mount | Volume |
|--------|-----------|--------|
| Где хранится | Любая директория хоста | Управляется Docker |
| Синтаксис | `-v /host/path:/container/path` | `-v volume_name:/container/path` |
| Когда использовать | Разработка, конфиги | Production, данные БД |
| Backup | Вручную | Через Docker |
| Производительность | Зависит от хоста | Оптимизировано Docker |

**Для разработки:** Используйте bind mounts  
**Для production:** Используйте volumes

---

## Часть 9: Сборка from scratch (опционально)

### Что такое scratch?

`scratch` - пустой образ (0 байт!)

**Когда использовать:**
- Статические бинарники (Go, Rust)
- Минимальный размер образа
- Максимальная безопасность (нет shell, утилит)

---

### Задание 9.1: Go приложение

**main.go:**
```go
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from Go in Docker!")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

**Dockerfile:**
```dockerfile
# Этап 1: Сборка
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY main.go .

# Компиляция статического бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app main.go

# Этап 2: Финальный образ
FROM scratch

# Копируем только бинарник
COPY --from=builder /app/app /app

# Запуск
CMD ["/app"]
```

**Сборка:**
```bash
docker build -t go-scratch .

# Размер образа
docker images go-scratch
# Размер: ~6-8MB!
```

**Запуск:**
```bash
docker run -p 8080:8080 go-scratch

curl http://localhost:8080
```

---

### Distroless образы (альтернатива scratch)

**Что такое distroless?**
- Минимальные образы от Google
- Только runtime, без shell/package manager
- Больше чем scratch, но проще в использовании

**Python с distroless:**

```dockerfile
FROM python:3.11-slim AS builder

WORKDIR /app
COPY requirements.txt .
RUN pip install --user --no-cache-dir -r requirements.txt
COPY app.py .

FROM gcr.io/distroless/python3-debian11

COPY --from=builder /root/.local /root/.local
COPY --from=builder /app /app

ENV PATH=/root/.local/bin:$PATH
WORKDIR /app

CMD ["python", "app.py"]
```

---

## Задание для самостоятельной работы

### Контейнеризация вашего проекта

**Задачи:** 
1. Возьмите ваш проект (из другой лабораторной или личный) и создайте для него Docker образ, описав Dockerfile
2. Выложите образ в Docker Hub
3. Создайте репозиторий в gitlab.mai.ru в вашей персональной группе и выложите приложение и Dockerfile c документацией по сборке и запуску

**Требования:**

#### 1. Dockerfile (20 баллов)

- ✅ Использовать подходящий базовый образ (slim/alpine)
- ✅ Правильно установить зависимости
- ✅ Оптимизировать размер образа (< 500MB)
- ✅ Использовать правильные CMD/ENTRYPOINT в JSON формате
- ✅ Правильный порядок инструкций (кэширование)

#### 2. .dockerignore (5 баллов)

- ✅ Исключить ненужные файлы
- ✅ Исключить .git, node_modules, __pycache__
- ✅ Исключить файлы разработки

#### 3. Docker Hub (20 баллов)

- ✅ Образ загружен на Docker Hub
- ✅ Минимум 2 тега (latest + версия)
- ✅ Публичный доступ
- ✅ README на Docker Hub с инструкциями по запуску

#### 4. Работающее приложение (25 баллов)

- ✅ Приложение запускается в контейнере
- ✅ Все зависимости установлены
- ✅ Порты правильно пробрасываются
- ✅ Функциональность работает корректно

#### 6. Bind mount для разработки (10 баллов)

- ✅ Можно запустить с bind mount
- ✅ Изменения кода применяются без пересборки
- ✅ Команды запуска с bind mount в README на Docker Hub

#### 7. Дополнительно (опционально, +10 баллов)

- Distroless или scratch образ
- Автоматическая сборка (GitHub Actions)
- Health check в Dockerfile
- Multi-architecture образ (amd64, arm64)

---

### Примеры проектов для контейнеризации

1. **Веб-приложение**
   - Flask/Django REST API
   - Express.js сервер
   - FastAPI приложение

2. **Телеграм бот**
   - Python telegram bot
   - Бот с базой данных

3. **CLI утилита**
   - Скрипт обработки данных
   - Автоматизация задач

4. **Микросервис**
   - Сервис обработки изображений
   - API для работы с файлами

---

### Критерии оценки

| Критерий | Баллы |
|----------|-------|
| Dockerfile корректный и оптимизированный | 20 |
| .dockerignore создан правильно | 5 |
| Образ на Docker Hub с тегами и README | 20 |
| Приложение работает в контейнере | 25 |
| Bind mount для разработки | 10 |
| Размер образа < 500MB | 10 |
| **ИТОГО** | **90** |
| Дополнительно: distroless/scratch/CI | +10 |
| **МАКСИМУМ** | **100** |

---

## Контрольные вопросы

### Базовые вопросы (обязательные)

1. **В чем разница между CMD и ENTRYPOINT?**
   - Когда использовать каждый из них?

2. **Когда использовать строковый формат CMD, а когда JSON?**
   - В чем преимущества JSON формата?

3. **Как уменьшить размер Docker образа?**
   - Назовите минимум 5 способов.

4. **Что такое слои (layers) в Docker образе?**
   - Как порядок инструкций влияет на размер?

5. **Зачем нужен файл .dockerignore?**
   - Что туда обычно добавляют?

6. **Как посмотреть логи контейнера?**
   - Как следить за логами в реальном времени?

7. **Как выполнить команду внутри запущенного контейнера?**
   - Какая разница между docker exec и docker run?

8. **Что делает команда `docker inspect`?**
   - Как получить конкретное поле из вывода?

9. **Как создать несколько тегов для одного образа?**
   - Зачем нужны разные теги?

10. **В чем разница между `docker run` и `docker exec`?**
    - Когда использовать каждый из них?

### Продвинутые вопросы (для самопроверки)

11. **Как монтировать локальную директорию в контейнер?**
    - В чем разница между bind mount и volume?

12. **Что такое bind mount и когда его использовать?**
    - Подходит ли он для production?

13. **Как отладить контейнер, который сразу падает?**
    - Какие команды использовать?

14. **Зачем нужен Docker Hub?**
    - Можно ли использовать другие registry?

15. **Что такое distroless образы и в чем их преимущество?**
    - Когда их использовать?

---

## Дополнительные материалы

### Шпаргалка по командам

```bash
# === ОБРАЗЫ ===
docker images                    # Список образов
docker pull image:tag            # Скачать образ
docker build -t name:tag .       # Собрать образ
docker rmi image:tag             # Удалить образ
docker image prune               # Удалить неиспользуемые

# === КОНТЕЙНЕРЫ ===
docker ps                        # Запущенные контейнеры
docker ps -a                     # Все контейнеры
docker run -d -p 8080:80 image   # Запустить контейнер
docker start container           # Запустить остановленный
docker stop container            # Остановить
docker restart container         # Перезапустить
docker rm container              # Удалить
docker container prune           # Удалить остановленные

# === ОТЛАДКА ===
docker logs -f container         # Логи
docker exec -it container bash   # Войти в контейнер
docker inspect container         # Подробная информация
docker stats                     # Использование ресурсов

# === DOCKER HUB ===
docker login                     # Войти
docker tag image user/repo:tag   # Создать тег
docker push user/repo:tag        # Загрузить
docker pull user/repo:tag        # Скачать

# === ОЧИСТКА ===
docker system prune              # Очистить все
docker system prune -a           # Очистить все + образы
docker volume prune              # Очистить volumes
```

### Полезные ссылки

#### Официальная документация
- [Docker Documentation](https://docs.docker.com/)
- [Dockerfile Reference](https://docs.docker.com/engine/reference/builder/)
- [Docker CLI Reference](https://docs.docker.com/engine/reference/commandline/cli/)
- [Best Practices](https://docs.docker.com/develop/dev-best-practices/)

#### Docker Hub
- [Docker Hub](https://hub.docker.com/)
- [Official Images](https://hub.docker.com/search?q=&type=image&image_filter=official)

#### Обучающие материалы
- [Docker для начинающих](https://github.com/allenov/Lessons/tree/main/2-lesson-docker)
- [Play with Docker](https://labs.play-with-docker.com/) - Онлайн песочница
- [Docker Curriculum](https://docker-curriculum.com/)

#### Инструменты
- [Dive](https://github.com/wagoodman/dive) - Анализ слоев образа
- [Hadolint](https://github.com/hadolint/hadolint) - Линтер для Dockerfile
- [Docker Slim](https://github.com/docker-slim/docker-slim) - Минимизация образов

#### Видео
- [Docker Tutorial for Beginners](https://www.youtube.com/watch?v=fqMOX6JJhGo)
- [Docker Mastery](https://www.udemy.com/course/docker-mastery/)

### Наши материалы

- [Quick Start Guide](./docs/quick-start.md) - Быстрый старт (5-30 минут)
- [Docker Cheat Sheet](./docs/docker-cheatsheet.md) - Шпаргалка по командам
- [Dockerfile Best Practices](./docs/dockerfile-best-practices.md) - Лучшие практики
- [FAQ](./docs/faq.md) - Часто задаваемые вопросы
- [Learning Resources](./docs/learning-resources.md) - 100+ обучающих ресурсов

### Дополнительные примеры

В папке `examples/` подробные примеры:
- `basic-commands/` - Все команды Docker с примерами
- `dockerfile/` - Примеры Dockerfile для разных языков
- `volumes/` - Работа с данными и volumes
- `security/` - Безопасность (опционально)
- `troubleshooting/` - Решение типичных проблем

---

## Следующие шаги

После успешного выполнения этой лабораторной работы вы будете готовы к:

**Лабораторная работа №3: Docker Compose и продвинутые возможности**
- Docker Compose для многоконтейнерных приложений
- Multi-stage builds для оптимизации
- Docker Buildx
- Docker Bake

**Дальнейшее изучение:**
- Kubernetes
- CI/CD с Docker
- Docker Swarm
- Production deployment

---

## Часто задаваемые вопросы

**В: Можно ли использовать Windows для выполнения работы?**  
О: Да, установите Docker Desktop с WSL 2 (инструкции в начале работы).

**В: Какой проект выбрать для контейнеризации?**  
О: Любой ваш проект: веб-приложение, API, телеграм бот, CLI утилиту.

**В: Образ получился > 500MB, что делать?**  
О: Используйте alpine/slim образы, объединяйте RUN команды, используйте .dockerignore.

**В: Как изменить код без пересборки образа?**  
О: Используйте bind mount: `docker run -v $(pwd):/app myapp`

**В: Контейнер сразу падает, что делать?**  
О: Смотрите логи: `docker logs <container_id>`, запустите с shell: `docker run -it image /bin/sh`

**В: Docker Hub требует платный аккаунт?**  
О: Нет, бесплатного аккаунта достаточно для публичных репозиториев.

---

**Успехов в контейнеризации ваших приложений!** 🐳

---