package main

import (
	"flag"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func main(){

	engine := gin.Default()
	//
	store := cookie.NewStore([]byte("userPlus"))// cookie的名字 叫做 userPlus 然后内容是不可见的。
	engine.Use(sessions.Sessions("sessionPlus",store))//preAction
	engine.Use(user)
	//
	route(engine)
	err := engine.Run(":666")
	glog.Error(err)
}
func init(){
	gin.SetMode(gin.ReleaseMode)
	flag.Parse()
}

