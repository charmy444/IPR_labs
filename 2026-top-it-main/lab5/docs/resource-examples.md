# Примеры работы с ресурсами Kubernetes

## Введение

В этом документе представлены подробные примеры работы с различными ресурсами Kubernetes, включая ConfigMap, Secret, PersistentVolume, PersistentVolumeClaim, Job, CronJob и другие. Эти ресурсы позволяют управлять конфигурацией, секретами, хранением данных и периодическими задачами в кластере Kubernetes.

Концептуальный разбор **базовых** объектов (Pod, Deployment, Service, метки, тома) с большим числом ссылок на официальную документацию — в серии файлов [Базовые ресурсы (подробно)](resources/README.md).

## ConfigMap

ConfigMap используется для хранения нечувствительных данных конфигурации в виде пар ключ-значение. Это позволяет отделить конфигурацию приложения от контейнерного образа.

### Создание ConfigMap

ConfigMap можно создать несколькими способами:

#### Из командной строки

```bash
# Создание ConfigMap из отдельных ключей
kubectl create configmap app-config \
  --from-literal=database.host=postgres \
  --from-literal=database.port=5432 \
  --from-literal=log.level=info

# Создание ConfigMap из файла
kubectl create configmap app-config \
  --from-file=application.properties
```

В PowerShell перенос строки удобно задавать символом обратной кавычки (`` ` ``) в конце строки либо выполнять команду в одну строку.

#### Из манифеста YAML

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: default
data:
  database.host: "postgres"
  database.port: "5432"
  log.level: "info"
  application.properties: |
    server.port=8080
    spring.datasource.url=jdbc:postgresql://postgres:5432/mydb
    spring.datasource.username=user
    spring.datasource.password=password
```

### Использование ConfigMap в Pod

ConfigMap можно использовать в Pod несколькими способами:

#### Как переменные окружения

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: nginx:1.20
    env:
    - name: DATABASE_HOST
      valueFrom:
        configMapKeyRef:
          name: app-config
          key: database.host
    - name: DATABASE_PORT
      valueFrom:
        configMapKeyRef:
          name: app-config
          key: database.port
    envFrom:
    - configMapRef:
        name: app-config
```

#### Как том (volume)


```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: nginx:1.20
    volumeMounts:
    - name: config-volume
      mountPath: /etc/config
  volumes:
  - name: config-volume
    configMap:
      name: app-config
```

## Secret

Secret используется для хранения чувствительных данных, таких как пароли, токены и ключи. Данные в Secret хранятся в base64-кодировке.

### Создание Secret

#### Из командной строки

```bash
# Создание Secret с учетными данными
kubectl create secret generic db-secret \
  --from-literal=username=admin \
  --from-literal=password='S3cureP@ssw0rd!'

# Создание Secret из файла (TLS)
kubectl create secret generic tls-secret \
  --from-file=tls.crt=/path/to/tls.crt \
  --from-file=tls.key=/path/to/tls.key
```

#### Из манифеста YAML

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
  namespace: default
type: Opaque
data:
  username: YWRtaW4=  # base64 encoded 'admin'
  password: UzNjdXJlUEBzc3cwcmQh  # base64 encoded 'S3cureP@ssw0rd!'
```

> **Примечание:** Для кодирования в base64 используйте команду `echo -n 'string' | base64`.

### Использование Secret в Pod

#### Как переменные окружения

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: nginx:1.20
    env:
    - name: DB_USERNAME
      valueFrom:
        secretKeyRef:
          name: db-secret
          key: username
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: db-secret
          key: password
```

#### Как том (volume)


```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: nginx:1.20
    volumeMounts:
    - name: secret-volume
      mountPath: "/etc/secret"
      readOnly: true
  volumes:
  - name: secret-volume
    secret:
      secretName: db-secret
```

## PersistentVolume и PersistentVolumeClaim

PersistentVolume (PV) — это ресурс в кластере, представляющий собой кусок хранилища. PersistentVolumeClaim (PVC) — это запрос на хранилище от пользователя.

### Создание PersistentVolume

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-example
  labels:
    type: local
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: manual
  hostPath:
    path: "/mnt/data"
    type: DirectoryOrCreate
```

### Создание PersistentVolumeClaim

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-example
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: manual
```

### Использование PVC в Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: nginx:1.20
    volumeMounts:
    - name: storage
      mountPath: "/usr/share/nginx/html"
  volumes:
  - name: storage
    persistentVolumeClaim:
      claimName: pvc-example
```

## Job

Job создает один или несколько Pods и гарантирует, что указанное количество Pods успешно завершится.

### Пример Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
```

### Запуск Job с параметрами

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: process-item
spec:
  completions: 5
  parallelism: 2
  template:
    spec:
      containers:
      - name: processor
        image: my-processor:1.0
        env:
        - name: ITEM
          value: "$(ITEM)"
      restartPolicy: OnFailure
```

## CronJob

CronJob создает Job по расписанию, аналогично cron в Unix-системах.

### Пример CronJob

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"  # Каждую минуту
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox:1.28
            args:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```

### Расписание Cron

| Выражение | Описание |
|-----------|----------|
| `* * * * *` | Каждую минуту |
| `0 * * * *` | Каждый час |
| `0 0 * * *` | Каждый день в полночь |
| `0 0 * * 0` | Каждую неделю в воскресенье |
| `0 0 1 * *` | Каждый месяц 1-го числа |

## Horizontal Pod Autoscaler (HPA)


HPA автоматически масштабирует количество реплик в Deployment, ReplicaSet или StatefulSet на основе наблюдаемых метрик.

### Пример HPA

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: php-apache
  namespace: default
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: php-apache
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 50
  - type: Resource
    resource:
      name: memory
      target:
        type: AverageValue
        averageValue: 200Mi
```

### Создание HPA из командной строки

```bash
# Автомасштабирование на основе использования CPU
kubectl autoscale deployment php-apache \
  --cpu-percent=50 \
  --min=1 \
  --max=10

# Просмотр статуса HPA
kubectl get hpa
```

## NetworkPolicy

NetworkPolicy определяет, как Pods могут взаимодействовать с другими Pods и сетевыми конечными точками.

### Пример NetworkPolicy

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend-to-backend
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 5000
```

### Запрет всего трафика

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

## PodDisruptionBudget

PodDisruptionBudget ограничивает количество Pods, которые могут быть одновременно недоступны из-за добровольных прерываний.

### Пример PodDisruptionBudget

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: my-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: my-app
```

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: my-pdb
spec:
  maxUnavailable: 25%
  selector:
    matchLabels:
      app: my-app
```

## Заключение

Эти примеры демонстрируют основные способы работы с различными ресурсами Kubernetes. Понимание и правильное использование этих ресурсов позволяет создавать надежные, масштабируемые и безопасные приложения в кластере Kubernetes.

## Официальная документация по разделам этого файла

- [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/) (en)
- [Secret](https://kubernetes.io/docs/concepts/configuration/secret/) (en)
- [Volumes](https://kubernetes.io/docs/concepts/storage/volumes/) · [PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) · [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/) (en)
- [Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/) · [CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cronjobs/) (en)
- [Horizontal Pod Autoscaling](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (en)
- [Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/) (en)
- [Pod Disruption Budget](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/#pod-disruption-budgets) (en)

Дополнительные ссылки (туториалы, GitLab, локальный кластер) — в конце [оглавления `docs/README.md`](README.md).