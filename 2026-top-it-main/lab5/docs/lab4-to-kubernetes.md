# От лабораторной №4 к Kubernetes: развёртывание своего приложения

В [лабораторной работе №4](../../lab4/README.md) вы настроили GitLab CI/CD, собирали Docker-образы и публиковали их в **GitLab Container Registry** (или Docker Hub). В лабораторной №5 цель та же по смыслу «доставки», но среда исполнения — **Kubernetes**: вы описываете желаемое состояние декларативно (YAML), а планировщик кластера поддерживает нужное число подов, сеть и обновления.

Этот документ связывает артефакты лаб. №4 с шагами деплоя в k8s.

## Что у вас должно быть после лаб. №4

- Репозиторий с рабочим `.gitlab-ci.yml`: сборка образа после тестов, публикация в registry.
- Один или несколько `Dockerfile` (например, отдельно для API и фронтенда — как в вашем проекте).
- Понимание переменных `$CI_REGISTRY`, `$CI_REGISTRY_IMAGE`, `$CI_REGISTRY_USER`, `$CI_JOB_TOKEN` / пароля для `docker login`.

Если образ ещё не публикуется в registry, сначала доведите пайплайн лаб. №4 до состояния «образ появляется в Container Registry проекта».

## От Docker Compose к объектам Kubernetes

| В `docker-compose` | В Kubernetes (минимальный набор) |
|--------------------|-----------------------------------|
| Сервис (контейнер) | `Deployment` (+ контейнер в `spec.template.spec.containers`) |
| `ports`, имя сервиса | `Service` (`ClusterIP` / `NodePort` / `LoadBalancer`) |
| `environment`, `.env` | `ConfigMap` и/или `Secret`, поле `env` / `envFrom` в Pod |
| Тома данных | `PersistentVolumeClaim` + монтирование в Pod (при необходимости) |
| Зависимость «БД поднялась» | Отдельный `Deployment`/`StatefulSet` для БД или внешняя БД; `initContainer` или порядок применения манифестов |

Важно: **DNS внутри кластера**. Поды обращаются друг к другу по имени **Service**, а не по имени пода. Имя вида `http://my-backend:8080` резолвится в виртуальный IP сервиса, который балансирует трафик между подами.

## Шаг 1. Выбор образов

### Вариант A: образы из GitLab Container Registry (основной для зачёта с лаб. №4)

1. В GitLab откройте проект → **Deploy → Container Registry** и скопируйте путь образа, например:
   - `registry.gitlab.mai.ru/group/subgroup/project/backend:main`
2. В манифесте `Deployment` укажите полный путь в поле `image:`.

Приватный registry: кластер при скачивании образа должен авторизоваться.

**Создание секрета для pull (тип `docker-registry`):**

```bash
kubectl create secret docker-registry gitlab-registry \
  --docker-server=registry.gitlab.mai.ru \
  --docker-username=<ваш_логин_или_ci_token_name> \
  --docker-password=<token_или_пароль> \
  --namespace=<ваш_namespace>
```

В `Deployment` (в `spec.template.spec`):

```yaml
spec:
  template:
    spec:
      imagePullSecrets:
        - name: gitlab-registry
      containers:
        - name: app
          image: registry.gitlab.mai.ru/your/group/project:tag
```

Подставьте реальный registry вашей инсталляции GitLab (в лаб. №4 используется `gitlab.mai.ru`).

### Вариант B: локальные образы (Docker Desktop Kubernetes)

Если вы собрали образ локально (`docker build -t myapp:1.0 .`) и кластер тот же Docker Desktop, образ доступен узлу без push. Укажите в манифесте то же имя тега, что при сборке. Для учебных заданий это удобно; для «как в проде» предпочтителен вариант A.

## Шаг 2. Namespace и изоляция

Создайте отдельный namespace под проект (как в основном README лаб. №5):

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: lab5-myapp
```

Все последующие ресурсы указывайте с `metadata.namespace: lab5-myapp` или применяйте с `-n lab5-myapp`.

## Шаг 3. Секреты и конфигурация

**Не кладите** токены и пароли в Git в открытом виде.

- Нечувствительное: URL API для фронта, уровень логов — `ConfigMap`.
- Пароли БД, API-ключи, токены ботов — `Secret` (данные в base64 в YAML или создание через `kubectl create secret generic ... --from-literal=...`).

Пример структуры для backend (идея, не копируйте секреты в репозиторий):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: lab5-myapp
type: Opaque
stringData:
  database-url: "postgres://user:pass@postgres:5432/dbname?sslmode=disable"
  api-key: "your-secret"
```

`stringData` удобнее, чем `data` + base64, при ручном создании манифестов.

Фронтенд (Next.js и аналоги): переменные с префиксом `NEXT_PUBLIC_` встраиваются на этапе сборки. В k8s часто либо собирают образ с нужным `ARG`/`ENV` в CI, либо задают публичный URL API, доступный браузеру (Ingress/NodePort), отдельно от внутреннего имени Service.

## Шаг 4. Deployment и Service для каждого компонента

Для **каждого** процесса из вашего `docker-compose` (backend, frontend, worker и т.д.):

1. **Deployment**: образ, `replicas`, `ports`, `env`, при необходимости `resources`, `livenessProbe` / `readinessProbe`.
2. **Service**: `selector` совпадает с `labels` подов; `port` / `targetPort` как в контейнере.

Имена Service используйте в переменных окружения других сервисов (как `BACKEND_URL=http://backend:8080` в compose → то же имя Service в k8s).

Минимальный каркас:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-backend
  namespace: lab5-myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
      component: backend
  template:
    metadata:
      labels:
        app: my-app
        component: backend
    spec:
      containers:
        - name: backend
          image: registry.gitlab.mai.ru/your/project/backend:main
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: my-backend
  namespace: lab5-myapp
spec:
  selector:
    app: my-app
    component: backend
  ports:
    - port: 8080
      targetPort: 8080
```

## Шаг 5. База данных

Учебные варианты:

1. **PostgreSQL/MySQL в кластере** — отдельный `Deployment` + `Service` + `PersistentVolumeClaim` (данные переживают перезапуск пода, но зависят от настроек тома). Для экзаменационного прототипа иногда допускают ephemeral БД без PVC (данные пропадут при удалении пода).
2. **Внешняя управляемая БД** — в `Secret` только строка подключения; в кластере нет StatefulSet для БД.

Убедитесь, что `DATABASE_URL` в backend указывает на хост **имени Service** БД в k8s (например `postgres://...@postgres:5432/...`), а не `localhost`.

## Шаг 6. Применение и проверка

```bash
kubectl apply -f namespace.yaml
kubectl apply -f secret.yaml   # или создайте секреты через kubectl
kubectl apply -f configmap.yaml
kubectl apply -f backend-deployment.yaml
kubectl apply -f backend-service.yaml
# ... остальные манифесты
kubectl get pods,svc -n lab5-myapp
kubectl logs -n lab5-myapp deploy/my-backend
kubectl describe pod -n lab5-myapp -l component=backend
```

Типичные проблемы:

| Симптом | Что проверить |
|---------|----------------|
| `ImagePullBackOff` | Правильность пути образа, `imagePullSecrets`, доступность registry |
| `CrashLoopBackOff` | `kubectl logs`, переменные окружения, доступность БД |
| Фронт не достучится до API | Имя Service, порт, CORS; для браузера — публичный URL, а не внутреннее DNS |
| Под `Running`, но не в `Ready` | `readinessProbe`, зависимости старта |

Доступ с хоста (Docker Desktop): часто используют `Service` типа `NodePort` или порт-форвард:

```bash
kubectl port-forward -n lab5-myapp svc/my-frontend 3000:3000
```

## Шаг 7. Связь с CI/CD (по желанию, продвинутый уровень)

Лаб. №4 заканчивается публикацией образа. Лаб. №5 — ручной или полуавтоматический деплой манифестов. Дальнейшее развитие (не обязательно для базовой лаб. №5):

- обновление тега образа в манифесте и `kubectl apply`;
- или GitOps (Argo CD, Flux);
- или отдельный job `deploy` в GitLab с `kubectl` / Helm на runner с доступом к кластеру.

## Эталон в репозитории лаб. №5

В [`examples/telegram-support-bot`](../examples/telegram-support-bot/k8s/) приведены манифесты **только приложения** (backend Go, frontend Next.js, секреты, сервисы): каталог `k8s/kustomization` и `k8s/helm`. PostgreSQL вынесен в отдельный учебный «репозиторий инфраструктуры» — [telegram-support-infra](../../lab6/examples/telegram-support-infra/README.md). Ingress и Istio — в `k8s/optional/`. Полный сценарий Kustomize/Helm и разделение app/infra — в [лабораторной №6](../../lab6/README.md). Сверяйте структуру лейблов, проб и переменных окружения.

## Критерии «лабораторная выполнена» в части своего приложения

- Образ приложения из лаб. №4 (или логичное продолжение того же репозитория) успешно запускается в Kubernetes.
- Есть хотя бы один `Deployment` и `Service` для основного сервиса; при многосервисной архитектуре — для ключевых компонентов.
- Секреты не хранятся в открытом виде в Git; для приватного registry настроен pull-secret при необходимости.
- Вы можете показать работающую проверку (`curl`, браузер, health-endpoint) и объяснить, как поды находят друг друга по DNS.

Детали базовых концепций — в [Основы Kubernetes](kubernetes-basics.md); углублённо по ресурсам и ссылкам на kubernetes.io — [Базовые ресурсы (подробно)](resources/README.md); справочник YAML по ConfigMap, PVC, Job — в [Примеры ресурсов](resource-examples.md).

## Внешние материалы по теме документа

- [Образы и приватные registry](https://kubernetes.io/docs/concepts/containers/images/) (en)
- [Создание `imagePullSecret` и использование в Pod](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) (en)
- [Переменные окружения в контейнере](https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/) (en) · [ConfigMap как env / том](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/) (en)
- [Docker Compose и Kubernetes (сопоставление концепций, статья Kubernetes)](https://kubernetes.io/docs/tasks/configure-pod-container/translate-compose-kubernetes/) (en)
- [GitLab Container Registry](https://docs.gitlab.com/ee/user/packages/container_registry/) (en)
- [Деплой в Kubernetes из GitLab (обзор возможностей продукта)](https://docs.gitlab.com/ee/user/project/clusters/deploy_to_cluster.html) (en) — для справки, в базовой лаб. №5 достаточно ручного `kubectl apply`
- [Helm — менеджер пакетов для Kubernetes](https://helm.sh/docs/intro/quickstart/) (en) — опционально, если захотите шаблонизировать манифесты

Общий список ссылок для лабораторной — в конце [README папки `docs`](README.md).
