package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type Merchant struct {
	Id           int    `json:"id"`
	MerchantName string `json:"merchant_name"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
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
	MerchantName    string `json:"merchant_name"`
	Description     string `json:"description"`
	OwnerEmail      string `json:"owner_email"`
	OwnerName       string `json:"owner_name"`
	DefaultPassword string `json:"default_password"`
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
	MerchantID   int    `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	OwnerEmail   string `json:"owner_email"`
	OwnerName    string `json:"owner_name"`
}

type MerchantEditReq struct {
	Id           int    `json:"-"`
	MerchantName string `json:"merchant_name"`
}

func (m MerchantEditReq) Validate() error {
	if err := validation.ValidateStruct(&m,
		validation.Field(&m.MerchantName, validation.Required),
	); err != nil {
		return err
	}
	return nil
}
