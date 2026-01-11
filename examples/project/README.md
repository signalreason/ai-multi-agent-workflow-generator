# Example Project: Note API

This is a minimal Go HTTP service used to test the multi-agent workflow manager.

## Run

```
go run ./cmd/noteapi
```

## Try it

```
curl -X POST http://localhost:8080/notes -d '{"body":"hello"}' -H 'Content-Type: application/json'
curl http://localhost:8080/notes
```

## Test

```
go test ./...
```
