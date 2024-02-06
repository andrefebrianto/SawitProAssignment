// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	SaveUser(ctx context.Context, input SaveUserInput) (int64, error)
	GetUserByID(ctx context.Context, input GetUserByIDInput) (output GetUserByIDOutput, err error)
	GetUserByPhone(ctx context.Context, input GetUserByPhoneInput) (output GetUserByPhoneOutput, err error)
	UpdateUserByID(ctx context.Context, input UpdateUserByIDInput) error
	IncreaseUserLoginCounterByID(ctx context.Context, input IncreaseUserLoginCounterByIDInput) error
}
