package enums

type Postgres string

const (
	Port     Postgres = "6432"
	Host     Postgres = "localhost"
	Username Postgres = "postgres"
	Password Postgres = "postgres"
	Database Postgres = "golang_db"
)
