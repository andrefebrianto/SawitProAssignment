// This file contains types that are used in the repository layer.
package repository

type SaveUserInput struct {
	Name     string
	Phone    string
	Password string
}

type GetUserByIDInput struct {
	ID int64
}

type GetUserByIDOutput struct {
	ID    int64
	Name  string
	Phone string
}

type GetUserByPhoneInput struct {
	Phone string
}

type GetUserByPhoneOutput struct {
	ID       int64
	Name     string
	Phone    string
	Password string
}

type UpdateUserByIDInput struct {
	ID    int64
	Name  string
	Phone string
}

type IncreaseUserLoginCounterByIDInput struct {
	ID int64
}
