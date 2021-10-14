package product_serv

import (
	"context"
	"fmt"
	"github.com/muchlist/mini_pos/dao/product_dao"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/logger"
	"github.com/muchlist/mini_pos/utils/mjwt"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"time"
)

type ProductServiceAssumer interface {
	ProductServiceModifier
	ProductServiceReader
}

type ProductServiceReader interface {
	Get(ctx context.Context, productID int, outletID int) (*dto.ProductModel, rest_err.APIError)
	FindProducts(ctx context.Context, search string, limit int, offset int, outletSpecific int) ([]dto.ProductModel, rest_err.APIError)
}

type ProductServiceModifier interface {
	CreateProduct(ctx context.Context, claims mjwt.CustomClaim, product dto.ProductModel) (int, rest_err.APIError)
	EditProduct(ctx context.Context, claims mjwt.CustomClaim, request dto.ProductEditRequest) (*dto.ProductModel, rest_err.APIError)
	DeleteProduct(ctx context.Context, claims mjwt.CustomClaim, productID int) rest_err.APIError
	SetCustomPrice(ctx context.Context, claims mjwt.CustomClaim, price dto.ProductPriceRequest) (*dto.ProductModel, rest_err.APIError)
	SetImagePath(ctx context.Context, productID int, path string) (*dto.ProductModel, rest_err.APIError)
}

func NewProductService(dao product_dao.ProductDaoAssumer) ProductServiceAssumer {
	return &productService{
		dao: dao,
	}
}

type productService struct {
	dao product_dao.ProductDaoAssumer
}

// CreateProduct melakukan register product oleh akun owner
func (u *productService) CreateProduct(ctx context.Context, claims mjwt.CustomClaim, product dto.ProductModel) (int, rest_err.APIError) {

	timeNow := time.Now().Unix()
	product.CreatedAt = timeNow
	product.UpdatedAt = timeNow
	product.MerchantID = claims.Merchant // merchant ID adalah sama dengan merchant id owner

	productID, err := u.dao.Insert(ctx, product)
	if err != nil {
		return 0, err
	}
	return productID, nil
}

// EditProduct
func (u *productService) EditProduct(ctx context.Context, claims mjwt.CustomClaim, request dto.ProductEditRequest) (*dto.ProductModel, rest_err.APIError) {
	editParams := dto.ProductEditModel{
		WhereID:         request.ID,
		WhereMerchantID: claims.Merchant, // <--- product yang diedit harus memiliki merchant id yang sama dengan pengedit
		Code:            dto.UppercaseString(request.Code),
		Name:            dto.UppercaseString(request.Name),
		MasterBuyPrice:  request.MasterBuyPrice,
		MasterSellPrice: request.MasterSellPrice,
	}

	result, err := u.dao.Edit(ctx, editParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// EditProduct
func (u *productService) SetImagePath(ctx context.Context, productID int, path string) (*dto.ProductModel, rest_err.APIError) {
	result, err := u.dao.SetImagePath(ctx, productID, path)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteProduct
func (u *productService) DeleteProduct(ctx context.Context, claims mjwt.CustomClaim, productID int) rest_err.APIError {
	err := u.dao.Delete(ctx, productID, claims.Merchant)
	if err != nil {
		return err
	}
	return nil
}

// CreateProduct melakukan register product oleh akun owner
func (u *productService) SetCustomPrice(ctx context.Context, claims mjwt.CustomClaim, price dto.ProductPriceRequest) (*dto.ProductModel, rest_err.APIError) {

	timeNow := time.Now().Unix()

	// verifikasi apakah product berasal dari merchant yang sama dengan user
	product, err := u.dao.Get(ctx, price.ProductID)
	if err != nil {
		return nil, err
	}
	if product.MerchantID != claims.Merchant {
		return nil, rest_err.NewBadRequestError("User tidak memeiliki hak akses untuk merubah harga product ini")
	}

	// generate ID dari product id dan outletID
	idGenerated := fmt.Sprintf("%d-%d", price.OutletID, price.ProductID)

	// periksa apakah price id tersebut exist
	existPrice, _ := u.dao.GetPriceDataWithID(ctx, idGenerated)
	var productResult *dto.ProductModel
	if existPrice != nil {
		// datanya ada, perlu dilakukan edit
		productResult, err = u.dao.EditCustomPrice(ctx, dto.ProductPriceModel{
			ID:        dto.UppercaseString(idGenerated),
			ProductID: price.ProductID,
			OutletID:  price.OutletID,
			BuyPrice:  price.BuyPrice,
			SellPrice: price.SellPrice,
			UpdatedAt: timeNow,
		})
	} else {
		// data tidak ada , lakukan insert
		productResult, err = u.dao.InsertCustomPrice(ctx, dto.ProductPriceModel{
			ID:        dto.UppercaseString(idGenerated),
			ProductID: price.ProductID,
			OutletID:  price.OutletID,
			BuyPrice:  price.BuyPrice,
			SellPrice: price.SellPrice,
			UpdatedAt: timeNow,
		})
	}
	if err != nil {
		return nil, err
	}
	return productResult, nil
}

// GetProductByID mendapatkan product dari database
func (u *productService) Get(ctx context.Context, productID int, outletID int) (*dto.ProductModel, rest_err.APIError) {
	var product *dto.ProductModel
	var err rest_err.APIError
	if outletID != 0 {
		// tampilkan harga dengan outlet spesifik
		product, err = u.dao.GetWithCustomPriceOutlet(ctx, productID, outletID)
	} else {
		product, err = u.dao.Get(ctx, productID)
	}
	if err != nil {
		return nil, err
	}

	return product, nil
}

// FindProducts
func (u *productService) FindProducts(ctx context.Context, search string, limit int, offset int, outletSpecific int) ([]dto.ProductModel, rest_err.APIError) {
	productList, err := u.dao.FindWithPagination(ctx, product_dao.FindParams{
		Search: search,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	if outletSpecific != 0 {
		customPrices, err := u.dao.FindCustomPriceOutlet(ctx, outletSpecific)
		if err != nil {
			logger.Info("Custom Price gagal didapatkan")
		}
		if len(customPrices) != 0 {
			priceMap := make(map[int]dto.ProductPriceModel)
			for _, price := range customPrices {
				priceMap[price.ProductID] = price
			}
			for i, product := range productList {
				val, ok := priceMap[product.ID]
				if ok {
					productList[i].BuyPrice = val.BuyPrice
					productList[i].SellPrice = val.SellPrice
				}
			}
		}
	}

	return productList, nil
}
