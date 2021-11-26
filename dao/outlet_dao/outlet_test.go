package outlet_dao

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"testing"
)

// SELECT id, merchant_id, outlet_name, address, created_at, updated_at
// FROM outlets
// WHERE (outlet_name ILIKE $1 AND merchant_id = $2)
// ORDER BY id ASC
// LIMIT 10
// OFFSET 20
func TestFindWithPagination(t *testing.T) {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlFrom := sb.Select(keyID, keyMerchantID, keyOutletName, keyAddress, keyCreatedAt, keyUpdatedAt).
		From(keyOutletTable)

	sqlFrom = sqlFrom.Where(sq.And{
		sq.ILike{keyOutletName: fmt.Sprint("%", "TEST", "%")},
		sq.Eq{keyMerchantID: 12345},
	})

	sqlStatement, args, err := sqlFrom.OrderBy(keyID + " ASC").
		Limit(uint64(10)).
		Offset(uint64(20)).
		ToSql()

	fmt.Println(sqlStatement)
	fmt.Printf("%v\n", args)
	assert.Nil(t, err)
}