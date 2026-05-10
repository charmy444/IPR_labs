# Инфраструктура PostgreSQL для Telegram Support Bot

Отдельный «репозиторий инфраструктуры» в учебных целях: **здесь только база данных** и сопутствующие манифесты. Репозиторий приложения ([telegram-support-bot](../../../lab5/examples/telegram-support-bot)) не должен дублировать StatefulSet PostgreSQL — так проще сопровождать релизы, секреты и бэкапы.

## Зачем отдельный каталог

- Команда платформы или SRE владеет жизненным циклом БД (обновления, PVC, бэкапы).
- Разработчики приложения подключаются по **контракту** (имя сервиса, порт, БД, способ получения пароля).
- В Git не хранятся реальные пароли продакшена; для CI используются секреты GitLab / External Secrets / Sealed Secrets.

## Предпосылки

- Кластер Kubernetes с рабочим **StorageClass** (Docker Desktop обычно создаёт `hostpath` или default SC).
- Установлены `kubectl` и (для варианта Helm) `helm`.

Проверка:

```bash
kubectl get storageclass
```

## Контракт для приложения (telegram-support-bot)

После установки этого чарта/манифестов в namespace `N`:

| Параметр | Значение |
|----------|----------|
| Хост (внутри кластера) | `postgres-0.postgres.<NAMESPACE>.svc.cluster.local` |
| Порт | `5432` |
| Пользователь БД | `postgres` (переопределяется в values — см. ниже) |
| Имя БД | `support_bot` |
| Пароль | Как в `values-dev.yaml` / `values-prod.yaml` (Helm) или в `Secret` (Kustomize) — **должен совпадать** с тем, что вы передаёте в `DATABASE_URL` приложения |

Пример строки подключения (подставьте свой namespace и пароль):

```text
postgres://postgres:PASSWORD@postgres-0.postgres.<NAMESPACE>.svc.cluster.local:5432/support_bot?sslmode=disable
```

Если приложение развёрнуто **в том же namespace**, часто достаточно короткого имени:

```text
postgres://postgres:PASSWORD@postgres-0.postgres:5432/support_bot?sslmode=disable
```

## Порядок деплоя (лабораторная работа №6)

1. Развернуть **инфраструктуру** (этот каталог) в выбранный namespace.
2. Дождаться готовности Pod `postgres-0`:

```bash
kubectl get pods -n <NAMESPACE> -l app=postgres
kubectl logs -n <NAMESPACE> postgres-0 -c postgres --tail=20
```

3. Развернуть **приложение** из репозитория `telegram-support-bot`, передав в Secret приложения `database-url`, согласованный с таблицей выше.

## Вариант A: Helm

```bash
cd lab6/examples/telegram-support-infra

# Dev (маленький том, слабые лимиты)
helm upgrade --install telegram-support-db ./k8s/helm/postgres-infra \
  --namespace telegram-demo --create-namespace \
  -f ./k8s/helm/postgres-infra/values-dev.yaml

# Prod-подобные значения (увеличьте пароль через --set или внешний секрет)
helm upgrade --install telegram-support-db ./k8s/helm/postgres-infra \
  --namespace telegram-demo \
  -f ./k8s/helm/postgres-infra/values-prod.yaml
```

Проверка:

```bash
kubectl get statefulset,pvc,pod -n telegram-demo -l app.kubernetes.io/name=postgres
```

Удаление:

```bash
helm uninstall telegram-support-db -n telegram-demo
```

**Важно:** при удалении StatefulSet PVC по умолчанию могут остаться — удалите их вручную, если нужен «чистый» стенд.

## Вариант B: Kustomize

```bash
cd lab6/examples/telegram-support-infra

kubectl apply -k k8s/kustomization/overlays/dev # namespace: telegram-demo
# или
kubectl apply -k k8s/kustomization/overlays/prod # namespace: telegram-prod (задайте пароль в secret до применения)
```

Просмотр сгенерированных манифестов без применения:

```bash
kubectl kustomize k8s/kustomization/overlays/dev
```

## StatefulSet и Headless Service

- **Headless Service** (`clusterIP: None`) нужен для стабильных сетевых идентификаторов Pod’ов StatefulSet.
- Учебный стенд — **одна реплика** PostgreSQL (не HA-кластер).

## Безопасность

- Не коммитьте файлы с реальными паролями продакшена.
- Для учёбы допустимы явные пароли в `values-dev.yaml`; для `values-prod.yaml` используйте плейсхолдер и подстановку в CI или внешний секрет-менеджер.

## Связанные материалы

- Приложение: [telegram-support-bot](../../../lab5/examples/telegram-support-bot)
- Лабораторная работа №6: [lab6/README.md](../../README.md)
