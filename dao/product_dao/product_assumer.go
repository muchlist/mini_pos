package product_dao

import (
	"context"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/rest_err"
)

type ProductDaoAssumer interface {
	ProductSaver
	ProductLoader
}

type ProductSaver interface {
	Insert(ctx context.Context, input dto.ProductModel) (int, rest_err.APIError)
	Edit(ctx context.Context, input dto.ProductEditModel) (*dto.ProductModel, rest_err.APIError)
	Delete(ctx context.Context, id int, filterMerchant int) rest_err.APIError
	EditCustomPrice(ctx context.Context, input dto.ProductPriceModel) (*dto.ProductModel, rest_err.APIError)
	InsertCustomPrice(ctx context.Context, input dto.ProductPriceModel) (*dto.ProductModel, rest_err.APIError)
	SetImagePath(ctx context.Context, productID int, path string) (*dto.ProductModel, rest_err.APIError)
}

type ProductLoader interface {
	Get(ctx context.Context, id int, merchantFilter int) (*dto.ProductModel, rest_err.APIError)
	GetWithCustomPriceOutlet(ctx context.Context, id int, outletID int) (*dto.ProductModel, rest_err.APIError)
	GetPriceDataWithID(ctx context.Context, priceID string) (*dto.ProductPriceModel, rest_err.APIError)
	FindWithPagination(ctx context.Context, opt FindParams, merchantFilter int) ([]dto.ProductModel, rest_err.APIError)
	FindCustomPriceOutlet(ctx context.Context, outletID int) ([]dto.ProductPriceModel, rest_err.APIError)
}
