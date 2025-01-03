# Arithmetic evaluator with a RESTful web API

> Русский [README](README.ru.md)

## Evaluator (calculator)

### Calabilities:

-   Uses +-\\*/ and () in PEMDAS order
-   Checks for syntax
-   Ignores non-syntax-breaking spaces (52 - 24 is ok, 52 - 2 4 is not)

### Limitations:

-   No exponents or other functions

## Running

-   `go run cmd/main.go` after installing dependencies in go.mod (gin)

### Testing

-   `go test internal/application/calcapi_test.go`
-   `go test pkg/calcapi/calculator_test.go`

### API

- RESTful API with JSON

#### "api/v1/calculate/" - POST

- In: {"expression": "..."}
- Out: {"result": ...}
- Error: 422 or 500 {"error": "..."}
- Example curl:

```bash
curl -L -X POST 'http://localhost:80/api/v1/calculate'\
-H 'Content-Type: application/json' -d '{"expression":"2+2"}'
```

TG: @neo536
