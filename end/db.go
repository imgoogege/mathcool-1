// 此文件，是关于数据的储存到mysql或者其它数据库中的一个文件。
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"sync"
)

type DB struct {
	lock  sync.Mutex
	Value map[string]interface{}
}

//返回一个新的*sql.DB
func NewDB() (*sql.DB, error) {
	return sql.Open("mysql", dbUrl)

}

// 查询db 返回查询的行数
func (d *DB) SearchDB(db *sql.DB, sql string) {
	if err := db.Ping(); err != nil {
		glog.Error("数据库无法打开，错误是:", err)
	}
	rows, err := db.Query("SELECT * FROM mathcoolEnd  LIMIT 30")
	if err != nil {
		glog.Error(err)
	}
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)// 写入数据
		d.Value["uid"] = uid
		d.Value["username"] = username
		d.Value["department"] = department
		d.Value["created"] = created
	}

}

// 插入DB 返回插入的自增的那个字段的id数字。
func (d *DB) InsertDB(db *sql.DB,tem string) int64 {
	if err := db.Ping(); err != nil {
		glog.Error("数据库无法打开，错误是:", err)
	}
	stmt, err := db.Prepare("INSERT mathcoolEnd SET username=?,department=?,created=?")
	if err != nil {
		glog.Error("数据库出错", err)
	}

	res, err := stmt.Exec(d.Value)
	if err != nil {
		glog.Error(err)
	}
	n, err := res.LastInsertId()
	if err != nil {
		glog.Error(err)
	}
	return n
}

// 更新DB 返回更新的行数
func (d *DB) UpdateDB(db *sql.DB)int64 {
	if err := db.Ping(); err != nil {
		glog.Error("数据库无法打开，错误是:", err)
	}
	stmt, err := db.Prepare("UPDATE userinfo SET username=? WHERE uid=?")
	if err != nil {
		glog.Error("数据库出错", err)
	}

	res, err := stmt.Exec(d.Value)
	if err != nil {
		glog.Error(err)
	}
	n, err := res.LastInsertId()
	if err != nil {
		glog.Error(err)
	}
	return n
}

// 删除DB 返回删除的行数
func (d *DB) DeleteDB(db *sql.DB) int64{
	if err := db.Ping(); err != nil {
		glog.Error("数据库无法打开，错误是:", err)
	}
	stmt, err := db.Prepare("DELETE FROM userinfo WHERE uid=?")
	if err != nil {
		glog.Error("数据库出错", err)
	}

	res, err := stmt.Exec(d.Value)
	if err != nil {
		glog.Error(err)
	}
	n, err := res.LastInsertId()
	if err != nil {
		glog.Error(err)
	}
	return n
}

// 自增的的id
