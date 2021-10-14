package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type Merchant struct {
	Id           int    `json:"id" example:"1"`
	MerchantName string `json:"merchant_name" example:"KUKUS-TOKO"`
	CreatedAt    int64  `json:"created_at" example:"1631341964"`
	UpdatedAt    int64  `json:"updated_at" example:"1631341964"`
}

func (m *Merchant) Prepare() {
	timeNow := time.Now().Unix()
	if m.CreatedAt == 0 {
		m.CreatedAt = timeNow
	}
	if m.UpdatedAt == 0 {
		m.UpdatedAt = timeNow
	}
}

type MerchantCreateReq struct {
	MerchantName    string `json:"merchant_name" example:"KUKUS-TOKO"`
	Description     string `json:"description" example:"penjualan barang barang tidak kasat mata"`
	OwnerEmail      string `json:"owner_email" example:"example@gmail.com"`
	OwnerName       string `json:"owner_name" example:"MUCHLIS"`
	DefaultPassword string `json:"default_password" example:"secret"`
}

func (m MerchantCreateReq) Validate() error {
	if err := validation.ValidateStruct(&m,
		validation.Field(&m.MerchantName, validation.Required),
		validation.Field(&m.OwnerName, validation.Required),
		validation.Field(&m.OwnerEmail, validation.Required, is.Email),
		validation.Field(&m.DefaultPassword, validation.Required, validation.Length(3, 20)),
	); err != nil {
		return err
	}
	return nil
}

type MerchantCreateRes struct {
	MerchantID   int    `json:"merchant_id" example:"1"`
	MerchantName string `json:"merchant_name" example:"KUKUS-TOKO"`
	OwnerEmail   string `json:"owner_email" example:"example@gmail.com"`
	OwnerName    string `json:"owner_name" example:"MUCHLIS"`
}

type MerchantEditReq struct {
	Id           int    `json:"-"`
	MerchantName string `json:"merchant_name" example:"KUKUS TOKO"`
}

func (m MerchantEditReq) Validate() error {
	if err := validation.ValidateStruct(&m,
		validation.Field(&m.MerchantName, validation.Required),
	); err != nil {
		return err
	}
	return nil
}
