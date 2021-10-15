package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type OutletModel struct {
	ID         int             `json:"id" example:"1"`
	MerchantID int             `json:"merchant_id" example:"1"`
	OutletName UppercaseString `json:"outlet_name" example:"BLOK B"`
	Address    string          `json:"address" example:"Jl Pangeran Samudera"`
	CreatedAt  int64           `json:"created_at" example:"1631341964"`
	UpdatedAt  int64           `json:"updated_at" example:"1631341964"`
}

type OutletCreateRequest struct {
	OutletName string `json:"outlet_name" example:"BLOK B"`
	Address    string `json:"address" example:"Jl Pangeran Samudera"`
}

func (o OutletCreateRequest) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.OutletName, validation.Required),
		validation.Field(&o.Address, validation.Required),
	)
}

type OutletEditRequest struct {
	ID         int    `json:"-"`
	OutletName string `json:"outlet_name" example:"BLOK B"`
	Address    string `json:"address" example:"Jl Pangeran Samudera"`
}

func (o OutletEditRequest) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.OutletName, validation.Required),
		validation.Field(&o.Address, validation.Required),
	)
}

type OutletEditModel struct {
	WhereID         int
	WhereMerchantID int
	OutletName      UppercaseString
	Address         string
	UpdatedAt       int64
}
