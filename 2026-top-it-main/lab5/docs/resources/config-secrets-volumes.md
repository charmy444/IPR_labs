# Конфигурация, секреты и тома

Цель — **не запекать** конфигурацию и секреты в образ контейнера: одни и те же образы должны работать в dev/stage/prod с разными параметрами.

## ConfigMap

**ConfigMap** хранит неконфиденциальные данные: ключи-значения, целые файлы конфигурации. Данные можно:

- пробросить в контейнер как **переменные окружения** (`env`, `envFrom`);
- смонтировать как **файлы** через volume (`configMap` volume);
- использовать в аргументах командной строки через `envsubst`/шаблоны (паттерн приложения).

### Важные детали

- При монтировании как volume изменения ConfigMap **могут** подтягиваться в под без пересоздания (с задержкой); для критичных приложений проверяйте поведение рантайма.
- С ConfigMap можно связать **`immutable: true`**, чтобы случайно не сломать работающие поды неожиданным обновлением (Kubernetes 1.21+).

- [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/) (en)
- [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/) (туториал, en)

## Secret

**Secret** предназначен для чувствительных данных: пароли, токены, TLS-сертификаты. В манифестах поле `data` — **base64**; удобнее для ручного редактирования **`stringData`** (Kubernetes сам закодирует при сохранении).

### Безопасность: что Secret не делает сам по себе

- Base64 — **не шифрование**; любой с доступом к API может прочитать Secret.
- Защита на уровне кластера: RBAC, [encryption at rest](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/) в etcd, политики организации.
- **Не коммить** реальные секреты в Git; для лабораторий используйте плейсхолдеры, SOPS, Sealed Secrets, External Secrets или секреты в CI/CD.

Типы: `Opaque` (общий), `kubernetes.io/dockerconfigjson` для `imagePullSecrets`, `kubernetes.io/tls` для Ingress и т.д.

- [Secret](https://kubernetes.io/docs/concepts/configuration/secret/) (en)
- [Pull an Image from a Private Registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) (en)

## Тома (Volumes): обзор для базового уровня

Том подключается к поду и виден указанным контейнерам через `volumeMounts`.

| Тип | Типичное использование |
|-----|-------------------------|
| `emptyDir` | Временные файлы, кэш; данные живут, пока жив под |
| `configMap` / `secret` | Файлы конфигурации и секретов в файловой системе |
| `persistentVolumeClaim` | Долговечное хранилище вне жизненного цикла пода |
| `downwardAPI` | Проброс метаданных пода/узла в файлы или env |

Полный список типов томов — в справочнике; для лаб. №5 достаточно понимать связку «манифест тома → mountPath в контейнере».

- [Volumes](https://kubernetes.io/docs/concepts/storage/volumes/) (en)
- [Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) (en)
- [Storage Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/) (en)

## Примеры YAML и продвинутые сценарии

Готовые фрагменты манифестов (env, volumeMount, PV/PVC) — в [Примеры ресурсов](../resource-examples.md). Там же — Job, CronJob, HPA и сетевые политики.

## Референсы

- [Configuration](https://kubernetes.io/docs/concepts/configuration/) (en) — раздел документации целиком
- [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/) · [Secrets в поде](https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets-as-files-from-a-pod) (en)
- [Projected volumes](https://kubernetes.io/docs/concepts/storage/volumes/#projected) (en) — несколько источников в одном томе
- [SubPath](https://kubernetes.io/docs/concepts/storage/volumes/#using-subpath) (en) — монтирование одного файла из тома

Наверх: [оглавление раздела «Базовые ресурсы»](README.md).
