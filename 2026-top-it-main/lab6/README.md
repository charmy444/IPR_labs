# Лабораторная работа №6: Kustomize и Helm, разделение приложения и инфраструктуры

## Цель работы

Научиться:

- разделять **манифесты приложения** и **инфраструктуру** (в первую очередь базу данных) так, как это делают в промышленной разработке;
- описывать приложение через **Kustomize** (база + overlays) и через **Helm** (chart + несколько values-файлов);
- понимать роль **StatefulSet**, **Headless Service** и **PVC** на примере PostgreSQL в отдельном «репозитории инфраструктуры».

## Требования

- Выполнена [лабораторная работа №5](../lab5/README.md): `kubectl`, кластер Kubernetes (например Docker Desktop).
- Установлен [Helm](https://helm.sh/docs/intro/install/) 3.x.
- Собраны образы примера [telegram-support-bot](../lab5/examples/telegram-support-bot) (или образы вашего проекта из [лаб. №4](../lab4/README.md)).

## Расширенная документация

Подробные справочники по инструментам и сравнению подходов — в папке [docs/](docs/README.md):

- [Kustomize](docs/kustomize.md)
- [Helm](docs/helm.md)
- [Helm и Kustomize: сравнение и выбор](docs/helm-vs-kustomize.md)

## Педагогический принцип: два репозитория

В одном репозитории живёт **код приложения** и манифесты **только** frontend/backend (и их конфигурация). **PostgreSQL** (и аналогичные системы) вносят в **отдельный** репозиторий или каталог, который владеет values, runbook’ом в README и политикой данных.

Почему так:

- разные команды и релизные циклы (приложение vs платформа);
- секреты и бэкапы БД не смешиваются с Dockerfile приложения;
- приложение подключается к БД по **контракту** (DNS-имя, порт, имя БД), описанному в README инфраструктуры.

Локально **docker-compose** примера бота остаётся монолитным для удобства предыдущих лаб; в Kubernetes мы явно учим **другое** разбиение.

## Эталоны в репозитории

| Каталог | Содержимое |
|---------|------------|
| [examples/telegram-support-infra](examples/telegram-support-infra/README.md) | Только PostgreSQL: `k8s/helm/postgres-infra`, `k8s/kustomization` (base + dev/prod), README с контрактом для приложения |
| [lab5/examples/telegram-support-bot/k8s](../lab5/examples/telegram-support-bot/k8s/README.md) | Только приложение: `kustomization/`, `helm/telegram-support-app`, опционально `optional/istio` и `optional/ingress` |

---

## Часть A. Инфраструктура и stateful-контур

### A.1. Теория (кратко)

**StatefulSet** даёт предсказуемые имена подов (`postgres-0`, …), упорядоченное масштабирование и устойчивые сетевые идентификаторы. Для стабильного DNS внутри кластера к нему подключают **Headless Service** (`clusterIP: None`). Данные БД хранят в **PVC**, задаваемых через `volumeClaimTemplates`.

Эти объекты относятся к **инфраструктуре**, а не к чарту приложения.

### A.2. Практика

1. Откройте [examples/telegram-support-infra/README.md](examples/telegram-support-infra/README.md).
2. Убедитесь, что в кластере есть StorageClass: `kubectl get storageclass`.
3. Установите БД **одним** из способов (или обоими по желанию преподавателя):

**Helm**

```bash
cd lab6/examples/telegram-support-infra
helm upgrade --install telegram-support-db ./k8s/helm/postgres-infra \
  --namespace telegram-demo --create-namespace \
  -f ./k8s/helm/postgres-infra/values-dev.yaml
```

**Kustomize**

```bash
cd lab6/examples/telegram-support-infra
kubectl apply -k k8s/kustomization/overlays/dev
```

4. Проверьте Pod и PVC:

```bash
kubectl get pods,pvc -n telegram-demo -l app=postgres
```

5. Зафиксируйте **контракт** для `DATABASE_URL` приложения (хост `postgres-0.postgres...`, порт `5432`, БД `support_bot`, пользователь и пароль из values/Secret инфраструктуры).

---

## Часть B. Kustomize — приложение (первым по порядку лабы)

В каталоге приложения эталонная структура:

```text
k8s/kustomization/
  base/ # Deployment, Service, ConfigMap (без Secret с токеном — см. overlays)
  overlays/
    dev/
    prod/
```

1. Изучите [base](../lab5/examples/telegram-support-bot/k8s/kustomization/base) и overlays `dev` / `prod`.
2. Убедитесь, что в **resources** overlays **нет** манифестов PostgreSQL.
3. В overlay `dev` задан namespace `telegram-demo` и `database-url` в Secret согласован с инфраструктурой в том же namespace.
4. В overlay `prod` приложение в namespace `telegram-prod`, а в Secret указан **FQDN** БД в `telegram-demo` (пример меж-namespace контракта).

Сборка и просмотр:

```bash
cd lab5/examples/telegram-support-bot
kubectl kustomize k8s/kustomization/overlays/dev
kubectl apply -k k8s/kustomization/overlays/dev
```

Перед применением задайте реальный `bot-token` в `overlays/*/secret.yaml` или примените патч; не коммитьте токены.

---

## Часть C. Helm — то же приложение

Эталон: [k8s/helm/telegram-support-app](../lab5/examples/telegram-support-bot/k8s/helm/telegram-support-app).

1. Разберите `values.yaml`, `values-dev.yaml`, `values-prod.yaml`: различаются реплики, ресурсы, способ задания `database.url` (короткий хост vs FQDN).
2. Проверьте шаблоны без установки:

```bash
cd lab5/examples/telegram-support-bot
helm template telegram-app ./k8s/helm/telegram-support-app \
  --namespace telegram-demo \
  -f ./k8s/helm/telegram-support-app/values-dev.yaml
```

3. Установка:

```bash
helm upgrade --install telegram-app ./k8s/helm/telegram-support-app \
  --namespace telegram-demo --create-namespace \
  -f ./k8s/helm/telegram-support-app/values-dev.yaml \
  --set telegram.botToken="YOUR_TOKEN"
```

Команды `helm upgrade`, `helm rollback`, `helm uninstall` отработайте на учебном namespace самостоятельно.

---

## Сравнение Kustomize и Helm

| | Kustomize | Helm |
|---|-----------|------|
| Подход | Патчи, слои, композиция YAML без собственного языка шаблонов | Шаблоны Go + values, релиз как единица установки |
| Удобно когда | Несколько окружений как «надстройки» над одними и теми же манифестами, GitOps без чартов | Параметризация, зависимости чартов, версионирование пакета |
| Минусы | Сложная логика в шаблонах неудобна | Нужно следить за версией chart и совместимостью values |

Оба инструмента часто сосуществуют; в лабе важно **освоить оба** на одном приложении.

---

## Промышленная норма

- Инфраструктура (БД, кластерные CRD, сетевые политики) — отдельный репозиторий или монорепо с жёстким разделением каталогов и Codeowners.
- Пароли продакшена не лежат в Git; используются переменные CI, Sealed Secrets, External Secrets Operator и т.д.
- Приложение получает строку подключения из Secret, собранного из values/CI, согласованного с README/platform.

---

## Задание для самостоятельной работы

1. **Инфраструктура:** по образцу [telegram-support-infra](examples/telegram-support-infra) оформите каталог с PostgreSQL (Helm и/или Kustomize), двумя наборами values/overlays (например dev и prod), README с контрактом для приложения и порядком деплоя.
2. **Приложение:** в своём (или форкнутом) репозитории с сервисами создайте структуру:

```text
k8s/
  kustomization/
    base/
    overlays/
      dev/
      prod/
  helm/
    <имя-чарта>/
      values-dev.yaml
      values-prod.yaml
      templates/
```

В **`k8s/kustomization` и `k8s/helm` не должно быть манифестов базы данных.** Строка `DATABASE_URL` (или эквивалент) задаётся через overlay/values и совпадает с контрактом инфраструктуры.

3. Поднимите сначала инфраструктуру, затем приложение; покажите работающие health-checks backend и (по возможности) доступ к frontend.

---

## Контрольные вопросы

1. Зачем выносить PostgreSQL из репозитория приложения?
2. Что такое контракт между приложением и БД в Kubernetes?
3. Зачем StatefulSet’у Headless Service?
4. Что происходит с PVC при удалении StatefulSet (какие нюансы)?
5. Чем `helm upgrade` принципиально отличается от `kubectl apply -k` с точки зрения модели «релиза»?

---

## Дополнительные материалы

- Материалы курса: [папка docs/](docs/README.md) — Kustomize, Helm, сравнение подходов
- [Документация Kustomize](https://kubectl.docs.kubernetes.io/references/kustomize/) (через kubectl)
- [Документация Helm](https://helm.sh/docs/)
- [StatefulSet](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) (официальная документация Kubernetes)
