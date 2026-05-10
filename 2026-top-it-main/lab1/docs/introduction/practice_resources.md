# Практические ресурсы для изучения Git и GitLab

## Интерактивные платформы для практики

### [Learn Git Branching](https://github.com/pcottle/learnGitBranching)
Интерактивная визуальная платформа для изучения Git, которая позволяет:
- Визуализировать ветки, коммиты и слияния
- Практиковаться в основных командах Git
- Проходить пошаговые уроки от базовых до продвинутых тем
- Использовать встроенный симулятор Git в браузере

### [Git Practice Labs](https://github.com/labex-labs/git-practice-labs)
Комплексные практические задания для отработки навыков работы с Git:
- Пошаговые лабораторные работы по основным концепциям Git
- Практика работы с ветками, слияниями и разрешением конфликтов
- Упражнения по настройке и использованию GitLab
- Сценарии совместной разработки

### [Introduction to GitHub](https://github.com/skills/introduction-to-github)
Официальные интерактивные курсы от GitHub:
- Базовые и продвинутые курсы по работе с GitHub
- Практические задания с автоматической проверкой
- Изучение Pull Requests, Issues, Actions и других функций GitHub
- Сертификация по завершении курсов

## Видеоуроки

### [Полный курс Git и GitHub](https://www.youtube.com/watch?v=zZBiln_2FhM)
Комплексный курс, охватывающий все основные аспекты работы с Git и GitHub:
- Основы системы контроля версий
- Работа с репозиториями, ветками и коммитами
- Слияние и разрешение конфликтов
- Работа с удаленными репозиториями
- Практические примеры и кейсы

### [Git за 1 час](https://www.youtube.com/watch?v=eMETcugEX_c)
Быстрый курс для начинающих, который позволяет освоить основы Git за короткое время:
- Установка и настройка Git
- Основные команды и рабочие процессы
- Практические примеры использования
- Работа с GitHub

### [Продвинутый Git](https://www.youtube.com/watch?v=O00FTZDxD0o)
Курс для тех, кто уже знаком с основами и хочет углубить свои знания:
- Продвинутые техники ветвления и слияния
- Интерактивный rebase и перезапись истории
- Работа с подмодулями
- Инструменты визуализации и отладки

## Командные инструменты CLI

### [GitHub CLI (gh)](https://github.com/cli/cli)
Официальный командный интерфейс для GitHub, который интегрирует концепции GitHub в терминал:
- Управление репозиториями, issues, pull requests и релизами
- Поддержка аутентификации через OAuth
- Возможность создания и управления CI/CD пайплайнами
- Поддержка расширений и кастомизации

#### Установка
- **macOS**: `brew install gh`
- **Windows**: `winget install GitHub.cli` или `choco install gh`
- **Linux**: `sudo apt install gh` (Debian/Ubuntu), `sudo dnf install gh` (Fedora), `sudo pacman -S github-cli` (Arch)

#### Установка пакетных менеджеров
**Homebrew (macOS)**:
```bash
# Установка Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Проверка установки
brew --version
```

**Chocolatey (Windows)**:
```powershell
# Установка Chocolatey
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Проверка установки
choco --version
```

Примеры использования:
```bash
# Аутентификация
gh auth login

# Создание репозитория
gh repo create my-project --public

# Работа с pull requests
gh pr create --title "Add new feature" --body "Implementation of new feature"
gh pr review --approve
gh pr merge

# Управление issues
gh issue list
gh issue create --title "Bug report"
```

### [GitLab CLI (glab)](https://gitlab.com/gitlab-org/cli)
Командный интерфейс для GitLab, позволяющий управлять GitLab-функциями прямо из терминала:
- Работа с issues, merge requests, CI/CD пайплайнами и релизами
- Интеграция с GitLab Duo AI/ML функциями
- Поддержка управления SSH-ключами и переменными
- Возможность выполнения raw API запросов к GitLab

#### Установка
- **macOS**: `brew install glab`
- **Windows**: `choco install glab`
- **Linux**: `sudo snap install glab` (Snapcraft), `sudo dnf install glab` (Fedora), `go install gitlab.com/gitlab-org/cli/cmd/glab@main` (все платформы)

Примеры использования:
```bash
# Аутентификация
glab auth login

# Управление merge requests
glab mr list --assignee=@me
glab mr create --title "New feature" --description "Implementation details"
glab mr approve
glab mr merge

# Работа с CI/CD
glab ci list
glab pipeline ci view
glab pipeline ci trace

# Управление issues
glab issue create -m release-2.0.0 -t "My title here" --label important
glab issue list --state opened
```

## Дополнительные инструменты и ресурсы

### Git Cheat Sheet
Шпаргалка с основными командами Git:
- Инициализация и клонирование репозиториев
- Работа с ветками и слияниями
- Управление изменениями и история коммитов
- Работа с удаленными репозиториями

### Git Workflow Visualizer
Инструменты для визуализации истории Git:
- `git log --oneline --graph --decorate --all` - визуализация в терминале
- GitKraken, SourceTree - графические клиенты с визуализацией
- Онлайн-визуализаторы для демонстрации концепций

## Рекомендации по использованию

1. Начните с **Learn Git Branching** для понимания основных концепций
2. Перейдите к **Git Practice Labs** для отработки практических навыков
3. Используйте **Introduction to GitHub** для изучения специфики работы с GitHub
4. Дополните знания с помощью видеоуроков
5. Освойте командные инструменты **gh** и **glab** для эффективной работы
6. Регулярно практикуйтесь с реальными проектами
7. Используйте шпаргалки и визуализаторы для закрепления знаний

Эти ресурсы помогут вам не только понять теорию, но и получить практический опыт работы с Git и GitLab, что является ключевым для успешного освоения систем контроля версий.