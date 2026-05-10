# Helm: подробный разбор

[Helm](https://helm.sh/) — **менеджер пакетов** для Kubernetes. Единица распространения — **chart**: каталог с шаблонами манифестов, значениями по умолчанию и метаданными. Установленный chart в кластере — это **release** (имя + namespace + версия + история обновлений).

Helm 3 не использует серверный компонент Tiller: шаблоны рендерятся на машине, куда установлен CLI, а в API-сервер уходят обычные объекты Kubernetes.

---

## Структура чарта

Минимальный набор (как в курсе):

```text
telegram-support-app/
  Chart.yaml          # имя, версия чарта, версия приложения
  values.yaml         # значения по умолчанию
  values-dev.yaml     # переопределения для dev (пример)
  values-prod.yaml    # переопределения для prod
  templates/          # шаблоны *.yaml (и вспомогательные *.tpl)
```

### `Chart.yaml`

Содержит `apiVersion: v2`, `name`, `version` (версия пакета), `appVersion` (ориентировочная версия приложения), при необходимости `dependencies` для subchart’ов.

### `values.yaml`

Иерархическая структура (образы, реплики, ресурсы, произвольные ключи для ваших шаблонов). В лаборатории те же оси, что и в Kustomize-overlays: **разные файлы values** под окружения вместо дублирования шаблонов.

При установке файлы мержатся: сначала `values.yaml`, затем `-f values-dev.yaml` и ключи `--set` / `--set-file`.

### `templates/`

Файлы с расширением `.yaml` (и `.tpl`) — это **шаблоны Go templates**. На выходе после рендеринга должен получиться валидный YAML (или несколько документов `---` в одном файле).

В шаблонах доступны:

- `.Values` — итоговые values;
- `.Release` — имя релиза, namespace и т.д.;
- `.Chart` — метаданные из `Chart.yaml`;
- функции Sprig и встроенные функции Helm (`include`, `required`, `toYaml`, `indent`).

Типичные приёмы:

- `{{ .Values.frontend.replicas }}` в `spec.replicas`;
- `{{- toYaml .Values.frontend.resources | nindent 12 }}` для вложенного блока `resources`;
- `metadata.namespace: {{ .Release.Namespace }}` для явной привязки к namespace релиза.

Файл `_helpers.tpl` (опционально) хранит именованные фрагменты `define` / `template` для меток и имён, чтобы не дублировать разметку.

---

## Жизненный цикл релиза

| Команда | Назначение |
|---------|------------|
| `helm install <release> ./chart` | Первичная установка |
| `helm upgrade --install <release> ./chart` | Обновление или установка, если релиза ещё нет |
| `helm uninstall <release>` | Удаление ресурсов релиза (аннотации Helm отслеживают принадлежность) |
| `helm rollback <release> <revision>` | Откат к предыдущей ревизии |
| `helm history <release>` | История версий values/манифестов |
| `helm template <release> ./chart` | Рендер в stdout **без** обращения к кластеру — отладка и CI |
| `helm lint ./chart` | Проверка чарта |

Ревизии хранятся в кластере (Secret’ы с историей релиза), поэтому откат возможен даже если локальные файлы уже изменились — ориентируйтесь на документацию по ограничениям и очистке.

---

## Несколько окружений

Практики:

1. **Несколько values-файлов:**  
   `helm upgrade --install app ./chart -f values.yaml -f values-prod.yaml`
2. **Разные release-имена** для одного чарта: `myapp-dev`, `myapp-prod` в разных namespace.
3. **`--set` и `--set-file`** для секретов из CI (не коммитить в репозиторий).

В лаборатории эталон: `values-dev.yaml` и `values-prod.yaml` для [telegram-support-app](../../lab5/examples/telegram-support-bot/k8s/helm/telegram-support-app).

---

## Секреты и чувствительные данные

- Не кладите пароли и токены в Git в открытом виде.
- Для учебы допустимы плейсхолдеры в `values-dev.yaml`; для prod — подстановка через `helm install ... --set telegram.botToken="$TOKEN"` или внешний секрет-менеджер.
- `helm secrets` (плагины) и интеграции с SOPS — вне базовой программы курса, но полезны в работе.

---

## Зависимости между чартами (кратко)

В `Chart.yaml` можно объявить `dependencies` на другие чарты (например community PostgreSQL). Приложение в лаб. 6 **намеренно** не тянет БД как subchart: инфраструктура вынесена в отдельный пример. В реальных проектах subchart — компромисс между «всё в одном релизе» и «раздельные репозитории».

---

## Отладка

1. `helm template ...` и сравнение с ожидаемым YAML.
2. `helm get manifest <release>` — что реально задеплоено.
3. Уменьшить чарт до минимального Deployment и нарастить сложность.

---

## Связь с лабораторной №6

- Чарт приложения: [telegram-support-app](../../lab5/examples/telegram-support-bot/k8s/helm/telegram-support-app).
- Чарт инфраструктуры PostgreSQL: [postgres-infra](../examples/telegram-support-infra/k8s/helm/postgres-infra).

---

## Внешние материалы

- [Документация Helm](https://helm.sh/docs/)
- [Шаблоны и функции](https://helm.sh/docs/chart_template_guide/)
- [Best practices](https://helm.sh/docs/chart_best_practices/)
