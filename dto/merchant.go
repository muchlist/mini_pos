package dto

import "time"

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

type MerchantCreateRes struct {
	MerchantID   int    `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	OwnerEmail   string `json:"owner_email"`
	OwnerName    string `json:"owner_name"`
}
