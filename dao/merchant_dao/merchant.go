package merchant_dao

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
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
	keyMerchantTable = "merchant"
	keyID            = "id"
	keyMerchantName  = "merchant_name"
	keyDescription   = "description"
	keyCreatedAt     = "created_at"
	keyUpdatedAt     = "updated_at"

	keyUserTable      = "users"
	keyUserMerchantID = "merchant_id"
	keyUserDefOutlet  = "def_outlet"
	keyUserName       = "name"
	keyUserEmail      = "email"
	keyUserPassword   = "password"
	keyUserRole       = "role"
)

type merchantDao struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func New(db *pgxpool.Pool) MerchantDaoAssumer {
	return &merchantDao{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (m *merchantDao) Insert(ctx context.Context, input dto.MerchantCreateReq) (*dto.MerchantCreateRes, rest_err.APIError) {

	// ------------------------------------------------------------- begin
	trx, err := m.db.Begin(ctx)
	defer func(trx pgx.Tx) {
		_ = trx.Rollback(context.Background())
	}(trx)

	timeNow := time.Now().Unix()
	response := dto.MerchantCreateRes{}

	// -------------------------------------------------------------- insert merchant data
	sqlStatement, args, err := m.sb.Insert(keyMerchantTable).
		Columns(keyMerchantName, keyDescription, keyCreatedAt, keyUpdatedAt).
		Values(input.MerchantName, input.Description, timeNow, timeNow).
		Suffix(dao.Returning(keyID, keyMerchantName)).
		ToSql()
	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	err = trx.QueryRow(ctx, sqlStatement, args...).Scan(&response.MerchantID, &response.MerchantName)
	if err != nil {
		logger.Error("error saat trx query merchant (Insert:0)", err)
		return nil, sql_err.ParseError(err)
	}

	// ------------------------------------------------------------- insert new for merchant
	sqlStatement, args, err = m.sb.Insert(keyUserTable).Columns(
		keyUserMerchantID,
		keyUserDefOutlet,
		keyUserName,
		keyUserEmail,
		keyUserPassword,
		keyCreatedAt,
		keyUpdatedAt,
		keyUserRole).
		Values(response.MerchantID, 0, input.OwnerName, input.OwnerEmail, input.DefaultPassword, timeNow, timeNow, "owner").
		Suffix(dao.Returning(keyUserEmail, keyUserName)).ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	err = trx.QueryRow(ctx, sqlStatement, args...).Scan(&response.OwnerEmail, &response.OwnerName)
	if err != nil {
		logger.Error("error saat trx query users(Insert:1)", err)
		return nil, rest_err.NewBadRequestError("Email tidak tersedia")
	}

	// ------------------------------------------------------------- commit
	if err := trx.Commit(ctx); err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrCommit, err)
	}

	return &response, nil
}

func (m *merchantDao) Edit(ctx context.Context, input dto.Merchant) (*dto.Merchant, rest_err.APIError) {
	timeNow := time.Now().Unix()
	sqlStatement, args, err := m.sb.Update(keyMerchantTable).
		SetMap(squirrel.Eq{
			keyMerchantName: input.MerchantName,
			keyUpdatedAt:    timeNow,
		}).
		Where(squirrel.Eq{
			keyID: input.Id,
		}).
		Suffix(dao.Returning(keyID, keyMerchantName, keyCreatedAt, keyUpdatedAt)).
		ToSql()

	if err != nil {
		logger.Error("error saat edit merchant(Edit:0)", err)
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var res dto.Merchant
	err = m.db.QueryRow(ctx, sqlStatement, args...).Scan(&res.Id, &res.MerchantName, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}

	return &res, nil
}

func (m *merchantDao) Delete(ctx context.Context, id int) rest_err.APIError {
	sqlStatement, args, err := m.sb.Delete(keyMerchantTable).
		Where(squirrel.Eq{keyID: id}).
		ToSql()
	if err != nil {
		return rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	res, err := db.DB.Exec(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat delete merchant(Delete:0)", err)
		return rest_err.NewInternalServerError("gagal saat penghapusan merchant", err)
	}

	if res.RowsAffected() == 0 {
		return rest_err.NewBadRequestError(fmt.Sprintf("Merchant dengan id %d tidak ditemukan", id))
	}

	return nil
}

func (m *merchantDao) Get(ctx context.Context, id int) (*dto.Merchant, rest_err.APIError) {
	sqlStatement, args, err := m.sb.Select(keyID, keyMerchantName, keyCreatedAt, keyUpdatedAt).
		From(keyMerchantTable).
		Where(squirrel.Eq{
			keyID: id,
		}).ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var res dto.Merchant
	err = db.DB.QueryRow(ctx, sqlStatement, args...).Scan(&res.Id, &res.MerchantName, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan data merchant", err)
	}

	return &res, nil
}

type FindParams struct {
	Search string
	Limit  int
	Cursor int
}

// FindWithCursor example : ?limit=10&cursor=last_id_from_previous_fetch
func (m *merchantDao) FindWithCursor(ctx context.Context, opt FindParams) ([]dto.Merchant, rest_err.APIError) {

	// ------------------------------------------------------------------------- find user
	sqlFrom := m.sb.Select(keyID, keyMerchantName, keyCreatedAt, keyUpdatedAt).
		From(keyMerchantTable)

	// where
	if len(opt.Search) > 0 {
		// search
		sqlFrom = sqlFrom.Where(squirrel.And{
			squirrel.Gt{keyID: opt.Cursor},
			squirrel.ILike{keyMerchantName: fmt.Sprint("%", opt.Search, "%")},
		})
	} else {
		// find
		sqlFrom = sqlFrom.Where(squirrel.Gt{keyID: opt.Cursor})
	}

	sqlStatement, args, err := sqlFrom.OrderBy(keyID + " ASC").
		Limit(uint64(opt.Limit)).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}
	rows, err := db.DB.Query(ctx, sqlStatement, args...)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar merchant", err)
	}
	defer rows.Close()

	merchants := make([]dto.Merchant, 0)
	for rows.Next() {
		merchant := dto.Merchant{}
		err := rows.Scan(&merchant.Id, &merchant.MerchantName, &merchant.CreatedAt, &merchant.UpdatedAt)
		if err != nil {
			return nil, sql_err.ParseError(err)
		}
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}
