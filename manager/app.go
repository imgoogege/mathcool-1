package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net/http"
)

func preAction(ctx *gin.Context){
}

func ifErr(err error,ctx *gin.Context,value interface{}){
	if err != nil {
		glog.Error(err,value)
		ctx.JSON(http.StatusOK,gin.H{
			"success":"error",
			"data":fmt.Sprint(err,value),
		})
	}
	return
}


