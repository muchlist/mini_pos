package dto

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/muchlist/mini_pos/configs/roles"
	"github.com/muchlist/mini_pos/utils/sfunc"
	"strings"
)

type UserModel struct {
	ID         int             `json:"id" example:"1"`
	Email      LowercaseString `json:"email" example:"example@example.com"`
	Name       UppercaseString `json:"name" example:"muchlis"`
	Password   string          `json:"-"`
	Role       LowercaseString `json:"role" example:"owner,employee"`
	CreatedAt  int64           `json:"created_at" example:"1631341964"`
	UpdatedAt  int64           `json:"updated_at" example:"1631341964"`
	MerchantID int             `json:"merchant_id" example:"1"`
	DefOutlet  int             `json:"def_outlet" example:"1"`
}

type UserRegisterRequest struct {
	Email     string `json:"email" example:"example@example.com"`
	Name      string `json:"name" example:"muchlis"`
	Role      string `json:"role" example:"owner,employee"`
	Password  string `json:"password" example:"password123"`
	DefOutlet int    `json:"def_outlet" example:"1"`
}

func (u UserRegisterRequest) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Role, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	); err != nil {
		return err
	}

	if !sfunc.InSlice(strings.ToLower(u.Role), roles.GetRolesAvailable()) {
		// assums role tidak cocok dengan enum yang tersedia di database
		return errors.New(fmt.Sprintf("Role yang dimasukkan salah, gunakan %v", roles.GetRolesAvailable()))
	}

	return nil
}

type UserEditRequest struct {
	ID        int    `json:"-"`
	Email     string `json:"email" example:"example@example.com"`
	Name      string `json:"name" example:"muchlis"`
	DefOutlet int    `json:"def_outlet" example:"1"`
	Role      string `json:"role" example:"employee"`
}

func (u UserEditRequest) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Role, validation.Required),
	); err != nil {
		return err
	}
	if !sfunc.InSlice(strings.ToLower(u.Role), roles.GetRolesAvailable()) {
		// assums role tidak cocok dengan enum yang tersedia di database
		return errors.New(fmt.Sprintf("Role yang dimasukkan salah, gunakan %v", roles.GetRolesAvailable()))
	}

	return nil
}

type UserEditModel struct {
	WhereID         int
	WhereMerchantID int
	Email           LowercaseString
	Name            UppercaseString
	Role            LowercaseString
	UpdatedAt       int64
	DefOutlet       int
}

// UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginResponse struct {
	ID           int    `json:"id" example:"1"`
	Email        string `json:"email" example:"example@example.com"`
	Name         string `json:"name" example:"muchlis"`
	MerchantID   int    `json:"merchant_id" example:"1"`
	DefOutlet    int    `json:"def_outlet" example:"1"`
	Role         string `json:"role" example:"owner,employee"`
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"`
	Expired      int64  `json:"expired" example:"1631341964"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"`
}

// UserRefreshTokenResponse mengembalikan token dengan claims yang
// sama dengan token sebelumnya dengan expired yang baru
type UserRefreshTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"`
	Expired     int64  `json:"expired" example:"1631341964"`
}
