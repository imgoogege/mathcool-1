package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"html/template"
	"net/http"
	"strconv"
)

// 首页显示
func signUpGET(ctx *gin.Context) {
	tem, err := template.ParseFiles("./view/layout.html", "./view/contentSignUp.html")
	ifErr(err, ctx, "无法渲染tem")
	fmt.Println(userMap["1234567"])
	tem.Execute(ctx.Writer, nil)
}

// 首页登陆// 要设置cookie
func signUpPOST(ctx *gin.Context) {
	defer dbHere.Close()
	defer glog.Flush()
	signUpUserName, _ := ctx.GetPostForm("signUpUserName")
	signUpPassword, _ := ctx.GetPostForm("signUpPassword")
	rows, err := dbHere.Query("SELECT salt,level,db,user_id,is_root  FROM managerUser WHERE user_name=?", signUpUserName)
	ifErr(err, ctx, "登陆取db出错")
	var salt string
	var level int
	var db string
	var user_id int64
	var isRoot int
	for rows.Next() {
		rows.Scan(&salt, &level, &db, &user_id,&isRoot)
		if salt == "" || level == 0 || db == "" || user_id == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "取db出错2",
			})
			return
		}
	}
	l := strconv.FormatInt(int64(level), 10)
	//level + username + password + salt
	if fmt.Sprintf("%x", md5.Sum([]byte(l+signUpUserName+signUpPassword+salt))) != db {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "无法验证密码",
			"success": "error",
		})
		return
	}
	mc := new(ManagerCookie)
	mc.name = "managerUser"
	s := new(ManagerSession)
	rows, err = dbHere.Query("SELECT session_plus FROM managerSession WHERE user_id=?", user_id)
	ifErr(err, ctx, "无法获取sessionplus")
	var sessionplus string
	for rows.Next() {
		rows.Scan(&sessionplus)
		if sessionplus == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法获取sessionPlus",
			})
			return
		}
	}
	mc.Set(sessionplus, ctx)
	s.userName = signUpUserName
	s.db = db
	s.level = level
	s.salt = salt
	s.isRoot = isRoot
	s.sessionPlus = sessionplus
	s.userID = user_id
	userMap[s.sessionPlus] = *s
	http.Redirect(ctx.Writer, ctx.Request, "/", 301)

}
func setUserGET(ctx *gin.Context) {
	data := make(map[string]interface{})
	v, _ := ctx.Cookie("managerUser")
	s := userMap[v]
	b, _ := ctx.Get("makeSureIsUser")
	if b.(bool) {
		if s.isRoot == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "错误root is null",
			})
		}
		fmt.Println("测试root", s.isRoot)
		if s.isRoot == 1 {
			// root用户
			data["isRoot"] = true

		} else if s.isRoot == 2 {
			//	普通用户
			data["isRoot"] = false

		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"data":    "数据库数据错误",
				"success": "error",
			})
		}
		tem, err := template.ParseFiles("./view/layout.html", "./view/contentIndex.html")
		ifErr(err, ctx, "无法渲染tem")
		tem.Execute(ctx.Writer, data)
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/signUp", 301)
	}
}

// 控制台
func consoleGET(ctx *gin.Context) {
	defer dbHere.Close()
	b, _ := ctx.Get("makeSureIsUser")
	fmt.Println("测试b", b)
	if b.(bool) {
		tem, err := template.ParseFiles("./view/layout.html", "./view/contentConsole.html")
		ifErr(err, ctx, "无法渲染tem")
		tem.Execute(ctx.Writer, nil)
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/signUp", 301)
	}
}

// 设置员工
func rootSetUserPOST(ctx *gin.Context) {
	defer glog.Flush()
	defer dbHere.Close()
	mc := new(ManagerCookie)
	mc.name = "managerUser"
	mc.Get("managerUser", ctx)
	userName, b := ctx.GetPostForm("indexRootSetUser")
	if !b {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "没有用户名",
			"success": "error",
		})
		return
	}
	level, b := ctx.GetPostForm("level")
	if !b {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "没有level",
			"success": "error",
		})
		return
	}
	b1, _ := ctx.Get("makeSureIsUser")
	if b1.(bool) {
		s := userMap[mc.value]
		fmt.Println("测试userMap", userMap, "测试mc", *mc)
		if s.isRoot != 1 {
			var err error
			ifErr(err, ctx, "不是root，有攻击的可能性")
			ctx.JSON(http.StatusOK, gin.H{
				"data":    "不是root，请停止访问",
				"success": "error",
			})
			return
		}
		fmt.Println("测试isroot", s.isRoot)
		st, err := dbHere.Prepare("INSERT managerUser SET user_name=?,level=?,is_root=?,db=?,salt=?")
		ifErr(err, ctx, "无法insertmanagerUser user_name")
		l, err := strconv.Atoi(level)
		ifErr(err, ctx, "level无法转化为int，怀疑有攻击")
		db, salt, password := encty(l, userName)
		r, err := st.Exec(userName, level, 2, db, salt)
		ifErr(err, ctx, "插入managerUser username时出错")
		last, err := r.LastInsertId()
		ifErr(err, ctx, last)
		row, err := r.RowsAffected()
		ifErr(err, ctx, row)

		rows, err := dbHere.Query("SELECT user_id FROM managerUser WHERE user_name=?", userName)
		ifErr(err, ctx, "无法获取user——id")
		var user_id int64
		for rows.Next() {
			rows.Scan(&user_id)
			if user_id == 0 {
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "出错",
				})
			}
		}
		//开始设置session
		st, err = dbHere.Prepare("INSERT managerSession SET user_id=?,session_plus=?")
		ifErr(err, ctx, "无法insertmanagerUser user_name")
		l, err = strconv.Atoi(level)
		ifErr(err, ctx, "level无法转化为int，怀疑有攻击")
		sessionPlus, _, _ := encty(l, userName)
		r, err = st.Exec(user_id, sessionPlus)
		ifErr(err, ctx, "插入managerUser username时出错")
		last, err = r.LastInsertId()
		ifErr(err, ctx, last)
		row, err = r.RowsAffected()
		ifErr(err, ctx, row)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "ok",
			"data":    password,
		})

	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "访问错误",
			"success": "error",
		})
	}
}

// 显示内容
func contentPOST(ctx *gin.Context) {
	defer dbHere.Close()
	defer glog.Flush()
	mc := new(ManagerCookie)
	mc.Get("managerUser", ctx)
	b, _ := ctx.Get("makeSureIsUser")
	if b.(bool) {
		tem, err := template.ParseFiles("./view/layout.html", "./view/contentContent.html")
		ifErr(err, ctx, "无法渲染tem")
		data := make(map[string]interface{})
		//sessionPlus := mc.value
		//ms := userMap[sessionPlus]
		result,e := ctx.GetPostForm("result")

		if e && result != ""{
			resultSlice := make([]map[string]interface{},0)
			var contentPlus,title string
			//query := fmt.Sprintf(`SELECT content_plus,title FROM content WHERE title LIKE '\%%s\%'`,result)
			//fmt.Println("SELECT content_plus,title FROM content WHERE title LIKE ?",result)
			rows,err := dbHere.Query("SELECT content_plus,title FROM content WHERE title LIKE ?","%"+result+"%")
			ifErr(err,ctx,"无法获取search数据")
			for rows.Next(){
				rows.Scan(&contentPlus,&title)
				if contentPlus == "" || title == "" {
					ctx.JSON(http.StatusOK,gin.H{
						"data":"无法获取search数据2",
						"success":"error",
					})
					return
				}
				resultSlice = append(resultSlice,gin.H{
					"contentPlus":contentPlus,
					"title":title,
				})
			}
			data["resultSlice"] = resultSlice

		}
		contentPlus,e := ctx.GetPostForm("contentPlus")
		if e && contentPlus != ""{
			resultSlice := make([]map[string]interface{},0)
			var title string
			rows,err := dbHere.Query("SELECT title FROM content WHERE content_plus=?",contentPlus)
			ifErr(err,ctx,"无法获取search数据")
			for rows.Next(){
				rows.Scan(&title)
				if title == "" {
					ctx.JSON(http.StatusOK,gin.H{
						"data":"无法获取search数据2",
						"success":"error",
					})
					return
				}
				resultSlice = append(resultSlice,gin.H{
					"contentPlus":contentPlus,
					"title":title,
				})
			}
			data["contentSlice"] = resultSlice

		}
		commentID,e := ctx.GetPostForm("commentID")
		var idd int
		if commentID == "" {
			idd = 0
		}
		idd,_ = strconv.Atoi(commentID)
		if e && idd != 0 {
			resultSlice := make([]map[string]interface{},0)
			var comment_value string
			rows,err := dbHere.Query("SELECT comment_value FROM comment WHERE comment_id=?",idd)
			ifErr(err,ctx,"无法获取search数据")
			for rows.Next(){
				rows.Scan(&comment_value)
				if comment_value == "" {
					ctx.JSON(http.StatusOK,gin.H{
						"data":"无法获取search数据2",
						"success":"error",
					})
					return
				}
				resultSlice = append(resultSlice,gin.H{
					"commentID":commentID,
					"commentValue":comment_value,
				})
			}
			fmt.Println("测试comment",resultSlice)
			data["commentSlice"] = resultSlice
		}
		userName,e := ctx.GetPostForm("userName")
		if e && userName != ""{
			resultSlice := make([]map[string]interface{},0)
			var user_plus string
			var user_id int64

			rows,err := dbHere.Query("SELECT user_plus,user_id FROM user WHERE user_name=?",userName)
			ifErr(err,ctx,"无法获取search数据")
			for rows.Next(){
				rows.Scan(&user_plus,&user_id)
				if user_plus == "" || user_id == 0 {
					ctx.JSON(http.StatusOK,gin.H{
						"data":"无法获取search数据2",
						"success":"error",
					})
					return
				}
				resultSlice = append(resultSlice,gin.H{
					"userPlus":user_plus,
					"userID":user_id,
					"userName":userName,
				})
			}
			data["userSlice"] = resultSlice
		}
		tem.Execute(ctx.Writer, data)
	} else {
	http.Redirect(ctx.Writer,ctx.Request,"/signUp",301)
	}
}

//删除文章
func deleteContentGET(ctx *gin.Context) {
	defer dbHere.Close()
	defer glog.Flush()
	mc := new(ManagerCookie)
	mc.Get("managerUser", ctx)
	b, _ := ctx.Get("makeSureIsUser")
	contentPlus := ctx.Query("contentPlus")
	if contentPlus == "" {
		ctx.JSON(http.StatusOK,gin.H{
			"success":"error",
			"data":"删除失败，无contentpLUS 的query",
		})
	}
	if b.(bool) {
		//sessionPlus := mc.value
		//ms := userMap[sessionPlus]
		st,err := dbHere.Prepare("DELETE FROM content WHERE content_plus=?")
		ifErr(err,ctx,"无法删除content")
		r,err := st.Exec(contentPlus)
		ifErr(err,ctx,"无法删除content")
		_,err = r.RowsAffected()
		ifErr(err,ctx,"无法删除content")
		_,err = r.LastInsertId()
		ifErr(err,ctx,"无法删除content")
		ctx.JSON(http.StatusOK,gin.H{"success":"ok","data":"删除成功",})
	} else {
		http.Redirect(ctx.Writer,ctx.Request,"/signUp",301)
	}
}

// 删除评论
func deleteCommentGET(ctx *gin.Context) {
	defer dbHere.Close()
	defer glog.Flush()
	mc := new(ManagerCookie)
	mc.Get("managerUser", ctx)
	b, _ := ctx.Get("makeSureIsUser")
	commentID := ctx.Query("commentID")
	if commentID == "" {
		ctx.JSON(http.StatusOK,gin.H{
			"success":"error",
			"data":"删除失败，无contentpLUS 的query",
		})
	}
	id,_ := strconv.Atoi(commentID)
	if b.(bool) {
		//sessionPlus := mc.value
		//ms := userMap[sessionPlus]
		st,err := dbHere.Prepare("DELETE FROM comment WHERE comment_id=?")
		ifErr(err,ctx,"无法删除comment1")
		r,err := st.Exec(id)
		ifErr(err,ctx,"无法删除comment2")
		_,err = r.RowsAffected()
		ifErr(err,ctx,"无法删除comment3")
		_,err = r.LastInsertId()
		ifErr(err,ctx,"无法删除comment4")
		ctx.JSON(http.StatusOK,gin.H{"success":"ok","data":"删除成功",})

	} else {
		http.Redirect(ctx.Writer,ctx.Request,"/signUp",301)
	}
}

//删除用户
func deleteUserGET(ctx *gin.Context) {
	defer dbHere.Close()
	defer glog.Flush()
	mc := new(ManagerCookie)
	mc.Get("managerUser", ctx)
	b, _ := ctx.Get("makeSureIsUser")
	userID := ctx.Query("userID")
	if  userID== "" {
		ctx.JSON(http.StatusOK,gin.H{
			"success":"error",
			"data":"删除失败，无contentpLUS 的query",
		})
	}
	id,_ := strconv.Atoi(userID)
	if b.(bool) {
		//sessionPlus := mc.value
		//ms := userMap[sessionPlus]
		st,err := dbHere.Prepare("DELETE FROM user WHERE user_id=?")
		ifErr(err,ctx,"无法删除user1")
		r,err := st.Exec(id)
		ifErr(err,ctx,"无法删除user2")
		_,err = r.RowsAffected()
		ifErr(err,ctx,"无法删除user3")
		_,err = r.LastInsertId()
		ifErr(err,ctx,"无法删除user4")
		ctx.JSON(http.StatusOK,gin.H{"success":"ok","data":"删除成功",})

	} else {
		http.Redirect(ctx.Writer,ctx.Request,"/signUp",301)
	}
}

//删除员工
func deleteManagerUserGET(ctx *gin.Context) {
	defer dbHere.Close()
	defer glog.Flush()
	mc := new(ManagerCookie)
	mc.Get("managerUser", ctx)
	b, _ := ctx.Get("makeSureIsUser")
	if b.(bool) {
		//sessionPlus := mc.value
		//ms := userMap[sessionPlus]


	} else {
		http.Redirect(ctx.Writer,ctx.Request,"/signUp",301)
	}
}

// 登出
func signOutGET(ctx *gin.Context) {
	defer dbHere.Close()
	defer glog.Flush()
	mc := new(ManagerCookie)
	mc.Get("managerUser", ctx)
	b, _ := ctx.Get("makeSureIsUser")
	if b.(bool) {
		//sessionPlus := mc.value
		//ms := userMap[sessionPlus]
		mc.Set("",ctx)
		ctx.JSON(http.StatusOK,gin.H{"success":"ok","data":"您已经退出"})
	}
}
