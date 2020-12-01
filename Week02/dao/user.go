package dao

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	xerrors "github.com/pkg/errors"

	"geektime/Go-000/Week02/model"
)

const (
	MYSQLSRC = "root:123456@tcp(127.0.0.1:3306)/testdb?charset=utf8"
)

var (
	db             *sql.DB
	ErrDaoNotFound = errors.New("Dao:No rows found")
)

func init() {
	var err error
	db, err = sql.Open("mysql", MYSQLSRC)
	if err != nil {
		panic(err)
	}
}

func DBConn() *sql.DB {
	return db
}

func GetUserInfo(name string) (*model.User, error) {
	user := &model.User{}

	stmt, err := DBConn().Prepare(
		"select user_name,user_age from tbl_user where user_name=? limit 1")
	if err != nil {
		return nil, xerrors.Wrap(err, "Prepare statement failed")
	}
	defer stmt.Close()

	err = stmt.QueryRow(name).Scan(&user.Name, &user.Age)
	if err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			err = xerrors.Wrap(ErrDaoNotFound, "no result found in sql")
		} else {
			err = xerrors.Wrap(err, "QueryRow failed")
		}
		return nil, err
	}
	return user, nil
}
