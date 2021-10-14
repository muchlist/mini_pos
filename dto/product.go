package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type ProductModel struct {
	ID              int             `json:"id" example:"1"`
	MerchantID      int             `json:"merchant_id" example:"20"`
	Code            UppercaseString `json:"code" example:"CAT-20"` // SKU
	Name            UppercaseString `json:"name" example:"JAM TANGAN"`
	MasterBuyPrice  int             `json:"master_buy_price" example:"1000000"`
	MasterSellPrice int             `json:"master_sell_price" example:"1050000"`
	BuyPrice        int             `json:"buy_price" example:"1000000"`  // berasal dari table lain
	SellPrice       int             `json:"sell_price" example:"1000000"` // berasal dari table lain
	Image           string          `json:"image" example:"image/products/121634211915.jpg"`
	CreatedAt       int64           `json:"created_at" example:"1631341964"`
	UpdatedAt       int64           `json:"updated_at" example:"1631341964"`
}

type ProductCreateRequest struct {
	Code            string `json:"code" example:"CAT-20"` // SKU
	Name            string `json:"name" example:"JAM TANGAN"`
	MasterBuyPrice  int    `json:"master_buy_price" example:"1000000"`
	MasterSellPrice int    `json:"master_sell_price" example:"1050000"`
}

func (p ProductCreateRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Code, validation.Required),
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.MasterBuyPrice, validation.Required),
		validation.Field(&p.MasterSellPrice, validation.Required),
	)
}

type ProductEditRequest struct {
	ID              int    `json:"-"`
	Code            string `json:"code" example:"CAT-20"` // SKU
	Name            string `json:"name" example:"JAM TANGAN"`
	MasterBuyPrice  int    `json:"master_buy_price" example:"1000000"`
	MasterSellPrice int    `json:"master_sell_price" example:"1050000"`
}

func (p ProductEditRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Code, validation.Required),
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.MasterBuyPrice, validation.Required),
		validation.Field(&p.MasterSellPrice, validation.Required),
	)
}

type ProductEditModel struct {
	WhereID         int
	WhereMerchantID int
	Code            UppercaseString
	Name            UppercaseString
	MasterBuyPrice  int
	MasterSellPrice int
}

type ProductPriceModel struct {
	ID        UppercaseString `json:"id"  example:"1-20"` // combine productID-outletID
	ProductID int             `json:"product_id"  example:"1"`
	OutletID  int             `json:"outlet_id" example:"20"`
	BuyPrice  int             `json:"buy_price" example:"1000000"`
	SellPrice int             `json:"sell_price" example:"1050000"`
	UpdatedAt int64           `json:"updated_at" example:"1631341964"`
}

type ProductPriceRequest struct {
	ProductID int `json:"product_id"  example:"1"`
	OutletID  int `json:"outlet_id" example:"20"`
	BuyPrice  int `json:"buy_price" example:"1000000"`
	SellPrice int `json:"sell_price" example:"1050000"`
}

func (p ProductPriceRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ProductID, validation.Required),
		validation.Field(&p.OutletID, validation.Required),
		validation.Field(&p.BuyPrice, validation.Required),
		validation.Field(&p.SellPrice, validation.Required),
	)
}
