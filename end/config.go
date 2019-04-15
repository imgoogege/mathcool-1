package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	serverURL = ""
	version   = ""
	dbUrl = ""
	port = ""
	portT = ""
)
// gin's middleware --config when your config.json is in ./file/config.json
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
	version = fastjson.GetString(rData, "version")
	dbUrl = fastjson.GetString(rData,"dbUrl")
	port = fastjson.GetString(rData,"port")
}
