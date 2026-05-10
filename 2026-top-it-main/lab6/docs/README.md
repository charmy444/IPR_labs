# Документация к лабораторной работе №6

Материалы дополняют основной [`README.md`](../README.md): здесь подробнее разобраны **Kustomize**, **Helm** и **сравнение подходов** для самостоятельного чтения и справки.

## Оглавление

| Документ | Для чего |
|----------|----------|
| [Kustomize](kustomize.md) | Модель base/overlays, `kustomization.yaml`, патчи, генераторы, команды `kubectl`, типичные ошибки |
| [Helm](helm.md) | Чарт, values, шаблоны, релизы, жизненный цикл, `helm template`, практика с секретами |
| [Helm и Kustomize: сравнение](helm-vs-kustomize.md) | Различия моделей, сценарии «что выбрать», совместное использование |

## Эталоны в репозитории

- Приложение: [telegram-support-bot](../../lab5/examples/telegram-support-bot/k8s/README.md) — каталоги `k8s/kustomization` и `k8s/helm/telegram-support-app`.
- Инфраструктура (PostgreSQL): [telegram-support-infra](../examples/telegram-support-infra/README.md) — зеркальная структура `k8s/kustomization` и `k8s/helm/postgres-infra`.

## Рекомендуемый порядок

1. Пройти практику по основному [`README.md`](../README.md) (части A–C).
2. При необходимости углубиться: сначала [Kustomize](kustomize.md), затем [Helm](helm.md).
3. Закрепить выбор инструмента через [сравнение](helm-vs-kustomize.md).
