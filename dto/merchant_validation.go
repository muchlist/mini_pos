package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

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
