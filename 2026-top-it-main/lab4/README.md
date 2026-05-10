# Лабораторная работа №4: Непрерывная интеграция и доставка (CI/CD) в GitLab

## Цель работы

Освоить основы непрерывной интеграции и доставки (CI/CD) с использованием GitLab CI/CD. Научиться создавать пайплайны для автоматической сборки, тестирования и публикации приложений.

## Требования

- Аккаунт на [gitlab.mai.ru](https://gitlab.mai.ru)
- Базовые знания Git
- Успешно выполненные лабораторные работы №1-3

---

## Содержание

- [Часть 1: Введение в GitLab CI/CD](#часть-1-введение-в-gitlab-cicd)
- [Часть 2: Основы конфигурации .gitlab-ci.yml](#часть-2-основы-конфигурации-gitlab-ciyml)
- [Часть 2.1: Настройка GitLab Runner](#часть-21-настройка-gitlab-runner)
- [Часть 3: Создание пайплайна с unit-тестами](#часть-3-создание-пайплайна-с-unit-тестами)
- [Часть 4: Сборка и публикация Docker-образов](#часть-4-сборка-и-публикация-docker-образов)
- [Часть 5: Работа с артефактами и кэшем](#часть-5-работа-с-артефактами-и-кэшем)
- [Часть 5.1: Оптимизация пайплайнов и работа с зависимостями](#часть-51-оптимизация-пайплайнов-и-работа-с-зависимостями)
- [Часть 6: Переменные и секреты в CI/CD](#часть-6-переменные-и-секреты-в-cicd)
- [Часть 7: Триггеры и расписание пайплайнов](#часть-7-триггеры-и-расписание-пайплайнов)
- [Часть 8: Продвинутые техники CI/CD](#часть-8-продвинутые-техники-cicd)
- [Самостоятельное задание](#самостоятельное-задание)
- [Контрольные вопросы](#контрольные-вопросы)

---

## Часть 1: Введение в GitLab CI/CD

### Что такое CI/CD?

**Непрерывная интеграция (Continuous Integration, CI)** - это практика разработки программного обеспечения, при которой все изменения кода регулярно объединяются в основном ветке проекта. Это позволяет выявлять ошибки на ранних стадиях разработки.

**Непрерывная доставка (Continuous Delivery, CD)** - это подход к разработке программного обеспечения, при котором команды создают программное обеспечение в коротких циклах, обеспечивая уверенность в том, что оно может быть выпущено в любое время.

### Преимущества CI/CD

- ✅ Быстрое обнаружение ошибок
- ✅ Автоматизация рутинных задач
- ✅ Повышение качества кода
- ✅ Ускорение выхода на рынок
- ✅ Снижение рисков при развертывании

### Архитектура GitLab CI/CD

GitLab CI/CD состоит из следующих компонентов:

1. **Runner** - агент, который выполняет задания пайплайна
   - **Shared runners** - общие runners, доступные всем проектам
   - **Group runners** - runners для конкретной группы проектов
   - **Specific runners** - runners, привязанные к конкретному проекту
   - **Tags** - метки для назначения определенных job на конкретные runners

2. **Pipeline** - полный цикл CI/CD процесса
   - **Parent pipeline** - основной пайплайн
   - **Child pipeline** - дочерний пайплайн, вызываемый из родительского
   - **Multi-project pipeline** - пайплайн, охватывающий несколько проектов

3. **Stage** - этап пайплайна (build, test, deploy)
   - Выполняются последовательно
   - Все job в рамках stage выполняются параллельно

4. **Job** - конкретная задача в рамках stage
   - **Needs** - для создания зависимостей между job
   - **Dependencies** - для передачи артефактов между job

5. **Artifact** - файлы, созданные в процессе выполнения job
   - **Reports** - специальные артефакты для тестов, покрытия кода и т.д.
   - **Expire in** - время хранения артефактов

---

## Часть 2: Основы конфигурации .gitlab-ci.yml

### Структура .gitlab-ci.yml

Файл `.gitlab-ci.yml` описывает конфигурацию пайплайна. Основные элементы:

```yaml
# Определение этапов выполнения
stages:
  - build
  - test

# Переменные окружения
variables:
  DOCKER_DRIVER: overlay2
  IMAGE_NAME: $CI_REGISTRY_IMAGE:$CI_COMMIT_REF_NAME

# Задание (job)
build-job:
  stage: build
  script:
    - echo "Building the application..."
    - npm install
    - npm run build
  artifacts:
    paths:
      - dist/

# Задание с условиями выполнения
test-job:
  stage: test
  script:
    - echo "Running tests..."
    - npm run test
  only:
    - branches
  except:
    - master
```

Более подробную информацию о синтаксисе и ключевых словах можно найти в [официальной документации GitLab CI/CD](https://docs.gitlab.com/ee/ci/yaml/).

### Основные ключевые слова

| Ключевое слово | Описание |
|----------------|----------|
| `stages` | Определяет этапы пайплайна |
| `stage` | Указывает, к какому этапу относится job |
| `script` | Команды, которые будут выполнены |
| `before_script` | Команды, выполняемые перед каждым job |
| `after_script` | Команды, выполняемые после каждого job |
| `artifacts` | Файлы, которые будут сохранены после выполнения job |
| `cache` | Файлы, которые будут кэшированы между job |
| `only` | Условия, при которых job будет выполнен |
| `except` | Условия, при которых job не будет выполнен |
| `rules` | Более гибкие условия выполнения job |
| `needs` | Создает зависимости между job |
| `tags` | Указывает, на каких runners выполнять job |
| `timeout` | Максимальное время выполнения job |
| `retry` | Количество попыток при неудаче |
| `parallel` | Количество параллельных экземпляров job |
| `environment` | Окружение для развертывания |
| `resource_group` | Группировка ресурсов для предотвращения конфликтов |

### Практическое задание 2.1: Создание первого пайплайна

1. Создайте новый проект в GitLab
2. Добавьте файл `.gitlab-ci.yml` со следующим содержимым:

```yaml
stages:
  - build
  - test

build-job:
  stage: build
  script:
    - echo "Building the application..."
    - mkdir -p build
    - echo "Build successful!" > build/status.txt
  artifacts:
    paths:
      - build/

test-job:
  stage: test
  script:
    - echo "Running tests..."
    - cat build/status.txt
    - echo "Tests passed!"
```

3. Закоммитьте и запушьте изменения
4. Перейдите в раздел CI/CD → Pipelines и наблюдайте за выполнением

---

## Часть 2.1: Настройка GitLab Runner

### Типы GitLab Runner

GitLab Runner может работать в различных режимах:

1. **Shared Runner** - предоставляется GitLab, доступен всем проектам
2. **Group Runner** - настраивается для группы проектов
3. **Specific Runner** - настраивается для конкретного проекта

### Установка GitLab Runner

#### Установка на Linux (Ubuntu/Debian)

```bash
# Добавляем официальный репозиторий GitLab
curl -L "https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh" | sudo bash

# Устанавливаем GitLab Runner
sudo apt-get install gitlab-runner

# Регистрируем runner
sudo gitlab-runner register
```

#### Установка с Docker

```bash
# Запускаем GitLab Runner в контейнере
docker run -d --name gitlab-runner --restart always \
  -v /srv/gitlab-runner/config:/etc/gitlab-runner \
  -v /var/run/docker.sock:/var/run/docker.sock \
  gitlab/gitlab-runner:latest

# Регистрируем runner
docker exec -it gitlab-runner gitlab-runner register
```

### Конфигурация GitLab Runner

Файл конфигурации `/etc/gitlab-runner/config.toml`:

```toml
concurrent = 4
check_interval = 0

[session_server]
  session_timeout = 1800

[[runners]]
  name = "docker-runner"
  url = "https://gitlab.com/"
  token = "YOUR_RUNNER_TOKEN"
  executor = "docker"
  [runners.docker]
    tls_verify = false
    image = "docker:24.0.5"
    privileged = false
    disable_entrypoint_overwrite = false
    oom_kill_disable = false
    disable_cache = false
    volumes = ["/cache"]
    shm_size = 0
  [runners.cache]
    [runners.cache.s3]
    [runners.cache.gcs]
```

### Теги и их использование

Теги позволяют назначать задания на конкретные runners:

```yaml
# Задание будет выполняться только на runner с тегом docker
build-job:
  stage: build
  tags:
    - docker
  script:
    - docker build -t app .

# Задание будет выполняться только на runner с тегом production
deploy-job:
  stage: deploy
  tags:
    - production
    - shell
  script:
    - ./deploy.sh
```

---

## Часть 3: Создание пайплайна с unit-тестами

### Поддержка различных языков программирования

GitLab CI/CD поддерживает большинство популярных языков программирования. Вот примеры для разных технологий:

#### Python с pytest и coverage

```yaml
stages:
  - test

pytest:
  stage: test
  image: python:3.12
  before_script:
    - pip install -r requirements.txt
    - pip install pytest pytest-cov
  script:
    - pytest --junitxml=report.xml --cov=. --cov-report=xml:coverage.xml
  coverage: '/TOTAL.+ ([0-9]{1,3}%)/'
  artifacts:
    when: always
    reports:
      junit: report.xml
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml
```

#### Node.js с Jest и coverage

```yaml
stages:
  - test

jest-test:
  stage: test
  image: node:20
  before_script:
    - npm ci
  script:
    - npm test -- --coverage --testResultsProcessor=jest-junit
  coverage: '/All files[^|]*\|[^|]*\s+([\d\.]+)/'
  artifacts:
    when: always
    reports:
      junit: junit.xml
      coverage_report:
        coverage_format: cobertura
        path: coverage/cobertura-coverage.xml
    paths:
      - coverage/
```

#### Java с Maven и JaCoCo

```yaml
stages:
  - test

maven-test:
  stage: test
  image: maven:3.8-openjdk-11
  script:
    - mvn test
  artifacts:
    when: always
    reports:
      junit:
        - target/surefire-reports/TEST-*.xml
        - target/failsafe-reports/TEST-*.xml
      coverage_report:
        coverage_format: cobertura
        path: target/site/cobertura.xml
```

#### Python с pytest

```yaml
stages:
  - test

pytest:
  stage: test
  image: python:3.12
  before_script:
    - pip install -r requirements.txt
    - pip install pytest
  script:
    - pytest --junitxml=report.xml
  artifacts:
    when: always
    reports:
      junit: report.xml
```

#### Node.js с Jest

```yaml
stages:
  - test

jest-test:
  stage: test
  image: node:20
  before_script:
    - npm ci
  script:
    - npm test -- --coverage --testResultsProcessor=jest-junit
  artifacts:
    when: always
    reports:
      junit: junit.xml
    paths:
      - coverage/
```

#### Java с Maven

```yaml
stages:
  - test

maven-test:
  stage: test
  image: maven:3.8-openjdk-11
  script:
    - mvn test
  artifacts:
    when: always
    reports:
      junit:
        - target/surefire-reports/TEST-*.xml
        - target/failsafe-reports/TEST-*.xml
```

#### Go с gotestsum

```yaml
stages:
  - test

golang-test:
  stage: test
  image: golang:1.21
  script:
    - go install gotest.tools/gotestsum@latest
    - gotestsum --junitfile report.xml --format testname
  artifacts:
    when: always
    reports:
      junit: report.xml
```

#### Ruby с RSpec

```yaml
stages:
  - test

rspec:
  stage: test
  image: ruby:3.2
  before_script:
    - bundle install
  script:
    - bundle exec rspec --format progress --format RspecJunitFormatter --out rspec.xml
  artifacts:
    when: always
    reports:
      junit: rspec.xml
```

Дополнительные примеры для различных языков программирования можно найти в [официальной документации GitLab](https://docs.gitlab.com/ee/ci/examples/).

### Практическое задание 3.1: Python приложение с тестами

Создайте простое Python приложение с unit-тестами:

**calculator.py:**
```python
def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

def multiply(a, b):
    return a * b

def divide(a, b):
    if b == 0:
        raise ValueError("Cannot divide by zero")
    return a / b
```

**test_calculator.py:**
```python
import unittest
from calculator import add, subtract, multiply, divide

class TestCalculator(unittest.TestCase):
    def test_add(self):
        self.assertEqual(add(2, 3), 5)
        self.assertEqual(add(-1, 1), 0)

    def test_subtract(self):
        self.assertEqual(subtract(5, 3), 2)
        self.assertEqual(subtract(0, 5), -5)

    def test_multiply(self):
        self.assertEqual(multiply(3, 4), 12)
        self.assertEqual(multiply(-2, 3), -6)

    def test_divide(self):
        self.assertEqual(divide(10, 2), 5)
        self.assertEqual(divide(9, 3), 3)
        
        with self.assertRaises(ValueError):
            divide(10, 0)

if __name__ == '__main__':
    unittest.main()
```

**requirements.txt:**
```
pytest==7.4.3
pytest-cov==4.1.0
```

**.gitlab-ci.yml:**
```yaml
stages:
  - test

variables:
  PIP_CACHE_DIR: "$CI_PROJECT_DIR/.cache/pip"

cache:
  paths:
    - .cache/pip
    - venv/

before_script:
  - python -V
  - pip install virtualenv
  - virtualenv venv
  - source venv/bin/activate
  - pip install -r requirements.txt

pytest:
  stage: test
  script:
    - python -m pytest --junitxml=report.xml --cov=. --cov-report=xml:coverage.xml
  coverage: '/TOTAL.+ ([0-9]{1,3}%)/'
  artifacts:
    when: always
    reports:
      junit: report.xml
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml
    paths:
      - htmlcov/
```

---

## Часть 4: Сборка и публикация Docker-образов

### Работа с Container Registry

GitLab предоставляет встроенный Container Registry для хранения Docker-образов. Для работы с ним используются предопределенные переменные:

- `$CI_REGISTRY` - адрес registry
- `$CI_REGISTRY_IMAGE` - путь к образу
- `$CI_REGISTRY_USER` - имя пользователя
- `$CI_REGISTRY_PASSWORD` - пароль (job token)

### Практическое задание 4.1: Сборка Docker-образа

Создайте простое веб-приложение и соберите для него Docker-образ:

**app.py:**
```python
from flask import Flask

app = Flask(__name__)

@app.route('/')
def hello():
    return 'Hello from GitLab CI/CD!'

@app.route('/health')
def health():
    return {'status': 'healthy'}

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

**requirements.txt:**
```
Flask==2.3.2
```

**Dockerfile:**
```dockerfile
FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 5000

CMD ["python", "app.py"]
```

**.gitlab-ci.yml:**
```yaml
stages:
  - build
  - test
  - publish

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"
  CONTAINER_IMAGE: $CI_REGISTRY_IMAGE:$CI_COMMIT_REF_SLUG

before_script:
  - docker info

build-image:
  stage: build
  image: docker:24.0.5
  services:
    - docker:24.0.5-dind
  script:
    - echo "Building Docker image..."
    - docker build -t $CONTAINER_IMAGE .
    - echo "Image built successfully"
  artifacts:
    reports:
      dotenv: build.env

test-image:
  stage: test
  image: docker:24.0.5
  services:
    - docker:24.0.5-dind
  script:
    - echo "Testing Docker image..."
    - docker run --rm $CONTAINER_IMAGE python -c "print('Image is working')"
    - echo "Image test passed"

publish-image:
  stage: publish
  image: docker:24.0.5
  services:
    - docker:24.0.5-dind
  before_script:
    - echo "Logging into GitLab Container Registry..."
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - echo "Publishing Docker image..."
    - docker push $CONTAINER_IMAGE
    - echo "Image published to $CONTAINER_IMAGE"
  only:
    - main
```

### Публикация в Docker Hub

Для публикации в Docker Hub необходимо добавить переменные окружения с учетными данными:

```yaml
publish-to-dockerhub:
  stage: publish
  image: docker:24.0.5
  services:
    - docker:24.0.5-dind
  before_script:
    - echo "Logging into Docker Hub..."
    - docker login -u $DOCKER_HUB_USER -p $DOCKER_HUB_PASSWORD
  script:
    - docker tag $CONTAINER_IMAGE $DOCKER_HUB_USER/my-app:$CI_COMMIT_REF_SLUG
    - docker push $DOCKER_HUB_USER/my-app:$CI_COMMIT_REF_SLUG
  only:
    - main
```

---

## Часть 5: Работа с артефактами и кэшем

### Артефакты (Artifacts)

Артефакты - это файлы, созданные в процессе выполнения job, которые можно скачать или передать между job.

```yaml
build-app:
  stage: build
  script:
    - npm run build
  artifacts:
    paths:
      - dist/
    expire_in: 1 week
```

Более подробную информацию об артефактах можно найти в [документации GitLab](https://docs.gitlab.com/ee/ci/pipelines/job_artifacts.html).

### Кэш (Cache)

Кэш используется для ускорения выполнения job за счет сохранения зависимостей между запусками.

```yaml
variables:
  NODE_MODULES_CACHE_KEY: $CI_COMMIT_REF_SLUG

cache:
  key: $NODE_MODULES_CACHE_KEY
  paths:
    - node_modules/

install-dependencies:
  script:
    - npm ci
```

### Продвинутые техники кэширования

#### Кэширование для Python

```yaml
variables:
  PIP_CACHE_DIR: "$CI_PROJECT_DIR/.cache/pip"

cache:
  key: $CI_COMMIT_REF_SLUG
  paths:
    - .cache/pip
    - venv/

before_script:
  - python -m venv venv
  - source venv/bin/activate
  - pip install --cache-dir $PIP_CACHE_DIR -r requirements.txt

test-python:
  stage: test
  script:
    - pytest
```

#### Кэширование для Node.js

```yaml
variables:
  NPM_CACHE_DIR: "$CI_PROJECT_DIR/.npm"
  NODE_MODULES_CACHE_KEY: $CI_COMMIT_REF_SLUG

cache:
  key: ${NODE_MODULES_CACHE_KEY}
  paths:
    - node_modules/
    - .npm/

install-dependencies:
  stage: build
  script:
    - npm ci --cache $NPM_CACHE_DIR --prefer-offline
```

#### Кэширование для Maven

```yaml
variables:
  MAVEN_OPTS: "-Dmaven.repo.local=$CI_PROJECT_DIR/.m2/repository"

cache:
  key: $CI_COMMIT_REF_SLUG
  paths:
    - .m2/repository/

maven-build:
  stage: build
  script:
    - mvn clean install
```

#### Кэширование для Docker слоев

```yaml
variables:
  DOCKER_BUILDKIT: 1
  BUILDKIT_INLINE_CACHE: 1

cache:
  key: $CI_COMMIT_REF_SLUG
  paths:
    - .docker-cache/

docker-build:
  stage: build
  script:
    - |
      docker build \
        --cache-from $CI_REGISTRY_IMAGE:cache \
        --cache-to $CI_REGISTRY_IMAGE:cache \
        --build-arg BUILDKIT_INLINE_CACHE=1 \
        -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA \
        .
```

Подробнее о кэшировании в GitLab CI/CD можно узнать из [официальной документации](https://docs.gitlab.com/ee/ci/caching/).

### Практическое задание 5.1: Оптимизация сборки с кэшем

```yaml
stages:
  - build
  - test

variables:
  NODE_VERSION: "20"
  CACHE_KEY: $CI_COMMIT_REF_SLUG-node-$NODE_VERSION

cache:
  key: $CACHE_KEY
  paths:
    - node_modules/
    - .npm/

before_script:
  - node -v
  - npm -v

install:
  stage: build
  script:
    - npm ci --cache .npm --prefer-offline
  cache:
    key: $CACHE_KEY
    paths:
      - node_modules/
      - .npm/
    policy: pull-push

build:
  stage: build
  script:
    - npm run build
  artifacts:
    paths:
      - dist/
    expire_in: 1 day

test:
  stage: test
  script:
    - npm run test
  artifacts:
    reports:
      junit: test-results.xml

deploy:
  stage: deploy
  script:
    - echo "Deploying application..."
  only:
    - main
```

---

## Часть 5.1: Оптимизация пайплайнов и работа с зависимостями

### Управление зависимостями между job

#### Использование needs для создания DAG (Directed Acyclic Graph)

```yaml
stages:
  - build
  - test
  - security
  - deploy

build-frontend:
  stage: build
  script:
    - echo "Building frontend"
    - npm run build
  artifacts:
    paths:
      - dist/

build-backend:
  stage: build
  script:
    - echo "Building backend"
    - mvn package
  artifacts:
    paths:
      - target/*.jar

test-frontend:
  stage: test
  needs: [build-frontend]
  script:
    - echo "Testing frontend"
    - npm test

test-backend:
  stage: test
  needs: [build-backend]
  script:
    - echo "Testing backend"
    - mvn test

security-scan:
  stage: security
  needs: [test-frontend, test-backend]
  script:
    - echo "Running security scan"
    - ./security-scan.sh

deploy-staging:
  stage: deploy
  needs: [security-scan]
  script:
    - echo "Deploying to staging"
    - ./deploy-staging.sh
  environment:
    name: staging
```

#### Параллельное выполнение с зависимостями

```yaml
stages:
  - validate
  - test
  - build
  - deploy

validate-code:
  stage: validate
  script:
    - echo "Running linters and formatters"
    - npm run lint
    - npm run format:check

unit-tests:
  stage: test
  needs: [validate-code]
  parallel: 3
  script:
    - echo "Running unit tests - shard $CI_NODE_INDEX"
    - npm run test:shard $CI_NODE_INDEX

integration-tests:
  stage: test
  needs: [validate-code]
  script:
    - echo "Running integration tests"
    - npm run test:integration

build-app:
  stage: build
  needs: [unit-tests, integration-tests]
  script:
    - echo "Building application"
    - npm run build
  artifacts:
    paths:
      - dist/

deploy-app:
  stage: deploy
  needs: [build-app]
  script:
    - echo "Deploying application"
    - npm run deploy
  only:
    - main
```

### Оптимизация времени выполнения

#### Использование Docker Layer Caching

```yaml
variables:
  DOCKER_BUILDKIT: 1
  BUILDKIT_INLINE_CACHE: 1

build-docker:
  stage: build
  image: docker:24.0.5
  services:
    - docker:24.0.5-dind
  script:
    - |
      docker build \
        --cache-from $CI_REGISTRY_IMAGE:latest \
        --cache-to type=inline \
        --build-arg BUILDKIT_INLINE_CACHE=1 \
        -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA \
        -t $CI_REGISTRY_IMAGE:latest \
        .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
    - docker push $CI_REGISTRY_IMAGE:latest
```

#### Умное кэширование с ключами

```yaml
variables:
  CACHE_KEY_PREFIX: "v1"

cache:
  key: ${CACHE_KEY_PREFIX}-${CI_COMMIT_REF_SLUG}-${CI_PROJECT_ID}
  paths:
    - node_modules/
    - .npm/
    - .cache/

# Разные ключи для разных сценариев
cache-main:
  cache:
    key: ${CACHE_KEY_PREFIX}-main-${CI_PROJECT_ID}
    paths:
      - node_modules/
  only:
    - main

cache-mr:
  cache:
    key: ${CACHE_KEY_PREFIX}-mr-${CI_MERGE_REQUEST_IID}-${CI_PROJECT_ID}
    paths:
      - node_modules/
  only:
    - merge_requests
```

#### Оптимизация тестирования

```yaml
# Разделение тестов по типам
test-unit:
  stage: test
  script:
    - npm run test:unit
  artifacts:
    reports:
      junit: junit-unit.xml
    expire_in: 1 week

test-integration:
  stage: test
  script:
    - npm run test:integration
  artifacts:
    reports:
      junit: junit-integration.xml
    expire_in: 1 week

test-e2e:
  stage: test
  script:
    - npm run test:e2e
  artifacts:
    reports:
      junit: junit-e2e.xml
    expire_in: 1 week
  when: manual  # Запускать вручную для экономии времени
```

### Мониторинг и метрики

#### Сбор метрик пайплайна

```yaml
collect-metrics:
  stage: .post
  script:
    - |
      echo "Pipeline duration: $CI_PIPELINE_DURATION_SECONDS seconds"
      echo "Job duration: $CI_JOB_DURATION_SECONDS seconds"
      echo "Total jobs: $(echo $CI_JOB_NAME | wc -l)"
      
      # Отправка метрик в систему мониторинга
      curl -X POST "$METRICS_ENDPOINT" \
        -H "Content-Type: application/json" \
        -d "{
          \"pipeline_id\": \"$CI_PIPELINE_ID\",
          \"project\": \"$CI_PROJECT_NAME\",
          \"branch\": \"$CI_COMMIT_REF_NAME\",
          \"duration\": $CI_PIPELINE_DURATION_SECONDS,
          \"status\": \"$CI_PIPELINE_SOURCE\"
        }"
  when: always
```

---

## Часть 6: Переменные и секреты в CI/CD

### Типы переменных

1. **Predefined Variables** - предопределенные переменные GitLab
2. **Project Variables** - переменные уровня проекта
3. **Group Variables** - переменные уровня группы
4. **Environment Variables** - переменные окружения

### Наиболее используемые предопределенные переменные

| Переменная | Описание | Пример |
|------------|----------|--------|
| `$CI_COMMIT_REF_NAME` | Имя ветки или тега | `main`, `feature/test` |
| `$CI_COMMIT_SHA` | Хеш коммита | `1a2b3c4d5e6f` |
| `$CI_PIPELINE_ID` | ID пайплайна | `12345` |
| `$CI_JOB_ID` | ID задания | `67890` |
| `$CI_REGISTRY` | Адрес Container Registry | `registry.gitlab.com` |
| `$CI_REGISTRY_IMAGE` | Путь к образу | `registry.gitlab.com/group/project` |
| `$CI_PROJECT_DIR` | Путь к проекту | `/builds/group/project` |
| `$CI_PROJECT_NAME` | Имя проекта | `my-project` |
| `$CI_PROJECT_NAMESPACE` | Пространство имен проекта | `my-group` |
| `$CI_RUNNER_TAGS` | Теги runner'а | `docker,linux` |

### Настройка переменных проекта

1. **Переменные проекта**:
   - Перейдите в Settings → CI/CD → Variables
   - Нажмите "Add variable"
   - Укажите ключ, значение и параметры:
     - **Protected** - доступно только для защищенных веток
     - **Masked** - скрыто в логах
     - **Environment scope** - окружение, для которого доступна переменная

2. **Переменные группы**:
   - Перейдите в Group Settings → CI/CD → Variables
   - Настройте аналогично переменным проекта

### Использование переменных в .gitlab-ci.yml

```yaml
variables:
  # Глобальные переменные
  NODE_VERSION: "20"
  APP_NAME: "my-application"

build-job:
  stage: build
  variables:
    # Переменные уровня job
    BUILD_ENV: "production"
  script:
    - echo "Building $APP_NAME with Node.js $NODE_VERSION"
    - echo "Build environment: $BUILD_ENV"
    - echo "Commit: $CI_COMMIT_SHA"
    - echo "Branch: $CI_COMMIT_REF_NAME"

deploy-job:
  stage: deploy
  script:
    - echo "Deploying to $DEPLOY_TARGET"
    - echo "Registry: $CI_REGISTRY"
    - echo "Image: $CI_REGISTRY_IMAGE:$CI_COMMIT_REF_SLUG"
  only:
    - main
```

### Переменные окружения (Environments)

```yaml
deploy-staging:
  stage: deploy
  script:
    - echo "Deploying to staging"
    - ./deploy.sh --env=staging
  environment:
    name: staging
    url: https://staging.example.com
  variables:
    DEPLOY_TARGET: "staging"
    API_URL: "https://api-staging.example.com"

deploy-production:
  stage: deploy
  script:
    - echo "Deploying to production"
    - ./deploy.sh --env=production
  environment:
    name: production
    url: https://example.com
  variables:
    DEPLOY_TARGET: "production"
    API_URL: "https://api.example.com"
  when: manual
  only:
    - main
```

Подробнее о переменных в GitLab CI/CD можно узнать из [официальной документации](https://docs.gitlab.com/ee/ci/variables/).

### Работа с секретами

Для хранения секретов используются masked variables:

```yaml
deploy-production:
  stage: deploy
  script:
    - echo "Deploying to production with API key: $PRODUCTION_API_KEY"
    - ./deploy.sh
  environment:
    name: production
  only:
    - main
```

### Практическое задание 6.1: Безопасное хранение конфигураций

1. Перейдите в Settings → CI/CD → Variables
2. Добавьте переменные:
   - `DATABASE_URL` (masked)
   - `API_KEY` (masked)
   - `DEPLOY_TARGET`

3. Используйте их в .gitlab-ci.yml:

```yaml
stages:

deploy-app:
  stage: deploy
  script:
    - echo "Deploying to $DEPLOY_TARGET"
    - ./deploy.sh --database-url=$DATABASE_URL --api-key=$API_KEY
  environment:
    name: $DEPLOY_TARGET
  only:
    - main
```

---

## Часть 7: Триггеры и расписание пайплайнов

### Триггеры пайплайнов

Пайплайны можно запускать по различным событиям:

```yaml
# Запуск при пуше в любую ветку
push-job:
  script:
    - echo "Push detected"
  only:
    - branches

# Запуск при создании merge request
mr-job:
  script:
    - echo "Merge request created"
  only:
    - merge_requests

# Запуск по расписанию
scheduled-job:
  script:
    - echo "Scheduled job running"
  only:
    - schedules
```

### Расписание пайплайнов

1. Перейдите в CI/CD → Schedules
2. Нажмите "New schedule"
3. Задайте параметры:
   - Interval Pattern (cron syntax)
   - Target branch
   - Variables (опционально)

### Практическое задание 7.1: Создание расписания

Создайте пайплайн, который ежедневно запускает тесты:

```yaml
daily-test:
  stage: test
  script:
    - echo "Running daily tests at $(date)"
    - python -m pytest
  only:
    - schedules
```

---

## Часть 8: Продвинутые техники CI/CD

### Multi-Project Pipelines

Multi-project pipelines позволяют создавать зависимости между пайплайнами разных проектов:

```yaml
# В проекте A (upstream)
trigger-downstream:
  stage: deploy
  trigger:
    project: group/project-b
    branch: main
    strategy: depend

# В проекте B (downstream)
build-from-upstream:
  stage: build
  script:
    - echo "Building from upstream pipeline"
    - echo "Upstream project: $CI_UPSTREAM_PROJECT"
    - echo "Upstream pipeline: $CI_UPSTREAM_PIPELINE_ID"
```

### Child Pipelines

Child pipelines позволяют разбивать сложные пайплайны на управляемые части:

```yaml
# Основной пайплайн
stages:
  - generate
  - deploy

generate-child-pipeline:
  stage: generate
  script:
    - |
      cat > child-pipeline.yml << EOF
      stages:
        - build
        - test
        - deploy
      
      build-app:
        stage: build
        script:
          - echo "Building application"
      
      test-app:
        stage: test
        script:
          - echo "Testing application"
      
      deploy-app:
        stage: deploy
        script:
          - echo "Deploying application"
      EOF
  artifacts:
    paths:
      - child-pipeline.yml

trigger-child-pipeline:
  stage: deploy
  trigger:
    include:
      - artifact: child-pipeline.yml
        job: generate-child-pipeline
```

### Security Scanning

GitLab предоставляет встроенные инструменты для сканирования безопасности:

```yaml
# SAST (Static Application Security Testing)
sast:
  stage: test
  image: docker:24.0.5
  services:
    - docker:24.0.5-dind
  variables:
    SEARCH_MAX_DEPTH: 4
  script:
    - echo "Running SAST scan"
    - |
      if [ -f "Dockerfile" ]; then
        docker run --rm -v "$PWD":/app -w /app \
          registry.gitlab.com/security-products/sast:latest \
          /app
      fi
  artifacts:
    reports:
      sast: gl-sast-report.json
  only:
    - branches
  except:
    - schedules

# Dependency Scanning
dependency-scanning:
  stage: test
  image: docker:24.0.5
  services:
    - docker:24.0.5-dind
  script:
    - echo "Running dependency scanning"
    - |
      if [ -f "package.json" ]; then
        docker run --rm -v "$PWD":/app -w /app \
          registry.gitlab.com/security-products/dependency-scanning:latest \
          /app
      fi
  artifacts:
    reports:
      dependency_scanning: gl-dependency-scanning-report.json
  only:
    - branches
  except:
    - schedules

# Container Scanning
container-scanning:
  stage: test
  image: docker:24.0.5
  services:
    - docker:24.0.5-dind
  script:
    - echo "Running container scanning"
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - |
      docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
        -v "$PWD":/app -w /app \
        registry.gitlab.com/security-products/container-scanning:latest \
        /app
  artifacts:
    reports:
      container_scanning: gl-container-scanning-report.json
  only:
    - branches
  except:
    - schedules
```

### Parallel Jobs

Параллельное выполнение заданий для ускорения пайплайна:

```yaml
test-parallel:
  stage: test
  parallel: 4
  script:
    - echo "Running test job $CI_NODE_INDEX of $CI_NODE_TOTAL"
    - |
      # Разделение тестов между параллельными job
      TOTAL_TESTS=$(find . -name "*_test.rb" | wc -l)
      TESTS_PER_JOB=$((TOTAL_TESTS / CI_NODE_TOTAL + 1))
      START=$((CI_NODE_INDEX * TESTS_PER_JOB - TESTS_PER_JOB + 1))
      END=$((CI_NODE_INDEX * TESTS_PER_JOB))
      
      # Запуск только своей части тестов
      find . -name "*_test.rb" | sed -n "${START},${END}p" | xargs rspec
```

### Resource Groups

Resource groups предотвращают одновременное выполнение конфликтующих операций:

```yaml
deploy-production:
  stage: deploy
  script:
    - echo "Deploying to production"
    - ./deploy.sh
  resource_group: production
  environment:
    name: production
    url: https://example.com
  only:
    - main
```

### Conditional Execution с Rules

Более гибкое управление выполнением job с использованием rules:

```yaml
job-with-rules:
  script:
    - echo "This job uses rules"
  rules:
    # Запускать для merge requests
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    # Запускать для main ветки
    - if: $CI_COMMIT_REF_NAME == "main"
    # Запускать для тегов
    - if: $CI_COMMIT_TAG
    # Запускать по расписанию
    - if: $CI_PIPELINE_SOURCE == "schedule"
    # Запускать вручную
    - when: manual
      allow_failure: true
```

### Matrix Strategy

Выполнение job с различными комбинациями параметров:

```yaml
test-matrix:
  stage: test
  parallel:
    matrix:
      - NODE_VERSION: [16, 18, 20]
        OS: [ubuntu, alpine]
  script:
    - echo "Testing with Node.js $NODE_VERSION on $OS"
    - |
      if [ "$OS" == "ubuntu" ]; then
        image="node:$NODE_VERSION"
      else
        image="node:$NODE_VERSION-alpine"
      fi
      docker run --rm $image npm test
```

### Blue-Green Deployment

Реализация blue-green развертывания:

```yaml
deploy-blue:
  stage: deploy
  script:
    - echo "Deploying to blue environment"
    - docker-compose -f docker-compose.blue.yml up -d
    - ./health-check.sh blue
  environment:
    name: blue
    url: https://blue.example.com
  only:
    - main

switch-to-blue:
  stage: deploy
  script:
    - echo "Switching traffic to blue environment"
    - ./switch-traffic.sh blue
  environment:
    name: production
    url: https://example.com
  when: manual
  only:
    - main
  dependencies:
    - deploy-blue

deploy-green:
  stage: deploy
  script:
    - echo "Deploying to green environment"
    - docker-compose -f docker-compose.green.yml up -d
    - ./health-check.sh green
  environment:
    name: green
    url: https://green.example.com
  only:
    - main

switch-to-green:
  stage: deploy
  script:
    - echo "Switching traffic to green environment"
    - ./switch-traffic.sh green
  environment:
    name: production
    url: https://example.com
  when: manual
  only:
    - main
  dependencies:
    - deploy-green
```

---

## Самостоятельное задание

Подготовьте репозиторий с пайплайном для тестирования приложения через unit тесты, а также сборку после удачного тестирования приложения и публикацию в GitLab Artifact Registry или DockerHub.

### Требования к заданию:

#### 1. Репозиторий (20 баллов)
- ✅ Создан репозиторий в GitLab
- ✅ Содержит README с описанием проекта
- ✅ Используется осмысленная структура файлов
- ✅ Есть .gitignore с нужными исключениями
- ✅ Есть лицензия и информация о разработчиках

#### 2. Unit-тесты (25 баллов)
- ✅ Реализованы unit-тесты для приложения
- ✅ Тесты покрывают основную функциональность (>70%)
- ✅ Пайплайн запускает тесты автоматически
- ✅ Результаты тестов отображаются в GitLab
- ✅ Используются отчеты о покрытии кода
- ✅ Есть интеграционные тесты (дополнительно)

#### 3. Сборка приложения (25 баллов)
- ✅ Приложение собирается в пайплайне
- ✅ Сборка запускается только после успешных тестов
- ✅ Используются артефакты для передачи сборки между этапами
- ✅ Оптимизирована с помощью кэширования
- ✅ Есть проверка работоспособности сборки
- ✅ Используются multi-stage builds для Docker (дополнительно)

#### 4. Публикация (20 баллов)
- ✅ Docker-образ публикуется в GitLab Container Registry или DockerHub
- ✅ Публикация происходит только после успешной сборки
- ✅ Используются теги для версионирования
- ✅ Образ содержит все необходимые зависимости
- ✅ Есть проверка работоспособности опубликованного образа
- ✅ Настроено автоматическое тегирование для релизов (дополнительно)

#### 5. Конфигурация пайплайна (10 баллов)
- ✅ Используются переменные окружения
- ✅ Есть разные этапы (build, test, deploy)
- ✅ Настроены условия запуска job
- ✅ Используются секреты для аутентификации
- ✅ Есть обработка ошибок и логирование
- ✅ Используются правила вместо only/except (дополнительно)

### Пример структуры проекта

```
my-application/
├── README.md
├── .gitignore
├── .gitlab-ci.yml
├── Dockerfile
├── requirements.txt  # для Python
├── package.json      # для Node.js
├── src/
│   ├── main.py
│   └── utils.py
├── tests/
│   ├── test_main.py
│   └── test_utils.py
└── docs/
    └── api.md
```

### Пример .gitlab-ci.yml для справки

```yaml
stages:
  - test
  - build
  - deploy

variables:
  IMAGE_TAG: $CI_REGISTRY_IMAGE:$CI_COMMIT_REF_SLUG

test:
  stage: test
  image: python:3.12
  before_script:
    - pip install -r requirements.txt
    - pip install pytest pytest-cov
  script:
    - pytest --junitxml=report.xml --cov=. --cov-report=xml
  coverage: '/TOTAL.+ ([0-9]{1,3}%)/'
  artifacts:
    reports:
      junit: report.xml
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml

build:
  stage: build
  image: docker:24.0.5
  services:
    - docker:24.0.5-dind
  script:
    - docker build -t $IMAGE_TAG .
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker push $IMAGE_TAG
  only:
    - main
    - develop

deploy:
  stage: deploy
  script:
    - echo "Deploying $IMAGE_TAG"
    # Здесь ваша логика развертывания
  only:
    - main
  when: manual
```

### Дополнительные бонусы (до 10 баллов)

- **Security scanning** (3 балла): Включите SAST или dependency scanning
- **Parallel execution** (2 балла): Используйте параллельное выполнение тестов
- **Environment management** (3 балла): Настройте staging и production окружения
- **Monitoring** (2 балла): Добавьте сбор метрик и уведомлений

### Критерии оценки

| Критерий | Баллы |
|----------|-------|
| Репозиторий | 20 |
| Unit-тесты | 25 |
| Сборка приложения | 25 |
| Публикация | 20 |
| Конфигурация пайплайна | 10 |
| **ИТОГО** | **100** |

---

## Контрольные вопросы

### Базовые вопросы

1. **Что такое GitLab CI/CD и какие преимущества он предоставляет?**
   - Опишите основные компоненты архитектуры

2. **Какова структура файла .gitlab-ci.yml?**
   - Назовите основные элементы конфигурации

3. **В чем разница между stages, jobs и pipeline?**
   - Как они взаимодействуют между собой?

4. **Какие типы переменных существуют в GitLab CI/CD?**
   - Где и как их правильно использовать?

5. **Что такое артефакты и для чего они нужны?**
   - Как настроить их хранение и передачу между job?

### Продвинутые вопросы

6. **Как оптимизировать пайплайн с помощью кэширования?**
   - Приведите примеры для разных языков программирования

7. **Как обеспечить безопасность при работе с секретами?**
   - Какие механизмы защиты предоставляет GitLab?

8. **Как настроить автоматическую публикацию Docker-образов?**
   - В чем разница между GitLab Registry и DockerHub?

9. **Какие способы запуска пайплайнов существуют?**
   - Когда использовать каждый из них?

10. **Как настроить расписание выполнения пайплайнов?**
    - Приведите пример cron-выражения для ежедневного запуска

11. **Как отладить пайплайн, который завершается с ошибкой?**
    - Какие инструменты диагностики предоставляет GitLab?

12. **Что такое dependency proxy и зачем он нужен?**
    - Как его настроить и использовать?

13. **Как реализовать blue-green deployment в GitLab CI/CD?**
    - Опишите архитектуру такого решения

14. **Как настроить автоматическое тестирование в merge requests?**
    - Как отобразить результаты в интерфейсе GitLab?

15. **Какие best practices существуют для написания .gitlab-ci.yml?**
    - Приведите 5 рекомендаций по оптимизации пайплайнов

### Экспертные вопросы

16. **Как работают multi-project pipelines и когда их следует использовать?**
    - Опишите сценарии использования и преимущества

17. **В чем разница между child pipelines и parent pipelines?**
    - Когда следует разбивать пайплайн на дочерние части?

18. **Как настроить security scanning в GitLab CI/CD?**
    - Какие типы сканирования поддерживаются и как их интегрировать?

19. **Как реализовать parallel execution для ускорения пайплайна?**
    - Приведите примеры разделения тестов между параллельными job

20. **Что такое resource groups и как они помогают избежать конфликтов?**
    - Опишите сценарии использования

21. **Как работает matrix strategy и когда ее применять?**
    - Приведите примеры тестирования на разных окружениях

22. **Как реализовать blue-green deployment в GitLab CI/CD?**
    - Опишите архитектуру и шаги реализации

23. **В чем разница между rules и only/except?**
    - Когда следует использовать каждый подход?

24. **Как оптимизировать время выполнения пайплайна?**
    - Приведите 10 конкретных рекомендаций

25. **Как настроить мониторинг и алерты для CI/CD пайплайнов?**
    - Какие метрики важно отслеживать?

---

## Итоги лабораторной работы

### Что было изучено

В ходе выполнения лабораторной работы вы освоили:

1. **Основы CI/CD**:
   - Понимание принципов непрерывной интеграции и доставки
   - Архитектура GitLab CI/CD
   - Преимущества автоматизации процессов разработки

2. **Конфигурация пайплайнов**:
   - Синтаксис `.gitlab-ci.yml`
   - Основные ключевые слова и их применение
   - Структура stages, jobs и pipelines

3. **GitLab Runner**:
   - Типы runners и их настройка
   - Установка и регистрация
   - Использование тегов для распределения задач

4. **Тестирование в CI/CD**:
   - Интеграция unit-тестов для различных языков
   - Отчеты о покрытии кода
   - Автоматический запуск тестов

5. **Работа с Docker**:
   - Сборка и публикация Docker-образов
   - Работа с Container Registry
   - Оптимизация Docker-сборок

6. **Артефакты и кэширование**:
   - Управление артефактами между job
   - Оптимизация с помощью кэша
   - Продвинутые техники кэширования

7. **Переменные и секреты**:
   - Типы переменных в GitLab
   - Безопасное хранение секретов
   - Управление окружениями

8. **Триггеры и расписание**:
   - Условия запуска пайплайнов
   - Настройка расписаний
   - Автоматизация процессов

9. **Продвинутые техники**:
   - Multi-project pipelines
   - Child pipelines
   - Security scanning
   - Parallel execution
   - Blue-green deployment

### Приобретенные навыки

После выполнения этой лабораторной работы вы сможете:

- ✅ Создавать и настраивать CI/CD пайплайны в GitLab
- ✅ Интегрировать автоматическое тестирование
- ✅ Оптимизировать производительность пайплайнов
- ✅ Работать с Docker-образами и registry
- ✅ Реализовывать безопасные практики CI/CD
- ✅ Настраивать мониторинг и алерты
- ✅ Применять продвинутые техники CI/CD

### Практическое применение

Полученные знания могут быть применены для:

- Автоматизации процессов разработки
- Повышения качества кода
- Ускорения доставки продуктов
- Снижения рисков при развертывании
- Улучшения командной работы

---

## Дополнительные материалы

### Официальная документация
- [GitLab CI/CD Documentation](https://docs.gitlab.com/ee/ci/)
- [.gitlab-ci.yml Reference](https://docs.gitlab.com/ee/ci/yaml/)
- [GitLab Container Registry](https://docs.gitlab.com/ee/user/packages/container_registry/)
- [GitLab CI/CD Examples](https://docs.gitlab.com/ee/ci/examples/)
- [GitLab CI/CD Variables](https://docs.gitlab.com/ee/ci/variables/)
- [GitLab CI/CD Artifacts](https://docs.gitlab.com/ee/ci/pipelines/job_artifacts.html)
- [GitLab CI/CD Caching](https://docs.gitlab.com/ee/ci/caching/)

### Обучающие ресурсы
- [GitLab CI/CD Tutorial](https://docs.gitlab.com/ee/ci/quick_start/)
- [CI/CD Best Practices](https://docs.gitlab.com/ee/ci/pipelines/pipeline_efficiency.html)
- [Docker Integration](https://docs.gitlab.com/ee/ci/docker/)

### Полезные инструменты
- [GitLab Runner](https://docs.gitlab.com/runner/) - агент для выполнения job
- [GitLab CLI (glab)](https://gitlab.com/gitlab-org/cli) - командная строка для GitLab
- [Docker Scout](https://docs.docker.com/scout/) - анализ безопасности образов
- [gitlab-ci-local](https://github.com/firecow/gitlab-ci-local) - инструмент для локального тестирования пайплайнов

---

## Следующие шаги

После успешного выполнения этой лабораторной работы вы будете готовы к:

**Дальнейшее изучение:**
- Kubernetes для оркестрации контейнеров
- Infrastructure as Code (Terraform)
- Advanced CI/CD patterns
- Security scanning in pipelines
- Performance monitoring

---

## Часто задаваемые вопросы

**В: Можно ли использовать GitLab CI/CD бесплатно?**  
О: Да, GitLab предоставляет бесплатный tier с достаточным количеством функций для обучения.

**В: Как отладить пайплайн, который не запускается?**  
О: Проверьте синтаксис .gitlab-ci.yml, логи runner'ов и настройки проекта в разделе CI/CD.

**В: Где хранятся артефакты пайплайна?**  
О: Артефакты хранятся в GitLab и доступны через веб-интерфейс или API.

**В: Можно ли запускать пайплайны локально?**  
О: Да, с помощью GitLab Runner или инструментов вроде gitlab-ci-local.

**В: Как настроить уведомления о результатах пайплайна?**  
О: В Settings → Integrations можно настроить webhook или интеграцию с Slack/Telegram.

---

**Удачи в освоении CI/CD!** 🚀