# E2E тестрвоание gRPC

Цель: проверь возможность end-to-end тестирования gRPC сервиса.

## Запуск E2E тестов

```bash
go test -tags=e2etest ./... -v
```

## gopls error

Чтобы не было ошибок от gopls достаточно просто указать в конфиге парамет buildFlags.

[Пример файла для zed](../.zed/settings.json)

```json
{
  "lsp": {
    "gopls": {
      "initialization_options": {
        "buildFlags": ["-tags=e2etest"],
      },
    },
  },
}
```
