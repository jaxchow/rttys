package sqlbuilder

import (
	"database/sql"
	"testing"

	"github.com/go-playground/assert/v2"
	_ "modernc.org/sqlite"
)

func TestSqlBuilder(t *testing.T) {

	sql := QueryBuilder(SQuery{})
	assert.Equal(t, sql, "SELECT * FROM telant WHERE 1=1")
}

func TestSqlBuilderWithLimiy(t *testing.T) {

	sql := QueryBuilder(SQuery{
		SPage: SPage{
			Limit: 10,
		},
	})
	assert.Equal(t, sql, "SELECT * FROM telant WHERE 1=1 limit 10")
}

func TestSqlBuilderWithOffset(t *testing.T) {

	sql := QueryBuilder(SQuery{
		SPage: SPage{
			Offset: 10,
		},
	})
	assert.Equal(t, sql, "SELECT * FROM telant WHERE 1=1 offset 10")
}

func TestInsertBuilder(t *testing.T) {
	val := STelant{
		Name:        "jaxteam",
		Owner:       "jaxchow",
		Description: "jaxteam org",
	}
	sql := InsertBuilder(val)
	assert.Equal(t, sql, "INSERT INTO telant(name,owner,description) VALUES('jaxteam','jaxchow','jaxteam org')")
}

func TestDeleteBuilder(t *testing.T) {
	val := SQuery{
		STelant: STelant{
			Name:        "jaxteam",
			Owner:       "jaxchow",
			Description: "jaxteam org",
		},
	}
	sql := DeleteBuilder(val)
	assert.Equal(t, sql, "DELETE FROM telant WHERE 1=1 and name='jaxteam' and owner='jaxchow' and description='jaxteam org'")
}

func TestSelect(t *testing.T) {
	dsn := "sqlite://rttys.db"
	db, err := sql.Open("sqlite", dsn)
	rs, err := Query(db, SQuery{})
	print(rs, err)
}
