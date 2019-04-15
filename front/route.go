package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func route(engin *gin.Engine) {
	engin.Use(Config)
	defer glog.Flush()
	// 静态路径设置
	pwd := pwdPbulic()
	if pwd != "" {
		engin.Static("/static", pwd)

	} else {
		glog.Error("静态路径路由执行错误")
	}

	//404
	engin.NoRoute(notFound)
	engin.GET("/noSign", noSignGET)
	engin.GET("/noContent", noContentGET)
	//主页面的GET方法
	engin.GET("/", indexGET)

	// search的 POST方法
	engin.POST("/search", searchPOST)

	// 捐赠
	engin.GET("/donate", donateGET)

	// 公司文化

	engin.GET("/culture", cultureGET)

	// 加入我们

	engin.GET("/join", joinGET)

	// 提出意见

	engin.GET("/advise", adviseGET)
	// 登陆
	engin.Any("/signIn", signIn)
	//engin.GET("/signIn", signInGET)
	//engin.POST("/signIn", signInPOST)

	// 注册
	engin.GET("/signUp", signUpGET)
	engin.POST("/signUp", signUpPOST)

	// 登出

	engin.GET("/signOut", signOutGET)

	// 联系我们

	engin.GET("/contact", contactGET)

	// 提出问题

	engin.GET("/question", questionGET)

	// test

	engin.GET("/test", test)

	// watch的综合路由
	engin.GET("/w", watchGET)

	// 个人信息
	engin.GET("/user", userGET)
	engin.POST("/user", userPOST)
	// user之下的个人的文章，消息和评论
	//公共个人页面
	engin.GET("/u/:userName", uGET)
	// 公式
	engin.GET("/formula", formulaGET)
	engin.POST("/formula", formulaPOST)

	// 试题
	engin.GET("/examQuestion", examQuestionGET)
	engin.POST("/examQuestion", examQuestionPOST)

	//修改信息
	engin.GET("/changeMS", changeMSGET)
	engin.POST("/changeMS", changeMSPOST)

	// 出题
	engin.GET("/makeExam", makeExamGET)
	engin.POST("/makeExam", makeExamPOST)

	// 我给大家出的题
	engin.GET("/myExam", myExamGET)

	// 我的排名
	engin.GET("/ranking", RankingGET)
	// 试题榜
	engin.GET("/testlist", testListGET)
	//all testlist

	//job

	engin.GET("/job", jobGET)

	//comment 评论
	engin.POST("/comment", commentPOST)
	// 删除评论
	engin.GET("deleteComment", deleteCommentGET)

	// 提交的容器
	engin.POST("/content", contentPOST)
	// 删除容器
	engin.GET("/deleteContent", deleteContentGET)

	//👍增加赞
	engin.GET("/addZan", addZanGET)

	// 增加image
	engin.POST("/addImage", addImagePOST)
	engin.GET("/deleteImage", deleteImageGET)

}
