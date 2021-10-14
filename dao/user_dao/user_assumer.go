package user_dao

import (
	"context"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/rest_err"
)

type UserDaoAssumer interface {
	UserSaver
	UserReader
}

type UserSaver interface {
	Insert(ctx context.Context, user dto.UserModel) (string, rest_err.APIError)
	Edit(ctx context.Context, userInput dto.UserEditModel) (*dto.UserModel, rest_err.APIError)
	Delete(ctx context.Context, id int, filterMerchant int) rest_err.APIError
	ChangePassword(ctx context.Context, input dto.UserModel) rest_err.APIError
}

type UserReader interface {
	GetByID(ctx context.Context, id int) (*dto.UserModel, rest_err.APIError)
	GetByEmail(ctx context.Context, email string) (*dto.UserModel, rest_err.APIError)
	FindWithPagination(ctx context.Context, opt FindPaginationParams) ([]dto.UserModel, rest_err.APIError)
}
