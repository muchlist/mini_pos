package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type ProductModel struct {
	ID              int             `json:"id"`
	MerchantID      int             `json:"merchant_id"`
	Code            UppercaseString `json:"code"` // SKU
	Name            UppercaseString `json:"name"`
	MasterBuyPrice  int             `json:"master_buy_price"`
	MasterSellPrice int             `json:"master_sell_price"`
	BuyPrice        int             `json:"buy_price"`  // berasal dari table lain
	SellPrice       int             `json:"sell_price"` // berasal dari table lain
	Image           string          `json:"image"`
	CreatedAt       int64           `json:"created_at"`
	UpdatedAt       int64           `json:"updated_at"`
}

type ProductCreateRequest struct {
	Code            string `json:"code"` // SKU
	Name            string `json:"name"`
	MasterBuyPrice  int    `json:"master_buy_price"`
	MasterSellPrice int    `json:"master_sell_price"`
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
	Code            string `json:"code"` // SKU
	Name            string `json:"name"`
	MasterBuyPrice  int    `json:"master_buy_price"`
	MasterSellPrice int    `json:"master_sell_price"`
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
	Code            UppercaseString `json:"code"` // SKU
	Name            UppercaseString `json:"name"`
	MasterBuyPrice  int             `json:"master_buy_price"`
	MasterSellPrice int             `json:"master_sell_price"`
}
