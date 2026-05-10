# Локальный стек наблюдаемости (лаб. 7)

Конфигурации для **Prometheus**, **Grafana** и **Grafana Tempo** при запуске через [docker-compose.observability.yml](../docker-compose.observability.yml).

## Запуск

Из корня примера `telegram-support-bot`:

```bash
docker compose -f docker-compose.yml -f docker-compose.observability.yml up -d
```

## Доступ

| Сервис | URL | Учётные данные |
|--------|-----|----------------|
| Grafana | http://localhost:3001 | `admin` / `admin` |
| Prometheus | http://localhost:9090 | — |
| Tempo (HTTP API) | http://localhost:3200 | — |
| OTLP (хост) | `localhost:4318` | для клиентов вне Docker |

Backend в overlay получает `OTEL_EXPORTER_OTLP_ENDPOINT=http://tempo:4318`.

## Проверка

1. http://localhost:8080/metrics — текст метрик Prometheus.
2. Prometheus → Status → Targets — job `telegram-support-backend` в состоянии **UP**.
3. Grafana → Dashboards → папка **Lab7** — дашборд эталона.
4. Grafana → Explore → datasource **Tempo** — трейсы после запросов к API (например `curl http://localhost:8080/api/stats`).

## Файлы

- `prometheus/prometheus.yml` — scrape `backend:8080`.
- `tempo/tempo.yaml` — monolithic Tempo, OTLP HTTP4318.
- `grafana/provisioning/` — datasources Prometheus + Tempo.
- `grafana/dashboards/lab7-backend.json` — эталонный дашборд.
