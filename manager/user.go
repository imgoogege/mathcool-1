package main

import (
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net/http"
)

const mysqlAddress  ="root:359258Ls!@tcp(localhost:3306)/mathcoolEnd"
func user(ctx *gin.Context){
	var err error
	ctx.Set("makeSureIsUser",false)
	dbHere, err = sql.Open("mysql", mysqlAddress)
	if err != nil {
		glog.Error(err)
		ctx.Set("makeSureIsUser",false)
		ctx.JSON(http.StatusOK,gin.H{
			"success":"error",
			"data":"无法open sql",
		})
		return
	}
	err = dbHere.Ping()
	if err != nil {
		glog.Error("数据库无法连接", err)
		ctx.Set("makeSureIsUser",false)
		ctx.JSON(http.StatusOK,gin.H{
			"success":"error",
			"data":"无法ping",
		})
		return
	}
	v,err  := ctx.Cookie("managerUser")
	if err != nil||v =="" {
		ctx.Set("makeSureIsUser",false)
		return
	}
	ifErr(err,ctx,"无cookie")
	s := new(ManagerSession)
	if _,ok := userMap[v];!ok { // 如果不存在这个map中
		// 无法在map中去得到这个信息
		rows,err := dbHere.Query("SELECT user_id FROM managerSession WHERE session_plus=?",v)
		if err != nil{
			ctx.Set("makeSureIsUser",false)
			ctx.JSON(http.StatusOK,gin.H{
				"success":"error",
				"data":"无法从managerSession中取得到user_id",
			})
			return
		}
		var user_id int64
		for rows.Next(){
			rows.Scan(&user_id)
			if user_id == 0 {
				ctx.Set("makeSureIsUser",false)
				ctx.JSON(http.StatusOK,gin.H{
					"success":"error",
					"data":"验证一下user_id是否是0",
				})
				return
			}

		}
		rows,err = dbHere.Query("SELECT user_name,is_root,level,db,salt FROM managerUser WHERE user_id=?",user_id)
		ifErr(err,ctx,"无法取得username")
		var level int
		var user_name string
		var is_root int
		var db ,salt string
		for rows.Next(){
			rows.Scan(&user_name,&is_root,&level,&db,&salt)
			if user_name =="" || is_root == 0 || level == 0 {
				ctx.Set("makeSureIsUser",false)
				ctx.JSON(http.StatusOK,gin.H{
					"success":"error",
					"data":"username空",
				})
				return
			}
		}
		s.db = db
		s.salt = salt
		s.level = level
		s.userName = user_name
		s.userID = user_id
		s.isRoot= is_root
		s.sessionPlus = v
		userMap[v]= *s
		ctx.Set("makeSureIsUser",true)
	}else {
		ctx.Set("makeSureIsUser",true)
	}
}

func(mc *ManagerCookie)Get(name string,ctx *gin.Context)  {
	value,err := ctx.Cookie(name)
	if err != nil  {
		ctx.JSON(http.StatusOK,gin.H{
			"success":"error",
			"data":"找不到cookie",
		})
		return
	}
	mc.value = value
	mc.name= name
}
func(mc *ManagerCookie)Set(value string,ctx *gin.Context){
	mc.value = value
	ctx.SetCookie(mc.name,mc.value,1*60*60*24*30,"","",false,true)
}
