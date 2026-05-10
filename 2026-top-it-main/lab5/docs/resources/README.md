# Базовые ресурсы Kubernetes (подробно)

Этот раздел **дополняет** краткий обзор в [Основы Kubernetes](../kubernetes-basics.md): здесь те же объекты разобраны глубже, с акцентом на связи между API, типичные ошибки и **официальные референсы** (Kubernetes, kubectl, сетевые плагины).

## Оглавление раздела

| Файл | О чём |
|------|--------|
| [Рабочие нагрузки: Pod, ReplicaSet, Deployment](workloads.md) | Планирование подов, жизненный цикл, стратегии обновления, пробы, ресурсы контейнера |
| [Сервисы и сеть внутри кластера](services-networking.md) | Типы Service, DNS, Endpoints, headless, Ingress в контексте базовой модели |
| [Namespace, метки и селекторы](namespaces-labels.md) | Изоляция имён, labels vs annotations, селекторы в манифестах |
| [Конфигурация, секреты и тома](config-secrets-volumes.md) | ConfigMap, Secret, способы монтирования, обзор типов томов |

После этого логично перейти к практическим YAML в [Примеры ресурсов](../resource-examples.md) и к сценарию деплоя из лаб. №4: [От лабораторной №4 к Kubernetes](../lab4-to-kubernetes.md).

## Справочник API и объектов

- [Kubernetes API Overview](https://kubernetes.io/docs/concepts/overview/kubernetes-api/) (en)
- [Обзор API (рус.)](https://kubernetes.io/ru/docs/concepts/overview/kubernetes-api/) — краткая версия
- [API Reference (все ресурсы)](https://kubernetes.io/docs/reference/kubernetes-api/) (en)

## kubectl: общие материалы

- [Overview kubectl](https://kubernetes.io/docs/reference/kubectl/) (en)
- [Шпаргалка kubectl](https://kubernetes.io/docs/reference/kubectl/cheatsheet/) (en) · [на русском](https://kubernetes.io/ru/docs/reference/kubectl/cheatsheet/)
- [kubectl explain](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_explain/) — разбор полей прямо в терминале, например: `kubectl explain deployment.spec.strategy`

## Сеть (концепции, не привязанные к одному YAML)

- [Концепции: Services, Load Balancing, Networking](https://kubernetes.io/docs/concepts/services-networking/) (en)
- [DNS для сервисов и подов](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/) (en)
- [Сетевые плагины (CNI)](https://kubernetes.io/docs/concepts/cluster-administration/networking/) (en)

## Безопасность и хранение секретов (точки входа)

- [Secrets — good practices](https://kubernetes.io/docs/concepts/configuration/secret/#good-practices) (en)
- [Шифрование данных в etcd (encryption at rest)](https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/) (en)

## Версии и обратная совместимость

- [Версионирование API](https://kubernetes.io/docs/reference/using-api/#api-versioning) (en)
- [Deprecated API Migration Guide](https://kubernetes.io/docs/reference/using-api/deprecation-guide/) (en) — полезно при обновлении кластера или чужих манифестов

В каждом файле ниже в конце есть **тематический блок ссылок** по разделу.
