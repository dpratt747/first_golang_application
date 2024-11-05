<!-- First steps -->

### Need to add a db migration tool

options are:

-   Golang-migrate
-   [Goose](https://github.com/pressly/goose)
-   Atlas

Going to use Goose for this project:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

```bash
go get github.com/pressly/goose/v3@latest
```

```bash
cd migrations
```

```bash
goose -s create new_user_table.sql
```

---

```bash
go mod init your_module_name
```

```bash
go get github.com/golang-migrate/migrate/v4
```

`go.sum` Is generated and updated automatically. It records the expected cryptographic checksums of the content of specific module versions, ensuring that future downloads of these modules are consistent and secure

## Bootstrap tool

[godev](https://github.com/zephinzer/godev)

Installation:

-   git clone https://github.com/zephinzer/godev.git
-   cd godev
-   go build -o $(go env GOPATH)/bin/godev
-   $(go env GOPATH)/bin/godev --version

Running Go-blueprint:

-   go-blueprint create --name db_access --framework gin --driver postgres

---

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

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

## Running Individual Tests:

All integration tests:

```bash
make itest
```

```bash
go test -v -run TestInsertNewUser ./integration_tests/database_test.go
```

Add postgres driver:

```bash
go get github.com/lib/pq
```

## Database migrations

```bash
goose -dir ./migrations postgres "user=postgres password=postgres port=5432 host=localhost dbname=golang_db sslmode=disable" up
```

## Undo migrations

```bash
goose -dir ./migrations postgres "user=postgres password=postgres port=5432 host=localhost dbname=golang_db sslmode=disable" down-to 0
```