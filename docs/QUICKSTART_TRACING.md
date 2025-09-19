# Быстрый старт с трейсингом

## Запуск с трейсингом

### 1. Запуск инфраструктуры
```bash
# Запуск всех сервисов включая трейсинг
docker-compose up -d

# Проверка статуса
docker-compose ps
```

### 2. Настройка переменных окружения
```bash
# Скопируйте пример конфигурации
cp .env.example .env

# Отредактируйте .env файл - установите TELEMETRY_ENABLED=true
echo "TELEMETRY_ENABLED=true" >> .env
echo "OTLP_ENDPOINT=http://localhost:4318" >> .env
echo "FACEIT_API_KEY=your_api_key_here" >> .env
```

### 3. Запуск приложения
```bash
# Сборка
go build -o faceit-cli main.go

# Запуск с трейсингом
./faceit-cli
```

## Просмотр трейсов

### Zipkin UI
- URL: http://localhost:9411
- Нажмите "Run Query" для просмотра всех трейсов
- Фильтруйте по сервису "faceit-cli"

### Jaeger UI (опционально)
```bash
# Запуск с Jaeger
docker-compose --profile jaeger up -d

# URL: http://localhost:16686
```

## Отключение трейсинга

```bash
# В .env файле
TELEMETRY_ENABLED=false

# Или через переменную окружения
TELEMETRY_ENABLED=false ./faceit-cli
```

## Что трейсится

- **app.run** - Основной трейс выполнения приложения
- **app.init_ui** - Инициализация UI
- **app.tui_execution** - Выполнение TUI программы
- **repository.get_player_by_nickname** - Поиск игрока по никнейму
- **repository.get_player_stats** - Получение статистики игрока
- **repository.get_player_recent_matches** - Получение последних матчей
- **repository.get_match_stats** - Получение статистики матча

## Полезные команды

```bash
# Просмотр логов OpenTelemetry Collector
docker-compose logs otel-collector

# Просмотр логов Zipkin
docker-compose logs zipkin

# Остановка всех сервисов
docker-compose down
```
