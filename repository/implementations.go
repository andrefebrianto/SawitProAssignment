package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SawitProRecruitment/UserService/constant"
)

func (r *Repository) SaveUser(ctx context.Context, input SaveUserInput) (int64, error) {
	var userID int64
	err := r.Db.QueryRowContext(ctx, "INSERT INTO users (name, phone, password) VALUES ($1, $2, $3) RETURNING id", input.Name, input.Phone, input.Password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *Repository) GetUserByID(ctx context.Context, input GetUserByIDInput) (output GetUserByIDOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, name, phone FROM users WHERE id = $1", input.ID).Scan(&output.ID, &output.Name, &output.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = constant.NotFoundErr
		}
		return
	}
	return
}

func (r *Repository) GetUserByPhone(ctx context.Context, input GetUserByPhoneInput) (output GetUserByPhoneOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, name, phone, password FROM users WHERE phone = $1", input.Phone).Scan(&output.ID, &output.Name, &output.Phone, &output.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = constant.NotFoundErr
		}
		return
	}
	return
}

func (r *Repository) UpdateUserByID(ctx context.Context, input UpdateUserByIDInput) error {
	_, err := r.Db.ExecContext(ctx, "UPDATE users SET name = $2, phone = $3 WHERE id = $1", input.ID, input.Name, input.Phone)
	return err
}

func (r *Repository) IncreaseUserLoginCounterByID(ctx context.Context, input IncreaseUserLoginCounterByIDInput) error {
	_, err := r.Db.ExecContext(ctx, "UPDATE users SET login_count = login_count + 1 WHERE id = $1", input.ID)
	return err
}
