package user_dao

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
	keyUserTable      = "users"
	keyUserID         = "id"
	keyUserMerchantID = "merchant_id"
	keyUserDefOutlet  = "def_outlet"
	keyUserName       = "name"
	keyUserEmail      = "email"
	keyUserPassword   = "password"
	keyUserRole       = "role"
	keyCreatedAt      = "created_at"
	keyUpdatedAt      = "updated_at"
)

type userDao struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func New(db *pgxpool.Pool) UserDaoAssumer {
	return &userDao{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (u userDao) Insert(ctx context.Context, user dto.UserModel) (string, rest_err.APIError) {
	timeNow := time.Now().Unix()

	sqlStatement, args, err := u.sb.Insert(keyUserTable).Columns(
		keyUserMerchantID,
		keyUserDefOutlet,
		keyUserName,
		keyUserEmail,
		keyUserPassword,
		keyCreatedAt,
		keyUpdatedAt,
		keyUserRole).
		Values(user.MerchantID, user.DefOutlet, user.Name, user.Email, user.Password, timeNow, timeNow, user.Role).
		Suffix(dao.Returning(keyUserEmail, keyUserName)).ToSql()

	if err != nil {
		return "", rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var email string
	var name string
	err = u.db.QueryRow(ctx, sqlStatement, args...).Scan(&email, &name)
	if err != nil {
		logger.Error("error saat query users(Insert:1)", err)
		return "", sql_err.ParseError(err)
	}

	return fmt.Sprintf("berhasil menambahkan user dengan nama %s - email %s", name, email), nil
}

func (u userDao) Edit(ctx context.Context, input dto.UserEditModel) (*dto.UserModel, rest_err.APIError) {
	sqlStatement, args, err := u.sb.Update(keyUserTable).
		SetMap(squirrel.Eq{
			keyUserName:      input.Name,
			keyUserRole:      input.Role,
			keyUserDefOutlet: input.DefOutlet,
			keyUpdatedAt:     input.UpdatedAt,
		}).
		Where(
			squirrel.And{
				squirrel.Eq{keyUserID: input.WhereID},
				squirrel.Eq{keyUserMerchantID: input.WhereMerchantID},
			}).
		Suffix(dao.Returning(
			keyUserID,
			keyUserMerchantID,
			keyUserDefOutlet,
			keyUserName,
			keyUserEmail,
			keyCreatedAt,
			keyUpdatedAt,
			keyUserRole)).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var user dto.UserModel
	err = u.db.QueryRow(
		ctx,
		sqlStatement, args...).Scan(
		&user.ID,
		&user.MerchantID,
		&user.DefOutlet,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role)
	if err != nil {
		logger.Error("error saat query users(Edit:1)", err)
		return nil, sql_err.ParseError(err)
	}
	return &user, nil
}

func (u userDao) Delete(ctx context.Context, id int, filterMerchant int) rest_err.APIError {
	sqlStatement, args, err := u.sb.Delete(keyUserTable).
		Where(squirrel.And{
			squirrel.Eq{keyUserID: id},
			squirrel.Eq{keyUserMerchantID: filterMerchant},
		}).ToSql()
	if err != nil {
		return rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	res, err := db.DB.Exec(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat exec users(Delete:0)", err)
		return rest_err.NewInternalServerError("gagal saat penghapusan user", err)
	}

	if res.RowsAffected() == 0 {
		return rest_err.NewBadRequestError(fmt.Sprintf("UserModel dengan username %d tidak ditemukan", id))
	}

	return nil
}

func (u userDao) ChangePassword(ctx context.Context, input dto.UserModel) rest_err.APIError {
	sqlStatement, args, err := u.sb.Update(keyUserTable).
		SetMap(squirrel.Eq{
			keyUserPassword: input.Password,
			keyUpdatedAt:    input.UpdatedAt,
		}).
		Where(keyUserID, input.ID).
		ToSql()

	if err != nil {
		return rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	res, err := db.DB.Exec(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat exec users(ChangePassword:0)", err)
		return sql_err.ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return rest_err.NewBadRequestError(fmt.Sprintf("UserModel dengan username %d tidak ditemukan", input.ID))
	}

	return nil
}

func (u userDao) GetByID(ctx context.Context, id int) (*dto.UserModel, rest_err.APIError) {
	sqlStatement, args, err := u.sb.Select(
		keyUserID,
		keyUserMerchantID,
		keyUserDefOutlet,
		keyUserName,
		keyUserEmail,
		keyCreatedAt,
		keyUpdatedAt,
		keyUserRole,
	).From(keyUserTable).Where(squirrel.Eq{keyUserID: id}).ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var user dto.UserModel
	err = db.DB.QueryRow(ctx, sqlStatement, args...).
		Scan(&user.ID, &user.MerchantID, &user.DefOutlet, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err != nil {
		logger.Error("error saat QueryRow users(Get:0)", err)
		return nil, rest_err.NewInternalServerError("gagal mendapatkan user", err)
	}

	return &user, nil
}

func (u userDao) GetByEmail(ctx context.Context, email string) (*dto.UserModel, rest_err.APIError) {
	sqlStatement, args, err := u.sb.Select(
		keyUserID,
		keyUserMerchantID,
		keyUserDefOutlet,
		keyUserName,
		keyUserEmail,
		keyUserPassword,
		keyCreatedAt,
		keyUpdatedAt,
		keyUserRole,
	).From(keyUserTable).Where(squirrel.Eq{keyUserEmail: email}).ToSql()

	if err != nil {
		logger.Error("error saat query builder users(GetByEmail:0)", err)
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var user dto.UserModel
	err = db.DB.QueryRow(ctx, sqlStatement, args...).
		Scan(&user.ID, &user.MerchantID, &user.DefOutlet, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err != nil {
		logger.Error("error saat QueryRow users(GetByEmail:1)", err)
		return nil, rest_err.NewInternalServerError("gagal mendapatkan user", err)
	}

	return &user, nil
}

type FindPaginationParams struct {
	Search string
	Limit  int
	Offset int
}

func (u *userDao) FindWithPagination(ctx context.Context, opt FindPaginationParams) ([]dto.UserModel, rest_err.APIError) {
	// ------------------------------------------------------------------------- find user
	sqlFrom := u.sb.Select(
		keyUserID,
		keyUserMerchantID,
		keyUserDefOutlet,
		keyUserName,
		keyUserEmail,
		keyCreatedAt,
		keyUpdatedAt,
		keyUserRole).
		From(keyUserTable)

	// where
	if len(opt.Search) > 0 {
		// search
		sqlFrom = sqlFrom.Where(squirrel.ILike{keyUserName: fmt.Sprint("%", opt.Search, "%")})
	}

	sqlStatement, args, err := sqlFrom.OrderBy(keyUserName + " ASC").
		Limit(uint64(opt.Limit)).
		Offset(uint64(opt.Offset)).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}
	rows, err := db.DB.Query(ctx, sqlStatement, args...)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar user", err)
	}
	defer rows.Close()

	users := make([]dto.UserModel, 0)
	for rows.Next() {
		var user dto.UserModel
		err := rows.Scan(&user.ID, &user.MerchantID, &user.DefOutlet, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role)
		if err != nil {
			return nil, sql_err.ParseError(err)
		}
		users = append(users, user)
	}

	return users, nil
}
