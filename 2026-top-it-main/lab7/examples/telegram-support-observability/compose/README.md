# Конфигурации для `docker compose`

Файлы для сервисов из [docker-compose.yml](../docker-compose.yml) в родительском каталоге.

| Каталог | Назначение |
|---------|------------|
| `prometheus/prometheus.yml` | Scrape: сам Prometheus + опционально backend на `host.docker.internal:8080` |
| `tempo/tempo.yaml` | Monolithic Tempo, OTLP HTTP 4318 |
| `grafana/provisioning/` | Datasources Prometheus + Tempo, провайдер дашбордов |
| `grafana/dashboards/` | Эталонный дашборд Lab 7 (метрики бота) |
