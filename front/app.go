package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func path1(path, file string) string {
	return path + "/" + file
}
func temp(a ...string) (result []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path := filepath.Join(dir, "view")
	result = []string{
		path1(path, "layout.html"),
		path1(path, "allCss.html"),
		path1(path, "allJs.html"),
		path1(path, "nav.html"),
		path1(path, "foot.html"),
		path1(path, "head.html"),
	}
	for _, v := range a {
		result = append(result, path1(path, v))
	}
	return
}

// 返回的是public的路径
func pwdPbulic() string {
	pwd, err := filepath.Abs(".")
	if err != nil {
		return ""
	} else {
		return filepath.Join(pwd, "public")
	}
}

// 返回的是file的路径
func pwdFile() string {
	pwd, err := filepath.Abs(".")
	if err != nil {
		return ""
	} else {
		return filepath.Join(pwd, "file")
	}
}

// 验证csrf防护机制,机制是 将这个东西 保存在前端的array。
func iscsfr(now time.Time, ctx gin.Context, formHiddenValue string) bool {
	hiddenValue, _ := ctx.GetPostForm(formHiddenValue) // 得到这个value值
	plus, err := ctx.Cookie("sessionID")
	if err != nil {
		glog.Error(err)
	}

	if hiddenValue == csrfMap[plus] {
		delete(csrfMap, plus) // 删除这个mapvalue，等待下一次的输入。
		return true
	}
	return false
}

// 如果set呢？首先 使用 time.now（） 将值导入，然后 立刻 去得到 sessionID 将value 给 csrf[sessionID] = time.Now().String()即可。

// 数字标准化 小于1000正常显示 大于1000 显示为 6.5k 去1位小数 大于等于10万 统一显示 10w+
func numberTransform(number int64) (transNumber string) {
	if number < 1000 {
		return fmt.Sprint(number)
	} else if number >= 1000 && number <= 100000 {
		return fmt.Sprintf("%.1f", number/1000)
	} else {
		return "10万+"
	}

}

// 页码转换。不需要什么复杂的算法，反正就10个数字，
func pageS(page int) [10]int {
	//输入一个页面，这个页面必须是第2位
	var data [10]int
	if page <= 0 {
		data[0] = page
	} else {
		data[0] = page - 1
	}
	data[1] = page
	data[2] = page + 1
	data[3] = page + 2
	data[4] = page + 3
	data[5] = page + 4
	data[6] = page + 5
	data[7] = page + 6
	data[8] = page + 7
	data[9] = page + 8
	return data
}

func pageConversion(pageString string, ctx *gin.Context) (pre int, pageArray []int, last int, this int) {
	var page int
	var pageResult [10]int
	if pageString == "" {
		page = 0
	} else {
		var err error
		page, err = strconv.Atoi(pageString)
		if err != nil {
			glog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{"success": "error", "data": "输入的page错误"})
			return
		}
	}
	pageResult = pageS(page)
	return pageResult[0], pageResult[2:8], pageResult[2], pageResult[1]
}

func nav(ctx *gin.Context, headTitle string, data map[string]interface{}) {
	// 判断是否登陆，然后以及nav的写法
	data["Head_title"] = headTitle  // 首页的title设置
	e, _ := ctx.Get("makeSureUser") // 是否登陆的标志
	if !e.(bool) {
		data["signUpStatus"] = "注册"
		data["signOutStatus"] = ""
		data["signInStatus"] = "登录"
	} else {
		data["signUpStatus"] = ""
		data["signInStatus"] = ""
		data["signOutStatus"] = "登出"
	}
	// nav 固定用法结束。
}
func tag(tagValue int) (value string) {
	switch tagValue {
	case 1:
		return "小学数学"
	case 2:
		return "大学本科数学"
	case 3:
		return "小学数学"
	case 4:
		return "初中数学"
	case 5:
		return "高中数学"
	case 6:
		return "研究生数学"
	default:
		return "高级数学"

	}
}

//  浏览量增加
func addView(contentPlus string, fc *FrontCookie) {
	fc.GetValueFromServerBySessionPlus(serverURL + "/addSeeNumber?contentPlus=" + contentPlus)
}

// 👍增加 TODO: 和 更改信息 和 更改密码 *我的排名*（暂时不写这个功能） 安全防护方面
//func addZan(contentPlus string,fc *FrontCookie) {
//	fc.GetValueFromServerBySessionPlus(serverURL + "?addZan?contentPlus="+contentPlus)
//}
func ifErrReturn(err error, ctx *gin.Context, value interface{}) {
	if err != nil {
		glog.Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"success": "error",
			"data":    value,
		})
		return
	}
}

func noSignIn(ctx *gin.Context, fc *FrontCookie) {
	data := make(map[string]interface{})
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	nav(ctx, "nudao.xyz-数学酷吗", data)
	tem, err := template.ParseFiles(temp("index.html", "text_center_notSignIn.html", "text.html", "text-left.html", "text-right-index.html")...)
	ifErrReturn(err, ctx, "无法显示nosign")
	err = tem.Execute(ctx.Writer, data)
	ifErrReturn(err, ctx, "无法渲染 nosign")

}

func noContent(ctx *gin.Context, fc *FrontCookie) {
	data := make(map[string]interface{})
	hot := make([]map[string]interface{}, 0)
	rightHot(fc, ctx, &hot)
	data["rightHot"] = hot
	nav(ctx, "nudao.xyz-数学酷吗", data)
	tem, err := template.ParseFiles(temp("index.html", "text_center_notFindContent.html", "text.html", "text-left.html", "text-right-index.html")...)
	ifErrReturn(err, ctx, "无法显示notfind content")
	err = tem.Execute(ctx.Writer, data)
	ifErrReturn(err, ctx, "无法渲染 notfind content")
}

func contentTypeList(value int) string {
	switch value {
	case 1:
		return "文章"
	case 2:
		return "问答"
	case 3:
		return "工作"
	case 6:
		return "公式配套试题"
	case 7:
		return "意见"
	default:
		return "无法识别"
	}
}

func contentSex(value int) string {
	switch value {
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return "无法识别"
	}
}
