package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/gomega"
)

var (
	errDBClosed = errors.New("sql: database is closed")
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestRepository_SaveUser(t *testing.T) {
	RegisterTestingT(t)

	gofakeit.Seed(time.Now().UnixNano())
	var input repository.SaveUserInput
	gofakeit.Struct(input)
	query := regexp.QuoteMeta("INSERT INTO users (name, phone, password) VALUES ($1, $2, $3) RETURNING id")

	t.Run("Failed to prepare staement", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		mock.ExpectPrepare(query).WillReturnError(errDBClosed)

		result, err := repo.SaveUser(ctx, input)
		Expect(result).Should(Equal(int64(0)))
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Failed to execute query", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)
		prep.ExpectQuery().WithArgs(input.Name, input.Phone, input.Password).WillReturnError(errDBClosed)

		result, err := repo.SaveUser(ctx, input)
		Expect(result).Should(Equal(int64(0)))
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Saved the data successfully", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)

		lastId := int64(1)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(lastId)
		prep.ExpectQuery().WithArgs(input.Name, input.Phone, input.Password).WillReturnRows(rows)

		result, err := repo.SaveUser(ctx, input)
		Expect(err).Should(BeNil())
		Expect(result).Should(Equal(lastId))
	})
}

func TestRepository_GetUserByID(t *testing.T) {
	RegisterTestingT(t)

	gofakeit.Seed(time.Now().UnixNano())
	var input repository.GetUserByIDInput
	gofakeit.Struct(input)
	query := regexp.QuoteMeta("SELECT id, name, phone FROM users WHERE id = $1")

	t.Run("Failed to prepare staement", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		mock.ExpectPrepare(query).WillReturnError(errDBClosed)

		result, err := repo.GetUserByID(ctx, input)
		Expect(result).Should(Equal(repository.GetUserByIDOutput{}))
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Failed to execute query", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)
		prep.ExpectQuery().WithArgs(input.ID).WillReturnError(errDBClosed)

		result, err := repo.GetUserByID(ctx, input)
		Expect(result).Should(Equal(repository.GetUserByIDOutput{}))
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Retrieved the data successfully", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)

		var output repository.GetUserByIDOutput
		gofakeit.Struct(&output)

		rows := sqlmock.NewRows([]string{"id", "name", "phone"}).AddRow(output.ID, output.Name, output.Phone)
		prep.ExpectQuery().WithArgs(input.ID).WillReturnRows(rows)

		result, err := repo.GetUserByID(ctx, input)
		Expect(err).Should(BeNil())
		Expect(result).Should(Equal(output))
	})
}

func TestRepository_GetUserByPhone(t *testing.T) {
	RegisterTestingT(t)

	gofakeit.Seed(time.Now().UnixNano())
	var input repository.GetUserByPhoneInput
	gofakeit.Struct(input)
	query := regexp.QuoteMeta("SELECT id, name, phone, password FROM users WHERE phone = $1")

	t.Run("Failed to prepare staement", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		mock.ExpectPrepare(query).WillReturnError(errDBClosed)

		result, err := repo.GetUserByPhone(ctx, input)
		Expect(result).Should(Equal(repository.GetUserByPhoneOutput{}))
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Failed to execute query", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)
		prep.ExpectQuery().WithArgs(input.Phone).WillReturnError(errDBClosed)

		result, err := repo.GetUserByPhone(ctx, input)
		Expect(result).Should(Equal(repository.GetUserByPhoneOutput{}))
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Retrieved the data successfully", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)

		var output repository.GetUserByPhoneOutput
		gofakeit.Struct(&output)

		rows := sqlmock.NewRows([]string{"id", "name", "phone", "password"}).AddRow(output.ID, output.Name, output.Phone, output.Password)
		prep.ExpectQuery().WithArgs(input.Phone).WillReturnRows(rows)

		result, err := repo.GetUserByPhone(ctx, input)
		Expect(err).Should(BeNil())
		Expect(result).Should(Equal(output))
	})
}

func TestRepository_UpdateUserByID(t *testing.T) {
	RegisterTestingT(t)

	gofakeit.Seed(time.Now().UnixNano())
	var input repository.UpdateUserByIDInput
	gofakeit.Struct(input)
	query := regexp.QuoteMeta("UPDATE users SET name = $2, phone = $3 WHERE id = $1")

	t.Run("Failed to prepare staement", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		mock.ExpectPrepare(query).WillReturnError(errDBClosed)

		err := repo.UpdateUserByID(ctx, input)
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Failed to execute query", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(input.ID, input.Name, input.Phone).WillReturnError(errDBClosed)

		err := repo.UpdateUserByID(ctx, input)
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Retrieved the data successfully", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)

		var output repository.GetUserByPhoneOutput
		gofakeit.Struct(&output)

		prep.ExpectExec().WithArgs(input.ID, input.Name, input.Phone).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserByID(ctx, input)
		Expect(err).Should(BeNil())
	})
}

func TestRepository_IncreaseUserLoginCounterByID(t *testing.T) {
	RegisterTestingT(t)

	gofakeit.Seed(time.Now().UnixNano())
	var input repository.IncreaseUserLoginCounterByIDInput
	gofakeit.Struct(input)
	query := regexp.QuoteMeta("UPDATE users SET login_count = login_count + 1 WHERE id = $1")

	t.Run("Failed to prepare staement", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		mock.ExpectPrepare(query).WillReturnError(errDBClosed)

		err := repo.IncreaseUserLoginCounterByID(ctx, input)
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Failed to execute query", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(input.ID).WillReturnError(errDBClosed)

		err := repo.IncreaseUserLoginCounterByID(ctx, input)
		Expect(err).Should(Equal(errDBClosed))
	})

	t.Run("Retrieved the data successfully", func(t *testing.T) {
		ctx := context.Background()
		db, mock := NewMock()
		repo := repository.Repository{Db: db}

		prep := mock.ExpectPrepare(query)

		var output repository.GetUserByPhoneOutput
		gofakeit.Struct(&output)

		prep.ExpectExec().WithArgs(input.ID).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.IncreaseUserLoginCounterByID(ctx, input)
		Expect(err).Should(BeNil())
	})
}
