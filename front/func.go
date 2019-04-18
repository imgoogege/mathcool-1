package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/googege/goo/imageDealWith"
	"github.com/valyala/fastjson"
	"html/template"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

// index获得得title只能是文章 问答 和工作。也就是 typelist 1 2 3, 测试成功。
func indexGET(ctx *gin.Context) {
	defer glog.Flush() // 将glog 传出
	page := ctx.Query("page")
	pre, pageArray, last, this := pageConversion(page, ctx)                                                                                  // 生成
	fc := NewFrontCookie("mathcool", "", "", "")                                                                                             // fc
	data := make(map[string]interface{})                                                                                                     // 这是一个要传递给前端的data
	err := fc.GetCookie("mathcool", ctx)                                                                                                     //先get到cookie这个value后，然后再从远程服务器得到数据
	glog.V(2).Infoln("/")                                                                                                                    // 作为数据的log
	tem, err := template.ParseFiles(temp("index.html", "text_center_index.html", "text.html", "text-left.html", "text-right-index.html")...) // 将文件 导出为完整的HTML
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "首页载入出错",
			"success": "error",
		})
		return
	}
	nav(ctx, "nudao.xyz-数学酷吗", data)
	typeList := ctx.Query("typeList")
	if typeList == "" {
		typeList = "1"
	}
	rangeValue := make([]map[string]interface{}, 0)
	query := fmt.Sprintf("/indexArticleTitleList?typeList=%s&page=%s", typeList, page)
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "无法从后台服务器获取数据",
			"success": "error",
		})
		return
	}
	v, err := fastjson.ParseBytes(value)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "无法从后台服务器获取数据",
			"success": "error",
		})
		return
	}
	for _, v := range v.GetArray() {
		var m = make(map[string]interface{})
		m["see_number"] = v.GetInt("see_number")
		tagValue := v.GetInt("tag")
		m["tag"] = tag(tagValue)
		m["contentPlus"] = string(v.GetStringBytes("contentPlus"))
		m["join_time"] = string(v.GetStringBytes("join_time"))
		m["title"] = string(v.GetStringBytes("title"))
		if len(v.GetStringBytes("userName")) >= 9 {
			m["userName"] = string(v.GetStringBytes("userName")[:8])
		} else {
			m["userName"] = string(v.GetStringBytes("userName"))
		}
		m["zan"] = v.GetInt("zan")
		rangeValue = append(rangeValue, m)
	}
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	data["pagePre"] = pre
	data["pageLast"] = last
	data["page"] = pageArray
	data["pageThis"] = this
	data["dataRange"] = rangeValue
	tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}
}

// 读取右边的那个热点信息 测试成功。
func rightHot(fc *FrontCookie, ctx *gin.Context, sliceMap *[]map[string]interface{}) {
	query := fmt.Sprintf("/rightHot?l=h")
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
	v, err := fastjson.ParseBytes(value)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "无法从后台服务器获取数据",
			"success": "error",
		})
		return
	}
	for _, v := range v.GetArray() {
		var m = make(map[string]interface{})
		m["contentPlus"] = string(v.GetStringBytes("contentPlus"))
		m["title"] = string(v.GetStringBytes("title"))
		*sliceMap = append(*sliceMap, m)
	}
}

// search 主页面的搜索
func searchPOST(ctx *gin.Context) {
	data := make(map[string]interface{})
	glog.V(2).Infoln("公式")
	defer glog.Flush()
	page := ctx.Query("page")
	searchResult, _ := ctx.GetPostForm("searchValue")
	pre, pageArray, last, this := pageConversion(page, ctx) // 生成
	tem, err := template.ParseFiles(temp("index.html", "text_center_search.html", "text.html", "text-left.html", "text-right-index.html")...)
	if err != nil {
		glog.Error(err)
	}
	fc := NewFrontCookie("mathcool", "", "", "") // fc
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	query := fmt.Sprintf("/search?&page=%s&searchResult=%s", page, searchResult)
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
	ifErrReturn(err, ctx, "无法从后端获取关于公式的信息")
	oneValue, err := fastjson.ParseBytes(value)
	ifErrReturn(err, ctx, "无法解析json")
	dataSlice := make([]map[string]interface{}, 0)
	for _, v := range oneValue.GetArray() {
		dataSlice = append(dataSlice, map[string]interface{}{
			"title":       string(v.GetStringBytes("title")),
			"contentPlus": string(v.GetStringBytes("contentPlus")),
		})
	}
	data["pagePre"] = pre
	data["pageLast"] = last
	data["page"] = pageArray
	data["pageThis"] = this
	data["searchData"] = dataSlice
	tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}

}

// 404的执行函数
func notFound(ctx *gin.Context) {
	defer glog.Flush()
	glog.V(2).Infoln("404")
	tem, err := template.ParseFiles(temp("index.html", "text_no.html", "text.html", "text-left.html", "text-right-index.html")...)
	if err != nil {
		glog.Error(err)
	}
	tem.Execute(ctx.Writer, nil)
	if err != nil {
		glog.Error(err)
	}

}

//捐赠的函数

func donateGET(ctx *gin.Context) {
	donate := new(SmallDonate)
	glog.V(2).Infoln("donate")
	defer glog.Flush()
	tem, err := template.ParseFiles(temp("small-donate.html", "small-nav.html", "index-donate.html")...)
	if err != nil {
		glog.Error(err)
	}
	donate.Value = []*Donate{
		{"https://github.com/googege", "https://raw.githubusercontent.com/realnudao/images/master/alipay.png"},
		{"https://coastroad.net", "https://raw.githubusercontent.com/realnudao/images/master/alipay.png"},
		{"https://github.com/googege", "https://raw.githubusercontent.com/realnudao/images/master/alipay.png"},
	}
	donate.Head_title = "捐赠"
	tem.Execute(ctx.Writer, donate)
	if err != nil {
		glog.Error(err)
	}

}

//加入我们

func joinGET(ctx *gin.Context) {
	glog.V(2).Infoln("join")
	defer glog.Flush()
	tem, err := template.ParseFiles(temp("small-nav.html", "index-join.html")...) //可以多但是不能少
	if err != nil {
		glog.Error(err)
	}
	tem.Execute(ctx.Writer, nil)
	if err != nil {
		glog.Error(err)
	}
}

// 企业文化
func cultureGET(ctx *gin.Context) {
	glog.V(2).Infoln("culture")
	defer glog.Flush()
	tem, err := template.ParseFiles(temp("small-nav.html", "index-culture.html")...)
	if err != nil {
		glog.Error(err)
	}
	tem.Execute(ctx.Writer, nil)
	if err != nil {
		glog.Error(err)
	}
}

// 登录 测试成功
func signIn(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		tem, err := template.ParseFiles(temp("small-nav.html", "index-signIn.html")...)
		if err != nil {
			glog.Error(err)
		}
		tem.Execute(ctx.Writer, nil)
		if err != nil {
			glog.Error(err)
		}
	} else if ctx.Request.Method == "POST" {
		// 得到form中的query, form中的name value
		fc := NewFrontCookie("mathcool", "", "", "")
		ctx.Request.ParseForm()
		email := ctx.Request.FormValue("signInEmail")
		pwd := ctx.Request.FormValue("signInPassword")
		query := fmt.Sprintf("/signIn?password=%s&email=%s", pwd, email)
		byteValue, _ := fc.PostValueToServerBySessionPlus(serverURL+query, "", nil)
		isOk := fastjson.GetString(byteValue, "success")
		if isOk != "ok" {
			data := fastjson.GetString(byteValue, "data")
			ctx.JSON(http.StatusOK, gin.H{"success": "error", "data": data})
			return
		}
		sessionPlus := fastjson.GetString(byteValue, "data")
		fc.value = sessionPlus
		fc.SetCookie(ctx)
		fmt.Fprint(ctx.Writer, "<script>window.location='/'</script>")
	}
}

// 注册 测试成功。
func signUpGET(ctx *gin.Context) {
	glog.V(2).Infoln("signin")
	defer glog.Flush()
	tem, err := template.ParseFiles(temp("small-nav.html", "index-signUp.html", "foot.html")...)
	if err != nil {
		glog.Error(err)
	}
	tem.Execute(ctx.Writer, nil)
	if err != nil {
		glog.Error(err)
	}
}

// 注册，的post 测试成功
func signUpPOST(ctx *gin.Context) {
	fc := NewFrontCookie("mathcool", "", "", "")
	email, b := ctx.GetPostForm("signUpEmail") // 获取到 邮箱信息
	if !b {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "请输入邮箱地址",
		})
		return
	}
	password, b := ctx.GetPostForm("signUpPassword") //密码
	if !b {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "请输入密码",
		})
		return
	}
	userName, b := ctx.GetPostForm("signUpUserName") //用户名
	if !b {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "请输入用户名",
		})
		return
	}
	Sex, b := ctx.GetPostForm("signUpSex") //性别
	if !b {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "请选择性别",
		})
		return
	}
	year, _ := ctx.GetPostForm("signUpUserYear")           //年龄
	phone, _ := ctx.GetPostForm("signUpPhoneNumber")       //电话
	description, _ := ctx.GetPostForm("signUpDescription") //简介
	query := fmt.Sprintf("/signUp?email=%s&password=%s&userName=%s&sex=%s&year=%s&phoneNumber=%s&description=%s", email, password, userName, Sex, year, phone, description)
	value, err := fc.PostValueToServerBySessionPlus(serverURL+query, "application/json", nil)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    err,
		})
		return
	}

	if su := fastjson.GetString(value, "success"); su != "ok" {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    err,
		})
		return
	}
	sessionPlus := fastjson.GetString(value, "sessionPlus")
	fc.value = sessionPlus
	fc.SetCookie(ctx)
	fmt.Fprint(ctx.Writer, fmt.Sprint("<script>window.location='/'</script>"))
}

// 登出 功能测试成功。
func signOutGET(ctx *gin.Context) {
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.DeleteValue(ctx)
	fmt.Fprint(ctx.Writer, fmt.Sprint("<script>window.location='/'</script>"))

}

// 提出意见 测试成功
func adviseGET(ctx *gin.Context) {
	data := make(map[string]interface{})
	glog.V(2).Infoln("advise")
	defer glog.Flush()
	tem, err := template.ParseFiles(temp("small-editor.html", "small-nav.html", "index-advise.html")...)
	if err != nil {
		glog.Error(err)
	}
	data["typeList"] = "7"
	tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}
}

// 联系我们
func contactGET(ctx *gin.Context) {
	glog.V(2).Infoln("advise")
	defer glog.Flush()
	tem, err := template.ParseFiles(temp("small-nav.html", "index-contact.html")...)
	if err != nil {
		glog.Error(err)
	}
	tem.Execute(ctx.Writer, nil)
	if err != nil {
		glog.Error(err)
	}
}

// 提出问题
func questionGET(ctx *gin.Context) {
	defer glog.Flush()
	fc := NewFrontCookie("mathcool", "", "", "")
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		data := make(map[string]interface{})
		data["motherContentPlus"] = ctx.Query("motherContentPlus")
		glog.V(2).Infoln("/")
		tem, err := template.ParseFiles(temp("small-editor.html", "index.html", "text_center_question.html", "text.html", "text-left.html", "text-right-index.html")...)
		if err != nil {
			glog.Error(err)
		}
		nav(ctx, "nudao.xyz---数学酷吗", data)
		hot := make([]map[string]interface{}, 0)
		rightHot(fc, ctx, &hot)
		data["rightHot"] = hot
		data["typeList"] = ctx.Query("typeList")
		data["advise"] = testDataRight
		tem.Execute(ctx.Writer, data)
		if err != nil {
			glog.Error(err)
		}
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", http.StatusMovedPermanently)
	}
}

//test

func test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"a":   "12",
		"1":   "332",
		"332": 333232,
	})
}

// watch 详细页面，这个页面就是描述的是观看页面，使用跟YouTube类似的模式，watch?v=fdfdfew2332fdf

// 规则1 ，如果 通过检索，该contentPlus存在comment那么就将这个comment打进去
// 规则2，如果这个contentPlus拥有 配套试题，那么将配套试题打进去
//  data["title"] jotin_time author 然后 data["comment"] = []map[stirng]interface{} [{username userplus commnet value jointime_time}]
// data[smallTest] = []map[string]interface{}{ username userplus value joti_toime }
func watchGET(ctx *gin.Context) {
	defer glog.Flush()
	data := make(map[string]interface{})
	fc := NewFrontCookie("mathcool", "", "", "")
	contentPlus := ctx.Query("contentPlus") // 获取contentPlus
	if contentPlus == "" {
		http.Redirect(ctx.Writer, ctx.Request, "/noContent", http.StatusMovedPermanently)
	}
	query := fmt.Sprintf("/w?contentPlus=%s", contentPlus)
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    "无法从服务器查询到数据",
		})
		return
	}
	// 每次浏览都增加一次浏览量
	addView(contentPlus, fc)
	// fastJson在parse 从后端传来的数据
	oneValue, err := fastjson.ParseBytes(value)
	ifErrReturn(err, ctx, "fastjson无法解析数据")
	data["contentValue"] = template.HTML(string(oneValue.GetStringBytes("contentValue")))
	data["title"] = string(oneValue.GetStringBytes("title"))
	if data["title"].(string) == "" {
		http.Redirect(ctx.Writer, ctx.Request, "/noContent", http.StatusMovedPermanently)
	}
	data["userName"] = string(oneValue.GetStringBytes("userName"))
	data["userPlus"] = string(oneValue.GetStringBytes("userPlus"))
	data["join_time"] = string(oneValue.GetStringBytes("join_time"))
	// 传入关于本文章的评论
	slice := make([]map[string]interface{}, 0)
	for _, v := range oneValue.GetArray("comment") {
		slice = append(slice, map[string]interface{}{
			"commentValue": template.HTML(string(v.GetStringBytes("commentValue"))),
			"userName":     string(v.GetStringBytes("userName")),
			"userPlus":     string(v.GetStringBytes("userPlus")),
		})
	}
	data["comment"] = slice
	// 评论数据完结。
	// 查看是否有小题，也就是将 这个seeesion 当做 motherid即可。然后找到一些 文章 返回 title username userPlus contentValue即可 不要返回 contentPlus
	query2 := fmt.Sprintf("/smallQ?contentPlus=%s", contentPlus)
	value2, err := fc.GetValueFromServerBySessionPlus(serverURL + query2)
	ifErrReturn(err, ctx, "寻找配套试题出错")
	oneValue2, err := fastjson.ParseBytes(value2)
	ifErrReturn(err, ctx, "无法找到配套试题")
	slice2 := make([]map[string]interface{}, 0)
	for _, v := range oneValue2.GetArray("data") {
		slice2 = append(slice2, map[string]interface{}{
			"contentPlus":  string(v.GetStringBytes("contentPlus")),
			"contentValue": template.HTML(string(v.GetStringBytes("contentValue"))),
			"userName":     string(v.GetStringBytes("userName")),
			"userPlus":     string(v.GetStringBytes("userPlus")),
		})
	}
	if len(slice2) == 0 {
		data["isHave"] = false
	} else {
		data["isHave"] = true
	}
	data["small"] = slice2
	tem, err := template.ParseFiles(temp("small-comment.html", "smallExam.html", "index.html", "text_center_article.html", "text.html", "text-left.html", "text-right-index.html")...) // 将文件 导出为完整的HTML
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "首页载入出错",
			"success": "error",
		})
		return
	}
	nav(ctx, "数学酷吗---"+string(oneValue.GetStringBytes("title")), data)
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	data["contentPlus"] = contentPlus
	if oneValue.GetInt("typeList") == 5 {
		data["isTypeList"] = true
	} else {
		data["isTypeList"] = false
	}
	err = tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}

}

// 个人用户的页面
func userGET(ctx *gin.Context) {
	defer glog.Flush()
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	b, _ := ctx.Get("makeSureUser")

	if b.(bool) {
		typeList := ctx.Query("typeList")
		if typeList == "" {
			typeList = "1"
		}
		tem, err := template.ParseFiles(temp("index.html", "text_center_user.html", "text.html", "text-left.html", "text-right-index.html")...) // 将文件 导出为完整的HTML
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"data":    "首页载入出错",
				"success": "error",
			})
			return
		}
		data := make(map[string]interface{})
		page := ctx.Query("page")
		pre, pageArray, last, this := pageConversion(page, ctx) // 生成
		nav(ctx, "nudao.xyz-数学酷吗", data)
		hot := make([]map[string]interface{}, 0)
		rightHot(fc, ctx, &hot)
		data["rightHot"] = hot
		data["pagePre"] = pre
		data["pageLast"] = last
		data["page"] = pageArray
		data["pageThis"] = this
		query := fmt.Sprintf("/user?typeList=%s&page=%s", typeList, page)
		value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
		ifErrReturn(err, ctx, "无法从服务端获取信息")
		data["userName"] = fastjson.GetString(value, "userName")
		data["sex"] = contentSex(fastjson.GetInt(value, "sex"))
		data["year"] = fastjson.GetString(value, "year")
		data["joinTime"] = fastjson.GetString(value, "joinTime")
		data["email"] = fastjson.GetString(value, "email")
		data["phoneNumber"] = fastjson.GetString(value, "phoneNumber")
		data["description"] = fastjson.GetString(value, "description")
		oneValue, err := fastjson.ParseBytes(value)
		ifErrReturn(err, ctx, "无法解析array")
		sliceData := make([]map[string]interface{}, 0)
		for _, v := range oneValue.GetArray("content") {
			sliceData = append(sliceData, gin.H{
				"title":              string(v.GetStringBytes("title")),
				"contentPlus":        string(v.GetStringBytes("content_plus")),
				"contentPlusComment": string(v.GetStringBytes("contentPlus")),
				"typeList":           contentTypeList(v.GetInt("type")),
				"commentValue":       template.HTML(v.GetStringBytes("commentValue")),
				"commentID":          v.GetInt64("commentID"),
			})
		}
		if typeList == "1" {
			data["css1"] = "btn-info"
			data["typeList"] = true
		} else {
			data["css3"] = "btn-info"
			data["typeList"] = false
		}
		data["sliceData"] = sliceData
		//user上传的图库
		sliceImageVa := make([]map[string]interface{}, 0)
		value1, err := fc.GetValueFromServerBySessionPlus(serverURL + "/readImage?l=h")
		ifErrReturn(err, ctx, "无法获取照片")
		oneValue1, err := fastjson.ParseBytes(value1)
		ifErrReturn(err, ctx, "无法解析value")
		for _, v := range oneValue1.GetArray("data") {
			sliceImageVa = append(sliceImageVa, gin.H{
				"imgValue": string(v.GetStringBytes("imgValue")),
				"imgID":    v.GetInt64("imgID"),
			})
		}
		data["imgSlice"] = sliceImageVa
		//
		err = tem.Execute(ctx.Writer, data)
		ifErrReturn(err, ctx, "无法聚集成html")
	} else {
		ctx.JSON(200, "没登录")
	}
}

// 目的是为了 修改密码等信息
func userPOST(ctx *gin.Context) {

}
func uGET(ctx *gin.Context) {
	defer glog.Flush()
	data := make(map[string]interface{})
	userName := ctx.Param("userName")
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + "/u/" + userName + "?l=h")
	ifErrReturn(err, ctx, "从后端获取信息错误")
	nav(ctx, "nudao.xyz-数学酷吗", data)
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	data["userName"] = fastjson.GetString(value, "userName")
	if fastjson.GetString(value, "sex") == "1" {
		data["sex"] = "男"
	} else if fastjson.GetString(value, "year") == "2" {
		data["sex"] = "女"
	} else {
		data["sex"] = "保密"
	}
	data["year"] = fastjson.GetString(value, "year")
	data["description"] = fastjson.GetString(value, "description")
	tem, err := template.ParseFiles(temp("index.html", "text_center_u.html", "text.html", "text-left.html", "text-right-index.html")...) // 将文件 导出为
	if err != nil {
		glog.Error(err)
	}
	tem.Execute(ctx.Writer, data)

}

// 公式 ，显示出来公式 显示的内容就是 全部的 公式 ,然后 search search后 还是这个页面 但是 显示的东西是跟搜索相关的东西。
func formulaGET(ctx *gin.Context) {
	data := make(map[string]interface{})
	glog.V(2).Infoln("公式")
	defer glog.Flush()
	page := ctx.Query("page")
	pre, pageArray, last, this := pageConversion(page, ctx) // 生成
	tem, err := template.ParseFiles(temp("index.html", "text_center_formula.html", "text.html", "text-left.html", "text-right-index.html")...)
	if err != nil {
		glog.Error(err)
	}
	fc := NewFrontCookie("mathcool", "", "", "") // fc
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	query := fmt.Sprintf("/formula?&page=%s", page)
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
	ifErrReturn(err, ctx, "无法从后端获取关于公式的信息")
	oneValue, err := fastjson.ParseBytes(value)
	ifErrReturn(err, ctx, "无法解析json")
	dataSlice := make([]map[string]interface{}, 0)
	for _, v := range oneValue.GetArray() {
		dataSlice = append(dataSlice, map[string]interface{}{
			"title":       string(v.GetStringBytes("title")),
			"contentPlus": string(v.GetStringBytes("contentPlus")),
		})
	}
	data["formulaData"] = dataSlice
	nav(ctx, "nudao.xyz---公式", data)
	data["pagePre"] = pre
	data["pageLast"] = last
	data["page"] = pageArray
	data["pageThis"] = this
	tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}
}

func formulaPOST(ctx *gin.Context) {
	data := make(map[string]interface{})
	glog.V(2).Infoln("公式")
	defer glog.Flush()
	page := ctx.Query("page")
	searchResult, _ := ctx.GetPostForm("formulaValue")
	pre, pageArray, last, this := pageConversion(page, ctx) // 生成
	tem, err := template.ParseFiles(temp("index.html", "text_center_formulaPOST.html", "text.html", "text-left.html", "text-right-index.html")...)
	if err != nil {
		glog.Error(err)
	}
	fc := NewFrontCookie("mathcool", "", "", "") // fc
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	query := fmt.Sprintf("/formulaSearch?page=%s&searchResult=%s", page, searchResult)
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
	ifErrReturn(err, ctx, "无法从后端获取关于公式的信息")
	oneValue, err := fastjson.ParseBytes(value)
	ifErrReturn(err, ctx, "无法解析json")
	dataSlice := make([]map[string]interface{}, 0)
	for _, v := range oneValue.GetArray() {
		dataSlice = append(dataSlice, map[string]interface{}{
			"title":       string(v.GetStringBytes("title")),
			"contentPlus": string(v.GetStringBytes("contentPlus")),
		})
	}
	data["formulaData"] = dataSlice
	nav(ctx, "nudao.xyz---公式", data)
	data["pagePre"] = pre
	data["pageLast"] = last
	data["page"] = pageArray
	data["pageThis"] = this
	tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}
}

//试题
func examQuestionGET(ctx *gin.Context) {
	data := make(map[string]interface{})
	glog.V(2).Infoln("试题")
	defer glog.Flush()
	page := ctx.Query("page")
	pre, pageArray, last, this := pageConversion(page, ctx) // 生成
	tem, err := template.ParseFiles(temp("index.html", "text_center_examQuestion.html", "text.html", "text-left.html", "text-right-index.html")...)
	if err != nil {
		glog.Error(err)
	}
	fc := NewFrontCookie("mathcool", "", "", "") // fc
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	query := fmt.Sprintf("/exam?&page=%s", page)
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
	ifErrReturn(err, ctx, "无法从后端获取关于试题的信息")
	oneValue, err := fastjson.ParseBytes(value)
	ifErrReturn(err, ctx, "无法解析json")
	dataSlice := make([]map[string]interface{}, 0)
	for _, v := range oneValue.GetArray() {
		dataSlice = append(dataSlice, map[string]interface{}{
			"title":       string(v.GetStringBytes("title")),
			"contentPlus": string(v.GetStringBytes("contentPlus")),
		})
	}
	data["examQuestion"] = dataSlice
	nav(ctx, "nudao.xyz--试题", data)
	data["pagePre"] = pre
	data["pageLast"] = last
	data["page"] = pageArray
	data["pageThis"] = this
	tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}
}

func examQuestionPOST(ctx *gin.Context) {
	data := make(map[string]interface{})
	glog.V(2).Infoln("试题")
	defer glog.Flush()
	page := ctx.Query("page")
	searchResult := ctx.Query("examQuestion")
	pre, pageArray, last, this := pageConversion(page, ctx) // 生成
	tem, err := template.ParseFiles(temp("index.html", "text_center_examQuestionPOST.html", "text.html", "text-left.html", "text-right-index.html")...)
	if err != nil {
		glog.Error(err)
	}
	fc := NewFrontCookie("mathcool", "", "", "") // fc
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	query := fmt.Sprintf("/examSearch?&page=%s&searchResult=%s", page, searchResult)
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
	ifErrReturn(err, ctx, "无法从后端获取关于试题的信息")
	oneValue, err := fastjson.ParseBytes(value)
	ifErrReturn(err, ctx, "无法解析json")
	dataSlice := make([]map[string]interface{}, 0)
	for _, v := range oneValue.GetArray() {
		dataSlice = append(dataSlice, map[string]interface{}{
			"title":       string(v.GetStringBytes("title")),
			"contentPlus": string(v.GetStringBytes("contentPlus")),
		})
	}
	data["examQuestion"] = dataSlice
	nav(ctx, "nudao.xyz--试题", data)
	data["pagePre"] = pre
	data["pageLast"] = last
	data["page"] = pageArray
	data["pageThis"] = this
	tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}
}

//修改你的密码

func changeMSGET(ctx *gin.Context) {
	defer glog.Flush()
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		data := make(map[string]interface{})
		tem, err := template.ParseFiles(temp("index.html", "text_center_changeMS.html", "text.html", "text-left.html", "text-right-index.html")...) // 将文件 导出为完整的HTML
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"data":    "首页载入出错",
				"success": "error",
			})
			return
		}
		nav(ctx, "nudao.xyz-数学酷吗", data)
		hot := make([]map[string]interface{}, 0)
		rightHot(fc, ctx, &hot)
		data["rightHot"] = hot
		tem.Execute(ctx.Writer, data)
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", http.StatusMovedPermanently)
	}

}
func changeMSPOST(ctx *gin.Context) {
	defer glog.Flush()
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		data := make(map[string]interface{})
		data["oldPassWord"], _ = ctx.GetPostForm("oldPassword")
		data["newPassWord"], _ = ctx.GetPostForm("newPassword")
		data["newYear"], _ = ctx.GetPostForm("newYear")
		data["newPhoneNumber"], _ = ctx.GetPostForm("newPhoneNumber")
		data["newDescription"], _ = ctx.GetPostForm("newDescription")
		query := fmt.Sprintf("/changeMS?l=h")
		oneValue, err := fc.PostValueToServerBySessionPlus(serverURL+query, "", data)
		ifErrReturn(err, ctx, "无法从后端获取信息")
		if fastjson.GetString(oneValue, "success") == "error" {
			ctx.JSON(http.StatusOK, gin.H{
				"data":    fastjson.GetString(oneValue, "data"),
				"success": "error",
			})
			return
		} else if fastjson.GetString(oneValue, "success") == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"data":    "你没有任何的动作",
				"success": "error",
			})
			return
		} else {
			fc.value = ""
			fc.SetCookie(ctx)
			http.Redirect(ctx.Writer, ctx.Request, "/signIn", http.StatusMovedPermanently)
		}
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", http.StatusMovedPermanently)
	}

}

// 出题
func makeExamGET(ctx *gin.Context) {
	glog.V(2).Infoln("advise")
	defer glog.Flush()
	tem, err := template.ParseFiles(temp("small-editor.html", "small-nav.html", "index-makeExam.html")...)
	if err != nil {
		glog.Error(err)
	}
	tem.Execute(ctx.Writer, nil)
	if err != nil {
		glog.Error(err)
	}

}

// 首先就是要找到 mother contentPLUS
func makeExamPOST(ctx *gin.Context) {

}

// 我给大家出的题
//todo ING
func myExamGET(ctx *gin.Context) {
	defer glog.Flush()
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		data := make(map[string]interface{})
		tem, err := template.ParseFiles(temp("index.html", "text_center_myExam.html", "text.html", "text-left.html", "text-right-index.html")...)
		if err != nil {
			glog.Error(err)
		}
		page := ctx.Query("page")
		pre, pageArray, last, this := pageConversion(page, ctx) // 生成
		hot := make([]map[string]interface{}, 0)
		rightHot(fc, ctx, &hot)
		data["rightHot"] = hot
		data["pagePre"] = pre
		data["pageLast"] = last
		data["page"] = pageArray
		data["pageThis"] = this
		query := fmt.Sprintf("/myTest?page=%s&typeList=6", page)
		value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
		ifErrReturn(err, ctx, "无法从后端获取数据")
		oneValue, err := fastjson.ParseBytes(value)
		ifErrReturn(err, ctx, "无法解析json数据")
		rangeValue := make([]map[string]interface{}, 0)
		for _, v := range oneValue.GetArray() {
			rangeValue = append(rangeValue, gin.H{
				"title":       string(v.GetStringBytes("title")),
				"contentPlus": string(v.GetStringBytes("contentPlus")),
			})
		}
		data["dataRange"] = rangeValue
		//传入结束
		//
		tem.Execute(ctx.Writer, data)
		if err != nil {
			glog.Error(err)
		}
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", http.StatusMovedPermanently)
	}

}

// 排名
func RankingGET(ctx *gin.Context) {
	fc := NewFrontCookie("mathcool", "", "", "") // fc
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		data := make(map[string]interface{}) // 这是一个要传递给前端的data
		err := fc.GetCookie("mathcool", ctx) //先get到cookie这个value后，然后再从远程服务器得到数据
		tem, err := template.ParseFiles(temp("index.html", "text_center_ranking.html", "text.html", "text-left.html", "text-right-index.html")...)
		if err != nil {
			glog.Error(err)
		}
		hot := make([]map[string]interface{}, 0)
		rightHot(fc, ctx, &hot)
		data["rightHot"] = hot
		nav(ctx, "nudao.xyz-数学酷吗", data)

		//传入结束
		//
		tem.Execute(ctx.Writer, data)
		if err != nil {
			glog.Error(err)
		}
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", http.StatusMovedPermanently)
	}
}

// 试题榜单

func testListGET(ctx *gin.Context) {
	tem, err := template.ParseFiles(temp("index.html", "text_center_testlist.html", "text.html", "text-left.html", "text-right-index.html")...)
	if err != nil {
		glog.Error(err)
	}
	data := make(map[string]interface{})
	nav(ctx, "nudao.xyz-数学酷吗", data)
	hot := make([]map[string]interface{}, 0)
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	oneValue, err := fc.GetValueFromServerBySessionPlus(serverURL + "/testList?l=h")
	ifErrReturn(err, ctx, "无法获取后端数据")
	v, err := fastjson.ParseBytes(oneValue)
	ifErrReturn(err, ctx, "解析错误")
	sliceData := make([]map[string]interface{}, 0)
	for _, v := range v.GetArray() {
		sliceData = append(sliceData, gin.H{
			"number":   v.GetInt("number"),
			"userName": string(v.GetStringBytes("userName")),
			"userPlus": string(v.GetStringBytes("userPlus")),
		})
	}
	rightHot(fc, ctx, &hot)
	data["dataSlice"] = sliceData
	data["rightHot"] = hot
	//传入结束
	//
	tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}
}

// 工作
func jobGET(ctx *gin.Context) {
	defer glog.Flush() // 将glog 传出
	page := ctx.Query("page")
	pre, pageArray, last, this := pageConversion(page, ctx)                                                                                // 生成
	fc := NewFrontCookie("mathcool", "", "", "")                                                                                           // fc
	data := make(map[string]interface{})                                                                                                   // 这是一个要传递给前端的data
	err := fc.GetCookie("mathcool", ctx)                                                                                                   //先get到cookie这个value后，然后再从远程服务器得到数据
	glog.V(2).Infoln("/")                                                                                                                  // 作为数据的log
	tem, err := template.ParseFiles(temp("index.html", "text_center_job.html", "text.html", "text-left.html", "text-right-index.html")...) // 将文件 导出为完整的HTML
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "首页载入出错",
			"success": "error",
		})
		return
	}
	nav(ctx, "nudao.xyz-工作", data)
	rangeValue := make([]map[string]interface{}, 0)
	query := fmt.Sprintf("/job?typeList=%s&page=%s", 3, page)
	value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "无法从后台服务器获取数据",
			"success": "error",
		})
		return
	}
	v, err := fastjson.ParseBytes(value)
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"data":    "无法从后台服务器获取数据",
			"success": "error",
		})
		return
	}
	for _, v := range v.GetArray() {
		var m = make(map[string]interface{})
		m["see_number"] = v.GetInt("see_number")
		tagValue := v.GetInt("tag")
		m["tag"] = tag(tagValue)
		m["contentPlus"] = string(v.GetStringBytes("contentPlus"))
		m["join_time"] = string(v.GetStringBytes("join_time"))
		m["title"] = string(v.GetStringBytes("title"))
		if len(v.GetStringBytes("userName")) >= 10 {
			m["userName"] = string(v.GetStringBytes("userName")[:9])
		} else {
			m["userName"] = string(v.GetStringBytes("userName"))
		}
		m["zan"] = v.GetInt("zan")
		rangeValue = append(rangeValue, m)
	}
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["pagePre"] = pre
	data["page"] = pageArray
	data["pageLast"] = last
	data["pageThis"] = this
	data["dataRange"] = rangeValue
	data["rightHot"] = hot
	tem.Execute(ctx.Writer, data)
	if err != nil {
		glog.Error(err)
	}
}

// 提交评论
func commentPOST(ctx *gin.Context) {
	defer glog.Flush()
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		contentPlus := ctx.Query("contentPlus")
		if contentPlus == "" {
			http.Redirect(ctx.Writer, ctx.Request, "/noContent", http.StatusMovedPermanently)
			return
		}
		body, bo := ctx.GetPostForm("comment")
		body = string(markdownValue([]byte(body)))
		data := make(map[string]interface{})
		if !bo {
			http.Redirect(ctx.Writer, ctx.Request, "/noContent", http.StatusMovedPermanently)
			return
		}
		data["data"] = body
		query := fmt.Sprintf("/addComment?contentPlus=%s", contentPlus)
		value, err := fc.PostValueToServerBySessionPlus(serverURL+query, "text/html", data)
		ifErrReturn(err, ctx, "从后端传来的命令异常comment")
		if fastjson.GetString(value, "success") == "error" {
			ctx.JSON(http.StatusOK, "无法提交评论"+fastjson.GetString(value, "data"))
			return
		}
		http.Redirect(ctx.Writer, ctx.Request, "/w?contentPlus="+contentPlus, http.StatusMovedPermanently)
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", http.StatusMovedPermanently)
	}
}

// 提交容器
func contentPOST(ctx *gin.Context) {
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		var typeList string
		var motherContentPlus string
		typeListQuery := ctx.Query("typeList")
		typeListForm, _ := ctx.GetPostForm("typeList")
		if typeListQuery == "" {
			typeList = typeListForm
		} else if typeListQuery != "" && typeListQuery != typeListForm {
			typeList = typeListQuery
		}
		if typeList == "6" {
			motherContentPlus = ctx.Query("motherContentPlus")
		}
		tag, _ := ctx.GetPostForm("tag")
		title, _ := ctx.GetPostForm("editorTitle")
		body, _ := ctx.GetPostForm("editorBody")
		datat := make(map[string]interface{})
		datat["title"] = title
		datat["contentValue"] = string(markdownValue([]byte(body)))
		//tag,_ := ctx.GetPostForm("tag-list")
		// 7 就是 提出的意见。
		query := fmt.Sprintf("/w?typeList=%s&tag=%s&motherContentPlus=%s", typeList, tag, motherContentPlus)
		// 传入一般的data就行了，这里就不用json了 因为 下面的那个函数已经用过了。
		v, _ := fc.PostValueToServerBySessionPlus(serverURL+query, "", datat)
		if fastjson.GetString(v, "success") == "ok" {
			data := "<script>window.location='/'</script>"
			fmt.Fprint(ctx.Writer, data)
		} else {
			ctx.JSON(http.StatusOK, "提交失败"+fastjson.GetString(v, "data"))
		}
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", http.StatusMovedPermanently)
	}
}

// 增加赞
func addZanGET(ctx *gin.Context) {
	defer glog.Flush()
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		fc.GetValueFromServerBySessionPlus(serverURL + "/addZan?contentPlus=" + ctx.Query("contentPlus"))
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", http.StatusMovedPermanently)
	}
}

func noSignGET(ctx *gin.Context) {
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	noSignIn(ctx, fc)
}
func noContentGET(ctx *gin.Context) {
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	noContent(ctx, fc)
}

// 删除评论
func deleteCommentGET(ctx *gin.Context) {
	defer glog.Flush()
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		contentPlus := ctx.Query("contentPlus")
		commentID := ctx.Query("commentID")
		if commentID == "" {
			http.Redirect(ctx.Writer, ctx.Request, "/noContent", http.StatusMovedPermanently)
			return
		}
		query := fmt.Sprintf("/deleteComment?contentPlus=%s&commentID=%s", contentPlus, commentID)
		value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
		ifErrReturn(err, ctx, "在删除评论的时候从服务器返回信息错误")
		if fastjson.GetString(value, "success") == "error" {
			http.Redirect(ctx.Writer, ctx.Request, "/noContent", http.StatusMovedPermanently)
			return
		}
		http.Redirect(ctx.Writer, ctx.Request, "/user?typeList=2", http.StatusMovedPermanently)
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "noSign", http.StatusMovedPermanently)
	}
}

// 删除容器
func deleteContentGET(ctx *gin.Context) {
	defer glog.Flush()
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	b, _ := ctx.Get("makeSureUser")
	if b.(bool) {
		contentPlus := ctx.Query("contentPlus")
		if contentPlus == "" {
			http.Redirect(ctx.Writer, ctx.Request, "/noContent", http.StatusMovedPermanently)
			return
		}
		query := fmt.Sprintf("/deleteW?contentPlus=%s", contentPlus)
		value, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
		ifErrReturn(err, ctx, "在删除容器的时候从服务器返回的信息错误")
		if fastjson.GetString(value, "success") == "error" {
			http.Redirect(ctx.Writer, ctx.Request, "/noContent", http.StatusMovedPermanently)
			return
		}
		http.Redirect(ctx.Writer, ctx.Request, "/user?typeList=1", http.StatusMovedPermanently)

	} else {
		http.Redirect(ctx.Writer, ctx.Request, "noSign", http.StatusMovedPermanently)
	}
}

// 增加image
// 算法是 1 用户点击数据上传，
//接受用户的数据，然后上传到github
//将用户的本地image删除，然后返回一个github的url
//给用户将这个url储存在数据库中然后将这个图片显示在user中，
//然后让用户点击复制然后就可以将这个图片加载到这个文章或者评论中。
func addImagePOST(ctx *gin.Context) {
	b, _ := ctx.Get("makeSureUser")
	r := rand.New(source)
	if b.(bool) {
		var errIM error
		fc := NewFrontCookie("mathcool", "", "", "")
		fc.GetCookie("mathcool", ctx)
		form, err := ctx.MultipartForm()
		ifErrReturn(err, ctx, "无法获取上传的文件")
		wait := make(chan struct{}, 10)
		group := sync.WaitGroup{}
		data := make(map[interface{}]interface{})
		lock := sync.Mutex{}
		group.Add(len(form.File["file"]))
		var i int
		for _, v := range form.File["file"] {
			i++
			go func(v *multipart.FileHeader, i int) {
				t := make([]byte, 10)
				r.Read(t)
				lock.Lock()
				defer lock.Unlock()
				defer group.Done()
				wait <- struct{}{}
				file, err := v.Open()
				ifErrReturn(err, ctx, "无法剖析file")
				typeValue := v.Filename[len(v.Filename)-3:]
				typeValue = strings.ToLower(typeValue)
				if typeValue == "png" || typeValue == "jpeg" || typeValue == "jpg" || typeValue == "peg" || typeValue == "gif" {
					defer file.Close()
					creatFile, err := os.Create("./public/img/" + fmt.Sprintf("%x", t) + v.Filename)
					defer creatFile.Close()
					ifErrReturn(err, ctx, "无法创建文件")
					fmt.Println("测试", typeValue)
					errIM = imageDealWith.Compression(typeValue, 600, file, creatFile)
					if errIM == nil {
						r := regexp.MustCompile("./public/img/")
						realName := r.ReplaceAllString(creatFile.Name(), "")
						data[i] = realName
					} else {
						err = os.Remove(creatFile.Name())
						ifErrReturn(err, ctx, "无法删除")
					}
				}
				//_, err = io.Copy(creatFile, file)
				<-wait

			}(v, i)
		}
		group.Wait()
		//
		//只有结束了才能上传数据库//
		// 上传传数据库
		if errIM == nil {
			fmt.Println("测试文件", data)
			wait2 := sync.WaitGroup{}
			speed := make(chan struct{}, 10)
			wait2.Add(len(data))
			for _, v := range data {
				go func(v interface{}) {
					speed <- struct{}{}
					defer wait2.Done()
					query := fmt.Sprintf("/addImage?l=h&imgValue=%s", v)
					fc.GetValueFromServerBySessionPlus(serverURL + query)
					<-speed
				}(v)
			}
			wait2.Wait()
			// 上传执行完毕。
		}
		http.Redirect(ctx.Writer, ctx.Request, "/user", 301)

	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", 301)
	}
}

// 删除这个image
func deleteImageGET(ctx *gin.Context) {
	b, _ := ctx.Get("makeSureUser")
	r := rand.New(source)
	t := make([]byte, 10)
	fc := NewFrontCookie("mathcool", "", "", "")
	fc.GetCookie("mathcool", ctx)
	r.Read(t)
	if b.(bool) {
		iv := ctx.Query("imgValue")
		err := os.Remove("./public/img/" + iv)
		ifErrReturn(err, ctx, "无法删除")
		id := ctx.Query("imgID")
		query := fmt.Sprintf("/deleteImage?imgID=%s", id)
		va, err := fc.GetValueFromServerBySessionPlus(serverURL + query)
		ifErrReturn(err, ctx, "接受数据出错")
		if fastjson.GetString(va, "success") == "error" {
			ctx.JSON(http.StatusOK, gin.H{
				"data":    fastjson.GetString(va, "data"),
				"success": "error",
			})
			return
		}
		http.Redirect(ctx.Writer, ctx.Request, "/user", 301)
	} else {
		http.Redirect(ctx.Writer, ctx.Request, "/noSign", 301)
	}

}

func weiboSignInGET(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(200, gin.H{
			"data":    "无法通过微博登陆",
			"success": "error",
		})
		return
	}
	reUrl := "https://127.0.0.1/weiboSignIn"
	query := fmt.Sprintf("https://api.weibo.com/oauth2/access_token?client_id=3266637437&client_secret=33dfb7df4d81fa187c9377af1a75adb9&grant_type=authorization_code&redirect_uri=%s&code=%s", reUrl, code)
	rc := NewFrontCookie("mathcool", "", "", "")
	value, err := rc.PostValueToServerBySessionPlus(query, "", nil)
	ifErrReturn(err, ctx, "获取信息失败,您无法登陆")
	accessToken := fastjson.GetString(value, "access_token") // sessionPlus
	query = fmt.Sprintf("https://api.weibo.com/oauth2/get_token_info?access_token=%s", accessToken)
	value1, err := rc.PostValueToServerBySessionPlus(query, "", nil)
	ifErrReturn(err, ctx, "无法获取信息")
	uid1 := fastjson.GetInt(value1, "uid")
	fmt.Println("测试accs",accessToken,uid1)
	query = fmt.Sprintf("https://api.weibo.com/2/users/show.json?access_token=%s&uid=%d", accessToken, uid1)
	value, err = rc.GetValueFromServerBySessionPlus(query)
	ifErrReturn(err, ctx, "无法获取信息")
	userName := fastjson.GetString(value, "screen_name")
	location := fastjson.GetString(value, "location")
	gender := fastjson.GetString(value, "gender")
	description := fastjson.GetString(value, "description")
	fmt.Println(string(value))
	var sex int
	if gender == "男" {
		sex = 1
	} else {
		sex = 2
	}
	query = fmt.Sprintf("/weiboSignIn?uid=%s%d", "weibo", uid1)
	_, err = rc.PostValueToServerBySessionPlus(serverURL+query, "", gin.H{
		"userName":    userName,
		"location":    location,
		"sex":         sex,
		"description": description,
	})
	if err == nil {
		rc.value = fmt.Sprintf("%s%d", "weibo", uid1)
		rc.SetCookie(ctx)
		http.Redirect(ctx.Writer, ctx.Request, "/", 200)
	} else {
		ctx.JSON(200, gin.H{
			"success": "error",
			"data":    "无法通过微博登陆",
		})
	}
	//id	int64	用户UID
	//idstr	string	字符串型的用户UID
	//screen_name	string	用户昵称
	//name	string	友好显示名称
	//province	int	用户所在省级ID
	//city	int	用户所在城市ID
	//location	string	用户所在地
	//description	string	用户个人描述
	//url	string	用户博客地址
	//profile_image_url	string	用户头像地址（中图），50×50像素
	//profile_url	string	用户的微博统一URL地址
	//domain	string	用户的个性化域名
	//weihao	string	用户的微号
	//gender
}
