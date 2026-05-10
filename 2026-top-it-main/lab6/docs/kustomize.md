# Kustomize: подробный разбор

[Kustomize](https://kustomize.io/) — инструмент для **настройки Kubernetes-манифестов без встраивания шаблонного языка в YAML**. Вы остаётесь на «чистом» API Kubernetes: те же `Deployment`, `Service`, `ConfigMap`, что и у `kubectl apply -f`, но с возможностью **накладывать слои** (overlays) и **собирать** итоговый набор файлов перед применением.

С версии Kubernetes 1.14 Kustomize встроен в `kubectl` (`kubectl kustomize`, `kubectl apply -k`). Отдельная установка бинарника `kustomize` нужна только если требуется версия новее, чем в вашем `kubectl`.

---

## Зачем это нужно

| Задача | Как Kustomize помогает |
|--------|-------------------------|
| Одинаковая «база» для dev/stage/prod | Каталог **base** + несколько **overlays**, каждый добавляет только отличия |
| Не плодить копии целых YAML | Патчи к тем же ресурсам (replicas, лимиты, образ) |
| Разные namespace и метки | Поля `namespace`, `commonLabels` / `labels` в `kustomization.yaml` |
| Подставить конфиг без ручного дублирования | `configMapGenerator`, `secretGenerator` (осторожно с секретами в Git) |

Kustomize **не** заменяет Helm полностью: у него нет понятия «релиза» как у Helm и нет зависимостей между чартами в том же виде. Зато манифесты остаются валидными YAML для API-сервера и их проще ревьюить в MR.

---

## Структура проекта: base и overlays

Типичная раскладка (как в курсе):

```text
k8s/kustomization/
  base/
    kustomization.yaml
    deployment.yaml
    service.yaml
    ...
  overlays/
    dev/
      kustomization.yaml
      patch-replicas.yaml
      secret.yaml
    prod/
      kustomization.yaml
      patch-resources.yaml
      secret.yaml
```

- **base** — общие манифесты без привязки к конкретному окружению (или с нейтральными значениями).
- **overlays/dev**, **overlays/prod** — только то, чем окружения **отличаются**: namespace, число реплик, размер ресурсов, отдельные `Secret`, иногда дополнительные ресурсы.

Overlay подключает base так:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: telegram-demo

resources:
  - ../../base
  - secret.yaml

patches:
  - path: patch-replicas.yaml
```

---

## Файл `kustomization.yaml` (основные поля)

### `resources`

Список файлов или каталогов с манифестами (и других `kustomization`, если подключаете чужой base).

### `namespace`

Добавляет **один и тот же** `metadata.namespace` ко всем ресурсам, где это уместно (не ко всем типам — исключения есть в документации).

### `commonLabels` и `labels`

Метки, которые нужно повесить на **все** (или почти все) ресурсы. В новых версиях Kustomize поле `commonLabels` помечено как устаревающее в пользу блока `labels` с опциями (в т.ч. `includeSelectors: false`, чтобы не менять селекторы `Service`/`Deployment` без необходимости).

Практическое правило: если после сборки **селекторы** перестали совпадать с **лейблами подов**, проверьте, не добавились ли лишние метки в `spec.selector` у `Service` или `matchLabels` у `Deployment`.

### `patches`

Список патчей к уже объявленным ресурсам:

- **стратегический merge** (фрагмент YAML с тем же `kind`, `metadata.name`, изменёнными полями);
- **JSON6902** для точечных операций;
- **патчи через target** (`target.kind`, `target.name`).

Пример стратегического патча — изменить только `spec.replicas` у `Deployment` с именем `telegram-backend`, не копируя весь файл.

### `images`

Массовая замена образа (`name`, `newName`, `newTag`) без правки всех YAML вручную — удобно для CI, который подставляет тег сборки.

### `configMapGenerator` и `secretGenerator`

Генерируют `ConfigMap`/`Secret` из литералов или файлов. Имена получаются с хэшем содержимого (если не отключить), что помогает при перезапуске подов при смене конфига.

**Секреты:** не храните реальные значения в Git. Для учебы допустимы плейсхолдеры; в работе — внешние системы (Sealed Secrets, External Secrets, переменные CI).

### `replacements` (Kustomize 4+)

Копирование значений из одного поля в другое (продвинутый сценарий; в базовой лабе можно не использовать).

---

## Команды

Просмотр результата **без** применения в кластер:

```bash
kubectl kustomize k8s/kustomization/overlays/dev
```

Применение:

```bash
kubectl apply -k k8s/kustomization/overlays/dev
```

Проверяйте вывод `kubectl kustomize` перед `apply`: так проще ловить ошибки композиции и лишние изменения.

---

## Типичные ошибки

1. **Имя ресурса в патче не совпадает** с именем в base — патч молча не применится или Kustomize выдаст ошибку.
2. **Селекторы и лейблы** после `commonLabels` не совпадают — сервис не находит поды.
3. **Дублирование Secret** в base и overlay — нужно явно решить, где живёт единственный источник правды.
4. **Порядок ресурсов** иногда важен для CRD; для обычных Deployment/Service обычно достаточно одного `apply -k`.

---

## Связь с лабораторной №6

- Инфраструктура PostgreSQL: [telegram-support-infra](../examples/telegram-support-infra/k8s/kustomization).
- Приложение: [telegram-support-bot](../../lab5/examples/telegram-support-bot/k8s/kustomization).

---

## Внешние материалы

- [Официальная документация Kustomize](https://kubectl.docs.kubernetes.io/references/kustomize/) (через kubectl)
- [Kustomize на kustomize.io](https://kustomize.io/)
