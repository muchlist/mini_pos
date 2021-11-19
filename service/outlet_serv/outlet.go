package outlet_serv

import (
	"context"
	"github.com/muchlist/mini_pos/dao/outlet_dao"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/mjwt"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"time"
)

type OutletServiceAssumer interface {
	OutletServiceModifier
	OutletServiceReader
}

type OutletServiceReader interface {
	GetOutletByID(ctx context.Context, claims mjwt.CustomClaim, outletID int) (*dto.OutletModel, rest_err.APIError)
	FindOutlets(ctx context.Context, claims mjwt.CustomClaim, search string, limit int, offset int) ([]dto.OutletModel, rest_err.APIError)
}

type OutletServiceModifier interface {
	CreateOutlet(ctx context.Context, claims mjwt.CustomClaim, outlet dto.OutletModel) (int, rest_err.APIError)
	EditOutlet(ctx context.Context, claims mjwt.CustomClaim, request dto.OutletEditRequest) (*dto.OutletModel, rest_err.APIError)
	DeleteOutlet(ctx context.Context, claims mjwt.CustomClaim, outletID int) rest_err.APIError
}

func NewOutletService(dao outlet_dao.OutletDaoAssumer) OutletServiceAssumer {
	return &outletService{
		dao: dao,
	}
}

type outletService struct {
	dao outlet_dao.OutletDaoAssumer
}

// CreateOutlet melakukan register outlet oleh akun owner
func (u *outletService) CreateOutlet(ctx context.Context, claims mjwt.CustomClaim, outlet dto.OutletModel) (int, rest_err.APIError) {

	timeNow := time.Now().Unix()
	outlet.CreatedAt = timeNow
	outlet.UpdatedAt = timeNow
	outlet.MerchantID = claims.Merchant // merchant ID adalah sama dengan merchant id owner

	outletID, err := u.dao.Insert(ctx, outlet)
	if err != nil {
		return 0, err
	}
	return outletID, nil
}

// EditOutlet
func (u *outletService) EditOutlet(ctx context.Context, claims mjwt.CustomClaim, request dto.OutletEditRequest) (*dto.OutletModel, rest_err.APIError) {
	editParams := dto.OutletEditModel{
		WhereID:         request.ID,
		WhereMerchantID: claims.Merchant, // <--- outlet yang diedit harus memiliki merchant id yang sama dengan pengedit
		OutletName:      dto.UppercaseString(request.OutletName),
		Address:         request.Address,
		UpdatedAt:       time.Now().Unix(),
	}

	result, err := u.dao.Edit(ctx, editParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteOutlet
func (u *outletService) DeleteOutlet(ctx context.Context, claims mjwt.CustomClaim, outletID int) rest_err.APIError {
	err := u.dao.Delete(ctx, outletID, claims.Merchant)
	if err != nil {
		return err
	}
	return nil
}

// GetOutletByID mendapatkan outlet dari database
func (u *outletService) GetOutletByID(ctx context.Context, claims mjwt.CustomClaim, outletID int) (*dto.OutletModel, rest_err.APIError) {
	outlet, err := u.dao.Get(ctx, outletID, claims.Merchant)
	if err != nil {
		return nil, err
	}
	return outlet, nil
}

// FindOutlets
func (u *outletService) FindOutlets(ctx context.Context, claims mjwt.CustomClaim, search string, limit int, offset int) ([]dto.OutletModel, rest_err.APIError) {
	outletList, err := u.dao.FindWithPagination(ctx, outlet_dao.FindParams{
		Search: search,
		Limit:  limit,
		Offset: offset,
	}, claims.Merchant)
	if err != nil {
		return nil, err
	}
	return outletList, nil
}
