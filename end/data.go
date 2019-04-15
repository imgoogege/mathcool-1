// 数据
package main

import (
	"database/sql"
	"github.com/golang/glog"
	"math/rand"
	"sync"
	"time"
)

// 储存session数据的地方 通过sessionID 来获取session对象的所有数据(也就是说 通过sessionID来检索一个人的基本不牵涉到机密的一些信息)
// 通过SessionID 来对应一个完整的Session 对象。作为缓存，后期使用redis代替。
var (
	// 全局的sessionID map
	SessionMap = make(map[string]*Session)
	// 通过 ID 来加密生成一个对应的Plus 来避免id被暴露，这里的是一个生成一个资源，全局一个资源即可。
	source         = rand.NewSource(time.Now().UnixNano())
	dbHere         *sql.DB
)

type Session struct {
	lock        sync.Mutex //  加入一个锁
	SessionID   int64      `json:"session_id" redis:"session_id"`     /// 自增的Session数据
	SessionPlus string     `json:"session_plus" redis:"session_plus"` // 外漏的session
	User
}
type comment struct {
	commentID int64 `json:"comment_id"`
	ContentID int64  `json:"content_id"` // 绑定文章的id
	commentValue     string `json:"comment_value"`      //实际的内容
	UserID    int64  `json:"user_id"`    // 绑定 用户的id。
}
type User struct {
	UserID      int64  `json:"user_id" redis:"user_id"`
	UserPlus    string `json:"user_plus" redis:"user_plus"`
	UserName    string `json:"user_name" redis:"user_name"`
	Sex         int    `json:"sex" redis:"sex"`
	Year        string `json:"year" redis:"year"`
	JoinTime    string `json:"join_time" redis:"join_time"`
	Email       string `json:"email" redis:"email"`
	PhoneNumber string `json:"phone_number" redis:"phone_number"`
	Description string `json:"description" redis:"description"`
	Salt        string `json:"salt" redis:"salt"`
	DBPassword  string `json:"db_password"`
}

// 所有的只要是内容，就是这个对象 ContentID 的生成 用 title + author + time.Now().month + time.now().Day() + 一个2位的随机数
type Content struct {
	ContentID      int64  `json:"content_id" redis:"content_id"`
	ContentPlus    string `json:"content_plus" redis:"content_plus"`
	Title          string `json:"title" redis:"title"`
	UserID         int64  `json:"user_id" redis:"user_id"`
	ContentValue   string `json:"content_value" redis:"content_value"`
	typeList       int    `json:"type_list"`        // 判定是什么类型的内容 1 文章 2 问答 3 工作 4 题 5 公式 6 公式配套试题 小试题
	motherContentID int64 `json:"mother_content_id"` //这个字段是为了给例如 公式 然后公式的配套试题的那种类型准备的// 1 是 2 不是
	zan            int64  `json:"zan"` //  赞的个数
	seeNumber int64 `json:"see_number"` // 浏览量
	joinTime string `json:"join_time"`
}

func open() {
	var err error
	dbHere, err = sql.Open("mysql", "root:359258Ls!@tcp(localhost:3306)/mathcoolEnd")
	dbHere.SetMaxOpenConns(2000)
	dbHere.SetMaxIdleConns(1000)
	if err != nil {
		glog.Error(err)
	}
	err = dbHere.Ping()
	if err != nil {
		glog.Error("数据库无法连接", err)
	}
}
