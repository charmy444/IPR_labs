# Установка Git на Windows, macOS и Linux

## Общие сведения

Git — это распределённая система контроля версий, используемая для отслеживания изменений в исходном коде. Перед началом работы с GitLab MAI необходимо установить Git на вашу операционную систему.

## Установка на Windows

### Способ 1: Установка через официальный установщик

1. Перейдите на официальный сайт Git: https://git-scm.com/download/win
2. Скачайте установщик `Git-<версия>-64-bit.exe`
3. Запустите скачанный файл
4. Следуйте инструкциям мастера установки:
   - **Select Components**: Оставьте все компоненты по умолчанию
   - **Choosing the default editor used by Git**: Выберите `Use Visual Studio Code as Git’s default editor` (если установлен VS Code) или `Use Vim as the default editor`
   - **Adjusting your PATH environment**: Выберите `Use Git from Windows Command Prompt`
   - **Choosing the SSH executable**: Оставьте `Use bundled OpenSSH`
   - **Configuring the line ending conversions**: Выберите `Checkout Windows-style, commit Unix-style line endings`
   - **Configuring the terminal emulator**: Выберите `Use MinTTY`
   - **Choosing extra options**: Оставьте все опции включёнными
5. Нажмите `Install` и дождитесь завершения установки
6. После установки нажмите `Finish`

### Способ 2: Установка через Chocolatey (пакетный менеджер)

Если у вас установлен Chocolatey, откройте командную строку от имени администратора и выполните:

```bash
choco install git
```

### Проверка установки

Откройте командную строку (cmd) или PowerShell и выполните:

```bash
git --version
```

Вы должны увидеть вывод вида: `git version 2.xx.x.windows.1`

## Установка на macOS

### Способ 1: Установка через Homebrew

1. Убедитесь, что у вас установлен Homebrew. Если нет — установите его:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

2. Установите Git:

```bash
brew install git
```

### Способ 2: Установка через официальный установщик

1. Перейдите на официальный сайт Git: https://git-scm.com/download/mac
2. Скачайте установщик `.pkg`
3. Запустите скачанный файл
4. Следуйте инструкциям мастера установки
5. После установки перезагрузите терминал

### Способ 3: Установка через Xcode Command Line Tools

Если у вас установлен Xcode, вы можете установить Git через командные инструменты:

```bash
xcode-select --install
```

Затем следуйте инструкциям на экране.

### Проверка установки

Откройте терминал и выполните:

```bash
git --version
```

Вы должны увидеть вывод вида: `git version 2.xx.x`

## Установка на Linux

### Ubuntu/Debian

```bash
sudo apt update
sudo apt install git
```

### CentOS/RHEL/Fedora

```bash
sudo yum install git
```

Или для новых версий:

```bash
sudo dnf install git
```

### Arch Linux

```bash
sudo pacman -S git
```

### Проверка установки

Откройте терминал и выполните:

```bash
git --version
```

Вы должны увидеть вывод вида: `git version 2.xx.x`

## Настройка Git после установки

После установки Git необходимо настроить ваше имя и email для идентификации коммитов:

```bash
git config --global user.name "Ваше Имя Фамилия"
git config --global user.email "your.email@mai.ru"
```

Для удобства работы рекомендуется также настроить цветной вывод и алиасы:

```bash
git config --global color.ui auto
git config --global alias.st status
git config --global alias.co checkout
git config --global alias.ci commit
git config --global alias.br branch
git config --global alias.tree "log --oneline --graph --decorate --all"
```

## Дополнительные ресурсы

- [Официальная документация Git](https://git-scm.com/doc)
- [Документация GitLab](https://docs.gitlab.com/)
## Установка на Windows (дополнительно)

### Способ 3: Установка через Git Bash

1. Скачайте установщик с официального сайта: https://gitforwindows.org/
2. Запустите установщик и следуйте инструкциям
3. В процессе установки:
   - Выберите редактор по умолчанию (рекомендуется использовать Notepad++ или Visual Studio Code)
   - Настройте PATH: выберите "Git from Windows Command Prompt"
   - Выберите SSH-клиент: "Use OpenSSH"
   - Выберите способ обработки окончаний строк: "Checkout Windows-style, commit Unix-style"
   - Выберите терминал: "Use MinTTY"
4. Завершите установку
- [Git Book на русском языке](https://git-scm.com/book/ru/v2)