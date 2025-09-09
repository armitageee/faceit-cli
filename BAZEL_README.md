# Bazel Build System для faceit-cli

Этот проект настроен для сборки с использованием Bazel - современной системы сборки от Google.

## Установка

### Bazelisk (рекомендуется)
```bash
# Установка через Homebrew
brew install bazelisk

# Или скачать с GitHub
# https://github.com/bazelbuild/bazelisk/releases
```

### Bazel (альтернатива)
```bash
# Установка через Homebrew
brew install bazel

# Или скачать с GitHub
# https://github.com/bazelbuild/bazel/releases
```

## Основные команды

### Сборка проекта
```bash
# Сборка основного бинарника
bazelisk build //:faceit-cli

# Сборка всех целей
bazelisk build //...
```

### Запуск тестов
```bash
# Запуск всех тестов
bazelisk test //...

# Запуск тестов конкретного пакета
bazelisk test //internal/logger:logger_test
```

### Очистка кэша
```bash
# Очистка кэша
bazelisk clean

# Полная очистка (включая загруженные зависимости)
bazelisk clean --expunge
```

### Управление зависимостями
```bash
# Обновление зависимостей из go.mod
bazelisk run //:gazelle -- update-repos -from_file=go.mod

# Обновление BUILD файлов
bazelisk run //:gazelle
```

## Структура проекта

```
faceit-cli/
├── WORKSPACE              # Конфигурация Bazel workspace
├── BUILD.bazel           # Основной BUILD файл
├── .bazelrc              # Настройки Bazel
├── internal/             # Внутренние пакеты
│   ├── app/
│   │   └── BUILD.bazel   # BUILD файл для app пакета
│   ├── cache/
│   │   └── BUILD.bazel   # BUILD файл для cache пакета
│   ├── config/
│   │   └── BUILD.bazel   # BUILD файл для config пакета
│   ├── entity/
│   │   └── BUILD.bazel   # BUILD файл для entity пакета
│   ├── logger/
│   │   └── BUILD.bazel   # BUILD файл для logger пакета
│   ├── repository/
│   │   └── BUILD.bazel   # BUILD файл для repository пакета
│   └── ui/
│       └── BUILD.bazel   # BUILD файл для ui пакета
└── go.mod                # Go модули (используется Gazelle)
```

## Конфигурация

### WORKSPACE
Основной файл конфигурации Bazel workspace, содержит:
- Настройку правил Go (rules_go)
- Настройку Gazelle для управления зависимостями
- Регистрацию Go toolchain

### .bazelrc
Файл с настройками Bazel:
- Режим компиляции (optimized)
- Настройки тестирования
- Настройки Go (GOPROXY, GOSUMDB)

### BUILD файлы
Каждый пакет имеет свой BUILD.bazel файл, который определяет:
- `go_library` - для библиотек
- `go_binary` - для исполняемых файлов
- `go_test` - для тестов
- Зависимости между пакетами

## Управление зависимостями

Проект использует Gazelle для автоматического управления зависимостями:

1. **Добавление новой зависимости**:
   ```bash
   go get github.com/example/package
   bazelisk run //:gazelle -- update-repos -from_file=go.mod
   bazelisk run //:gazelle
   ```

2. **Обновление зависимостей**:
   ```bash
   go mod tidy
   bazelisk run //:gazelle -- update-repos -from_file=go.mod
   bazelisk run //:gazelle
   ```

## Отладка

### Подробный вывод
```bash
# Подробный вывод сборки
bazelisk build //:faceit-cli --verbose_failures

# Подробный вывод тестов
bazelisk test //... --verbose_failures
```

### Анализ зависимостей
```bash
# Показать все зависимости
bazelisk query --output=graph //:faceit-cli

# Показать зависимости конкретного пакета
bazelisk query --output=graph //internal/logger:logger
```

### Очистка и пересборка
```bash
# При проблемах с кэшем
bazelisk clean --expunge
bazelisk build //:faceit-cli
```

## Преимущества Bazel

1. **Инкрементальная сборка** - пересобираются только измененные части
2. **Кэширование** - результаты сборки кэшируются
3. **Параллельная сборка** - автоматическое распараллеливание
4. **Воспроизводимость** - одинаковые результаты на разных машинах
5. **Масштабируемость** - поддержка больших монорепозиториев

## Troubleshooting

### Проблемы с зависимостями
```bash
# Очистить кэш и пересобрать
bazelisk clean --expunge
bazelisk run //:gazelle -- update-repos -from_file=go.mod
bazelisk run //:gazelle
bazelisk build //:faceit-cli
```

### Проблемы с Go версией
Убедитесь, что в WORKSPACE указана правильная версия Go:
```python
go_register_toolchains(version = "1.23.0")
```

### Проблемы с сетью
Если есть проблемы с загрузкой зависимостей, проверьте настройки в .bazelrc:
```
build --action_env=GOPROXY=direct
build --action_env=GOSUMDB=off
```

## Интеграция с IDE

### VS Code
Установите расширение "Bazel" для поддержки синтаксиса и автодополнения.

### GoLand/IntelliJ
Установите плагин "Bazel" для поддержки Bazel проектов.

## Дополнительные ресурсы

- [Официальная документация Bazel](https://bazel.build/)
- [Bazel Go правила](https://github.com/bazelbuild/rules_go)
- [Gazelle документация](https://github.com/bazelbuild/bazel-gazelle)
- [Bazelisk документация](https://github.com/bazelbuild/bazelisk)
