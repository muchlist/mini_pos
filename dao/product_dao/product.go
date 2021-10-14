package product_dao

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
	keyProductTable = "products"
	keyProID        = "id"
	keyProMerchID   = "merchant_id"
	keyProCode      = "code"
	keyProName      = "name"
	keyProDefBuy    = "def_buy_price"
	keyProDefSell   = "def_sell_price"
	keyProImage     = "image"
	keyCreatedAt    = "created_at"
	keyUpdatedAt    = "updated_at"

	keyProductPriceTable     = "product_price"
	keyProductPriceProductID = "product_id"
	keyProductPriceID        = "id"
	keyProductPriceBuy       = "buy_price"
	keyProductPriceSell      = "sell_price"
	keyProductPriceOutletID  = "outlet_id"
)

type ProductDaoAssumer interface {
	Insert(ctx context.Context, input dto.ProductModel) (int, rest_err.APIError)
	Edit(ctx context.Context, input dto.ProductEditModel) (*dto.ProductModel, rest_err.APIError)
	Delete(ctx context.Context, id int, filterMerchant int) rest_err.APIError
	EditCustomPrice(ctx context.Context, input dto.ProductPriceModel) (*dto.ProductModel, rest_err.APIError)
	InsertCustomPrice(ctx context.Context, input dto.ProductPriceModel) (*dto.ProductModel, rest_err.APIError)
	SetImagePath(ctx context.Context, productID int, path string) (*dto.ProductModel, rest_err.APIError)
	Get(ctx context.Context, id int) (*dto.ProductModel, rest_err.APIError)
	GetWithCustomPriceOutlet(ctx context.Context, id int, outletID int) (*dto.ProductModel, rest_err.APIError)
	GetPriceDataWithID(ctx context.Context, priceID string) (*dto.ProductPriceModel, rest_err.APIError)
	FindWithPagination(ctx context.Context, opt FindParams) ([]dto.ProductModel, rest_err.APIError)
	FindCustomPriceOutlet(ctx context.Context, outletID int) ([]dto.ProductPriceModel, rest_err.APIError)
}

type productDao struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func New(db *pgxpool.Pool) ProductDaoAssumer {
	return &productDao{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (p *productDao) Insert(ctx context.Context, input dto.ProductModel) (int, rest_err.APIError) {
	timeNow := time.Now().Unix()
	// -------------------------------------------------------------- insert merchant data
	sqlStatement, args, err := p.sb.Insert(keyProductTable).
		Columns(keyProMerchID, keyProCode, keyProName, keyProDefBuy, keyProDefSell, keyProImage, keyCreatedAt, keyUpdatedAt).
		Values(input.MerchantID, input.Code, input.Name, input.MasterBuyPrice, input.MasterSellPrice, input.Image, timeNow, timeNow).
		Suffix(dao.Returning(keyProID)).
		ToSql()
	if err != nil {
		return 0, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var createdID int
	err = p.db.QueryRow(ctx, sqlStatement, args...).Scan(&createdID)
	if err != nil {
		logger.Error("error saat queryRow product (Insert:0)", err)
		return 0, sql_err.ParseError(err)
	}

	return createdID, nil
}

func (p *productDao) Edit(ctx context.Context, input dto.ProductEditModel) (*dto.ProductModel, rest_err.APIError) {
	timeNow := time.Now().Unix()
	sqlStatement, args, err := p.sb.Update(keyProductTable).
		SetMap(squirrel.Eq{
			keyProCode:    input.Code,
			keyProName:    input.Name,
			keyProDefBuy:  input.MasterBuyPrice,
			keyProDefSell: input.MasterSellPrice,
			keyUpdatedAt:  timeNow,
		}).
		Where(squirrel.And{
			squirrel.Eq{keyProID: input.WhereID},
			squirrel.Eq{keyProMerchID: input.WhereMerchantID}}).
		Suffix(dao.Returning(keyProID, keyProMerchID, keyProCode, keyProName, keyProDefBuy, keyProDefSell, keyProImage, keyCreatedAt, keyUpdatedAt)).
		ToSql()

	if err != nil {
		logger.Error("error saat edit product(Edit:0)", err)
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var res dto.ProductModel
	err = p.db.QueryRow(ctx, sqlStatement, args...).
		Scan(&res.ID, &res.MerchantID, &res.Code, &res.Name, &res.MasterBuyPrice, &res.MasterSellPrice, &res.Image, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	// set price to master if 0
	if res.BuyPrice == 0 {
		res.BuyPrice = res.MasterBuyPrice
	}
	if res.SellPrice == 0 {
		res.SellPrice = res.MasterSellPrice
	}

	return &res, nil
}

func (p *productDao) Delete(ctx context.Context, id int, filterMerchant int) rest_err.APIError {
	sqlStatement, args, err := p.sb.Delete(keyProductTable).
		Where(squirrel.And{
			squirrel.Eq{keyProID: id},
			squirrel.Eq{keyProMerchID: filterMerchant},
		}).
		ToSql()
	if err != nil {
		return rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	res, err := db.DB.Exec(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat delete product(Delete:0)", err)
		return sql_err.ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return rest_err.NewBadRequestError(fmt.Sprintf("Product dengan id %d tidak ditemukan", id))
	}

	return nil
}

func (p *productDao) SetImagePath(ctx context.Context, productID int, path string) (*dto.ProductModel, rest_err.APIError) {
	timeNow := time.Now().Unix()
	sqlStatement, args, err := p.sb.Update(keyProductTable).
		SetMap(squirrel.Eq{
			keyProImage:  path,
			keyUpdatedAt: timeNow,
		}).
		Where(squirrel.Eq{keyProID: productID}).
		Suffix(dao.Returning(keyProID, keyProMerchID, keyProCode, keyProName, keyProDefBuy, keyProDefSell, keyProImage, keyCreatedAt, keyUpdatedAt)).
		ToSql()

	if err != nil {
		logger.Error("error saat edit product(SetImagePath:0)", err)
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var res dto.ProductModel
	err = p.db.QueryRow(ctx, sqlStatement, args...).
		Scan(&res.ID, &res.MerchantID, &res.Code, &res.Name, &res.MasterBuyPrice, &res.MasterSellPrice, &res.Image, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	// set price to master if 0
	if res.BuyPrice == 0 {
		res.BuyPrice = res.MasterBuyPrice
	}
	if res.SellPrice == 0 {
		res.SellPrice = res.MasterSellPrice
	}

	return &res, nil
}

func (p *productDao) Get(ctx context.Context, id int) (*dto.ProductModel, rest_err.APIError) {
	sqlStatement, args, err := p.sb.Select(
		keyProID,
		keyProMerchID,
		keyProCode,
		keyProName,
		keyProDefBuy,
		keyProDefSell,
		keyProImage,
		keyCreatedAt,
		keyUpdatedAt,
	).
		From(keyProductTable).
		Where(squirrel.Eq{keyProID: id}).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var res dto.ProductModel
	err = db.DB.QueryRow(ctx, sqlStatement, args...).
		Scan(&res.ID, &res.MerchantID, &res.Code, &res.Name, &res.MasterBuyPrice, &res.MasterSellPrice, &res.Image, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		logger.Error("error saat get product(Get:0)", err)
		return nil, sql_err.ParseError(err)
	}

	// set price to master if 0
	if res.BuyPrice == 0 {
		res.BuyPrice = res.MasterBuyPrice
	}
	if res.SellPrice == 0 {
		res.SellPrice = res.MasterSellPrice
	}

	return &res, nil
}

func (p *productDao) GetWithCustomPriceOutlet(ctx context.Context, id int, outletID int) (*dto.ProductModel, rest_err.APIError) {
	sqlStatement, args, err := p.sb.Select(
		dao.A(keyProID),
		dao.A(keyProMerchID),
		dao.A(keyProCode),
		dao.A(keyProName),
		dao.A(keyProDefBuy),
		dao.A(keyProDefSell),
		dao.A(keyProImage),
		dao.A(keyCreatedAt),
		dao.A(keyUpdatedAt),
		dao.CoalesceInt(dao.B(keyProductPriceBuy), 0),
		dao.CoalesceInt(dao.B(keyProductPriceSell), 0),
	).
		From(keyProductTable + " A").
		LeftJoin(keyProductPriceTable + " B ON A.id = B.product_id").
		Where(squirrel.And{
			squirrel.Eq{dao.A(keyProID): id},
			squirrel.Eq{dao.B(keyProductPriceOutletID): outletID},
		}).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var res dto.ProductModel
	err = db.DB.QueryRow(ctx, sqlStatement, args...).
		Scan(&res.ID, &res.MerchantID, &res.Code, &res.Name, &res.MasterBuyPrice, &res.MasterSellPrice, &res.Image, &res.CreatedAt, &res.UpdatedAt, &res.BuyPrice, &res.SellPrice)
	if err != nil {
		logger.Error("error saat get product(GetWithCustomPriceOutlet:0)", err)
		return nil, sql_err.ParseError(err)
	}

	// set price to master if 0
	if res.BuyPrice == 0 {
		res.BuyPrice = res.MasterBuyPrice
	}
	if res.SellPrice == 0 {
		res.SellPrice = res.MasterSellPrice
	}

	return &res, nil
}

type FindParams struct {
	Search string
	Limit  int
	Offset int
}

// FindWithPagination example : ?limit=10&offset=10
func (p *productDao) FindWithPagination(ctx context.Context, opt FindParams) ([]dto.ProductModel, rest_err.APIError) {

	// ------------------------------------------------------------------------- find user
	sqlFrom := p.sb.Select(
		keyProID,
		keyProMerchID,
		keyProCode,
		keyProName,
		keyProDefBuy,
		keyProDefSell,
		keyProImage,
		keyCreatedAt,
		keyUpdatedAt).
		From(keyProductTable)

	// where
	if len(opt.Search) > 0 {
		// search
		sqlFrom = sqlFrom.Where(squirrel.ILike{keyProName: fmt.Sprint("%", opt.Search, "%")})
	}

	sqlStatement, args, err := sqlFrom.OrderBy(keyProName + " ASC").
		Limit(uint64(opt.Limit)).
		Offset(uint64(opt.Offset)).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}
	rows, err := db.DB.Query(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat query product(FindWithPagination:0)", err)
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar product", err)
	}
	defer rows.Close()

	products := make([]dto.ProductModel, 0)
	for rows.Next() {
		product := dto.ProductModel{}
		err := rows.Scan(&product.ID, &product.MerchantID, &product.Code, &product.Name, &product.MasterBuyPrice, &product.MasterSellPrice, &product.Image, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			logger.Error("error saat parsing product(FindWithPagination:1)", err)
			return nil, sql_err.ParseError(err)
		}

		// set price to master if 0
		if product.BuyPrice == 0 {
			product.BuyPrice = product.MasterBuyPrice
		}
		if product.SellPrice == 0 {
			product.SellPrice = product.MasterSellPrice
		}
		products = append(products, product)
	}

	return products, nil
}

func (p *productDao) InsertCustomPrice(ctx context.Context, input dto.ProductPriceModel) (*dto.ProductModel, rest_err.APIError) {
	timeNow := time.Now().Unix()
	// -------------------------------------------------------------- insert merchant data
	sqlStatement, args, err := p.sb.Insert(keyProductPriceTable).
		Columns(keyProductPriceID, keyProductPriceProductID, keyProductPriceOutletID, keyProductPriceBuy, keyProductPriceSell, keyUpdatedAt).
		Values(input.ID, input.ProductID, input.OutletID, input.BuyPrice, input.SellPrice, timeNow).
		Suffix(dao.Returning(keyProductPriceID)).
		ToSql()
	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var createdID string
	err = p.db.QueryRow(ctx, sqlStatement, args...).Scan(&createdID)
	if err != nil {
		logger.Error("error saat queryRow product (InsertCustomPrice:0)", err)
		return nil, sql_err.ParseError(err)
	}

	res, apiErr := p.GetWithCustomPriceOutlet(ctx, input.ProductID, input.OutletID)
	if err != nil {
		logger.Error("error saat GetWithCustomPriceOutlet (InsertCustomPrice:1)", err)
		return nil, apiErr
	}

	return res, nil
}

func (p *productDao) EditCustomPrice(ctx context.Context, input dto.ProductPriceModel) (*dto.ProductModel, rest_err.APIError) {
	timeNow := time.Now().Unix()
	sqlStatement, args, err := p.sb.Update(keyProductPriceTable).
		SetMap(squirrel.Eq{
			keyProductPriceBuy:  input.BuyPrice,
			keyProductPriceSell: input.SellPrice,
			keyUpdatedAt:        timeNow,
		}).
		Where(squirrel.Eq{keyProductPriceID: input.ID}).
		Suffix(dao.Returning(keyProductPriceID)).
		ToSql()

	if err != nil {
		logger.Error("error saat edit product price(EditCustomPrice:0)", err)
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var createdID string
	err = p.db.QueryRow(ctx, sqlStatement, args...).Scan(&createdID)
	if err != nil {
		logger.Error("error saat queryRow product (EditCustomPrice:1)", err)
		return nil, sql_err.ParseError(err)
	}

	res, apiErr := p.GetWithCustomPriceOutlet(ctx, input.ProductID, input.OutletID)
	if err != nil {
		logger.Error("error saat GetWithCustomPriceOutlet (EditCustomPrice:2)", err)
		return nil, apiErr
	}

	return res, nil
}

func (p *productDao) GetPriceDataWithID(ctx context.Context, priceID string) (*dto.ProductPriceModel, rest_err.APIError) {
	sqlStatement, args, err := p.sb.Select(
		keyProductPriceID,
		keyProductPriceProductID,
		keyProductPriceOutletID,
		keyProductPriceBuy,
		keyProductPriceSell,
		keyUpdatedAt,
	).
		From(keyProductPriceTable).
		Where(squirrel.Eq{keyProID: priceID}).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var res dto.ProductPriceModel
	err = db.DB.QueryRow(ctx, sqlStatement, args...).
		Scan(&res.ID, &res.ProductID, &res.OutletID, &res.BuyPrice, &res.SellPrice, &res.UpdatedAt)
	if err != nil {
		logger.Error("error saat queryRow(GetPriceWithID:0)", err)
		return nil, sql_err.ParseError(err)
	}

	return &res, nil
}

func (p *productDao) FindCustomPriceOutlet(ctx context.Context, outletID int) ([]dto.ProductPriceModel, rest_err.APIError) {

	// ------------------------------------------------------------------------- find user
	sqlStatement, args, err := p.sb.Select(
		keyProductPriceID,
		keyProductPriceProductID,
		keyProductPriceOutletID,
		keyProductPriceBuy,
		keyProductPriceSell,
		keyUpdatedAt,
	).
		From(keyProductPriceTable).
		Where(squirrel.Eq{keyProductPriceOutletID: outletID}).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	rows, err := db.DB.Query(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat query custom price(FindCustomPriceOutlet:0)", err)
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar custom price", err)
	}
	defer rows.Close()

	customPrices := make([]dto.ProductPriceModel, 0)
	for rows.Next() {
		price := dto.ProductPriceModel{}
		err := rows.Scan(&price.ID, &price.ProductID, &price.OutletID, &price.BuyPrice, &price.SellPrice, &price.UpdatedAt)
		if err != nil {
			logger.Error("error saat parsing product(FindCustomPriceOutlet:1)", err)
			return nil, sql_err.ParseError(err)
		}

		customPrices = append(customPrices, price)
	}

	return customPrices, nil
}
