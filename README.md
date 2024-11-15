<!-- First steps -->

## <ins>Running locally</ins>

```bash
docker-compose up -d postgres
```

```bash
goose -dir ./migrations postgres "user=postgres password=postgres port=6432 host=localhost dbname=golang_db sslmode=disable" up
```

```bash
make run
```

### <ins>Undo migrations</ins>

```bash
goose -dir ./migrations postgres "user=postgres password=postgres port=6432 host=localhost dbname=golang_db sslmode=disable" down-to 0
```

---

## <ins>Set up</ins>

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

```bash
go get github.com/pressly/goose/v3@latest
```

```bash
go mod init <your_module_name>
```

```bash
cd migrations
```

```bash
goose -s create new_user_table.sql
```

`go.sum` Is generated and updated automatically. It records the expected cryptographic checksums of the content of specific module versions, ensuring that future downloads of these modules are consistent and secure

## Bootstrap tool

[godev](https://github.com/zephinzer/godev)

Installation:

```bash
git clone https://github.com/zephinzer/godev.git
cd godev
```

-   godev init <project_name>
-   godev run

---

## <ins>MakeFile</ins>

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```

---

## <ins>Running Individual Tests:</ins>

### Unit tests:

```bash
go test -v -run <Test Name> ./tests/...
```

e.g.

```bash
go test -v -run TestHelloWorldHandler ./tests/...
```

### Integration tests:

```bash
go test -v -run <Test Name> ./integration_tests/...
```

e.g.

```bash
go test -v -run TestInsertNewUser ./integration_tests/...
```

### Running all unit tests:

```bash
make test
```

### Running all integration tests:

```bash
make itest
```
