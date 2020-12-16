package data

import (
	"database/sql"
	"errors"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"
	xerrors "github.com/pkg/errors"

	"geektime/Go-000/Week04/internal/biz"
)

const (
	MYSQLSRC = "root:123456@tcp(192.168.141.180:3306)/test?charset=utf8"
)

var ErrRecordNotFound = errors.New("record not found")

var Provider = wire.NewSet(NewDB, NewUserRepo)

func NewUserRepo(db *sql.DB) biz.UserRepo {
	return &userRepo{db: db}
}

func NewDB() (db *sql.DB, cf func(), err error) {
	//TODO: 需要从配置文件中加载
	db, err = sql.Open("mysql", MYSQLSRC)
	cf = func() {
		db.Close()
	}
	return
}

type userRepo struct {
	db *sql.DB
}

func (u *userRepo) GetUserById(id int64) (*biz.User, error) {
	user := &biz.User{}

	stmt, err := u.db.Prepare("select id, user_name, mobile from tbl_user where id=? limit 1")
	if err != nil {
		return nil, xerrors.Wrap(err, "prepare statement failed")
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Mobile)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = xerrors.Wrap(ErrRecordNotFound, "no result found in sql")
		} else {
			err = xerrors.Wrap(err, "queryRow failed")
		}
		return nil, err
	}
	return user, nil
}
