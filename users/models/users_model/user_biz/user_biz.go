package user_biz

import (
	"4vndevs/users/models/users_model/user_dot"
	"context"
	"errors"
)

type UserStorage interface {
	CreateUser(ctx context.Context, data *user_dot.UserDOT) (*user_dot.UserResponseDOT, error)
}
type usersModel struct {
	db UserStorage
}

func NewUsersModel(storage UserStorage) *usersModel {
	return &usersModel{db: storage}
}

func (model *usersModel) CreateNewUser(ctx context.Context, data *user_dot.UserDOT) (*user_dot.UserResponseDOT, error) {
	if data.Email == "" || data.Password == "" {
		return nil, errors.New("title can not be blank")
	}

	userResponse, err := model.db.CreateUser(ctx, data)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}
