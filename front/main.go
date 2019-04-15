// glog v  正式 1 debug 2。 -v 大于等于这个数字 才会显示，所以 正式部署 v=1 debug v=2
package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func init() {
	flag.Parse()                 // 为了golang/glog来设置的 flag分析。
	gin.SetMode(gin.ReleaseMode) // 设置 gin的生产 模式。
	glog.V(1).Infoln("release: coastroad is runing.")
	glog.V(2).Infoln("debug:coastroad is runing.") // glog v 1 是调试阶段 2 是正式部署阶段。
}
func main() {

	engin := gin.Default()
	// 检测cookie日期
	engin.Use(AddCookieTime) // 每次都检测cookie的maxage然后不行就续订。
	//然后续订成功
	// 前端用户判断
	engin.Use(USER)
	//
	route(engin)
	glog.Flush()
	err := engin.RunTLS(":443", pwdFile()+"/nudao.crt", pwdFile()+"/nudao.key")
	if err != nil {
		fmt.Println(err)
	}
}
