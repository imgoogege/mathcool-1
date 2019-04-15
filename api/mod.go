package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"sync"
)

// 处理数据库
type DB interface {
	Select(arg ...interface{}) (result []interface{}, err error) // 查询
	Delete(arg ...interface{}) (err error)                       //删除
	Update(arg ...interface{}) (err error)                       //更新
	Insert(arg ...interface{}) (err error)                       // 插入
}
type DbStruct struct {
	dbValue *sql.DB                // 数据库它本身。
	value   map[string]interface{} // 实际的内容
	lock    sync.Mutex             // 加锁
}

//OpenDB 打开一个数据库，并且是new一个新的数据库。
func (d *DbStruct) OpenDB(username, password, dbname string) {
	dataSourceName := username + ":" + password + "@/" + dbname
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		glog.Error(err)
	}
	defer db.Close()
	d.dbValue = db
}

// 查找
func (d *DbStruct) Select(arg ...interface{}) (result []interface{}, err error) {
	result = make([]interface{}, 0)
	if d.dbValue == nil {
		return nil, fmt.Errorf("没有初始化db数据库，请先使用d.OpenDB()")
	}
	rows, err := d.dbValue.Query("SELECT ? FROM ? WHERE age=? LIMIT 20", arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name interface{}
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil

}

// 删除
func (d *DbStruct) Delete(arg ...interface{}) (err error) {
	_, err = d.dbValue.Exec("DELETE FROM values WHERE ???", arg)
	return err
}

// 更新
func (d *DbStruct) Update(arg ...interface{}) (err error) {
	_, err = d.dbValue.Exec("UPDATE ? SET ?=?,?=? WHERE ? >?", arg)
	return err
}

// 插入
func (d *DbStruct) Insert(arg ...interface{}) (err error) {
	_, err = d.dbValue.Exec("INSERT INTO values (?,?) VALUES (?,?)", arg)
	return err
}
