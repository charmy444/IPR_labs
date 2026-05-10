# Namespace, метки и селекторы

## Namespace

**Namespace** — область имён для большинства ресурсов Kubernetes. Объекты с одинаковыми `kind` и `metadata.name` могут сосуществовать в **разных** namespace.

### Зачем это нужно

- Логическое разделение команд, окружений (dev/stage) или приложений.
- Ограничения по ресурсам через **ResourceQuota** и политики (RBAC привязывают права к namespace).
- Уменьшение риска случайного `kubectl delete` «чужого» деплоймента с тем же именем.

Системные компоненты обычно в `kube-system`; служебные объекты могут быть в `kube-public`, `kube-node-lease` и др.

### Практика

```bash
kubectl create namespace lab5
kubectl config set-context --current --namespace=lab5
kubectl get pods -A   # все namespace
```

По умолчанию контекст часто указывает `default`; явно задавайте `-n` или контекст, чтобы не применить манифест «не туда».

- [Namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) (en)
- [Имена и идентификаторы объектов](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/) (en)

## Labels

**Labels** — пары `ключ: значение` на метаданных объекта. Ими пользуются **селекторы** в Service, Deployment, NetworkPolicy, PDB и др.

### Рекомендуемые метки

Проект **Kubernetes** рекомендует стандартные ключи `app.kubernetes.io/name`, `app.kubernetes.io/instance`, `app.kubernetes.io/version` и др. — это упрощает интеграцию с Helm и общие дашборды.

- [Recommended Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/) (en)

### Селекторы

- **Equality-based**: `app=nginx`, `tier!=cache`.
- **Set-based**: `app in (web, api)`, `!canary`.

В манифестах Deployment/RS обычно `matchLabels`; в NetworkPolicy и части CLI — более богатый синтаксис.

- [Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) (en)

### Аннотации (annotations)

**Annotations** — произвольные метаданные **без** участия в селекторах: версия git, время деплоя, идентификаторы интеграций, подсказки Ingress-контроллеру. Не путать с labels: аннотации не должны решать, куда идёт трафик.

## Связь с учебным проектом

В Helm/Kustomize из `examples/telegram-support-bot` метки связывают Deployment → Service → (при необходимости) ServiceMonitor. Если после правки манифеста сервис «не видит» поды, первым делом проверьте **совпадение labels** между `spec.selector` Service и `template.metadata.labels` в Deployment.

## Референсы

- [Namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) (en)
- [Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) (en)
- [Annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) (en)
- [Field Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/field-selectors/) (en) — отбор по полям, например `status.phase=Running`
- [Well-Known Labels, Annotations and Taints](https://kubernetes.io/docs/reference/labels-annotations-taints/) (en)

Наверх: [оглавление раздела «Базовые ресурсы»](README.md).
