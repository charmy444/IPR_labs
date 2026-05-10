# Инфраструктура наблюдаемости для Telegram Support Bot (лаб. 7)

Отдельный учебный «репозиторий платформы»: **Prometheus Operator (kube-prometheus-stack)**, **Grafana** и **Grafana Tempo**. Код приложения и манифесты frontend/backend по-прежнему в [telegram-support-bot](../../../lab5/examples/telegram-support-bot); PostgreSQL — в [telegram-support-infra](../../../lab6/examples/telegram-support-infra/README.md).

## Зачем отдельный каталог

- Те же границы, что в лаб. 6: приложение ≠ платформенный стек мониторинга.
- Команда приложения не дублирует релизы Prometheus/Grafana в своём чарте.
- Секреты и политики хранения метрик/трейсов живут у владельцев платформы.

## Предпосылки

- **Локально:** Docker и Docker Compose (см. раздел [Docker Compose](#docker-compose-локальная-проверка-стека)).
- **Kubernetes:** кластер (например Docker Desktop), Helm 3.x; CRD `ServiceMonitor` (поставляется вместе с kube-prometheus-stack).

## Docker Compose (локальная проверка стека)

Тот же набор **Prometheus + Grafana + Tempo**, что и в примере бота, но в каталоге инфраструктуры — удобно проверять конфиги без k8s.

```bash
cd lab7/examples/telegram-support-observability
docker compose up -d
```

| Сервис | URL | Учётные данные |
|--------|-----|----------------|
| Grafana | http://localhost:13001 | `admin` / `admin` |
| Prometheus | http://localhost:19090 | — |
| Tempo (HTTP) | http://localhost:13200 | — |
| OTLP с хоста | `localhost:14318` (HTTP), `localhost:14317` (gRPC) | — |

Порты на хосте выбраны так, чтобы не пересекаться с [docker-compose.observability.yml](../../../lab5/examples/telegram-support-bot/docker-compose.observability.yml) из примера бота (`9090`, `3001`, `3200`, `4317–4318`).

Конфигурации лежат в [compose/](compose/README.md). Prometheus по умолчанию:

- scrape **самого себя** (цель всегда **UP**);
- дополнительный job **`telegram-support-backend`** — `host.docker.internal:8080` (если API бота слушает `8080` на хосте или проброшен с хоста, target станет **UP**; иначе **DOWN** — для проверки стека это нормально).

Чтобы трейсы шли в этот Tempo с хоста, задайте `OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:14318` для процесса backend.

Остановка: `docker compose down`.

**Проверка (после `up`):** подождите **15–20 с** — у Tempo endpoint `/ready` может временно отвечать неготовностью. Затем:

```bash
curl http://localhost:19090/-/healthy
curl http://localhost:13200/ready
```

В Prometheus → Status → Targets ожидается **UP** у job `prometheus`; у `telegram-support-backend` — **UP**, если на хосте доступен `http://localhost:8080/metrics` (например запущен backend бота).

## Порядок развёртывания в Kubernetes (типичный сценарий)

1. **PostgreSQL** (если нужна полная работа бота): [telegram-support-infra](../../../lab6/examples/telegram-support-infra/README.md).
2. **Приложение** `telegram-support-bot` в namespace, например `telegram-demo` ([k8s/README](../../../lab5/examples/telegram-support-bot/k8s/README.md)).
3. **Namespace и Tempo** из этого каталога:

```bash
cd lab7/examples/telegram-support-observability
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/tempo-configmap.yaml
kubectl apply -f k8s/tempo-deployment.yaml
kubectl apply -f k8s/tempo-service.yaml
```

4. **kube-prometheus-stack** (релиз в namespace `observability`):

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
  --namespace observability --create-namespace \
  -f helm/kube-prometheus-stack/values-dev.yaml
```

5. **ServiceMonitor** для backend приложения — в namespace, где лежит `Service` `telegram-backend`:

```bash
kubectl apply -f k8s/servicemonitor-backend.yaml -n telegram-demo
```

Убедитесь, что селекторы `matchLabels` совпадают с метками вашего Service (эталон: `app: telegram-support-bot`, `component: backend`).

## Контракт для приложения (OTLP)

Для отправки трейсов backend должен иметь переменные окружения (см. [лаб. 7](../../README.md)):

| Переменная | Пример значения |
|------------|-----------------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `http://tempo.observability.svc.cluster.local:4318` |
| `OTEL_SERVICE_NAME` | `telegram-support-backend` |

Добавьте их в Deployment приложения (Helm values, Kustomize patch или Secret).

## Доступ к Grafana и Prometheus

Проброс портов (имена pod/сервисов могут отличаться — проверьте `kubectl get svc -n observability`):

```bash
kubectl port-forward -n observability svc/prometheus-grafana 3000:80
kubectl port-forward -n observability svc/prometheus-kube-prometheus-prometheus 9090:9090
```

Логин Grafana по умолчанию в `values-dev.yaml`: **admin** / **admin** (смените для чего-либо кроме учебной машины).

## Проверка

1. Prometheus UI → Status → Targets — endpoint `telegram-support-backend` **UP**.
2. Grafana → Explore → **Tempo** — трейсы после HTTP-запросов к API приложения.
3. Метрики `/metrics` содержат `http_requests_total`, `http_request_duration_seconds_*`, бизнес-счётчики эталона.

## Файлы

| Путь | Назначение |
|------|------------|
| [docker-compose.yml](docker-compose.yml) | Локальный Prometheus, Grafana, Tempo |
| [compose/](compose/README.md) | Конфиги для compose |
| [k8s/tempo-*.yaml](k8s/tempo-deployment.yaml) | Monolithic Tempo + Service |
| [helm/kube-prometheus-stack/values-dev.yaml](helm/kube-prometheus-stack/values-dev.yaml) | Урезанный kube-prometheus-stack + datasource Tempo |
| [k8s/servicemonitor-backend.yaml](k8s/servicemonitor-backend.yaml) | Scrape `/metrics` у backend |

Подробнее о задании и скриншотах для отчёта — в [лабораторной №7](../../README.md).
