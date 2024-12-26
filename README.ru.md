# Арифметический вычислитель с RESTFUL веб-API

> English [README](README.md)

## Калькулятор

### Возможности:

- Использует +-\\*/ и () в порядке PEMDAS.
- Проверяет синтаксис
- Игнорирует пробелы, не нарушающие синтаксис (52 – 24 – нормально, 52 – 2 4 – нет)

### Ограничения:

- Никаких экспонентов или других функций.

## Запуск

-   `go run cmd/main.go` после установки зависимостей в go.mod (gin)

### Тесты

-   `go test internal/application/calcapi_test.go`
-   `go test pkg/calcapi/calculator_test.go`

### API

- RESTful API с JSON

#### "api/v1/calculate/" - POST

- Ввод: {"expression": "..."}
- Вывод: {"result": ...}
- Ошибка: 422 или 500 {"error": "..."}
- Примерный curl:

```bash
curl -L -X POST 'http://localhost:80/api/v1/calculate'\
-H 'Content-Type: application/json' -d '{"expression":"2+2"}'
```


TG: @neo536
