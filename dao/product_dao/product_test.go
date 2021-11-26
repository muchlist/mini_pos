package product_dao

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/muchlist/mini_pos/dao"
	"github.com/stretchr/testify/assert"
)

// Hanya untuk ingin melihat hasil querynya saja
// SELECT A.id, A.merchant_id, A.code, A.name, A.def_buy_price, A.def_sell_price, A.image, A.created_at, A.updated_at, Coalesce(B.buy_price,0), Coalesce(B.sell_price,0)
// FROM products A LEFT JOIN product_price B ON A.id = B.product_id
// WHERE (A.id = $1 AND B.outlet_id = $2)
func TestGetWithCustomPrice(t *testing.T) {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlStatement, args, err := sb.Select(
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
		Where(sq.And{
			sq.Eq{dao.A(keyProID): 1},
			sq.Eq{dao.B(keyProductPriceOutletID): 2},
		}).
		ToSql()

	println(sqlStatement)
	fmt.Printf("%v\n", args)
	assert.Nil(t, err)
}

// UPDATE products SET code =
// (CASE WHEN id = $1 THEN A WHEN id = $2 THEN $3 WHEN id = $3 THEN 0 WHEN id = $4 THEN 1 WHEN id = $5 THEN 2 END),
// updated_at = $6 WHERE (id = $7 AND merchant_id = $8)
// [123 125 it's true! 0 1 2 123123 1 2]
func TestUpdateBulk(t *testing.T) {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	caseStmt := sq.Case().
		When(sq.Eq{keyProID: "123"}, "AZ").
		When(sq.Eq{keyProID: "125"}, sq.Expr("?", "it's true!"))

	for i := 0; i < 3; i++ {
		caseStmt = caseStmt.When(sq.Eq{keyProID: i}, fmt.Sprintf("%d", i))
	}

	sqlStatement, args, err := sb.Update(keyProductTable).
		SetMap(sq.Eq{
			keyProCode:   Sub(caseStmt),
			keyUpdatedAt: 123123,
		}).
		Where(sq.And{
			sq.Eq{keyProID: 1},
			sq.Eq{keyProMerchID: 2}}).
		ToSql()

	println(sqlStatement)
	fmt.Printf("%v\n", args)
	assert.Nil(t, err)
}

func Sub(sb sq.CaseBuilder) sq.Sqlizer {
	sql, params, _ := sb.ToSql()
	return sq.Expr("("+sql+")", params...)
}
