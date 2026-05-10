# Метрики, Prometheus и Grafana

## Зачем нужны метрики

**Метрики** — это числовые измерения, собранные во времени: сколько запросов обработано, какова задержка, сколько ошибок, какова загрузка CPU. Они отвечают на вопросы **агрегированно** («как система вела себя за последний час для всех пользователей»), а не только на единичный запрос.

Преимущества при правильном использовании:

- **Низкая стоимость** хранения по сравнению с полным логированием каждого запроса.
- Удобство для **дашбордов** и **оповещений** (правила в Prometheus, Alertmanager, облачные аналоги).
- Естественная опора для **SLI/SLO** (например доля запросов быстрее 300 ms).

Официальная документация и концепции:

- [Prometheus — What is Prometheus?](https://prometheus.io/docs/introduction/overview/)
- [Prometheus — Data model](https://prometheus.io/docs/concepts/data_model/) (имя метрики + labels = временной ряд).
- [Prometheus — Metric types](https://prometheus.io/docs/concepts/metric_types/) (counter, gauge, histogram, summary).

## Модель Prometheus

- **Pull-модель:** сервер Prometheus **опрашивает** HTTP-эндпоинты (`/metrics`) у целей по расписанию. Это упрощает обнаружение «живых» инстансов и единый контроль частоты scrape (альтернатива — push-модель, см. Pushgateway для особых случаев).

 - [Prometheus — Pull model](https://prometheus.io/docs/introduction/faq/#why-do-you-pull-rather-than-push)

- **Метки (labels):** каждая временная серия идентифицируется именем метрики и набором `key=value`. Нестабильные или слишком детальные labels (например уникальный `user_id` на каждый запрос) раздувают **cardinality** и нагрузку на TSDB.

  - [Grafana blog — Cardinality management](https://grafana.com/blog/2022/02/15/what-is-cardinality-and-how-does-it-impact-metrics-and-cost/)

- **Запросы:** язык **PromQL** для выборки и агрегации в Grafana и в правилах записи/алертинга.

  - [Prometheus — Querying basics](https://prometheus.io/docs/prometheus/latest/querying/basics/)
  - [Grafana — Prometheus data source](https://grafana.com/docs/grafana/latest/datasources/prometheus/)

## Типы метрик (кратко)

| Тип | Назначение | Примеры |
|-----|------------|---------|
| **Counter** | Монотонно неубывающий счётчик (сброс при рестарте) | Число обработанных запросов, ошибок |
| **Gauge** | Значение может расти и падать | Длина очереди, число открытых соединений |
| **Histogram** | Распределение наблюдений по **buckets** (корзинам) | Задержка HTTP:0.1s, 0.5s,1s, … |

Histogram позволяет в PromQL строить квантили (например p95) через функцию `histogram_quantile`:

- [Prometheus — Histograms and summaries](https://prometheus.io/docs/practices/histograms/)

## RED и золотые сигналы (связь с метриками)

Методика **RED** для сервиса, обслуживающего запросы:

- **Rate** — скорость запросов (из counter).
- **Errors** — доля или скорость ошибок.
- **Duration** — распределение задержек (histogram).

Первоисточник по «золотым сигналам» (latency, traffic, errors, saturation) — глава SRE-книги Google:

- [SRE Book — The Four Golden Signals](https://sre.google/sre-book/monitoring-distributed-systems/#xref_monitoring_alerting_latency-traffic-errors-saturation)

RED удобно маппится на метрики вроде `http_requests_total` + `http_request_duration_seconds_bucket`.

Дополнительно для **ресурсов** применяют **USE** (utilization, saturation, errors):

- [Brendan Gregg — USE Method](https://www.brendangregg.com/usemethod.html)

## Grafana

**Grafana** — визуализация и исследование данных: дашборды, алерты (в зависимости от редакции), Explore. Datasource **Prometheus** использует **PromQL**.

- [Grafana documentation](https://grafana.com/docs/grafana/latest/)
- [Provisioning — datasources and dashboards](https://grafana.com/docs/grafana/latest/administration/provisioning/) (как в учебном примере `telegram-support-bot/observability/grafana/`)

## Типичные ошибки

- **Взрыв cardinality:** слишком много уникальных комбинаций labels (ID пользователя, полный URL с UUID в path и т.д.). Предпочтительно использовать **шаблон маршрута** (в Gin: `c.FullPath()` после `Next()`).
- Путать **counter** и **gauge** для «количества запросов» — для запросов нужен counter (или histogram с подсчётом).
- Слишком **узкие** buckets гистограммы для latency — плохо оцениваются квантили; слишком широкие — теряется деталь.

## Дополнительные материалы

- [Google SRE — Monitoring distributed systems (full chapter)](https://sre.google/sre-book/monitoring-distributed-systems/) — философия мониторинга, белый шум, алерты.
- [Prometheus — Best practices](https://prometheus.io/docs/practices/) — именование, labels, histograms.
- [CNCF — Prometheus graduated project](https://www.cncf.io/projects/prometheus/) — статус проекта в CNCF.
