package merchant_serv

import (
	"context"
	"github.com/muchlist/mini_pos/dao/merchant_dao"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/mcrypt"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"time"
)

type MerchantServiceAssumer interface {
	CreateMerchant(ctx context.Context, req dto.MerchantCreateReq) (*dto.MerchantCreateRes, rest_err.APIError)
	Edit(ctx context.Context, req dto.MerchantEditReq) (*dto.Merchant, rest_err.APIError)
	Delete(ctx context.Context, id int) rest_err.APIError
	Get(ctx context.Context, id int) (*dto.Merchant, rest_err.APIError)
	FindMerchant(ctx context.Context, search string, limit int, cursor int) ([]dto.Merchant, rest_err.APIError)
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

func (m *merchantService) Edit(ctx context.Context, req dto.MerchantEditReq) (*dto.Merchant, rest_err.APIError) {
	return m.dao.Edit(ctx, dto.Merchant{
		Id:           req.Id,
		MerchantName: req.MerchantName,
		UpdatedAt:    time.Now().Unix(),
	})
}

func (m *merchantService) Delete(ctx context.Context, id int) rest_err.APIError {
	return m.dao.Delete(ctx, id)
}

func (m *merchantService) Get(ctx context.Context, id int) (*dto.Merchant, rest_err.APIError) {
	return m.dao.Get(ctx, id)
}

func (m *merchantService) FindMerchant(ctx context.Context, search string, limit int, cursor int) ([]dto.Merchant, rest_err.APIError) {
	merchantList, err := m.dao.FindWithCursor(ctx, merchant_dao.FindParams{
		Search: search,
		Limit:  limit,
		Cursor: cursor,
	})
	if err != nil {
		return nil, err
	}
	return merchantList, nil
}
