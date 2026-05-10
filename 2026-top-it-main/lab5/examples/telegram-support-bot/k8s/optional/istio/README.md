# Istio (опционально)

Манифесты рассчитаны на кластер с установленным [Istio](https://istio.io/). Укажите namespace при применении (тот же, что у приложения):

```bash
kubectl apply -f gateway.yaml -n telegram-demo
kubectl apply -f virtualservice.yaml -n telegram-demo
```

Сервисы `telegram-frontend` и `telegram-backend` должны уже существовать.
