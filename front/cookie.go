//本文件有以下几个功能：
//
// 1 设置cookie
// 2 得到cookie的内容也就是sessionID
// 3 从远程内部服务器得到想要的cookie的sessionID值.
//4 是否登录

//TODO:预防xss攻击 （过滤掉输入中的所有特殊字符）和csrf攻击给每个用户生成一个特殊的formID，这样就避免了csrf攻击 以及对用户输入的信息，进行js过滤和服务器过滤。 sql注入的问题。
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

func NewFrontCookie(name, value string, path string, domain string) *FrontCookie {
	return &FrontCookie{
		name:     name,
		value:    value,
		maxAge:   1 * 60 * 60 * 24 * 7,
		path:     path,
		domain:   domain,
		secure:   true,
		httpOnly: true,
	}
}

// 在前端的cookie主题的对象
type FrontCookie struct {
	lock     sync.Mutex
	name     string    // cookie的名称
	value    string    // cookie的内容
	expires  time.Time // 例子：Date: Wed, 21 Oct 2015 07:28:00 GMT 就是失效的时间，不过，如果你上面都不设置，那么就是退出这个网站的时候就会消失了。
	maxAge   int       // cookie 的最长有效时间 跟expires类似，在 cookie 失效之前需要经过的秒数，它的优先级比expires高
	path     string    //指定一个 URL 路径，这个路径必须出现在要请求的资源的路径中才可以发送 Cookie 首部。
	domain   string    //指定 cookie 可以送达的主机名。假如没有指定，那么默认值为当前文档访问地址中的主机部分（但是不包含子域名）
	secure   bool      //一个带有安全属性的 cookie 只有在请求使用SSL和HTTPS协议的时候才会被发送到服务器。
	httpOnly bool      //设置了 HttpOnly 属性的 cookie 不能使用 JavaScript 经由  Document.cookie 属性、XMLHttpRequest 和  Request APIs 进行访问
}

// 设置cookie
func (cookie *FrontCookie) SetCookie(ctx *gin.Context) {
	cookie.lock.Lock()
	defer cookie.lock.Unlock()
	ctx.SetCookie(cookie.name, cookie.value, cookie.maxAge, cookie.path, cookie.domain, cookie.secure, cookie.httpOnly)
}

//从浏览器得到cookie的内容
//返回一个只拥有cookie.Value的FrontCookie指针
func (cookie *FrontCookie) GetCookie(name string, ctx *gin.Context) error {
	cookie.lock.Lock()
	defer cookie.lock.Unlock()
	value, err := ctx.Cookie(name)
	cookie.value = value
	cookie.name = name
	return err
}

// 将这个cookie中的value(如 SessionID)删除
func (cookie *FrontCookie) DeleteValue(ctx *gin.Context) {
	cookie.value = ""
	cookie.SetCookie(ctx)
}

// 改变这个cookie的有效时间,如果想让这个cookie改变为
func (cookie *FrontCookie) ChangeExpires(newExpires time.Time) {
	cookie.lock.Lock()
	defer cookie.lock.Unlock()
	cookie.expires = newExpires
}

// 续订这个有效时间，如果发现已经经过的时间，已经超过了一半，那么给他补充它缺失的时间。
func (cookie *FrontCookie) RenewExpipres() {
	cookie.lock.Lock()
	defer cookie.lock.Unlock()
	now := time.Now()
	setTime := time.Now().Add(time.Duration(cookie.maxAge) * time.Second)
	subTime := setTime.Sub(now)
	if subTime.Hours() <= float64(24*time.Nanosecond) {
		cookie.maxAge = int(time.Nanosecond * 60 * 60 * 24)
	}
}

// 通过sessionID 从后端取到该有的数据.
func (cookie *FrontCookie) GetValueFromServerBySessionPlus(serverUrl string) ([]byte, error) {
	res, err := http.Get(serverUrl + "&sessionPlus=" + cookie.value)
	defer res.Body.Close()
	if err != nil {
		return []byte(""), fmt.Errorf("在从后端得到数据的时候发生了错误，错误信息是: %v", err)
	}
	return ioutil.ReadAll(res.Body)

}

// 把 这个body中的数据，传入到post方法中。
func (cookie *FrontCookie) PostValueToServerBySessionPlus(serverUrl string, contentType string, data interface{}) ([]byte, error) {
	if contentType == "" {
		contentType = "application/json"
	}
	v, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(v))
	res, err := http.Post(serverUrl+"&sessionPlus="+cookie.value, contentType, reader)
	if err != nil {
		return []byte(""), err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

// 中间件 就是每次请求都会来判断 cookie的时间。
func AddCookieTime(ctx *gin.Context) {
	c, err := ctx.Request.Cookie("mathcool")
	if err != nil {
		return
	}
	fc := NewFrontCookie(c.Name, c.Value, c.Path, c.Domain)
	fc.RenewExpipres()
	fc.SetCookie(ctx)
}
