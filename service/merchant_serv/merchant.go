package merchant_serv

import (
	"context"
	"github.com/muchlist/mini_pos/dao/merchant_dao"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/mcrypt"
	"github.com/muchlist/mini_pos/utils/rest_err"
)

type MerchantServiceAssumer interface {
	CreateMerchant(ctx context.Context, req dto.MerchantCreateReq) (*dto.MerchantCreateRes, rest_err.APIError)
}

type merchantService struct {
	dao    merchant_dao.MerchantDaoAssumer
	crypto mcrypt.BcryptAssumer
}

func NewMerchantService(mDao merchant_dao.MerchantDaoAssumer, mCrypt mcrypt.BcryptAssumer) MerchantServiceAssumer {
	return &merchantService{
		dao:    mDao,
		crypto: mCrypt,
	}
}

func (m *merchantService) CreateMerchant(ctx context.Context, req dto.MerchantCreateReq) (*dto.MerchantCreateRes, rest_err.APIError) {
	// hashing default password yang diberikan
	hashPw, err := m.crypto.GenerateHash(req.DefaultPassword)
	if err != nil {
		return nil, err
	}
	req.DefaultPassword = hashPw
	return m.dao.Insert(ctx, req)
}
