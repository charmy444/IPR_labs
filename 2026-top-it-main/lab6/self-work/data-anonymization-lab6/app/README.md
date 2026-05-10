# SDA application contour

Здесь находятся только манифесты приложения:

- FastAPI backend;
- Next.js frontend;
- ConfigMap/Secret приложения;
- Services, Ingress, HPA.

PVC создаются в `../infra`. В app-манифестах есть только ссылки на существующие PVC по контракту:

- `sda-upload-store`;
- `sda-analysis-store`.

