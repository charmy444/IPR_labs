# Самостоятельная работа lab6: SDA, Kustomize и Helm

Это отдельная сдаваемая папка для lab6 на базе `data-anonymization-and-synthesis-tool228`.

У проекта нет PostgreSQL: backend хранит временные CSV upload/session данные в файловой системе. Поэтому инфраструктурный контур в этой работе — Kubernetes-хранилище (`PersistentVolumeClaim`) для:

- `/tmp/sda_upload_store`;
- `/tmp/sda_similar_analysis_store`.

Приложение и инфраструктура разделены:

```text
infra/   # PVC-хранилище, dev/prod overlays, Helm chart
app/     # backend/frontend, dev/prod overlays, Helm chart
```

## Сборка локальных образов

Из корня `2026-top-it-main`:

```bash
docker build -t sda-backend:local \
  -f data-anonymization-and-synthesis-tool228/dockerfiles/backend.Dockerfile \
  data-anonymization-and-synthesis-tool228

docker build -t sda-frontend:local \
  --build-arg NEXT_PUBLIC_API_BASE_URL=http://localhost:8000/api/v1 \
  -f data-anonymization-and-synthesis-tool228/dockerfiles/frontend.Dockerfile \
  data-anonymization-and-synthesis-tool228
```

Если npm внутри Docker падает, соберите frontend из уже готового Next.js standalone:

```bash
docker build -t sda-frontend:local \
  -f data-anonymization-and-synthesis-tool228/dockerfiles/frontend.prebuilt.Dockerfile \
  data-anonymization-and-synthesis-tool228
```

## Вариант A: Kustomize

Сначала инфраструктура:

```bash
kubectl apply -k lab6/self-work/data-anonymization-lab6/infra/k8s/kustomization/overlays/dev
kubectl get pvc -n sda-dev
```

Потом приложение:

```bash
kubectl apply -k lab6/self-work/data-anonymization-lab6/app/k8s/kustomization/overlays/dev
kubectl get pods,svc -n sda-dev
```

Проверка:

```bash
kubectl port-forward -n sda-dev svc/sda-backend 8000:8000
curl http://localhost:8000/api/v1/health
curl http://localhost:8000/metrics
```

В другом терминале для frontend:

```bash
kubectl port-forward -n sda-dev svc/sda-frontend 3000:3000
```

Открыть `http://localhost:3000`.

## Вариант B: Helm

Сначала инфраструктура:

```bash
helm upgrade --install sda-storage \
  lab6/self-work/data-anonymization-lab6/infra/k8s/helm/sda-storage-infra \
  --namespace sda-dev --create-namespace \
  -f lab6/self-work/data-anonymization-lab6/infra/k8s/helm/sda-storage-infra/values-dev.yaml
```

Потом приложение:

```bash
helm upgrade --install sda-app \
  lab6/self-work/data-anonymization-lab6/app/k8s/helm/sda-app \
  --namespace sda-dev --create-namespace \
  -f lab6/self-work/data-anonymization-lab6/app/k8s/helm/sda-app/values-dev.yaml
```

Проверка такая же:

```bash
kubectl get pods,svc,pvc -n sda-dev
kubectl port-forward -n sda-dev svc/sda-backend 8000:8000
curl http://localhost:8000/api/v1/health
```

## Контракт между app и infra

Приложение ожидает, что в namespace уже есть PVC:

- `sda-upload-store`;
- `sda-analysis-store`.

Эти PVC создаются только в `infra/`. В `app/` нет манифестов создания хранилища, только ссылки на уже существующие claims.
