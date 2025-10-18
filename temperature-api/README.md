# Temperature API (Go)

Простое приложение на Go для имитации работы удаленного датчика температуры.

## Описание

Temperature API - это простое веб-приложение на Go с использованием Gin framework, которое имитирует работу датчика температуры. При запросе возвращает случайные значения температуры для указанных местоположений.

## API Endpoints

### GET /temperature
Получить температуру для указанного местоположения или датчика.

**Параметры:**
- `location` (опционально) - название комнаты ("Living Room", "Bedroom", "Kitchen")
- `sensorId` (опционально) - ID датчика ("1", "2", "3")

**Примеры запросов:**
```
GET /temperature?location=Living Room
GET /temperature?sensorId=1
GET /temperature?location=Kitchen&sensorId=3
```

**Ответ:**
```json
{
  "sensorId": "1",
  "sensorType": "temperature",
  "location": "Living Room",
  "value": 22.5,
  "unit": "°C",
  "status": "Comfortable",
  "timestamp": "2024-01-15T10:30:00Z",
  "description": "Temperature in Living Room: 22.5°C (Comfortable)"
}
```

### GET /temperature/{sensorId}
Получить температуру по ID датчика.

**Пример:**
```
GET /temperature/2
```

### GET /temperature/health
Проверка здоровья сервиса.

**Ответ:**
```json
{
  "status": "ok",
  "message": "Temperature API is running!",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Логика mapping location/sensorId

Приложение автоматически определяет недостающие параметры:

- Если `location` не указан, используется значение на основе `sensorId`:
  - sensorId "1" → "Living Room"
  - sensorId "2" → "Bedroom"
  - sensorId "3" → "Kitchen"
  - другие → "Unknown"

- Если `sensorId` не указан, генерируется на основе `location`:
  - "Living Room" → sensorId "1"
  - "Bedroom" → sensorId "2"
  - "Kitchen" → sensorId "3"
  - другие → sensorId "0"

## Как запустить

### Локально (для разработки)
```bash
# Убедитесь, что Go установлен
go version

# Установите зависимости
go mod tidy

# Запустите приложение
go run main.go
```

Сервис будет доступен по адресу: http://localhost:8081

### Через Docker
```bash
# Соберите образ
docker build -t temperature-api .

# Запустите контейнер
docker run -p 8081:8081 temperature-api
```

### Через Docker Compose (вместе с другими сервисами)
```bash
# Убедитесь, что Docker Desktop запущен
cd apps
docker-compose up --build
```

## Тестирование API

### Используя curl
```bash
# Проверка здоровья
curl http://localhost:8081/temperature/health

# Получить температуру для Living Room
curl "http://localhost:8081/temperature?location=Living%20Room"

# Получить температуру для sensorId=1
curl "http://localhost:8081/temperature?sensorId=1"

# Получить температуру для sensorId=2
curl http://localhost:8081/temperature/2
```

### Используя Postman
Импортируйте коллекцию `smarthome-api.postman_collection.json` и используйте:
- Create Sensor
- Get All Sensors

Каждый вызов должен показывать разные значения температуры.

## Технические детали

- **Язык:** Go 1.21
- **Framework:** Gin
- **Порт:** 8081
- **Температурный диапазон:** 15-30°C
- **Статусы:** Cold (<18°C), Comfortable (18-25°C), Hot (>25°C)
- **CORS:** Разрешен для всех источников
- **Логирование:** Встроенное логирование Gin

## Структура проекта

```
temperature-api/
├── main.go          # Основной файл приложения
├── go.mod           # Go модуль
├── go.sum           # Зависимости
├── Dockerfile       # Docker конфигурация
└── README.md        # Документация
```

## Зависимости

- `github.com/gin-gonic/gin` - HTTP веб-фреймворк для Go