package main

import "github.com/gin-gonic/gin"

func route(engine *gin.Engine){
//	登陆
engine.GET("/signUp",signUpGET)
engine.POST("/signUp",signUpPOST)
//console
engine.GET("/",consoleGET)
//root set user
engine.GET("/setUser",setUserGET)
engine.POST("/rootSetUser",rootSetUserPOST)
// 删除文章
engine.GET("/deleteContent",deleteContentGET)
//删除用户
engine.GET("/deleteUser",deleteUserGET)
// 删除评论
engine.GET("/deleteComment",deleteCommentGET)
//删除员工
engine.GET("/deleteManagerUser",deleteManagerUserGET)
// 登出
engine.GET("/signOut",signOutGET)
// 显示内容

engine.POST("/content",contentPOST)

}