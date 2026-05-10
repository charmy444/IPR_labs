# Базовая конфигурация Git

## Введение

После установки Git необходимо выполнить первоначальную настройку. Эти настройки сохраняются в конфигурационных файлах и используются для идентификации автора коммитов и настройки поведения Git.

## Уровни конфигурации Git

Git предоставляет три уровня конфигурации:

1. **--system** — настройки для всей системы (все пользователи)
2. **--global** — настройки для текущего пользователя
3. **--local** — настройки для конкретного репозитория

Настройки применяются в порядке приоритета: local > global > system.

## Основные настройки

### Имя пользователя и email

Эти параметры обязательны для идентификации автора коммитов:

```bash
git config --global user.name "Ваше Имя Фамилия"
git config --global user.email "your.email@mai.ru"
```

### Текстовый редактор

Выберите текстовый редактор для написания сообщений коммитов и описаний Merge Request:

```bash
# Для Visual Studio Code
git config --global core.editor "code --wait"

# Для Sublime Text
git config --global core.editor "subl -n -w"

# Для Vim (по умолчанию)
git config --global core.editor "vim"
```

### Настройка цветов

Включите цветной вывод для лучшей читаемости:

```bash
git config --global color.ui auto
```

## Просмотр текущих настроек

Чтобы просмотреть все текущие настройки:

```bash
git config --list
```

Для просмотра настроек определенного уровня:

```bash
# Глобальные настройки
git config --global --list

# Локальные настройки (в текущем репозитории)
git config --local --list
```

Для просмотра значения конкретного параметра:

```bash
git config user.name
```

## Полезные дополнительные настройки

### Настройка переносов строк

Для корректной работы с различными операционными системами:

```bash
# На Windows
git config --global core.autocrlf true

# На Mac/Linux
git config --global core.autocrlf input
```

### Настройка ветки по умолчанию

Установите имя ветки по умолчанию для новых репозиториев:

```bash
git config --global init.defaultBranch main
```

### Настройка pager для вывода

Если вывод команд Git не помещается на экран, можно настроить pager:

```bash
git config --global core.pager "less -R"
```

### Алиасы для часто используемых команд

Создайте алиасы для ускорения работы:

```bash
# Сокращение для status
git config --global alias.st status

# Сокращение для checkout
git config --global alias.co checkout

# Сокращение для commit
git config --global alias.ci commit

# Сокращение для branch
git config --global alias.br branch

# Просмотр дерева коммитов
git config --global alias.tree "log --oneline --graph --decorate --all"
```

После создания этих алиасов вы можете использовать сокращенные команды:

```bash
git st    # вместо git status
git co    # вместо git checkout
git ci    # вместо git commit
git br    # вместо git branch
git tree  # для просмотра дерева коммитов
```

## Настройка файла .gitconfig

Все глобальные настройки сохраняются в файле `~/.gitconfig`. Вы можете редактировать этот файл напрямую для более сложных настроек:

```ini
[user]
  name = Ваше Имя Фамилия
  email = your.email@mai.ru

[core]
  editor = code --wait
  autocrlf = input
  pager = less -R

[init]
  defaultBranch = main

[color]
  ui = auto

[alias]
  st = status
  co = checkout
  ci = commit
  br = branch
  tree = log --oneline --graph --decorate --all
```

## Проверка конфигурации

После настройки проверьте, что все параметры установлены правильно:

```bash
git config --global user.name
git config --global user.email
git config --global core.editor
```

## Основные ссылки

- [Официальная документация Git по настройке](https://git-scm.com/book/ru/v2/Введение-Первоначальная-настройка-Git)
- [Документация Git Config](https://git-scm.com/docs/git-config)
- [Git Book на русском языке](https://git-scm.com/book/ru/v2)