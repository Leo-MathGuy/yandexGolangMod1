# Арифметический вычислитель с RESTFUL веб-API

> English [README](README.ru)

## Калькулятор

### Возможности:

- Использует +-\\*/ и () в порядке PEMDAS.
- Проверяет синтаксис
- Игнорирует пробелы, не нарушающие синтаксис (52 – 24 – нормально, 52 – 2 4 – нет)

### Ограничения:

- Никаких экспонентов или других функций.

### Запуск

-   `go run cmd/main.go` после установки зависимостей в go.mod (gin)

### Testing

-   `go test internal/application/calcapi_test.go`
-   `go test pkg/calcapi/calculator_test.go`

Поставь 100 пж
