package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SawitProRecruitment/UserService/constant"
)

func (r *Repository) SaveUser(ctx context.Context, input SaveUserInput) (userID int64, err error) {
	query := "INSERT INTO users (name, phone, password) VALUES ($1, $2, $3) RETURNING id"

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, input.Name, input.Phone, input.Password).Scan(&userID)
	if err != nil {
		return
	}

	return userID, nil
}

func (r *Repository) GetUserByID(ctx context.Context, input GetUserByIDInput) (output GetUserByIDOutput, err error) {
	query := "SELECT id, name, phone FROM users WHERE id = $1"

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, input.ID).Scan(&output.ID, &output.Name, &output.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = constant.NotFoundErr
		}
		return
	}
	return
}

func (r *Repository) GetUserByPhone(ctx context.Context, input GetUserByPhoneInput) (output GetUserByPhoneOutput, err error) {
	query := "SELECT id, name, phone, password FROM users WHERE phone = $1"

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, input.Phone).Scan(&output.ID, &output.Name, &output.Phone, &output.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = constant.NotFoundErr
		}
		return
	}
	return
}

func (r *Repository) UpdateUserByID(ctx context.Context, input UpdateUserByIDInput) (err error) {
	query := "UPDATE users SET name = $2, phone = $3 WHERE id = $1"

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, input.ID, input.Name, input.Phone)
	return
}

func (r *Repository) IncreaseUserLoginCounterByID(ctx context.Context, input IncreaseUserLoginCounterByIDInput) (err error) {
	query := "UPDATE users SET login_count = login_count + 1 WHERE id = $1"

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, input.ID)
	return
}
