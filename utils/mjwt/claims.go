package mjwt

import (
	"time"
)

type TokenType string

const (
	Access  TokenType = "Access"
	Refresh TokenType = "Refresh"
)

type CustomClaim struct {
	Identity    int
	Name        string
	Exp         int64
	ExtraMinute time.Duration
	Type        TokenType
	Fresh       bool
	Role        string
	Merchant    int
	Outlet      int
}
