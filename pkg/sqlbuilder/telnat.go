package sqlbuilder

import (
	"database/sql"
	"fmt"
	"strings"
)

type SPage struct {
	Limit  int8 `json:"limit"`
	Offset int8 `json:"offet"`
}

type SQuery struct {
	STelant
	SPage
}
type STelant struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
}

func RowsToArray(rows *sql.Rows) []STelant {
	result := make([]STelant, 0)
	for rows.Next() {
		telant := &STelant{}
		rows.Scan(&telant.Name, &telant.Owner, &telant.Description)
		result = append(result, *telant)
	}
	return result
}

func setBuilder(sqlBuilder *strings.Builder, value STelant) {
	if value.Name != "" {
		sqlBuilder.WriteString(fmt.Sprintf("set name=%s", value.Name))
	}

	if value.Owner != "" {
		sqlBuilder.WriteString(fmt.Sprintf("set owner=%s", value.Owner))
	}

	if value.Description != "" {
		sqlBuilder.WriteString(fmt.Sprintf("set name=%s", value.Description))
	}
}

func whereBuilder(sqlBuilder *strings.Builder, query SQuery) {
	if query.Name != "" {
		sqlBuilder.WriteString(fmt.Sprintf(" and name='%s'", query.Name))
	}
	if query.Owner != "" {
		sqlBuilder.WriteString(fmt.Sprintf(" and owner='%s'", query.Owner))
	}
	if query.Description != "" {
		sqlBuilder.WriteString(fmt.Sprintf(" and description='%s'", query.Description))
	}
	if query.Offset != 0 {
		sqlBuilder.WriteString(fmt.Sprintf(" offset %d", query.Offset))
	}
	if query.Limit != 0 {
		sqlBuilder.WriteString(fmt.Sprintf(" limit %d", query.Limit))
	}
}

func QueryBuilder(query SQuery) string {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("SELECT * FROM telant WHERE 1=1")
	whereBuilder(&sqlBuilder, query)
	return sqlBuilder.String()
}

func UpdateBuilder(value STelant, query SQuery) string {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("UPDATE telant")
	setBuilder(&sqlBuilder, value)
	whereBuilder(&sqlBuilder, query)
	return sqlBuilder.String()
}

func DeleteBuilder(query SQuery) string {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("DELETE FROM telant WHERE 1=1")
	whereBuilder(&sqlBuilder, query)
	return sqlBuilder.String()
}

func InsertBuilder(value STelant) string {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(fmt.Sprintf("INSERT INTO telant(name,owner,description) VALUES('%s','%s','%s')", value.Name, value.Owner, value.Description))
	return sqlBuilder.String()
}

func Query(db *sql.DB, query SQuery) ([]STelant, error) {
	rows, err := db.Query(QueryBuilder(query))
	if err != nil {
		return nil, err
	}
	result := RowsToArray(rows)
	return result, nil
}

func Update(db *sql.DB, value STelant, query SQuery) (int64, error) {
	result, err := db.Exec(UpdateBuilder(value, query))
	if err != nil {
		return 0, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func Delete(db *sql.DB, query SQuery) (int64, error) {
	result, err := db.Exec(DeleteBuilder(query))
	if err != nil {
		return 0, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func Insert(db *sql.DB, value STelant) (int64, error) {
	result, err := db.Exec(InsertBuilder(value))
	if err != nil {
		return 0, err
	}
	insertId, _ := result.LastInsertId()
	return insertId, nil
}
