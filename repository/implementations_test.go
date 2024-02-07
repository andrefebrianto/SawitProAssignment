package repository_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestRepository_SaveUser(t *testing.T) {
	t.Run("Failed to save user data", func(t *testing.T) {
		// ctx := context.Background()
		// db, mock := NewMock()
		// repo := repository.Repository{Db: db}

		// var input repository.SaveUserInput
		// gofakeit.Struct(input)

		// query := "INSERT INTO users (name, phone, password) VALUES ($1, $2, $3) RETURNING id"

		// // mock.ExpectQuery()
		// rows := sqlmock.Row

		// userID, err := repo.SaveUser(ctx, input)
	})

	t.Run("Saved the data successfully", func(t *testing.T) {

	})
}

func TestRepository_GetUserByID(t *testing.T) {

}

func TestRepository_GetUserByPhone(t *testing.T) {

}

func TestRepository_UpdateUserByID(t *testing.T) {

}

func TestRepository_IncreaseUserLoginCounterByID(t *testing.T) {

}

func TestRepository_Close(t *testing.T) {

}
