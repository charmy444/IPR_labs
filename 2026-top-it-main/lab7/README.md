# Лабораторная работа №7: Observability (Prometheus, Grafana, Grafana Tempo)

## Цель работы

Научиться:

- экспортировать **метрики** в формате Prometheus и настраивать их **скрейпинг**;
- строить **дашборды** в Grafana с запросами **PromQL**;
- внедрять **распределённый трейсинг** (OpenTelemetry, OTLP) и принимать трейсы в **Grafana Tempo**;
- разделять **код приложения** (инструментирование) и **платформенный стек** наблюдаемости по аналогии с лабораторной №6.

**Связь с лабораторными №5–6:** практика опирается на [telegram-support-bot](../lab5/examples/telegram-support-bot) и разделение приложения / инфраструктуры: PostgreSQL остаётся в [telegram-support-infra](../lab6/examples/telegram-support-infra/README.md), а Prometheus, Grafana и Tempo для Kubernetes выносятся в [examples/telegram-support-observability](examples/telegram-support-observability/README.md).

## Требования

- Выполнены [лабораторная работа №5](../lab5/README.md) (Kubernetes, `kubectl`) и [лабораторная работа №6](../lab6/README.md) (Helm, Kustomize, пример с БД и приложением) либо эквивалентные навыки.
- Установлены [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/) (версия 2.x).
- Для части с Kubernetes: кластер (например Docker Desktop), [Helm](https://helm.sh/docs/intro/install/) 3.x.
- Собраны образы примера [telegram-support-bot](../lab5/examples/telegram-support-bot) или вашего проекта из [лаб. №4](../lab4/README.md).

## Расширенная документация

Теория наблюдаемости, золотые сигналы, ссылки на внешние источники — в папке [docs/](docs/README.md):

- [Обзор и три столпа](docs/README.md)
- [Метрики, Prometheus и Grafana](docs/prometheus-grafana.md)
- [Трейсинг и Tempo](docs/tracing-tempo.md)
- [Логи в наблюдаемости](docs/logs-observability.md)

## Педагогический принцип: приложение и платформа

- В **репозитории приложения** — код, Dockerfile, `/metrics`, OpenTelemetry SDK, манифесты **только** frontend/backend (как в лаб. 5–6).
- В **каталоге/репозитории наблюдаемости** — установка Prometheus Operator, Grafana, Tempo, `ServiceMonitor`, политики хранения; это ближе к роли команды платформы / SRE.
- Секреты (токены бота, пароли Grafana в проде) в Git не кладутся; для учебы допустимы явные пароли в compose/values **только** на локальной машине.

## Эталоны в репозитории

| Каталог | Содержимое |
|---------|------------|
| [lab5/examples/telegram-support-bot](../lab5/examples/telegram-support-bot) | Backend с `/metrics` и OTLP, `docker-compose.observability.yml`, каталог [observability/](../lab5/examples/telegram-support-bot/observability/) |
| [examples/telegram-support-observability](examples/telegram-support-observability/README.md) | Локальный `docker-compose.yml`, манифесты Tempo, Helm `kube-prometheus-stack`, `ServiceMonitor` |

---

## Содержание

- [Часть A: Локальный стек наблюдаемости (Docker Compose)](#часть-a-локальный-стек-наблюдаемости-docker-compose)
  - [A.1. Два варианта эталона](#a1-два-варианта-эталона)
  - [A.2. Практика: приложение и overlay в каталоге бота](#a2-практика-приложение-и-overlay-в-каталоге-бота)
  - [A.3. Практика: только стек в каталоге observability](#a3-практика-только-стек-в-каталоге-observability)
  - [A.4. Проверка](#a4-проверка)
- [Часть B: Kubernetes](#часть-b-kubernetes)
  - [B.1. Порядок развёртывания](#b1-порядок-развёртывания)
  - [B.2. OTLP для приложения в кластере](#b2-otlp-для-приложения-в-кластере)
- [Задание для самостоятельной работы](#задание-для-самостоятельной-работы)
- [Контрольные вопросы](#контрольные-вопросы)
- [Дополнительные материалы](#дополнительные-материалы)

---

## Часть A: Локальный стек наблюдаемости (Docker Compose)

### A.1. Два варианта эталона

| Вариант | Когда удобен | Порты на хосте (типично) |
|---------|----------------|---------------------------|
| **1.** `docker-compose.yml` + `docker-compose.observability.yml` в [telegram-support-bot](../lab5/examples/telegram-support-bot) | Полный сценарий: Postgres, backend, frontend и стек в **одной** Docker-сети; scrape по имени сервиса `backend` | Prometheus **9090**, Grafana **3001**, Tempo **3200** |
| **2.** Только [docker-compose.yml](examples/telegram-support-observability/docker-compose.yml) в [telegram-support-observability](examples/telegram-support-observability/README.md) | Проверка Prometheus/Grafana/Tempo **без** пересборки бота; порты не конфликтуют с вариантом 1 | Prometheus **19090**, Grafana **13001**, Tempo **13200**, OTLP **14317–14318** |

Подробности конфигов варианта 2 — в [examples/telegram-support-observability/README.md](examples/telegram-support-observability/README.md).

### A.2. Практика: приложение и overlay в каталоге бота

1. Скопируйте `.env.example` в `.env`, задайте `BOT_TOKEN` (см. README примера).
2. Запустите приложение вместе со стеком наблюдаемости:

```bash
cd lab5/examples/telegram-support-bot
docker compose -f docker-compose.yml -f docker-compose.observability.yml up -d
```

3. Убедитесь, что поднялись сервисы `backend`, `postgres`, `prometheus`, `grafana`, `tempo` (см. `docker compose ps`).

### A.3. Практика: только стек в каталоге observability

1. Из корня репозитория курса:

```bash
cd lab7/examples/telegram-support-observability
docker compose up -d
```

2. Подождите **15–20 с** и проверьте готовность Tempo (ingester). Команды проверки и учётные данные Grafana — в [README примера](examples/telegram-support-observability/README.md).
3. Чтобы Prometheus собирал метрики с backend на хосте, API должен слушать **8080** на машине (`http://localhost:8080/metrics`). Для OTLP с хоста используйте `OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:14318`.

### A.4. Проверка

| Действие | Ожидание (вариант 1 / вариант 2) |
|----------|----------------------------------|
| Health API | `http://localhost:8080/health` (только при запущенном боте) |
| Метрики | `http://localhost:8080/metrics` — текст Prometheus |
| Prometheus UI | `http://localhost:9090` или `http://localhost:19090` → **Status → Targets**, job `prometheus` **UP**; job приложения **UP**, если scrape настроен и сервис доступен |
| Grafana | `http://localhost:3001` или `http://localhost:13001` (логин/пароль см. README observability-каталога или [observability/README.md](../lab5/examples/telegram-support-bot/observability/README.md)) |
| Трейсы | После запросов к API — **Explore → Tempo** в Grafana |

Сгенерируйте нагрузку, например:

```bash
curl http://localhost:8080/api/stats
```

Повторите запрос несколько раз и убедитесь, что на дашборде (папка **Lab7**) меняются метрики, а в Tempo появляются трейсы.

---

## Часть B: Kubernetes

### B.1. Порядок развёртывания

Пошаговые команды, Helm-релиз `kube-prometheus-stack`, манифесты Tempo и `ServiceMonitor` — в [examples/telegram-support-observability/README.md](examples/telegram-support-observability/README.md). Краткий порядок:

1. При необходимости развернуть БД: [telegram-support-infra](../lab6/examples/telegram-support-infra/README.md).
2. Развернуть приложение: [telegram-support-bot/k8s](../lab5/examples/telegram-support-bot/k8s/README.md).
3. Развернуть namespace, Tempo и Prometheus stack из каталога `telegram-support-observability`; применить `ServiceMonitor` в namespace приложения.

### B.2. OTLP для приложения в кластере

Задайте в Deployment backend переменные (значения — из README observability, обычно URL сервиса Tempo в namespace `observability`):

| Переменная | Пример |
|------------|--------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `http://tempo.observability.svc.cluster.local:4318` |
| `OTEL_SERVICE_NAME` | `telegram-support-backend` |

---

## Задание для самостоятельной работы

1. **Метрики:** в своём сервисе (на базе эталона или проекта курса) реализуйте эндпоинт `/metrics`, метрики HTTP (**counter** и **histogram**) с умеренной кардинальностью labels (шаблон маршрута, не произвольный path с ID). Добавьте минимум одну **бизнес-метрику** (например число ответов поддержки или входящих обращений).
2. **Скрейпинг:** в Docker Compose — цель Prometheus в статусе **UP**, корректный `metrics_path`; в Kubernetes — **ServiceMonitor** (или эквивалент), цель видна в Prometheus.
3. **Трейсинг:** настройте экспорт OTLP в **Tempo**; в Grafana найдите trace по HTTP-запросу к вашему API.
4. **Grafana:** подключите datasources Prometheus и Tempo; оформите дашборд (можно взять эталонный JSON и добавить **не менее одной** панели с собственным PromQL).
5. **Отчёт:** в репозитории должны быть конфигурации **и скриншоты** (не заменять только текстом README):
   - Prometheus: страница **Targets**, ваши jobs **UP**;
   - Grafana: дашборд с метриками backend;
   - Grafana: **Explore → Tempo** с развёрнутым trace (span-ы по запросу к API).

Рекомендуемая структура для скриншотов в форке: `docs/screenshots/lab7/`.

### Критерии приёмки

- Приложение стабильно работает, если OTLP endpoint **не** задан (трейсинг опционален на уровне конфигурации).
- Scrape метрик без ошибок в сценарии, заданном преподавателем (локально и/или в кластере).
- В Grafana видны метрики и хотя бы один полный trace для API.
- Есть скриншоты по списку выше и краткая инструкция «как воспроизвести».

**Дополнительно (по желанию):** Loki + Promtail, **exemplars** в Prometheus для перехода с графика задержки к trace.

---

## Контрольные вопросы

1. Чем **метрики** принципиально отличаются от **логов** с точки зрения агрегации и стоимости хранения?
2. Что такое **cardinality** в Prometheus и почему опасно ставить `user_id` в label на каждый запрос?
3. Зачем в распределённой системе единый **trace id** и приёмник трейсов (например Tempo)?
4. Для чего в Kubernetes используется ресурс **ServiceMonitor** в связке с Prometheus Operator?
5. Чем локальный запуск стека в `telegram-support-observability/docker-compose.yml` удобен по сравнению с немедленным развёртыванием только в кластер?

---

## Дополнительные материалы

- Материалы курса: [папка docs/](docs/README.md)
- [Prometheus — документация](https://prometheus.io/docs/introduction/overview/)
- [Grafana — документация](https://grafana.com/docs/grafana/latest/)
- [Grafana Tempo — документация](https://grafana.com/docs/tempo/latest/)
- [OpenTelemetry](https://opentelemetry.io/docs/)
- [Google SRE Book — мониторинг распределённых систем](https://sre.google/sre-book/monitoring-distributed-systems/)
