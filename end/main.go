package main

import (
	"database/sql"
	"flag"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	flag.Parse()                 // 为了golang/glog来设置的 flag分析。
	gin.SetMode(gin.ReleaseMode) // 设置 gin的生产 模式。
	glog.V(1).Infoln("release: coastroad is runing.")
	glog.V(2).Infoln("debug:coastroad is runing.") // glog v 1 是调试阶段 2 是正式部署阶段。
	var err error
	dbHere, err = sql.Open("mysql", "root:359258Ls@tcp(localhost:3306)/mathcoolEnd")
	dbHere.SetMaxOpenConns(200)
	dbHere.SetMaxIdleConns(2)
	if err != nil {
		glog.Error(err)
	}
	err = dbHere.Ping()
	if err != nil {
		glog.Error("数据库无法连接", err)
	}
}



func main() {
	go func() {
		fmt.Println(http.ListenAndServe(":6060", nil))
	}()
	engine := gin.Default()
	// 中间件布置区域
	engine.Use(Config) // 配置需要每次访问都会经过。
	engine.Use(isUser) // 在每次请求中判断是否是用户，并且只需要设置makesureIsUser即可。
	engine.Use(PreAction)
	engine.Use(cores())// cores攻击预警
	//
	route(engine)
	engine.Run(":520")
}

