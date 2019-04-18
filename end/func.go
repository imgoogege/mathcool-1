package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// 目的是 为了得到数据, 测试成功
func wGET(ctx *gin.Context) {

	var contentID string
	var contentPlus string
	var title string
	var userID int64
	var contentValue string
	var userName string
	var userPlus string
	//var commentID int64
	var typeList int
	data := make(map[string]interface{}) // 设置一个map，为了生成json
	v := ctx.Query("contentPlus")
	if v == "" { //判断一下，如果这个query不存在的话，那么就
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "内容不存在",
		})
		return
	}
	// 直接从数据库中查找
	rows, err := dbHere.Query("SELECT content_id,content_plus,title,user_id,content_value,join_time,type_list FROM content  WHERE content_plus=?", v)
	defer rows.Close()
	if err != nil {
		glog.Error("从content表中查询数据的时候出错", err)
	}
	var join_time string
	for rows.Next() {
		rows.Scan(&contentID, &contentPlus, &title, &userID, &contentValue, &join_time, &typeList)
	}
	defer rows.Close()
	data["typeList"] = typeList
	data["join_time"] = join_time
	data["title"] = title
	data["contentValue"] = contentValue
	rows, err = dbHere.Query("SELECT user_name,user_plus FROM user WHERE user_id=?", userID)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
	}
	for rows.Next() {
		err = rows.Scan(&userName, &userPlus)
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法查询到用户",
			})
			return
		}
	}
	data["userName"] = userName
	data["userPlus"] = userPlus
	r, err := dbHere.Query("SELECT comment_id FROM comment WHERE content_id=?", contentID)
	defer r.Close()
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法查询到评论",
		})
		return
	}
	commentSlice := make([]map[string]interface{}, 0)
	for r.Next() { // 评论的slice
		var commentID int64
		if err = r.Scan(&commentID); err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法读取commentID",
			})
			return
		}
		ro, err := dbHere.Query("SELECT comment_value,user_id FROM comment WHERE comment_id=?", commentID)
		defer ro.Close()
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法查询到评论的实际内容",
			})
			return
		}
		commentMap := make(map[string]interface{})
		for ro.Next() { //
			vale := ""
			var id int64
			err := ro.Scan(&vale, &id)
			if err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "无法写入评论内容",
				})
				return
			}
			rww, err := dbHere.Query("SELECT user_name,user_plus FROM user where user_id=?", id)
			defer rww.Close()
			if err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "无法查询到评论的username",
				})
				return
			}
			for rww.Next() {
				var userName, userPlus string
				rww.Scan(&userName, &userPlus)
				commentMap["userName"] = userName
				commentMap["userPlus"] = userPlus
			}
			commentMap["commentValue"] = vale
		}
		commentSlice = append(commentSlice, commentMap)
	}
	data["comment"] = commentSlice
	ctx.JSON(http.StatusOK, data)
}

// 传入文章等内容 测试成功
func wPOST(ctx *gin.Context) {
	// 首先要如何生成一个w?v 这个v如何生成以及储存是个问题
	// 生成这个content_plus 使用  题目 + user_plus
	tl := ctx.Query("typeList")
	if tl == "" {
		tl = "0"
	}
	t, err := strconv.Atoi(tl)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法获取typelist",
		})
		return
	}
	tg := ctx.Query("tag")
	if tg == "" {
		tg = "0"
	}
	tgInt, err := strconv.Atoi(tg)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法获取tag",
		})
		return
	}

	// 前端 就是这样就ok了。
	//a := make(map[string]interface{})
	//a["title"] = "测试我的这个titleのO__O "
	//a["contentValue"] = "这个人生我知道的，所以我要来测试一下，所以就可以测试喽~！@#@~！~~！"
	//v1,_ := json.Marshal(a)
	//r := strings.NewReader(string(v1))
	//res,_ := http.Post("http://localhost:520/w?sessionPlus=4343&typeList=3","application/json",r)
	//
	m, _ := ctx.Get("makeUserIsUser")
	if m.(bool) {
		value, err := ioutil.ReadAll(ctx.Request.Body) // 得到从前端传入来的value值
		valueString := fastjson.GetString(value, "contentValue")
		title := fastjson.GetString(value, "title") // 得到 title
		var user_id int64
		user_plus := ctx.Query("sessionPlus") // 从前端去得到sessionID 也就是后端的sessionPLus
		s := SessionMap[user_plus]
		user_id = s.UserID
		Content_plus, _ := Encryption(s.UserID, title) // 得到content_plus
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法制作plus",
			})
			return
		}

		var mother_content_id int64
		mother_content_plus := ctx.Query("motherContentPlus")
		if mother_content_plus != "" {
			ro, err := dbHere.Query("SELECT content_id FROM content WHERE content_plus=?", mother_content_plus)
			defer ro.Close()
			if err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "无法找到content_id",
				})
				return
			}
			for ro.Next() {
				err = ro.Scan(&mother_content_id)
				if err != nil {
					glog.Error(err)
					ctx.JSON(http.StatusOK, gin.H{
						"success": "error",
						"data":    "无法找到contentid",
					})
					return
				}
			}
		}
		stmt, err := dbHere.Prepare("INSERT content SET content_plus=?,title=?,user_id=?,content_value=?,type_list=?,tag=?,mother_content_id=?")
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法插入数据",
			})
			return
		}

		r, err := stmt.Exec(Content_plus, title, user_id, valueString, t, tgInt, mother_content_id) // 将数据传入
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "插入数据失败",
			})
			return
		}
		rid, err := r.LastInsertId()
		if err != nil || rid == 0 {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"data":    "数据插入失败",
				"success": "error",
			})

			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "成功",
			"success": "ok",
		})
		// j接下来将这些数据传递到数据库中
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "请登录,系统检测您尚未登录",
		})
	}
}

// 删除文章 测试成功
func deleteWGET(ctx *gin.Context) {

	// 如何删除一篇文章呢？根据什么来删除呢？那肯定是1 首先要验证用户是否是本人，2 通过文章的content_plus将文章删除掉。
	//1 验证是否是用户本人。
	m, _ := ctx.Get("makeUserIsUser")
	if m.(bool) {
		var user int64
		session_plus := ctx.Query("sessionPlus")
		content_plus := ctx.Query("contentPlus")
		s := SessionMap[session_plus]
		userID := s.UserID
		rows, err := dbHere.Query("SELECT user_id FROM content WHERE content_plus=?", content_plus)
		defer rows.Close()
		if err != nil {
			glog.Error(err)
		}
		for rows.Next() {
			if err = rows.Scan(&user); err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "",
				})
				return
			}

		}
		if user != userID {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			glog.Error("攻击警告⚠️")
			return
		}
		// 已经验证完毕了，是主人本身，那么根据这个content_plus 删除文章即可
		stmt, err := dbHere.Prepare("DELETE FROM content WHERE content_plus=? and user_id=?")
		if err != nil {
			glog.Error(err)
		}
		_, err = stmt.Exec(content_plus, s.UserID)
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			return
		}

	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "",
		})
	}
}

// index 的文章title等信息的列表 测试成功
func indexArticleTitleListGET(ctx *gin.Context) {
	defer glog.Flush()

	// 首先排列的方式有几种 默认的是1 按照时间顺序 2 浏览量 3 按照被赞的个数
	// 要输出的格式是一个slice 那么如何每个slice中需要的数据有 1.文章的title 2 文章的content_plus （赋值给a的href）3author也就是user_name
	//4 被赞的个数 5 浏览量 6 日期
	sortType := ctx.Query("typeList") // 判断类型是什么
	page := ctx.Query("page")
	i, _ := strconv.Atoi(sortType)
	if i <= 1 {
		i = 1
	}
	query := "type_list=1 OR type_list=2"
	indexTypeOne(ctx, page, i, query)
}

// 排列文章列表 测试成功
func indexTypeOne(ctx *gin.Context, page string, typeValue int, query string) {
	data := make([]map[string]interface{}, 0) // 设置要输入到前端的这个data slice类型。
	var typeValueS string
	if typeValue == 1 {
		typeValueS = `join_time`
	} else if typeValue == 2 {
		typeValueS = `see_number`
	} else if typeValue == 3 {
		typeValueS = `zan`
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "typeList 出错",
		})
		return
	}
	if page == "" {
		page = "0"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		glog.Error(err, "出错页码:", pageInt)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "page出错",
		})
		return
	}
	pageInt *= 66
	queryValue := fmt.Sprintf("SELECT title,content_plus,user_id,zan,join_time,see_number,tag FROM content WHERE %s ORDER BY %s DESC LIMIT 66 OFFSET ?", query, typeValueS)
	rows, err := dbHere.Query(queryValue, pageInt)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法在数据库中查找到数据",
		})
	}
	for rows.Next() {
		var user_id int64
		dataMap := make(map[string]interface{})
		// 通过user_id 找到user_name
		var title, plus, join string
		var zan, see, tag int64
		err = rows.Scan(&title, &plus, &user_id, &zan, &join, &see, &tag)
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法将数据充值",
			})
			return
		}
		dataMap["tag"] = tag
		dataMap["title"] = title
		dataMap["contentPlus"] = plus
		dataMap["zan"] = zan
		dataMap["join_time"] = join
		dataMap["see_number"] = see
		r, err := dbHere.Query("SELECT user_name FROM user WHERE user_id=?", user_id)
		defer r.Close()
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			return
		}
		for r.Next() {
			var userName string
			if err = r.Scan(&userName); err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "",
				})
				return
			}
			dataMap["userName"] = userName
		}
		data = append(data, dataMap)
	}
	ctx.JSON(http.StatusOK, data)
}

// job

func jobGET(ctx *gin.Context) {
	defer glog.Flush()

	// 首先排列的方式有几种 默认的是1 按照时间顺序 2 浏览量 3 按照被赞的个数
	// 要输出的格式是一个slice 那么如何每个slice中需要的数据有 1.文章的title 2 文章的content_plus （赋值给a的href）3author也就是user_name
	//4 被赞的个数 5 浏览量 6 日期
	sortType := ctx.Query("typeList") // 判断类型是什么
	page := ctx.Query("page")
	i, _ := strconv.Atoi(sortType)
	if i <= 1 {
		i = 1
	}
	query := "type_list=3"
	indexTypeOne(ctx, page, 1, query)
}

// 公式
func formulaGET(ctx *gin.Context) {
	defer glog.Flush()

	// 首先排列的方式有几种 默认的是1 按照时间顺序 2 浏览量 3 按照被赞的个数
	// 要输出的格式是一个slice 那么如何每个slice中需要的数据有 1.文章的title 2 文章的content_plus （赋值给a的href）3author也就是user_name
	//4 被赞的个数 5 浏览量 6 日期
	sortType := ctx.Query("typeList") // 判断类型是什么
	page := ctx.Query("page")
	i, _ := strconv.Atoi(sortType)
	if i <= 1 {
		i = 1
	}
	query := "type_list=5"
	indexTypeOne(ctx, page, 2, query)
}

//试题
func examGET(ctx *gin.Context) {
	defer glog.Flush()

	// 首先排列的方式有几种 默认的是1 按照时间顺序 2 浏览量 3 按照被赞的个数
	// 要输出的格式是一个slice 那么如何每个slice中需要的数据有 1.文章的title 2 文章的content_plus （赋值给a的href）3author也就是user_name
	//4 被赞的个数 5 浏览量 6 日期
	sortType := ctx.Query("typeList") // 判断类型是什么
	page := ctx.Query("page")
	i, _ := strconv.Atoi(sortType)
	if i <= 1 {
		i = 1
	}
	query := "type_list=4"
	indexTypeOne(ctx, page, 2, query)
}

// 注册 将数据从前端搞过来 然后 先设置 user表 再设置session表 再加入到map中。 测试成功。
func signUpPOST(ctx *gin.Context) {

	use := new(User)
	session := new(Session)
	// 获取到客户端获取的url的query
	password := ctx.Query("password")    // 密码
	email := ctx.Query("email")          // E-mail
	verificationPassword(&password)      // 将密码进行过滤
	verificationEmail(&email)            //将email进行过滤
	use.UserName = ctx.Query("userName") //获取username
	Sex := ctx.Query("sex")
	sexInt, _ := strconv.Atoi(Sex)
	use.Email = email
	use.Sex = sexInt
	use.Year = ctx.Query("year")
	use.PhoneNumber = ctx.Query("phoneNumber")
	use.Description = ctx.Query("description")
	//  注册成功后要1 设定User这个数据库 2 设置 session这个数据库 3 将sessionID 提取出来，将这个值传入全局的map或者是redis缓存中
	// 设定 密码的plus
	dbPassword, saltValue := Encryption(int64(sexInt), password)
	// 将user这个数据库搞定
	// 要通过客户端将数据 设置为post的body数据传输json过来，不然通过query无法传送，有长度显示
	// 将user inset 进去数据库
	stmt, err := dbHere.Prepare("INSERT user SET user_plus=?,user_name=?,email=?,db_password=?,salt=?,sex=?,year=?,phone_number=?,description=?") // 将数据加入到数据库中
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法注册",
		})
		return
	}
	use.UserPlus, _ = Encryption(int64(use.Sex), use.UserName)
	_, err = stmt.Exec(use.UserPlus, use.UserName, use.Email, dbPassword, saltValue, use.Sex, use.Year, use.PhoneNumber, use.Description)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法再注册是入住信息",
		})
		return
	}
	// 拿出来user_id
	rows, err := dbHere.Query("SELECT user_id,join_time FROM user WHERE email=?", email)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法去得到user——id",
		})
		return
	}
	for rows.Next() {
		var id int64
		var time string
		if err = rows.Scan(&id, &time); err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法获取id",
			})
			return
		}
		use.UserID = id
		use.JoinTime = time

	}
	// 将session这个数据库搞定
	// 设置session_plus
	session_plus, _ := Encryption(use.UserID, use.UserName)
	//
	session.SessionPlus = session_plus
	stmt, err = dbHere.Prepare("INSERT session SET session_plus=?,user_id=?")
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法注入session",
		})
		return
	}
	_, err = stmt.Exec(session.SessionPlus, use.UserID)
	if err != nil {
		glog.Error(err, err)
		return
	}
	//取出来session_id
	rows, err = dbHere.Query("SELECT session_id FROM session WHERE user_id=?", use.UserID)
	defer rows.Close()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法获取sessionid",
		})
		glog.Error(err)
		return
	}
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法注入id",
			})
			return
		}
		session.SessionID = id
	}
	//完成user
	//
	session.SessionPlus = session_plus
	session.User = *use
	SessionMap[session.SessionPlus] = session
	ctx.JSON(http.StatusOK, gin.H{
		"success":     "ok",
		"sessionPlus": session.SessionPlus,
	})
}

// 登陆 测试成功
func signInPOST(ctx *gin.Context) { //要看
	// 登陆就是将sessionMap中的这个session_plus给前端即可。因为 注册的时候已经将这个map赋值了，
	// 但是 如果发生了故障就需要重新进行赋值给这个sessionID
	// 首先验证先验证邮箱是否是真的用户 然后如果是 那么开始验证是否密码正确，然后再开始验证在map中是否有值，如果没有就给一个。
	var salt string
	var sex int
	var username string
	var userID int64
	var sessionPlus string
	password := ctx.Query("password")
	email := ctx.Query("email")
	var dbPassword string

	rows, err := dbHere.Query("SELECT salt,db_password,sex,user_name,user_id FROM user WHERE  email=?", email)
	defer rows.Close()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "email没有注册",
		})
		glog.Error(err, "非法email:", email)
		return
	}
	for rows.Next() {
		if err = rows.Scan(&salt, &dbPassword, &sex, &username, &userID); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法获取注册的信息" + fmt.Sprintf("%v", err),
			})
			glog.Error(err)
			return
		}
	}
	idValue := strconv.Itoa(sex)
	if UpPassword := md5.Sum([]byte(password + idValue + salt)); fmt.Sprintf("%x", UpPassword) != dbPassword {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "密码错误",
		})
		return
	}
	// 这个时候已经知道是正确的 那么该怎么将 session_plus抛出去呢？

	rows, err = dbHere.Query("SELECT session_plus FROM session WHERE  user_id=?", userID)
	defer rows.Close()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "找不到session",
		})
		glog.Error(err)
		return
	}
	for rows.Next() {
		if err = rows.Scan(&sessionPlus); err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法将sessionPlus传递出去",
			})
			glog.Error(err)
			return
		}
	}
	if sessionPlus == "" {
		ctx.JSON(http.StatusOK, gin.H{"data": "无法获取sessionPlus", "success": "error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": sessionPlus, "success": "ok"}) // 全部验证完毕了将这个数据也就是这个后端的
}

// 增加评论// 测试成功。
func addCommentPOST(ctx *gin.Context) {
	// 首先要锁定几件事 1 哪篇文章 2 哪个用户 3 哪个评论
	// 先判断用户，然后找到文章，然后增加评论，评论里记录这个文章的id。（用content_plus找到content_id）
	defer glog.Flush()

	var userID, contentID int64
	sessionPlus := ctx.Query("sessionPlus") // 获取一个人的session_plus
	contentPlus := ctx.Query("contentPlus") // 获取一个文章的a标签指向的东西
	commentBody := ctx.Request.Body
	commentValueByte, err := ioutil.ReadAll(commentBody)
	commentValue := fastjson.GetString(commentValueByte, "data")
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法读取评论",
		})
		return
	}
	m, _ := ctx.Get("makeUserIsUser")
	fmt.Println("展开测试", m)
	if m.(bool) {
		s := SessionMap[sessionPlus]
		fmt.Println("测试s", *s)
		userID = s.UserID
		fmt.Println("测试评论：：：：", userID)
		row, err := dbHere.Query("SELECT content_id FROM content WHERE content_plus=?", contentPlus)
		defer row.Close()
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "找不到数据库里的contentid",
			})
			return
		}
		for row.Next() {
			err = row.Scan(&contentID)
			if err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "无法插入contentID1",
				})
				return
			}
		}
		smt, err := dbHere.Prepare("INSERT comment SET user_id=?,content_id=?,comment_value=?")
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "误差插入user——id",
			})
			return
		}
		_, err = smt.Exec(userID, contentID, commentValue)
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法插入userid2",
			})
			return
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "请登录",
		})
	}

}

// 读取评论 测试成功
func readCommentGET(ctx *gin.Context) {
	// 先通过content_plus找到content_id 然后通过content_id 找到comment
	defer glog.Flush()

	var contentId string
	var commentID int64
	contentPlus := ctx.Query("contentPlus")
	rows, err := dbHere.Query("SELECT content_id FROM content WHERE content_plus=?", contentPlus)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "",
		})
		return
	}
	for rows.Next() {
		if err = rows.Scan(&contentId); err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			return
		}
	}
	// 找到content
	rows, err = dbHere.Query("SELECT comment_value,user_id,comment_id FROM comment WHERE content_id=?", contentId)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "",
		})
		return
	}
	data := make([]map[string]interface{}, 0)
	var userid int64
	for rows.Next() {
		var commentValue, userName, userPlus string
		if err = rows.Scan(&commentValue, &userid, &commentID); err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			return
		}
		rows, err := dbHere.Query("SELECT user_name,user_plus FROM user WHERE user_id=?", userid)
		defer rows.Close()
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			return
		}
		for rows.Next() {
			if err = rows.Scan(&userName, &userPlus); err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "",
				})
				return
			}
		}
		data = append(data, map[string]interface{}{
			"userName":     userName,
			"userPlus":     userPlus,
			"commentValue": commentValue,
			"commentID":    commentID,
		})

	}
	ctx.JSON(http.StatusOK, data)

}

//删除评论 每个人删除自己的文章和评论只能在自己的那个已经登录的user才能删除。这么写简单呀。 测试成功。
func deleteCommentGET(ctx *gin.Context) {
	// 先判断人和评论是不是一个，然后再删除这一条评论
	//使用人和文章的id来删除comment 那么首先先获取到user_id和content_id
	defer glog.Flush()

	var userID, contentID int64
	sessionPlus := ctx.Query("sessionPlus") // 获取一个人的session_plus
	contentPlus := ctx.Query("contentPlus") // 获取一个文章的a标签指向的东西
	commentID := ctx.Query("commentID")
	if commentID == "" {
		ctx.JSON(http.StatusOK, "评论id为零无法删除")
		return
	}
	commentIdInt, err := strconv.Atoi(commentID)
	ifErrReturn(err, ctx, "无法将前端送来的string commentid转码")
	m, _ := ctx.Get("makeUserIsUser")
	if m.(bool) {
		s := SessionMap[sessionPlus]
		userID = s.UserID
		row, err := dbHere.Query("SELECT content_id FROM content WHERE content_plus=?", contentPlus)
		defer row.Close()
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "文章找不到",
			})
			return
		}
		for row.Next() {
			err = row.Scan(&contentID)
			if err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "无法找到文章的id",
				})
				return
			}
		}
		stmt, err := dbHere.Prepare("DELETE FROM comment WHERE user_id=? AND comment_id=?")
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法删除",
			})
			return
		}
		if _, err = stmt.Exec(userID, commentIdInt); err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法删除",
			})
			return
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "未登录",
		})
	}
}

// 用户 都有什么？那肯定是先登录状态才可以。 先登录 然后 测试成功。
func userGET(ctx *gin.Context) {
	defer glog.Flush()

	data := make(map[string]interface{})
	m, _ := ctx.Get("makeUserIsUser")
	if m.(bool) { // 已经登录了。
		typeList := ctx.Query("typeList")
		page := ctx.Query("page")
		if page == "" {
			page = "0"
		}
		sessionPlus := ctx.Query("sessionPlus") //得到sessionPlus
		s := SessionMap[sessionPlus]
		data["userName"] = s.UserName
		data["sex"] = s.Sex
		data["year"] = s.Year
		data["joinTime"] = s.JoinTime
		data["email"] = s.Email
		data["phoneNumber"] = s.PhoneNumber
		data["description"] = s.Description
		//获取文章title和plus
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			glog.Error(err, "出错页码:", pageInt)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "page出错",
			})
			return
		}
		pageInt *= 66
		queryValue := fmt.Sprintf("SELECT title,content_plus,type_list FROM content WHERE user_id=?  ORDER BY %s DESC LIMIT 66 OFFSET ?", "join_time")
		rows, err := dbHere.Query(queryValue, s.UserID, pageInt)
		defer rows.Close()
		//rows, err := dbHere.Query("SELECT title,content_plus,type_list FROM content WHERE user_id=?", s.UserID)
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法找到title",
			})
			return
		}
		contentData := make([]map[string]interface{}, 0)
		for rows.Next() {
			var title, plus string
			var typeList int
			rows.Scan(&title, &plus, &typeList)
			contentData = append(contentData, map[string]interface{}{
				"title":        title,
				"content_plus": plus,
				"type":         typeList,
			})
		}

		//开始获取comment
		row, err := dbHere.Query("SELECT content_id, comment_value,comment_id  FROM comment WHERE user_id=?", s.UserID)
		defer row.Close()
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法获取content——id",
			})
			return
		}
		commentData := make([]map[string]interface{}, 0)
		for row.Next() {
			var content_id, comment_id int64
			var commentValue string
			var contentPlus string
			row.Scan(&content_id, &commentValue, &comment_id)
			s, err := dbHere.Query("SELECT content_plus FROM content WHERE content_id=?", content_id)
			defer s.Close()
			if err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "无法获取contentpLUS",
				})
				return
			}
			for s.Next() {
				s.Scan(&contentPlus)
			}
			commentData = append(commentData, gin.H{
				"commentValue": commentValue,
				"commentID":    comment_id,
				"contentPlus":  contentPlus,
			})
		}

		if typeList == "" || typeList == "1" {
			typeList = "1"
		} else {
			typeList = "2"
		}
		if typeList == "1" {
			data["content"] = contentData
		} else {
			data["content"] = commentData
		}
		//
		ctx.JSON(http.StatusOK, data)
		//
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "没有登录",
		})
	}

}

// 修改 或者是 增加 用户的值。  测试成功。
func userPOST(ctx *gin.Context) {
	defer glog.Flush()

	sex := ctx.Query("sex")
	sexInt, _ := strconv.Atoi(sex)
	year := ctx.Query("year")
	phone := ctx.Query("phoneNumber")
	description := ctx.Query("description")
	sessionPlus := ctx.Query("sessionPlus")
	s := SessionMap[sessionPlus]
	m, _ := ctx.Get("makeUserIsUser")
	if m.(bool) { // 已经登录了。
		st, err := dbHere.Prepare("UPDATE user SET sex=?,year=?,phone_number=?,description=? WHERE user_id=?")
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			return
		}
		if _, err = st.Exec(sexInt, year, phone, description, s.UserID); err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			return
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": "error",
		"data":    "",
	})
}

// 热点 title 其实就是根据 评论值 将前 55位的内容 文章 试题  工作 等传递出去, 测试成功。
func rightHotGET(ctx *gin.Context) {
	defer glog.Flush()

	var contentPlus, title string
	data := make([]map[string]interface{}, 0)
	rows, err := dbHere.Query("SELECT content_plus, title FROM content WHERE type_list=1 OR type_list=2 OR type_list=3 OR type_list=4 OR type_list=5  ORDER BY see_number DESC LIMIT 50")
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "",
		})
		return
	}
	for rows.Next() {
		if err = rows.Scan(&contentPlus, &title); err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			return
		}
		data = append(data, map[string]interface{}{
			"contentPlus": contentPlus,
			"title":       title,
		})
	}
	ctx.JSON(http.StatusOK, data)
}

// 搜索 根据type来搜索 type都有什么？    1. 文章 2. 问答 3. 试题 4. 公式 5. 工作 6. 公式的配套小试题
// 输出的是什么？只有title和content_plus  测试成功。
func searchGET(ctx *gin.Context) {
	defer glog.Flush()

	//wType := ctx.Query("typeList") // 来访的内容type
	searchResult := ctx.Query("searchResult")
	page := ctx.Query("page")
	query := fmt.Sprintf("(type_list=%d OR type_list=%d OR type_list=%d OR type_list=%d OR type_list=%d)", 1, 2, 3, 4, 5)
	search(ctx, query, page, searchResult)

}

// 公式的search结果
func formulaSearchGET(ctx *gin.Context) {
	defer glog.Flush()

	//wType := ctx.Query("typeList") // 来访的内容type
	searchResult := ctx.Query("searchResult")
	page := ctx.Query("page")
	query := fmt.Sprintf("type_list=%d", 5)
	search(ctx, query, page, searchResult)
}

// 试题

func examSearchGET(ctx *gin.Context) {
	defer glog.Flush()

	//wType := ctx.Query("typeList") // 来访的内容type
	searchResult := ctx.Query("searchResult")
	page := ctx.Query("page")
	query := fmt.Sprintf("type_list=%d", 4)
	search(ctx, query, page, searchResult)
}
func search(ctx *gin.Context, wType string, page string, searchResult string) {
	pageInt, _ := strconv.Atoi(page)
	if pageInt <= 0 { // 为了米面 offset过0点。
		pageInt = 0
	}
	pageInt *= 66
	// sql的模糊搜索
	data := make([]map[string]interface{}, 0)
	content := func(data *[]map[string]interface{}, typeValue string, page int, searchResult string) {
		q := fmt.Sprintf("SELECT title,content_plus FROM content WHERE %s AND title LIKE ? ORDER BY see_number DESC LIMIT 66 OFFSET ?", typeValue)
		rows, err := dbHere.Query(q, "%"+searchResult+"%", page)
		defer rows.Close()
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "",
			})
			return
		}
		for rows.Next() {
			var title, contentPlus string
			if err = rows.Scan(&title, &contentPlus); err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "",
				})
				return
			}
			*data = append(*data, map[string]interface{}{
				"title":       title,
				"contentPlus": contentPlus,
			})

		}
	}
	content(&data, wType, pageInt, searchResult) // 将值 复制到 data里，为什么这里是使用&data呢？原因就是 append 每次都是返回给一个新的slice所以如果是
	// 传值那么 最后的这个slice是跟原来的data不一致的，所以穿指针就行了。跟slice实际上是指向底层array的指针没关系。
	ctx.JSON(http.StatusOK, data)
}

// 增加👍找到这个文章 然后 每次请求增加一次， 测试成功。
func addZanGET(ctx *gin.Context) {
	// 1. 使用query去得到这个文章的content_plus然后根据这个唯一值，将这个字段找出来，然后将zan +1 即可
	contentPlus := ctx.Query("contentPlus") // 去得到contentPlus 这个是 名目的id，必须是不能重复。
	defer glog.Flush()

	stmt, err := dbHere.Prepare("UPDATE content SET zan=zan+1 WHERE content_plus=?")
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "",
		})
		return
	}
	_, err = stmt.Exec(contentPlus)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "",
		})
		return
	}
}

// 增加 浏览量 直接增加
// 首先 增加的量增加的是see_number这个contenttable的字段值, 测试成功。
func addSeeNumberGET(ctx *gin.Context) {
	// 1. 使用query去得到这个文章的content_plus然后根据这个唯一值，将这个字段找出来，然后将see_numeber +1 即可
	contentPlus := ctx.Query("contentPlus") // 去得到contentPlus 这个是 名目的id，必须是不能重复。 // 得到这个views的值
	defer glog.Flush()

	stmt, err := dbHere.Prepare("UPDATE content SET see_number=see_number+1 WHERE content_plus=?")
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "",
		})
		return
	}
	_, err = stmt.Exec(contentPlus)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "",
		})
		return
	}

}

// 不用登陆，可以看到每个人的主页，这种主页，就是没有相关内容，没有任何操作，只读，并且只读信息不牵涉到机密。 测试成功。
func uUserNameGET(ctx *gin.Context) {
	defer glog.Flush()

	var userName, description, sex, year string
	username := ctx.Param("userName") // parm: 判断到底是哪个username
	rows, err := dbHere.Query("SELECT user_name,description,sex,year FROM user where user_name=?", username)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "查无此人",
		})
		return
	}
	for rows.Next() {
		if err = rows.Scan(&userName, &description, &sex, &year); err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法传入信息",
			})
			return
		}
	}
	var data = make(map[string]interface{})
	data["userName"] = userName
	data["description"] = description
	data["sex"] = sex
	data["year"] = year
	ctx.JSON(http.StatusOK, data)
}

// 得到配套试题
// 这种类型的小试题一律没有 title 当然也是储存在content中 // 测试成功
func smallQGET(ctx *gin.Context) {
	defer glog.Flush()

	var contentID int64
	var userID int64
	data := make([]map[string]interface{}, 0)
	c := new(Content)
	contentPlus := ctx.Query("contentPlus")
	rows, err := dbHere.Query("SELECT content_id FROM content WHERE content_plus=?", contentPlus)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法找到content_id",
		})
		return
	}
	for rows.Next() {
		err = rows.Scan(&contentID)
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法找到contentid",
			})
			return
		}
	}
	// 拉取子信息，客户端只要判断是不是空就ok了。不是空就显示呗。空就是没有，就不显示。 这个小试题跟其它的东西一样但是它没有title。
	rows, err = dbHere.Query("SELECT user_id,content_plus,content_value FROM content WHERE mother_content_id=?", contentID)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法找到user——id",
		})
		return
	}
	for rows.Next() {
		var contentPlus, contentValue string
		err = rows.Scan(&userID, &contentPlus, &contentValue)
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法匹配userid",
			})
			return
		}
		c.ContentPlus = contentPlus
		c.ContentValue = contentValue
		var userName, userPlus string
		rows, err := dbHere.Query("SELECT user_name,user_plus FROM user WHERE user_id=?", userID)
		defer rows.Close()
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "无法匹配username",
			})
			return
		}
		for rows.Next() {
			if err = rows.Scan(&userName, &userPlus); err != nil {
				glog.Error(err)
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "无法匹配username",
				})
				return
			}

			data = append(data, map[string]interface{}{
				"motherContentID": c.motherContentID,
				"contentPlus":     c.ContentPlus,
				"contentValue":    c.ContentValue,
				"userPlus":        userPlus,
				"userName":        userName,
			})
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":    data,
		"success": "ok",
	})

}

func isEmailGET(ctx *gin.Context) {

	defer glog.Flush()
	email := ctx.Query("email")
	rows, _ := dbHere.Query("SELECT user_name FROM user WHERE email=?", email)
	defer rows.Close()
	for rows.Next() {
		var t string
		rows.Scan(&t)
		if t != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "邮箱已经注册",
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "ok",
				"data":    "邮箱没有被使用，您可以进行注册",
			})
		}

	}

}

func isUserNameGET(ctx *gin.Context) {

	defer glog.Flush()
	userName := ctx.Query("username")
	rows, _ := dbHere.Query("SELECT user_id FROM user WHERE user_name=?", userName)
	defer rows.Close()
	for rows.Next() {
		var t int64
		rows.Scan(&t)
		if t != 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "用户名已经注册",
			})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "ok",
				"data":    "用户名没有被使用，您可以进行注册",
			})
		}

	}
}

// testList 出题人的一个输出
func testListGET(ctx *gin.Context) {
	defer glog.Flush()

	testList := make([]map[string]interface{}, 0)
	rows, err := dbHere.Query("SELECT user_id FROM content WHERE type_list=6")
	defer rows.Close()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法从数据库中获取testlist",
		})
		return
	}
	t := make(map[int]int, 0)
	for rows.Next() {
		var userID int
		rows.Scan(&userID)
		t[userID]++
	}
	userID, _ := sortMap(t)
	wait := sync.WaitGroup{}
	if len(userID) >= 20 {
		userID = userID[:19]
	}
	wait.Add(len(userID))
	var lock sync.Mutex
	for _, v := range userID {
		speed := make(chan struct{}, 10)
		go func(v int) {
			var userName, userPlus string
			speed <- struct{}{}
			defer wait.Done()
			lock.Lock()
			defer lock.Unlock()
			defer func() {
				if r := recover(); err != nil {
					glog.Error(r)
				}
			}()
			rows, err := dbHere.Query("SELECT user_name,user_plus FROM user WHERE user_id=?", v)
			defer rows.Close()
			if err != nil {
				glog.Error(err)
			}
			for rows.Next() {
				rows.Scan(&userName, &userPlus)
			}
			testList = append(testList, gin.H{
				"userName": userName,
				"userPlus": userPlus,
				"userID":   v,
			})
			<-speed

		}(v)

	}
	wait.Wait()
	data := make([]map[string]interface{}, 0)
	for _, v := range testList {
		data = append(data, gin.H{
			"number":   t[v["userID"].(int)],
			"userName": v["userName"],
			"userPlus": v["userName"],
		})
	}

	ctx.JSON(http.StatusOK, data)
	// 从slice中取出来数据。
}

// myTest 关于我自己的 试题的全部输出
func myTestGET(ctx *gin.Context) {
	defer glog.Flush()

	// 首先排列的方式有几种 默认的是1 按照时间顺序 2 浏览量 3 按照被赞的个数
	// 要输出的格式是一个slice 那么如何每个slice中需要的数据有 1.文章的title 2 文章的content_plus （赋值给a的href）3author也就是user_name
	//4 被赞的个数 5 浏览量 6 日期
	page := ctx.Query("page")
	b, _ := ctx.Get("makeUserIsUser")
	if b.(bool) {
		sesionPlus := ctx.Query("sessionPlus")
		s := SessionMap[sesionPlus]
		query := fmt.Sprintf("type_list=6 AND user_id=%d", s.UserID)
		// 这里的typeValue是指的按照上面样式排列 这里默认按照时间顺序。
		indexTypeOne(ctx, page, 1, query)
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "未登录",
		})
	}
}

//修改密码等信息
func changeMSPOST(ctx *gin.Context) {
	defer glog.Flush()

	b, _ := ctx.Get("makeUserIsUser")
	if b.(bool) {
		session := ctx.Query("sessionPlus")
		s := SessionMap[session]
		fmt.Println("测试email", s.Email)
		value, err := ioutil.ReadAll(ctx.Request.Body)
		ifErrReturn(err, ctx, "无法从前端获取获取信息")
		oldPassWord := fastjson.GetString(value, "oldPassWord")
		if oldPassWord != "" {
			p0 := strconv.FormatInt(int64(s.Sex), 10)
			if fmt.Sprintf("%x", md5.Sum([]byte(oldPassWord+p0+s.Salt))) != s.DBPassword {
				ctx.JSON(http.StatusOK, gin.H{
					"success": "error",
					"data":    "旧的密码不正确，请重新输入",
				})
				return
			}
			newPassWord := fastjson.GetString(value, "newPassWord")
			result, salt := Encryption(int64(s.Sex), newPassWord)
			st, err := dbHere.Prepare("UPDATE user SET salt=?,db_password=? WHERE user_id=?")
			ifErrReturn(err, ctx, "更新密码错误1")
			re, err := st.Exec(salt, result, s.UserID)
			ifErrReturn(err, ctx, "更新密码错误2")
			ls, err := re.LastInsertId()
			ifErrReturn(err, ctx, ls)
			ld, err := re.RowsAffected()
			ifErrReturn(err, ctx, ld)
			s.DBPassword = result
			s.Salt = salt
		}
		NewYear := fastjson.GetString(value, "newYear")
		if NewYear != "" {
			st, err := dbHere.Prepare("UPDATE user SET year=? where user_id=?")
			ifErrReturn(err, ctx, "更新年龄错误1")
			re, err := st.Exec(NewYear, s.UserID)
			ls, err := re.LastInsertId()
			ifErrReturn(err, ctx, ls)
			ld, err := re.RowsAffected()
			ifErrReturn(err, ctx, ld)
			ifErrReturn(err, ctx, "更新年龄错误2")
			s.Year = NewYear
		}
		newDescription := fastjson.GetString(value, "newDescription")
		if newDescription != "" {
			st, err := dbHere.Prepare("UPDATE user SET description=? WHERE user_id=?")
			ifErrReturn(err, ctx, "更新简介错误1")
			re, err := st.Exec(newDescription, s.UserID)
			ifErrReturn(err, ctx, "更新简介错误2")
			ls, err := re.LastInsertId()
			ifErrReturn(err, ctx, ls)
			ld, err := re.RowsAffected()
			ifErrReturn(err, ctx, ld)
			s.Description = newDescription
		}
		newPhoneNumber := fastjson.GetString(value, "newPhoneNumber")
		if newPhoneNumber != "" {
			st, err := dbHere.Prepare("UPDATE user SET phone_number=? where user_id=?")
			ifErrReturn(err, ctx, "更新电话号码错误1")
			result, err := st.Exec(newPhoneNumber, s.UserID)
			ls, err := result.LastInsertId()
			ifErrReturn(err, ctx, ls)
			ld, err := result.RowsAffected()
			ifErrReturn(err, ctx, ld)
			ifErrReturn(err, ctx, "更新电话号码错误2")
			s.PhoneNumber = newPhoneNumber
		}
		SessionMap[s.SessionPlus] = s
		ctx.JSON(http.StatusOK, gin.H{
			"success": "ok",
			"data":    "修改成功",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "未登录",
		})
	}
}

// 添加照片
func addImageGET(ctx *gin.Context) {
	defer glog.Flush()

	b, _ := ctx.Get("makeUserIsUser")
	if b.(bool) {
		sp := ctx.Query("sessionPlus")
		s := SessionMap[sp]
		iv := ctx.Query("imgValue")
		if iv == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "value是零",
			})
			return
		}
		st, err := dbHere.Prepare("INSERT img SET user_id=?,img_value=?")
		ifErrReturn(err, ctx, "无法插入img")
		r, err := st.Exec(s.UserID, iv)
		ifErrReturn(err, ctx, "无法插入")
		_, err = r.LastInsertId()
		ifErrReturn(err, ctx, "无法插入2")
		_, err = r.RowsAffected()
		ifErrReturn(err, ctx, "无法插入3")
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "未登录",
			"success": "error",
		})
	}
}

// 删除照片
func deleteImageGET(ctx *gin.Context) {
	defer glog.Flush()

	b, _ := ctx.Get("makeUserIsUser")
	if b.(bool) {
		sp := ctx.Query("sessionPlus")
		s := SessionMap[sp]
		iv := ctx.Query("imgID")
		if iv == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "error",
				"data":    "value是零",
			})
			return
		}
		id, err := strconv.Atoi(iv)
		ifErrReturn(err, ctx, "从客户端得到的imgid错误")
		st, err := dbHere.Prepare("DELETE FROM img WHERE user_id=? AND img_id=?")
		ifErrReturn(err, ctx, "无法插入img")
		r, err := st.Exec(s.UserID, id)
		ifErrReturn(err, ctx, "无法插入")
		_, err = r.LastInsertId()
		ifErrReturn(err, ctx, "无法插入2")
		_, err = r.RowsAffected()
		ifErrReturn(err, ctx, "无法插入3")
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "未登录",
			"success": "error",
		})
	}
}

// 读取文章

func readImageGET(ctx *gin.Context) {
	defer glog.Flush()

	b, _ := ctx.Get("makeUserIsUser")
	if b.(bool) {
		data := make(map[string]interface{})
		sp := ctx.Query("sessionPlus")
		s := SessionMap[sp]
		rows, err := dbHere.Query("SELECT img_value,img_id FROM img WHERE user_id=?", s.UserID)
		defer rows.Close()
		ifErrReturn(err, ctx, "无法取得数据")
		sliceImg := make([]map[string]interface{}, 0)
		for rows.Next() {
			var value string
			var img_id int64
			rows.Scan(&value, &img_id)
			sliceImg = append(sliceImg, gin.H{
				"imgID":    img_id,
				"imgValue": value,
			})
		}
		if len(sliceImg) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"success": "error", "data": "无法提取数据 "})
			return
		}
		data["data"] = sliceImg
		data["success"] = "ok"
		ctx.JSON(http.StatusOK, data)
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "未登录",
			"success": "error",
		})
	}
}

// 微博登陆
func weiboSignInGET(ctx *gin.Context) {

}
func isFirst(ctx *gin.Context) {
	plus := ctx.Query("uid")
	fmt.Println("测试uid", plus)
	b, _ := ctx.Get("makeUserIsUser")
	if b.(bool) {

	} else {
		v, _ := ioutil.ReadAll(ctx.Request.Body)
		defer ctx.Request.Body.Close()
		userName := fastjson.GetString(v, "userName")
		sex := fastjson.GetInt(v, "sex")
		description := fastjson.GetString(v, "description")
		s := NewSession()
		email := "weibo@weibo.weibo" + plus
		user_plus, _ := Encryption(1, userName)
		st, err := dbHere.Prepare("INSERT user SET email=?,user_name=?,sex=?,user_plus=?,salt=?,db_password=?,description=?")
		ifErrReturn(err, ctx, gin.H{"success": "error", "data": "无法insertusername"})
		_, err = st.Exec(email, userName, sex, user_plus, "1", "weibo", description)
		if err != nil {
			fmt.Println(err)
			return
		}
		rows, err := dbHere.Query("SELECT user_id FROM user WHERE user_name=?", userName)
		ifErrReturn(err, ctx, gin.H{"success": "error", "data": "query err"})
		defer rows.Close()
		var user_id int64
		for rows.Next() {
			err = rows.Scan(&user_id)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		st, err = dbHere.Prepare("INSERT session SET user_id=?,session_plus=?")
		ifErrReturn(err, ctx, gin.H{"success": "error", "data": "query err"})
		_, err = st.Exec(user_id, plus)
		if err != nil {
			fmt.Println(err)
			return
		}
		s.UserID = user_id
		s.Email = email
		s.SessionPlus = plus
		s.Sex = sex
		s.Description = description
		s.UserName = userName
		s.UserPlus = user_plus
		SessionMap[plus] = s
	}

}

// 微博登出
func weiboSignOutGET(ctx *gin.Context) {
	plus := ctx.Query("sessionPlus")
	if _, ok := SessionMap[plus]; !ok {

	} else {
		delete(SessionMap, plus)
	}
}
