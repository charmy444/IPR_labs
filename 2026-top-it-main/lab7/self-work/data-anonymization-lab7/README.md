# Самостоятельная работа lab7: SDA Observability

Это отдельная сдаваемая папка для lab7 на базе `data-anonymization-and-synthesis-tool228`.

В backend добавлены:

- `/metrics` в формате Prometheus;
- HTTP counter и latency histogram;
- бизнес-метрика `sda_business_events_total`;
- опциональный OTLP export в Tempo через `OTEL_EXPORTER_OTLP_ENDPOINT`.

Если OTLP endpoint не задан, backend продолжает работать без tracing.

## Вариант A: локально через Docker Compose

Из корня `2026-top-it-main`:

```bash
docker compose -f lab7/self-work/data-anonymization-lab7/docker-compose.yml up --build
```

Compose использует `dockerfiles/frontend.prebuilt.Dockerfile`, чтобы frontend запускался локально без `npm ci` внутри Docker. Это полезно, если npm падает с сообщением про `/root/.npm/_logs/...`.

Порты:

- frontend: `http://localhost:3000`
- backend docs: `http://localhost:8000/docs`
- backend health: `http://localhost:8000/api/v1/health`
- backend metrics: `http://localhost:8000/metrics`
- Prometheus: `http://localhost:19090`
- Grafana: `http://localhost:13001` (`admin` / `admin`)
- Tempo: `http://localhost:13200`

Сгенерируйте запросы:

```bash
curl http://localhost:8000/api/v1/health
curl http://localhost:8000/api/v1/generate/templates
curl http://localhost:8000/metrics
```

В Grafana:

- datasource `Prometheus`;
- datasource `Tempo`;
- dashboard `SDA Backend / Lab7`.

Остановить:

```bash
docker compose -f lab7/self-work/data-anonymization-lab7/docker-compose.yml down
```

## Вариант B: локальный Kubernetes

Сначала запустите приложение из lab5:

```bash
kubectl apply -f lab5/self-work/data-anonymization-k8s/manifests/
```

Затем включите OTLP endpoint для backend:

```bash
kubectl set env deployment/sda-backend -n sda-lab5 \
  OTEL_EXPORTER_OTLP_ENDPOINT=http://tempo.sda-observability.svc.cluster.local:4318/v1/traces \
  OTEL_SERVICE_NAME=sda-backend
```

Разверните observability:

```bash
kubectl apply -f lab7/self-work/data-anonymization-lab7/k8s/
```

Проверка:

```bash
kubectl get pods,svc,configmap -n sda-observability
kubectl port-forward -n sda-observability svc/prometheus 19090:9090
kubectl port-forward -n sda-observability svc/grafana 13001:3000
```

Prometheus будет скрейпить backend по адресу `sda-backend.sda-lab5.svc.cluster.local:8000/metrics`.

## Скриншоты для сдачи

Положите в `docs/screenshots/lab7/`:

- `prometheus-targets.png` — target `sda-backend` в статусе UP;
- `grafana-dashboard.png` — dashboard с метриками backend;
- `tempo-trace.png` — Grafana Explore -> Tempo с trace HTTP-запроса.
