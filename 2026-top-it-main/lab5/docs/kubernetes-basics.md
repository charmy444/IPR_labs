# Основы Kubernetes

## Зачем нужен Kubernetes

**Kubernetes (k8s)** — система оркестрации контейнеров. Она отвечает на вопросы:

- как запустить **много копий** сервиса и держать их число стабильным;
- как **обновить** версию без полной остановки;
- как дать сервисам **стабильные имена и балансировку**, хотя отдельные контейнеры постоянно пересоздаются;
- как **ограничить ресурсы** и проверять «жив ли» контейнер.

Вы описываете **желаемое состояние** (YAML или API), а компоненты control plane приводят фактическое состояние кластера к этому описанию.

## Декларативная модель

Типичный цикл работы разработчика:

1. Собрать **образ** приложения (`docker build` / CI из лаб. №4).
2. Описать в манифесте **Deployment**: какой образ, сколько реплик, переменные окружения, пробы.
3. Описать **Service**: на какие поды направлять трафик и на каких портах.
4. Выполнить `kubectl apply -f ...`.

Kubernetes хранит объекты в **etcd**; **контроллеры** сравнивают «как надо» и «как есть» и создают/удаляют поды, endpoints и т.д.

Полезно помнить: манифест почти всегда содержит четыре логических уровня:

- `apiVersion` — версия API (`v1`, `apps/v1`, …).
- `kind` — тип ресурса (`Pod`, `Deployment`, `Service`, …).
- `metadata` — имя, namespace, **labels**, аннотации.
- `spec` — желаемая конфигурация именно этого типа.

## Архитектура кластера (упрощённо)

### Control plane (бывший «master»)

- **kube-apiserver** — единая точка входа для `kubectl` и внутренних компонентов.
- **etcd** — хранилище состояния кластера.
- **scheduler** — выбирает узел для нового пода.
- **controller-manager** — набор контроллеров (в т.ч. для Deployment, ReplicaSet, Service).

### Узлы (nodes)

- **kubelet** — запускает контейнеры на узле по заданию API.
- **kube-proxy** — участвует в реализации сетевых правил для Service.
- **container runtime** (containerd и др.) — фактический запуск контейнеров.

В **Docker Desktop** вы обычно видите один узел `docker-desktop`, на котором совмещены control plane и рабочая нагрузка — этого достаточно для учебных задач.

## Ключевые объекты (в порядке изучения)

Ниже — **сжатый** конспект. Развёрнутое объяснение с типичными ошибками и расширенным списком ссылок на [kubernetes.io](https://kubernetes.io/docs/) вынесено в серию файлов: [Базовые ресурсы (подробно)](resources/README.md) ([рабочие нагрузки](resources/workloads.md), [сервисы и DNS](resources/services-networking.md), [namespace и метки](resources/namespaces-labels.md), [конфигурация и тома](resources/config-secrets-volumes.md)).

### Pod

**Pod** — минимальная единица планирования: один или несколько контейнеров с общим сетевым namespace (один IP на под), общими томами при необходимости. В продакшене поды чаще **не создают вручную** — ими управляет Deployment.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
    - name: nginx
      image: nginx:1.20
      ports:
        - containerPort: 80
```

### Deployment и ReplicaSet

**Deployment** задаёт шаблон пода и желаемое число реплик; обеспечивает **rolling update** и откаты. Фактически он управляет **ReplicaSet**, который следит за числом работающих подов.

Важно: `spec.selector.matchLabels` в Deployment **должен совпадать** с `spec.template.metadata.labels`, иначе манифест будет отклонён.

### Service

**Service** — стабильный **виртуальный IP** (ClusterIP) и DNS-имя для группы подов, отобранных по **labels**. Типы:

| Тип | Назначение |
|-----|------------|
| `ClusterIP` | Доступ только **изнутри** кластера (по умолчанию). |
| `NodePort` | Пробрасывает порт на каждый узел — удобно для доступа с хоста в Docker Desktop. |
| `LoadBalancer` | Внешний балансировщик (в облаке); локально часто ведёт себя как NodePort. |

**DNS внутри кластера:** сервис `my-svc` в namespace `lab5` доступен как `my-svc.lab5.svc.cluster.local` и часто как короткое имя `my-svc` из того же namespace.

### Namespace

Логическая изоляция имён ресурсов и квот. Ресурсы с одинаковым `kind`+`name` могут существовать в разных namespace. Системные поды обычно в `kube-system`.

### ConfigMap и Secret

Вынесение конфигурации из образа. **ConfigMap** — для неконфиденциальных данных. **Secret** — для паролей и ключей (в etcd они должны быть защищены политиками кластера; в YAML — не коммитить реальные значения в публичный Git).

Подключение: `env`, `envFrom`, том `configMap` / `secret`.

### Ingress (опционально на первом этапе)

**Ingress** — правила маршрутизации HTTP(S) к Service; нужен **Ingress Controller** (nginx, traefik и т.д.). В минимальной лабораторной часто достаточно NodePort или `kubectl port-forward`.

## Labels и selectors

**Labels** — пары `ключ: значение` на объектах. **selector** в Service и Deployment выбирает поды с нужными метками. Это основа связи «какой трафик куда идёт».

Рекомендация: договоритесь о наборе меток для проекта, например `app`, `component`, `version`.

## Пробы: liveness и readiness

- **livenessProbe** — «процесс завис?» Если неуспешна, kubelet **перезапускает** контейнер.
- **readinessProbe** — «готов ли принимать трафик?» Если неуспешна, под **исключается** из endpoints Service (на него не пойдёт балансировка).

Типично для HTTP-сервиса:

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 15
  periodSeconds: 10
readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

`initialDelaySeconds` подбирают так, чтобы приложение успело стартовать до первых проверок.

## Ресурсы CPU и памяти (кратко)

В `containers[].resources` задают **requests** (планировщик резервирует место на узле) и **limits** (верхняя граница). Без limits контейнер может выжать память узла. Для учебного кластера можно начать с небольших значений и смотреть `kubectl top pod` (если установлен metrics-server).

## Обновление версии и откаты

При смене образа или шаблона Deployment запускается **rolling update**: создаются новые поды, старые постепенно удаляются.

```bash
kubectl set image deployment/my-deploy app=my-image:2.0 -n lab5
kubectl rollout status deployment/my-deploy -n lab5
kubectl rollout undo deployment/my-deploy -n lab5
```

## Основные команды kubectl

| Команда | Описание |
|---------|----------|
| `kubectl get pods` | Список подов (добавьте `-n namespace`) |
| `kubectl get svc,deploy` | Сервисы и деплойменты |
| `kubectl describe pod <name>` | События и причины ошибок |
| `kubectl logs <pod>` | Логи; `-f` — поток |
| `kubectl logs deploy/<name>` | Логи одного из подов деплоймента |
| `kubectl exec -it <pod> -- sh` | Оболочка в контейнере |
| `kubectl apply -f dir/` | Применить все манифесты из каталога |
| `kubectl delete -f file.yaml` | Удалить ресурсы из файла |

Шпаргалка: [Kubernetes Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/).

## Минимальный сценарий «поднять сервис»

Императивно (для эксперимента):

```bash
kubectl create deployment nginx --image=nginx:1.20
kubectl expose deployment nginx --port=80 --type=NodePort
kubectl get svc nginx
```

Декларативно (как в лабораторной): отдельные YAML для Namespace, Deployment, Service, Secret/ConfigMap — повторяемо и версионируемо в Git.

## Связь с лабораторной №4

Образы, которые вы собираете и пушите в GitLab Registry в лаб. №4, указываются в поле `image` в Deployment. Для приватного registry нужен `imagePullSecrets`. Пошаговый перенос приложения описан в [От лабораторной №4 к Kubernetes](lab4-to-kubernetes.md).

## Где читать дальше (внешние материалы)

Материалы по темам этого файла:

- [Что такое Kubernetes](https://kubernetes.io/ru/docs/concepts/overview/what-is-kubernetes/) (ru) · [What is Kubernetes](https://kubernetes.io/docs/concepts/overview/) (en)
- [Архитектура кластера](https://kubernetes.io/docs/concepts/architecture/) (en)
- [Pods](https://kubernetes.io/docs/concepts/workloads/pods/) · [Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) · [ReplicaSet](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/) (en)
- [Services](https://kubernetes.io/docs/concepts/services-networking/service/) · [DNS для сервисов и подов](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/) (en)
- [Namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) (en)
- [ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/) · [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/) (en)
- [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) · [Ingress controllers](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/) (en)
- [Назначение CPU и памяти контейнерам](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) (en)
- [Пробы (liveness, readiness, startup)](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes) (en)
- [Стратегии обновления Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy) (en) · [Rollouts и откаты](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#rolling-back-a-deployment) (en)

Сводный список ссылок (включая GitLab, Docker Desktop, интерактивные площадки) — в конце [оглавления документации лаб. №5](README.md).

## Дополнительно в репозитории

Углублённая теория по базовым ресурсам — в [docs/resources/](resources/README.md). Практические фрагменты YAML (тома, Job, HPA, сетевые политики) — в [Примеры ресурсов](resource-examples.md).
