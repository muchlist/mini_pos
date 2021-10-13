package merchant_dao

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/mini_pos/dao"
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

type MerchantDaoAssumer interface {
	MerchantSaver
}

type MerchantSaver interface {
	Insert(ctx context.Context, input dto.MerchantCreateReq) (*dto.MerchantCreateRes, rest_err.APIError)
}

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

//
//func (u *userDao) Edit(ctx context.Context, input dto.User) (*dto.User, rest_err.APIError) {
//
//	if len(input.Roles) == 0 {
//		return nil, rest_err.NewBadRequestError("role tidak boleh kosong")
//	}
//
//	// ------------------------------------------------------------------------- begin
//	trx, err := u.db.Begin(ctx)
//	defer func(trx pgx.Tx) {
//		_ = trx.Rollback(context.Background())
//	}(trx)
//
//	// ------------------------------------------------------------------------- user edit
//	sqlStatement, args, err := u.sb.Update(keyUserTable).
//		SetMap(squirrel.Eq{
//			keyEmail:     input.Email,
//			keyName:      input.Name,
//			keyUpdatedAt: input.UpdatedAt,
//		}).
//		Where(squirrel.Eq{
//			keyID: input.ID,
//		}).
//		Suffix(dao.Returning(keyID, keyEmail, keyName, keyCreatedAt, keyUpdatedAt)).
//		ToSql()
//
//	if err != nil {
//		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
//	}
//
//	var user dto.User
//	err = trx.QueryRow(
//		ctx,
//		sqlStatement, args...).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
//	if err != nil {
//		return nil, sql_err.ParseError(err)
//	}
//
//	// ------------------------------------------------------------------------- role delete
//	sqlStatement, args, err = u.sb.Delete(keyUsersRolesTable).
//		Where(squirrel.Eq{
//			keyUsersID: input.ID,
//		}).ToSql()
//	if err != nil {
//		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
//	}
//
//	_, err = db.DB.Exec(ctx, sqlStatement, args...)
//	if err != nil {
//		logger.Error("error saat trx exec usersRoles(ChangeRole:0)", err)
//		return nil, sql_err.ParseError(err)
//	}
//
//	// ------------------------------------------------------------------------- role insert
//	sqlInsert := u.sb.Insert(keyUsersRolesTable).Columns(keyRolesName, keyUsersID)
//	for _, roleName := range input.Roles {
//		sqlInsert = sqlInsert.Values(roleName, user.ID)
//	}
//	sqlStatement, args, err = sqlInsert.ToSql()
//	if err != nil {
//		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
//	}
//
//	_, err = trx.Exec(ctx, sqlStatement, args...)
//	if err != nil {
//		logger.Error("error saat trx query usersRoles(Insert:1)", err)
//		return nil, sql_err.ParseError(err)
//	}
//
//	// ------------------------------------------------------------------------- commit
//	if err := trx.Commit(ctx); err != nil {
//		return nil, rest_err.NewInternalServerError(dao.ErrCommit, err)
//	}
//
//	user.Roles = input.Roles
//
//	return &user, nil
//}
