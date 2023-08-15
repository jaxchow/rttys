/*
 * @Author: 周家建
 * @Mail: zhou_0611@163.com
 * @Date: 2021-07-27 19:02:39
 * @Description:
 */

package main

import (
	"database/sql"
	"log"
	"strings"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	_ "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func instanceDB(str string) (*sql.DB, error) {
	sp := strings.Split(str, "://")
	// if len(sp) == 2 {
	// 	return sql.Open(sp[0], sp[1])
	// } else {
	return sql.Open("mysql", sp[1])
	// }
}

func initGORM(str string) (*gorm.DB, error) {
	sp := strings.Split(str, "://")
	var dialector gorm.Dialector
	log.Default().Println("dialector dns", sp[1])
	dialector = mysql.Open(sp[1])

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CasbinAdapter(db *sql.DB) *sqladapter.Adapter {
	a, err := sqladapter.NewAdapter(db, "sqlite3", "")
	if err != nil {
		panic(err)
	}
	return a
}
