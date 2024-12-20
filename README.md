## <ins>Running locally</ins>

```bash
docker-compose up -d db
```

```bash
goose -dir ./migrations postgres "user=postgres password=postgres port=5432 host=localhost dbname=golang_db sslmode=disable" up
```

```bash
make run
```

---

```bash
make docker-up
make migrate-up
```

and to close things:

```bash
make docker-down
```

---

## Request Examples:

### Insert User:

```bash
curl --request POST \
  --url http://127.0.0.1:8080/user \
  --header 'Content-Type: application/json' \
  --data '{
	"username": "1",
	"email": "11211@email.com"
}'
```

### Get All Users:

```bash
curl --request GET \
  --url http://127.0.0.1:8080/users \
  --header 'Content-Type: application/json'
```

### Delete User:

```bash
curl --request DELETE \
  --url http://127.0.0.1:8080/user/2 \
  --header 'Content-Type: application/json'
```

---

### <ins>Undo migrations</ins>

```bash
goose -dir ./migrations postgres "user=postgres password=postgres port=5432 host=localhost dbname=golang_db sslmode=disable" down-to 0
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

goose -s create new_user_table sql
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
