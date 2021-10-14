package outlet_dao

import (
	"context"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/rest_err"
)

type OutletDaoAssumer interface {
	OutletSaver
	OutletLoader
}

type OutletSaver interface {
	Insert(ctx context.Context, input dto.OutletModel) (int, rest_err.APIError)
	Edit(ctx context.Context, input dto.OutletEditModel) (*dto.OutletModel, rest_err.APIError)
	Delete(ctx context.Context, id int, filterMerchant int) rest_err.APIError
}

type OutletLoader interface {
	Get(ctx context.Context, id int) (*dto.OutletModel, rest_err.APIError)
	FindWithPagination(ctx context.Context, opt FindParams) ([]dto.OutletModel, rest_err.APIError)
}
