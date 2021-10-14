package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type OutletModel struct {
	ID         int             `json:"id"`
	MerchantID int             `json:"merchant_id"`
	OutletName UppercaseString `json:"outlet_name"`
	Address    string          `json:"address"`
	CreatedAt  int64           `json:"created_at"`
	UpdatedAt  int64           `json:"updated_at"`
}

type OutletCreateRequest struct {
	OutletName string `json:"outlet_name"`
	Address    string `json:"address"`
}

func (o OutletCreateRequest) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.OutletName, validation.Required),
		validation.Field(&o.Address, validation.Required),
	)
}

type OutletEditRequest struct {
	ID         int    `json:"-"`
	OutletName string `json:"outlet_name"`
	Address    string `json:"address"`
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
