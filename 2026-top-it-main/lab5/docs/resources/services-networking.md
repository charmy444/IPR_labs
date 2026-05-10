# Сервисы и сеть внутри кластера

## Service: стабильная точка доступа к подам

Поды **эфемерны**: при пересоздании меняется IP. **Service** даёт **стабильный виртуальный IP** (ClusterIP) и **DNS-имя**, а трафик распределяется на поды, попадающие под `spec.selector` (обычно совпадает с метками из Pod template в Deployment).

Важно: Service **не** «поднимает» контейнеры — он только маршрутизирует к уже существующим подам.

### Типы Service (кратко, но с опорой на модель)

| Тип | Назначение |
|-----|------------|
| `ClusterIP` | Доступ **только изнутри** кластера (значение по умолчанию). |
| `NodePort` | Открывает порт на каждом узле, проксирует в Service — удобно для локальных кластеров (Docker Desktop, minikube). |
| `LoadBalancer` | Запрос внешнего балансировщика у облачного провайдера; локально часто ведёт себя как NodePort + «pending» внешнего IP. |
| `ExternalName` | CNAME на внешний DNS — без прокси-подов; полезно для постепенного переноса зависимостей. |

Подробнее: [Service](https://kubernetes.io/docs/concepts/services-networking/service/) (en).

### Endpoints и EndpointSlice

Для Service контроллер поддерживает объекты **Endpoints** (legacy) и **EndpointSlice** (современная нарезка по группам): это фактический список IP:port подов, готовых принимать трафик. Если у пода падает readiness, его адрес **исчезает** из списка — балансировка перестаёт слать туда запросы.

Просмотр:

```bash
kubectl get endpoints <service-name>
kubectl get endpointslice
```

### Headless Service (`clusterIP: None`)

Используется, когда нужен **DNS по каждому поду** (StatefulSet, peer-to-peer) или когда балансировку делает клиент. Для такого Service DNS обычно возвращает **все** поды, попадающие под селектор.

- [Headless Services](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services) (en)

### Сессионная привязность (sessionAffinity)

По умолчанию трафик распределяется между подами. `sessionAffinity: ClientIP` может «липнуть» к одному поду в пределах таймаута — полезно для некоторых legacy-приложений, но осторожно с обновлениями и балансировкой.

### kube-proxy

Реализация правил маршрутизации для Service на узлах — зона ответственности **kube-proxy** (iptables, ipvs или другие режимы в зависимости от кластера). Прикладному разработчику достаточно понимать модель «Service → поды по селектору»; детали — при отладке сети.

- [Virtual IPs and Service Proxies](https://kubernetes.io/docs/reference/networking/virtual-ips/) (en)

## DNS в кластере

Короткое имя: в том же namespace Service `backend` доступен как `backend`. FQDN: `backend.<namespace>.svc.cluster.local`.

Для подов записи DNS зависят от политики и настроек CoreDNS; базовая схема описана в документации.

- [DNS for Services and Pods](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/) (en)

## Ingress vs Service

- **Service** (ClusterIP/NodePort/LoadBalancer) работает на **L4** (TCP/UDP) в типичной модели «порт сервиса → порты подов».
- **Ingress** — **L7 HTTP(S)** маршрутизация по хосту и пути к Service; нужен установленный **Ingress Controller** (nginx, traefik, …).

Для лаборатории часто хватает NodePort или `kubectl port-forward`; Ingress подключают, когда нужен нормальный HTTP-маршрут с одного IP.

- [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) (en)
- [Ingress Controllers](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/) (en)

Пример nginx Ingress в репозитории: [`examples/telegram-support-bot/k8s/optional/ingress/`](../../examples/telegram-support-bot/k8s/optional/ingress/).

## Референсы

- [Service](https://kubernetes.io/docs/concepts/services-networking/service/) (en)
- [Connecting Applications with Services](https://kubernetes.io/docs/tutorials/services/connect-applications-service/) (туториал, en)
- [EndpointSlice](https://kubernetes.io/docs/concepts/services-networking/endpoint-slices/) (en)
- [Публикация сервиса (типы)](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types) (en)
- [Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/) (en) — следующий уровень после базовых Service; примеры в [Примеры ресурсов](../resource-examples.md)

Наверх: [оглавление раздела «Базовые ресурсы»](README.md).
