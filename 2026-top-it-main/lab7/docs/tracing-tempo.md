# Трейсинг и Grafana Tempo

## Зачем нужны трейсы

**Трейс (trace)** описывает **один логический запрос** (или задание) сквозь систему. Он состоит из **span**-ов: отрезков работы с началом, концом, атрибутами и связью «родитель — потомок». Так видно:

- **какие сервисы** участвовали в обработке;
- **сколько времени** занял каждый шаг;
- где возникла **ошибка** (статус span, события).

Без трейсов при пяти микросервисах вы видите пять разрозненных картин по логам и метрикам; с трейсом — **одну временную шкалу** для всего запроса.

Введение в концепции:

- [OpenTelemetry — Traces](https://opentelemetry.io/docs/concepts/signals/traces/)
- [OpenTelemetry — Distributed tracing](https://opentelemetry.io/docs/concepts/observability-primer/#distributed-traces) (связь контекста между процессами)

## Единый бэкенд (коллектор / хранилище) и сквозная отладка

В распределённой системе каждый сервис может порождать span-ы. Чтобы собрать их в **одно целое**, нужны:

1. **Общий идентификатор трейса** (`trace_id`) и связи span-ов — стандарт **W3C Trace Context** (заголовки вроде `traceparent` для HTTP).

   - [W3C Trace Context](https://www.w3.org/TR/trace-context/)

2. **Экспорт** span-ов в одну точку приёма. В OpenTelemetry это чаще всего протокол **OTLP** (gRPC **4317** или HTTP **4318**).

   - [OTLP specification](https://opentelemetry.io/docs/specs/otlp/)
   - [OTLP Exporter (Go и др.)](https://opentelemetry.io/docs/languages/go/exporters/)

3. **Хранилище и UI** для поиска и отображения: в лабораторной — **Grafana Tempo**; альтернативы — Jaeger, Zipkin, коммерческие APM.

   - [Grafana Tempo — Introduction](https://grafana.com/docs/tempo/latest/getting-started/)

**Зачем единый коллектор или единый endpoint для всех сервисов:**

- приложения конфигурируются **одинаково** (`OTEL_EXPORTER_OTLP_ENDPOINT`), проще онбординг новых сервисов;
- **OpenTelemetry Collector** может принимать несколько форматов, обрезать PII, семплировать трейсы (снизить стоимость), маршрутизировать в несколько бэкендов;

  - [OpenTelemetry Collector — Introduction](https://opentelemetry.io/docs/collector/)
  - [Collector deployment patterns](https://opentelemetry.io/docs/collector/deployment/)

- в **Grafana** один datasource **Tempo** даёт доступ ко **всем** трейсам системы: поиск по сервису, по длительности, по trace id из лога или метрики.

**Как это помогает дебажить сквозные вызовы (практический сценарий):**

1. Пользователь сообщает: «ошибка в 14:05» или вы видите всплеск **ошибок в метрике**.
2. Находите **медленный или ошибочный** trace в Tempo (по времени, сервису, тегам).
3. В дереве span-ов видно: например **800 ms** ушло на вызов **БД**, **50 ms** — на сервис **уведомлений**; ошибка — **timeout** на исходящем HTTP.
4. При необходимости переходите в **логи** того же пода, подставив **trace id** (корреляция).

Связка метрика ↔ трейс через **exemplars** (в Prometheus и Grafana):

- [Prometheus — Exemplars](https://prometheus.io/docs/prometheus/latest/feature_flags/#exemplars-storage)
- [Grafana — Intro to exemplars](https://grafana.com/docs/grafana/latest/fundamentals/exemplars/)

## OpenTelemetry (OTel) в лабораторной

Приложение создаёт span-ы через **SDK** и отправляет их **OTLP** в Tempo. Типичные переменные окружения (см. `.env.example` эталона):

| Переменная | Назначение |
|------------|------------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Базовый URL OTLP (например `http://tempo:4318` для HTTP) |
| `OTEL_SERVICE_NAME` | Имя сервиса в UI (например `telegram-support-backend`) |

Семантические соглашения по атрибутам (имя HTTP-маршрута, кода ответа и т.д.) описаны в **OpenTelemetry semantic conventions**:

- [Semantic conventions — HTTP](https://opentelemetry.io/docs/specs/semconv/http/)

Если endpoint не задан, эталонное приложение использует **no-op** tracer: метрики работают, трейсы не отправляются.

## Grafana Tempo

Tempo ориентирован на хранение трейсов с опорой на **объектное/локальное** хранилище и интеграцию с Grafana. Запросы и поиск — через документацию по API и Grafana Explore.

- [Tempo — Architecture](https://grafana.com/docs/tempo/latest/architecture/)
- [Tempo — Best practices](https://grafana.com/docs/tempo/latest/operations/best-practices/)

## Корреляция: трейсы, логи, метрики

Идеальная цепочка расследования:

1. **Метрика** или алерт указывает на проблему (золотые сигналы / RED).
2. **Трейс** локализует узкое место в цепочке.
3. **Логи** с тем же `trace_id` дают деталь сообщения и stack trace.

Подробнее про логи в этом курсе: [Логи в наблюдаемости](logs-observability.md).

Обзорный документ с тремя столпами и ссылками: [README документации](README.md).

## Дополнительные материалы

- [Distributed Systems Observability (Cindy Sridharan) — free O’Reilly ebook](https://www.oreilly.com/library/view/distributed-systems-observability/9781492033438/) — целостная картина метрик, логов, трейсов.
- [Jaeger — Documentation](https://www.jaegertracing.io/docs/) — альтернативная связка, те же идеи контекста.
- [OpenTelemetry — Getting started by language](https://opentelemetry.io/docs/languages/) — инструментирование сервисов.
