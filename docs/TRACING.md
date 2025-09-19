# OpenTelemetry Tracing

Этот документ описывает настройку и использование OpenTelemetry для трейсинга в faceit-cli.

## Обзор

Приложение интегрировано с OpenTelemetry для сбора и отправки трейсов. Трейсы помогают отслеживать производительность и отлаживать проблемы в CLI утилите.

## Компоненты

### 1. OpenTelemetry Collector
- **Порт**: 4317 (gRPC), 4318 (HTTP)
- **Конфигурация**: `otel-collector-config.yaml`
- **Функция**: Собирает трейсы от приложения и отправляет их в Zipkin

### 2. Zipkin
- **Порт**: 9411
- **URL**: http://localhost:9411
- **Функция**: Визуализация трейсов

### 3. Jaeger (опционально)
- **Порт**: 16686
- **URL**: http://localhost:16686
- **Функция**: Альтернативная визуализация трейсов

## Запуск с трейсингом

### 1. Запуск инфраструктуры
```bash
# Запуск всех сервисов включая трейсинг
docker-compose up -d

# Или только трейсинг сервисы
docker-compose up -d otel-collector zipkin
```

### 2. Настройка переменных окружения
```bash
# Скопируйте пример конфигурации
cp .env.example .env

# Отредактируйте .env файл
# Убедитесь что TELEMETRY_ENABLED=true
```

### 3. Запуск приложения
```bash
# С трейсингом
TELEMETRY_ENABLED=true ./faceit-cli

# Или используйте .env файл
./faceit-cli
```

## Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `TELEMETRY_ENABLED` | Включить трейсинг | `false` |
| `OTLP_ENDPOINT` | OTLP HTTP endpoint (без /v1/traces) | `http://localhost:4318` |
| `ZIPKIN_ENDPOINT` | Zipkin endpoint | `http://localhost:9411/api/v2/spans` |
| `SERVICE_NAME` | Имя сервиса | `faceit-cli` |
| `SERVICE_VERSION` | Версия сервиса | `dev` |
| `ENVIRONMENT` | Окружение | `development` |
| `OTEL_LOG_LEVEL` | Уровень логов OpenTelemetry | `fatal` (подавляет OTLP логи) |

### Трейсы

Приложение создает следующие трейсы:

1. **app.run** - Основной трейс выполнения приложения
2. **app.init_ui** - Инициализация UI
3. **app.tui_execution** - Выполнение TUI программы
4. **repository.get_player_by_nickname** - Поиск игрока по никнейму
5. **repository.get_player_stats** - Получение статистики игрока
6. **repository.get_player_recent_matches** - Получение последних матчей
7. **repository.get_match_stats** - Получение статистики матча

## Просмотр трейсов

### Zipkin
1. Откройте http://localhost:9411
2. Нажмите "Run Query" для просмотра всех трейсов
3. Используйте фильтры для поиска конкретных трейсов

### Jaeger (если используется)
1. Запустите с профилем jaeger: `docker-compose --profile jaeger up -d`
2. Откройте http://localhost:16686
3. Выберите сервис "faceit-cli" и нажмите "Find Traces"

## Отладка

### Проверка статуса сервисов
```bash
# Проверка статуса контейнеров
docker-compose ps

# Логи OpenTelemetry Collector
docker-compose logs otel-collector

# Логи Zipkin
docker-compose logs zipkin
```

### Тестирование OTLP endpoint
```bash
# Проверка доступности OTLP HTTP endpoint
curl http://localhost:4318/v1/traces

# Проверка Zipkin API
curl http://localhost:9411/api/v2/services
```

## Производительность

Трейсинг добавляет минимальные накладные расходы:
- Время выполнения увеличивается на ~1-2ms на операцию
- Память: ~1-2MB дополнительно
- Сетевой трафик: ~1-5KB на трейс

## Отключение трейсинга

Для отключения трейсинга установите:
```bash
export TELEMETRY_ENABLED=false
```

Или в .env файле:
```
TELEMETRY_ENABLED=false
```

## Подавление OTLP логов

По умолчанию OTLP логи подавляются, чтобы не засорять stdout. Все трейсы отправляются только в Zipkin через OTLP коллектор.

Если нужно включить OTLP логи для отладки:
```bash
export OTEL_LOG_LEVEL=debug
```

Или в .env файле:
```
OTEL_LOG_LEVEL=debug
```
