<!-- First steps -->

### Need to add a db migration tool

options are:

-   Golang-migrate
-   Goose
-   Atlas

Going to use golang-migrate for this project:

-   Create a go.mod file for dependencies. Equivilant is build.sbt/requirements file

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
