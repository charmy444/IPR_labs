# Рабочие нагрузки: Pod, ReplicaSet, Deployment

## Pod: минимальная единица планирования

**Pod** — это группа из одного или нескольких контейнеров с общим сетевым namespace (один IP на под при стандартной модели CNI), общими томами и общим контекстом Linux (PID namespace по умолчанию не общий для всех контейнеров — см. документацию по `shareProcessNamespace`).

### Зачем в поде несколько контейнеров

Типичные паттерны: **sidecar** (прокси, сбор логов), **adapter** (нормализация формата), **init**-контейнеры для подготовки тома до старта основного процесса. В учебных лабораториях чаще один контейнер на под; в продакшене sidecar встречается постоянно.

### Жизненный цикл и политика перезапуска

- Поле `spec.restartPolicy`: `Always` (по умолчанию), `OnFailure`, `Never`. Для **Job** обычно `OnFailure` или `Never`; для долгоживущих сервисов поды создаёт **Deployment** с `Always`.
- Фазы пода и условия готовности описаны в документации по жизненному циклу.

### Init-контейнеры

Контейнеры из `spec.initContainers` выполняются **по порядку** до старта `containers`. Удобно для миграций, ожидания доступности БД, скачивания артефактов. См. официальный гайд по init-контейнерам.

### Пробы: liveness, readiness, startup

- **livenessProbe** — «процесс в бесконечном зависании?» Неуспех → **перезапуск** контейнера.
- **readinessProbe** — «можно слать трафик?» Неуспех → под **убирается из Endpoints** Service.
- **startupProbe** — даёт медленно стартующему приложению время до того, как начнут действовать liveness/readiness (актуально для JVM, тяжёлых миграций при старте).

Подробности и типы проб (HTTP, TCP, exec): см. раздел про пробы в жизненном цикле пода.

### Ресурсы: requests и limits

- **requests** — то, что планировщик учитывает при выборе узла; гарантированный минимум по CPU/memory (модель зависит от версии и настроек кластера).
- **limits** — верхняя граница; для CPU обычно throttling, для memory — OOMKill при превышении.

Практика: начните с небольших значений, смотрите `kubectl top pod` (нужен [metrics-server](https://github.com/kubernetes-sigs/metrics-server)) или метрики из лаб. по наблюдаемости.

### ReplicaSet

**ReplicaSet** поддерживает заданное число подов по шаблону и селектору. Пользователи почти всегда создают **Deployment**, который сам управляет ReplicaSet. Прямое редактирование ReplicaSet для обновления образа — антипаттерн: для этого есть Deployment.

### Deployment: декларативное масштабирование и обновления

Ключевые поля:

- `spec.replicas` — желаемое число подов.
- `spec.selector` — **должен совпадать** с метками в `spec.template.metadata.labels` (иначе API отклонит манифест).
- `spec.strategy` — чаще всего `RollingUpdate` с `maxUnavailable` и `maxSurge`; альтернатива `Recreate` (короткий даунтайм, но проще модель).

Откаты: `kubectl rollout undo`, история ревизий в ReplicaSet-ах предыдущих поколений (см. документацию Deployment).

### Полезные команды

```bash
kubectl get pods -o wide
kubectl describe pod <name>
kubectl logs deploy/<name> --tail=100
kubectl rollout status deployment/<name>
kubectl rollout history deployment/<name>
```

## Типичные ошибки

- Несовпадение `selector` Deployment и `template.metadata.labels`.
- Слишком агрессивный **livenessProbe** без **startupProbe** — контейнер убивают до окончания прогрева.
- Отсутствие **readinessProbe** — трафик попадает на под, который ещё не готов отвечать.

## Референсы

- [Pods](https://kubernetes.io/docs/concepts/workloads/pods/) (en) · [Жизненный цикл пода](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/) (en)
- [Init Containers](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) (en)
- [Container probes](https://kubernetes.io/docs/concepts/configuration/liveness-readiness-startup-probes/) (en)
- [Управление ресурсами контейнеров](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) (en)
- [ReplicaSet](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/) (en)
- [Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) (en) · [Стратегии обновления](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy) (en)
- [Pod QoS](https://kubernetes.io/docs/concepts/workloads/pods/pod-qos/) (en) — как requests/limits влияют на классы качества обслуживания

Наверх: [оглавление раздела «Базовые ресурсы»](README.md).
