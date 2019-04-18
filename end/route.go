package main

import (
	"github.com/gin-gonic/gin"
)

func route(engine *gin.Engine){
	// 关于内容
	engine.GET("/w",wGET)//读
	engine.POST("/w",wPOST)//写
	engine.GET("/deleteW",deleteWGET)//删除文章
	// 登陆 注册 登出
	engine.POST("/signUp",signUpPOST)// 注册
	engine.POST("/signIn",signInPOST)//  登陆
	// 登出 客户端将session_Plus 从cookie中删除即可.
	// index
	//1 按照时间顺序 2 按照 评论数字 3 按照被赞的数字
	engine.GET("/indexArticleTitleList",indexArticleTitleListGET)
	// 设置评论(没有读取评论，因为读取文章的时候读取了),删除评论。
	engine.POST("/addComment",addCommentPOST)
	engine.GET("/readComment",readCommentGET)
	engine.GET("/deleteComment",deleteCommentGET)
	// USER,在user中要发送的数据除了基本的 用户信息，还有 用户的文章title 用户发出的评论，
	engine.GET("/user",userGET)
	engine.POST("/user",userPOST)
	engine.GET("/u/:userName",uUserNameGET)
	// index-right也就是热点推荐 只需要发送出title即可，而已不需要经过用户验证。后期加入用户验证，然后定向给用户推荐他喜欢的东西，这个再说。
	engine.GET("/rightHot",rightHotGET)
	// 搜索 发出 搜索的结果，甭管是 搜索什么 总归是 搜索 的内容，可以从 发来的 query上来判断即可。
	engine.GET("/search",searchGET)
	// 👍
	engine.GET("/addZan",addZanGET)
	// 增加浏览量
	engine.GET("/addSeeNumber",addSeeNumberGET)
	// 得到 配套试题

	engine.GET("/smallQ",smallQGET)

	// 邮箱是否被注册
	engine.GET("/isEmail",isEmailGET)
	// 用户名是否被注册
	engine.GET("isUserName",isUserNameGET)
	// job

	engine.GET("/job",jobGET)

	// 公式
	engine.GET("/formula",formulaGET)
	engine.GET("/formulaSearch",formulaSearchGET)
	// 试题
	engine.GET("/exam",examGET)
	engine.GET("/examSearch",examSearchGET)
	// 修改个人信息 比如密码 比如 登陆邮箱等
	engine.POST("/changeMS",changeMSPOST)

	// 获取 出题的人按照数量进行排行榜
	// 获取关于我的所有的出的题
	engine.GET("/testList",testListGET)
	engine.GET("/myTest",myTestGET)
	// 添加图片
	engine.GET("/addImage",addImageGET)
	// 删除图片
	engine.GET("/deleteImage",deleteImageGET)
	// 读取图片
	engine.GET("/readImage",readImageGET)
	// 微博登陆
	engine.POST("/weiboSignIn",isFirst,weiboSignInGET)
	engine.GET("/weiboSignOut",weiboSignOutGET)

	// 微博登出
}