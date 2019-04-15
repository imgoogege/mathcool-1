package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

var (
	redirectUrl string = "https://nudao.xyz"
)

func main() {
	engin := gin.Default()
	engin.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, redirectUrl)
	})
	engin.NoRoute(func(ctx *gin.Context) {
		// 跳转到请求的地方就ok了。
		ctx.Redirect(http.StatusMovedPermanently, path.Join(redirectUrl, ctx.Request.URL.Path))
	})
	err := engin.Run(":80")
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	flag.Parse()
}
