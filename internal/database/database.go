package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"

	"db_access/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type DatabaseService interface {
	InsertNewUser(user domain.User) (int, error)

	GetAllUsers() ([]domain.User, error)

	SoftDeleteUser(userId int) error
}

type service struct {
	db *sql.DB
}

var (
	dbInstance *service
)

func New(connectionString string) DatabaseService {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) SoftDeleteUser(userId int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return &domain.DatabaseTransactionError{Message: err.Error()}
	}
	statement := "INSERT INTO user_deletes(user_id) VALUES($1)"

	query, err := tx.Prepare(statement)
	if err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("Failed to prepare the SQL statement. [Reason]: %v", err)
		log.Println(errorMessage)
		return err
	}

	_, err = query.Exec(userId)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				log.Println("Unique constraint violation:", pqErr.Message)
				return &domain.UniqueConstraintDatabaseError{Message: pqErr.Message}
			case "23503":
				log.Println("User does not exist cannot delete", pqErr.Message)
				return &domain.UniqueConstraintDatabaseError{Message: pqErr.Message}
			default:
				log.Println("Database error:", pqErr.Code.Name())
				return &domain.UnmappedDatabaseError{Message: pqErr.Message}
			}
		}
		return &domain.UnmappedDatabaseError{Message: err.Error()}

	}
	defer query.Close()

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("Failed to commit the prepared SQL statement. [Reason]: %v", err)
		log.Println(errorMessage)
		return &domain.DatabaseTransactionError{Message: err.Error()}
	}

	log.Println("SQL query:", statement)

	return nil
}

func (s *service) GetAllUsers() ([]domain.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, &domain.DatabaseTransactionError{Message: err.Error()}
	}

	statement :=
		`
	SELECT u.id, u.username, u.email
	FROM users u 
	LEFT JOIN user_deletes ud on u.id = ud.user_id
	WHERE ud.user_id is NULL
	`

	query, err := tx.Prepare(statement)
	if err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("Failed to prepare the SQL statement. [Reason]: %v", err)
		log.Println(errorMessage)
		return nil, &domain.DatabaseTransactionError{Message: err.Error()}
	}

	rows, err := query.Query()
	if err != nil {
		return nil, &domain.UnmappedDatabaseError{Message: err.Error()}
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		users = append(users, user)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("Failed to commit the prepared SQL statement. [Reason]: %v", err)
		log.Println(errorMessage)
		return nil, &domain.DatabaseTransactionError{Message: err.Error()}
	}

	log.Println("SQL query:", statement)
	return users, nil
}

func (s *service) InsertNewUser(user domain.User) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, &domain.DatabaseTransactionError{Message: err.Error()}
	}

	statement := "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id"

	query, err := tx.Prepare(statement)
	if err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("Failed to prepare the SQL statement. [Reason]: %v", err)
		log.Println(errorMessage)
		return 0, &domain.DatabaseTransactionError{Message: err.Error()}
	}
	defer query.Close()

	err = query.QueryRow(user.Username, user.Email).Scan(&user.ID)
	if err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("Failed to execute the prepared SQL statement. [Reason]: %v", err)
		log.Println(errorMessage)

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				log.Println("Unique constraint violation:", pqErr.Message)
				return 0, &domain.UniqueConstraintDatabaseError{Message: pqErr.Message}
			default:
				log.Println("Database error:", pqErr.Code.Name())
				return 0, &domain.UnmappedDatabaseError{Message: pqErr.Message}
			}
		}
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("Failed to commit the prepared SQL statement. [Reason]: %v", err)
		log.Println(errorMessage)
		return 0, &domain.DatabaseTransactionError{Message: err.Error()}
	}

	log.Println("SQL query:", statement)

	return user.ID, nil
}
