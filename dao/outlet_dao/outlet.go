package outlet_dao

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/mini_pos/dao"
	"github.com/muchlist/mini_pos/db"
	"github.com/muchlist/mini_pos/dto"
	"github.com/muchlist/mini_pos/utils/logger"
	"github.com/muchlist/mini_pos/utils/rest_err"
	"github.com/muchlist/mini_pos/utils/sql_err"
	"time"
)

const (
	keyOutletTable = "outlets"
	keyID          = "id"
	keyMerchantID  = "merchant_id"
	keyOutletName  = "outlet_name"
	keyAddress     = "address"
	keyCreatedAt   = "created_at"
	keyUpdatedAt   = "updated_at"
)

type outletDao struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func New(db *pgxpool.Pool) OutletDaoAssumer {
	return &outletDao{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (o *outletDao) Insert(ctx context.Context, input dto.OutletModel) (int, rest_err.APIError) {
	timeNow := time.Now().Unix()

	// -------------------------------------------------------------- insert merchant data
	sqlStatement, args, err := o.sb.Insert(keyOutletTable).
		Columns(keyMerchantID, keyOutletName, keyAddress, keyCreatedAt, keyUpdatedAt).
		Values(input.MerchantID, input.OutletName, input.Address, timeNow, timeNow).
		Suffix(dao.Returning(keyID)).
		ToSql()
	if err != nil {
		return 0, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var createdID int
	err = o.db.QueryRow(ctx, sqlStatement, args...).Scan(&createdID)
	if err != nil {
		logger.Error("error saat query outlet (Insert:0)", err)
		return 0, sql_err.ParseError(err)
	}

	return createdID, nil
}

func (o *outletDao) Edit(ctx context.Context, input dto.OutletEditModel) (*dto.OutletModel, rest_err.APIError) {
	timeNow := time.Now().Unix()
	sqlStatement, args, err := o.sb.Update(keyOutletTable).
		SetMap(squirrel.Eq{
			keyOutletName: input.OutletName,
			keyUpdatedAt:  timeNow,
		}).
		Where(squirrel.And{
			squirrel.Eq{keyID: input.WhereID},
			squirrel.Eq{keyMerchantID: input.WhereMerchantID}}).
		Suffix(dao.Returning(keyID, keyMerchantID, keyOutletName, keyAddress, keyCreatedAt, keyUpdatedAt)).
		ToSql()

	if err != nil {
		logger.Error("error saat edit outlet(Edit:0)", err)
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var res dto.OutletModel
	err = o.db.QueryRow(ctx, sqlStatement, args...).
		Scan(&res.ID, &res.MerchantID, &res.OutletName, &res.Address, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}

	return &res, nil
}

func (o *outletDao) Delete(ctx context.Context, id int, filterMerchant int) rest_err.APIError {
	sqlStatement, args, err := o.sb.Delete(keyOutletTable).
		Where(squirrel.And{
			squirrel.Eq{keyID: id},
			squirrel.Eq{keyMerchantID: filterMerchant},
		}).
		ToSql()
	if err != nil {
		return rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	res, err := db.DB.Exec(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat delete outlet(Delete:0)", err)
		return sql_err.ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return rest_err.NewBadRequestError(fmt.Sprintf("Outlet dengan id %d tidak ditemukan", id))
	}

	return nil
}

func (o *outletDao) Get(ctx context.Context, id int, merchantIDifSpecific int) (*dto.OutletModel, rest_err.APIError) {
	sqlStatement, args, err := o.sb.Select(keyID, keyMerchantID, keyOutletName, keyAddress, keyCreatedAt, keyUpdatedAt).
		From(keyOutletTable).
		Where(squirrel.Eq{
			keyID:         id,
			keyMerchantID: merchantIDifSpecific,
		}).ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var res dto.OutletModel
	err = db.DB.QueryRow(ctx, sqlStatement, args...).
		Scan(&res.ID, &res.MerchantID, &res.OutletName, &res.Address, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		logger.Error("error saat query outlet(Get:0)", err)
		return nil, sql_err.ParseError(err)
	}

	return &res, nil
}

type FindParams struct {
	Search string
	Limit  int
	Offset int
}

// FindWithPagination example : ?limit=10&offset=10
func (o *outletDao) FindWithPagination(ctx context.Context, opt FindParams, merchantFilter int) ([]dto.OutletModel, rest_err.APIError) {

	// ------------------------------------------------------------------------- find user
	sqlFrom := o.sb.Select(keyID, keyMerchantID, keyOutletName, keyAddress, keyCreatedAt, keyUpdatedAt).
		From(keyOutletTable)

	// where
	if len(opt.Search) > 0 {
		// search
		sqlFrom = sqlFrom.Where(squirrel.And{
			squirrel.ILike{keyOutletName: fmt.Sprint("%", opt.Search, "%")},
			squirrel.Eq{keyMerchantID: merchantFilter},
		})
	} else {
		sqlFrom = sqlFrom.Where(squirrel.Eq{keyMerchantID: merchantFilter})
	}

	sqlStatement, args, err := sqlFrom.OrderBy(keyID + " ASC").
		Limit(uint64(opt.Limit)).
		Offset(uint64(opt.Offset)).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}
	rows, err := db.DB.Query(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat query outlet(FindWithPagination:0)", err)
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar outlet", err)
	}
	defer rows.Close()

	outlets := make([]dto.OutletModel, 0)
	for rows.Next() {
		outlet := dto.OutletModel{}
		err := rows.Scan(&outlet.ID, &outlet.MerchantID, &outlet.OutletName, &outlet.Address, &outlet.CreatedAt, &outlet.UpdatedAt)
		if err != nil {
			logger.Error("error saat parsing outlet(FindWithPagination:1)", err)
			return nil, sql_err.ParseError(err)
		}
		outlets = append(outlets, outlet)
	}

	return outlets, nil
}
