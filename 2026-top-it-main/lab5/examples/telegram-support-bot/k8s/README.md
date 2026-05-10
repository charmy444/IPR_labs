# Kubernetes: только приложение

Манифесты **не содержат PostgreSQL**. База развёртывается отдельно из примера инфраструктуры курса: [telegram-support-infra](../../../../lab6/examples/telegram-support-infra/README.md).

## Структура

```text
k8s/
  README.md                 # этот файл
  kustomization/
    base/                   # общие манифесты приложения
    overlays/
      dev/                  # namespace telegram-demo, меньше реплик
      prod/                 # усиленные ресурсы, другой URL API при необходимости
  helm/
    telegram-support-app/   # тот же стек через Helm + values-dev / values-prod
  optional/
    istio/                  # Gateway + VirtualService (нужен установленный Istio)
```

## Порядок работы

1. Развернуть БД из `lab6/examples/telegram-support-infra` в тот же namespace, что и overlay приложения (для `dev` это `telegram-demo`), либо скорректируйте хост в `DATABASE_URL` по [контракту](../../../../lab6/examples/telegram-support-infra/README.md).
2. Создать образы backend и frontend и загрузить их в кластер (или registry из лаб. 4).
3. Применить Kustomize или установить Helm chart.

### Kustomize

```bash
kubectl apply -k k8s/kustomization/overlays/dev
```

### Helm

```bash
helm upgrade --install telegram-app ./k8s/helm/telegram-support-app \
  --namespace telegram-demo --create-namespace \
  -f ./k8s/helm/telegram-support-app/values-dev.yaml
```

Подробности и задания: [лабораторная работа №6](../../../../lab6/README.md).
