package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	serverURL = ""
	dbname    = ""
	version   = ""
)

// 一个gin的中间件，所谓中间件，就是每次请求都会先经过这里，当然你也可以结合group来让某些经过某些不经过
func Config(ctx *gin.Context) {
	defer glog.Flush()
	pwdFile := func() string {
		pwd, err := filepath.Abs(".")
		if err != nil {
			return ""
		} else {
			return filepath.Join(pwd, "file")
		}
	}
	config, err := os.Open(filepath.Join(pwdFile(), "config.json"))
	defer config.Close()
	if err != nil {
		glog.Error("在打开config的时候出现错误", err)
	}
	rData, err := ioutil.ReadAll(config)
	if err != nil {
		glog.Error("在config的时候发生了错误", err)
	}
	serverURL = fastjson.GetString(rData, "serverURL")
	dbname = fastjson.GetString(rData, "dbname")
	version = fastjson.GetString(rData, "version")
}
