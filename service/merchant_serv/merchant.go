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
	mdao   merchant_dao.MerchantDaoAssumer
	mCrypt mcrypt.BcryptAssumer
}

func NewMerchantService(
	mdao merchant_dao.MerchantDaoAssumer,
	mCrypt mcrypt.BcryptAssumer) MerchantServiceAssumer {
	return &merchantService{
		mdao:   mdao,
		mCrypt: mCrypt,
	}
}

func (m *merchantService) CreateMerchant(ctx context.Context, req dto.MerchantCreateReq) (*dto.MerchantCreateRes, rest_err.APIError) {
	hashPw, err := m.mCrypt.GenerateHash(req.DefaultPassword)
	if err != nil {
		return nil, err
	}
	req.DefaultPassword = hashPw
	return m.mdao.Insert(ctx, req)
}
