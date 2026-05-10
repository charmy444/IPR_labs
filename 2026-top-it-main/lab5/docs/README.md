# Документация к лабораторной работе №5

Материалы в этой папке дополняют основной [`README.md`](../README.md) курса. Их удобно читать по порядку, если вы впервые работаете с Kubernetes.

## Оглавление

| Документ | Для чего |
|----------|----------|
| [Основы Kubernetes](kubernetes-basics.md) | Модель платформы, объекты API, сеть внутри кластера, пробы, типичный цикл «образ → Deployment → Service» |
| [Базовые ресурсы (подробно)](resources/README.md) | Углублённый разбор Pod, Deployment, Service, DNS, namespace, меток, ConfigMap/Secret и томов — **несколько тематических файлов** и расширенный список ссылок на документацию Kubernetes |
| [От лабораторной №4 к Kubernetes](lab4-to-kubernetes.md) | Как перенести приложение с CI/CD и Docker-образами из лаб. №4 в кластер: registry, секреты, манифесты, проверка |
| [Примеры ресурсов](resource-examples.md) | ConfigMap, Secret, тома, Job, CronJob, HPA, сетевые политики (справочник с примерами YAML и `kubectl`) |

## Рекомендуемый порядок работы

1. Прочитать [Основы Kubernetes](kubernetes-basics.md) и выполнить практику из основного `README.md` (пример frontend/backend), чтобы освоить `kubectl` и базовые объекты.
2. При необходимости углубиться в теорию по разделу [Базовые ресурсы (подробно)](resources/README.md) (workloads, сеть, метки, конфигурация).
3. Изучить [От лабораторной №4 к Kubernetes](lab4-to-kubernetes.md) и задеплоить **своё** приложение из лаб. №4 (или расширенный вариант из `examples/telegram-support-bot`).
4. Для готовых фрагментов YAML и смежных объектов (Job, HPA, NetworkPolicy) — [Примеры ресурсов](resource-examples.md).
5. Продолжить с [лабораторной №6](../../lab6/README.md): Kustomize, Helm и отдельный пример инфраструктуры PostgreSQL — [telegram-support-infra](../../lab6/examples/telegram-support-infra/README.md).

## Пример полного приложения в репозитории

В каталоге [`examples/telegram-support-bot`](../examples/telegram-support-bot) лежит многосервисное приложение (Go + Next.js + PostgreSQL) с `docker-compose` и каталогом `k8s/` (**только приложение**: Kustomize и Helm; БД — в [telegram-support-infra](../../lab6/examples/telegram-support-infra/README.md)). Используйте его как эталон после учебного примера из основного README.

## Внешние материалы для самостоятельного изучения

Ниже — проверенные точки входа: официальная документация, интерактив и справочники. Часть ссылок дана на русском (`/ru/`), часть — на английском (там обычно быстрее обновляют текст).

### Официальная документация Kubernetes

- [Главная страница документации](https://kubernetes.io/docs/home/) (en)
- [Документация на русском](https://kubernetes.io/ru/docs/) — не все разделы переведены; при расхождении смотрите en-версию
- [Обзор концепций](https://kubernetes.io/docs/concepts/) (en)
- [Туториалы](https://kubernetes.io/docs/tutorials/) (en)
- [Интерактив: основы за 6 модулей](https://kubernetes.io/docs/tutorials/kubernetes-basics/) (в браузере)
- [Справочник kubectl](https://kubernetes.io/docs/reference/kubectl/) (en)
- [Шпаргалка kubectl](https://kubernetes.io/docs/reference/kubectl/quick-reference/) (en) · [на русском](https://kubernetes.io/ru/docs/reference/kubectl/cheatsheet/)

### Кластер локально и инструменты

- [Kubernetes в Docker Desktop](https://docs.docker.com/desktop/kubernetes/) (en)
- [Установка kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) (en) · [вариант на русском](https://kubernetes.io/ru/docs/tasks/tools/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) — альтернатива одноузловому кластеру (en)
- [kind (Kubernetes in Docker)](https://kind.sigs.k8s.io/) — кластер в контейнерах, удобно для экспериментов (en)

### Практика в браузере (без своего кластера)

- [Killercoda: сценарии по Kubernetes](https://killercoda.com/kubernetes) — пошаговые упражнения в терминале в браузере
- [Play with Kubernetes](https://labs.play-with-k8s.com/) — временный кластер в браузере (ограничения по времени)

### Связка с контейнерами и образами

- [Образы контейнеров](https://kubernetes.io/docs/concepts/containers/images/) — теги, `imagePullPolicy`, приватные registry (en)
- [Pull образа из приватного registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) (en)

### GitLab (лаб. №4 → образы в k8s)

- [GitLab Container Registry](https://docs.gitlab.com/ee/user/packages/container_registry/) (en)
- [Аутентификация в registry из CI](https://docs.gitlab.com/ee/user/packages/container_registry/authenticate_with_container_registry.html) (en)
- [Переменные окружения предопределённые для CI/CD](https://docs.gitlab.com/ee/ci/variables/predefined_variables.html) — `CI_REGISTRY`, `CI_REGISTRY_IMAGE` и др. (en)

### Книги и углублённый материал (по желанию)

- [Kubernetes Patterns (Red Hat Developer, e-book)](https://developers.redhat.com/e-books/kubernetes-patterns) — паттерны проектирования под k8s; [сайт проекта](https://kubernetespatterns.io/)
- [CNCF: сертификация и курсы](https://www.cncf.io/training/certification/) — CKA/CKAD и др., если захотите системно закрепить знания

Тематические ссылки по объектам API для **базовых** ресурсов — в каталоге [resources/](resources/README.md); для ConfigMap, Secret, Job, HPA и др. с YAML — в конце [Примеры ресурсов](resource-examples.md).
