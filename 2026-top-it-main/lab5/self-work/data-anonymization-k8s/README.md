# Самостоятельная работа lab5: SDA в Kubernetes

Это отдельная сдаваемая папка для lab5 на базе проекта `data-anonymization-and-synthesis-tool228`.

Сценарий рассчитан на локальный Docker Desktop Kubernetes. В проекте нет PostgreSQL или другой внешней БД: backend хранит временные upload/session данные в файловом хранилище. Поэтому для lab5 добавлены `PersistentVolumeClaim`, `ConfigMap`, `Secret`, `Deployment`, `Service`, `Ingress` и `HPA`.

## Состав

```text
manifests/
  00-namespace.yaml
  01-configmap.yaml
  02-secret.yaml
  03-storage.yaml
  04-backend.yaml
  05-frontend.yaml
  06-ingress.yaml
  07-hpa.yaml
```

## Сборка локальных образов

Команды выполняются из корня репозитория `2026-top-it-main`:

```bash
docker build -t sda-backend:local \
  -f data-anonymization-and-synthesis-tool228/dockerfiles/backend.Dockerfile \
  data-anonymization-and-synthesis-tool228

docker build -t sda-frontend:local \
  --build-arg NEXT_PUBLIC_API_BASE_URL=http://localhost:8000/api/v1 \
  -f data-anonymization-and-synthesis-tool228/dockerfiles/frontend.Dockerfile \
  data-anonymization-and-synthesis-tool228
```

Если frontend-сборка падает на `npm ci` / `npm run build` и выводит только строку вида `A complete log of this run can be found in: /root/.npm/_logs/...`, используйте локальный prebuilt-вариант. Он берёт уже собранный Next.js standalone из `frontend/.next` и не запускает npm:

```bash
docker build -t sda-frontend:local \
  -f data-anonymization-and-synthesis-tool228/dockerfiles/frontend.prebuilt.Dockerfile \
  data-anonymization-and-synthesis-tool228
```

Проверка:

```bash
docker images | grep 'sda-'
```

Если Docker не может скачать базовые образы с Docker Hub, сначала проверьте сеть Docker Desktop:

```bash
docker pull python:3.12-slim
docker pull node:22-slim
```

## Запуск в Kubernetes

```bash
kubectl apply -f lab5/self-work/data-anonymization-k8s/manifests/
kubectl get pods,svc,pvc,ingress,hpa -n sda-lab5
```

Если images были пересобраны после первого запуска:

```bash
kubectl rollout restart deployment/sda-backend deployment/sda-frontend -n sda-lab5
```

## Локальная проверка

Backend:

```bash
kubectl port-forward -n sda-lab5 svc/sda-backend 8000:8000
curl http://localhost:8000/api/v1/health
curl http://localhost:8000/metrics
```

Frontend:

```bash
kubectl port-forward -n sda-lab5 svc/sda-frontend 3000:3000
```

Откройте:

```text
http://localhost:3000
```

Важно: frontend image собран с `NEXT_PUBLIC_API_BASE_URL=http://localhost:8000/api/v1`, поэтому для полной работы UI backend port-forward на `8000` должен быть запущен одновременно.

## Опционально: Ingress

Ingress работает только при установленном Ingress Controller, например `ingress-nginx`.

Локальное имя:

```text
127.0.0.1 sda.local
```

После этого:

```text
http://sda.local/
http://sda.local/api/v1/health
```

## Диагностика

```bash
kubectl describe pod -n sda-lab5
kubectl logs -n sda-lab5 deployment/sda-backend
kubectl logs -n sda-lab5 deployment/sda-frontend
```

Если HPA показывает `cpu: <unknown>`, значит в локальном кластере не установлен `metrics-server`. Это не мешает запуску приложения.
