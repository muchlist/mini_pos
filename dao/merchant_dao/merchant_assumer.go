package merchant_dao

import (
	"context"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/rest_err"
)

type MerchantDaoAssumer interface {
	MerchantSaver
	MerchantLoader
}

type MerchantSaver interface {
	Insert(ctx context.Context, input dto.MerchantCreateReq) (*dto.MerchantCreateRes, rest_err.APIError)
	Edit(ctx context.Context, input dto.Merchant) (*dto.Merchant, rest_err.APIError)
	Delete(ctx context.Context, id int) rest_err.APIError
}

type MerchantLoader interface {
	Get(ctx context.Context, id int) (*dto.Merchant, rest_err.APIError)
	FindWithCursor(ctx context.Context, opt FindParams) ([]dto.Merchant, rest_err.APIError)
}
