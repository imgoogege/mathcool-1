package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *User) GetUserPlus() {
	plus, salt := Encryption(u.UserID, u.JoinTime)
	u.UserPlus = plus
	u.Salt = salt
}

// makeSureIsUser 这个变量 除了在 signUp signIn这俩函数（也就是注册登陆）不判断外，其余都判断
func isUser(ctx *gin.Context) {
	ctx.Set("makeUserIsUser",false)
	u := new(User)
	s := new(Session)
	//1 j检查全局的map中是否有这个得到的sessionID 也就是数据库中的session_plus
	// 2 如果有 就将 makeSureIsUser 设置为true
	// 3 否则就开始启动数据库进行访问，使用sessionID 看看是否可以取出来值，如果可以就将这个sessionID 加入到这个全局的map中。

	//1
	var thing string
	sessionPlus := ctx.Query("sessionPlus") // 得到这个前端传入的sessionID 也就是后端的sessionPlus
	uid := ctx.Query("uid")
	if sessionPlus == "" && uid != "" {
		thing = uid
	}else if sessionPlus != "" && uid == "" {
		thing = sessionPlus
	}else {
		return
	}
	if _, ok := SessionMap[thing]; !ok {
		// 2
		if rows, err := dbHere.Query("SELECT session_id,user_id FROM session WHERE session_plus=? ", sessionPlus);  err!= nil  {
			defer rows.Close()
			ctx.Set("makeUserIsUser",false)
			return
		} else {
			for rows.Next() {
				var id ,id2 int64
				rows.Scan(&id, &id2)
				if id <=0 || id2 <= 0 {
					ctx.Set("makeUserIsUser",false)
					return
				}
				fmt.Println("测试id1，id2",id,id2)
				s.SessionID = id
				u.UserID = id2
			}
			if u.UserID <=0 {
				ctx.Set("makeUserIsUser",false)
				return
			}
			defer rows.Close()
			if rows, err = dbHere.Query("SELECT user_plus,user_name,sex,year,join_time,email,phone_number,description,salt,db_password FROM user WHERE user_id=?", u.UserID);  err!= nil{
				defer rows.Close()
				ctx.Set("makeUserIsUser",false)
				ctx.JSON(http.StatusOK,"无法登陆"+fmt.Sprintf("%v",err))
				fmt.Println(err)
				return
			}
			for rows.Next() {
				var plus,name,year,join,email,phone,descrip,salt,dbPassword string
				var sex int
				rows.Scan(&plus, &name, &sex, &year, &join, &email, &phone, &descrip,&salt,&dbPassword)
				if plus=="" {
					ctx.Set("makeUserIsUser",false)
					return
				}
				u.Salt = salt
				u.DBPassword = dbPassword
				u.UserPlus= plus
				u.UserName = name
				u.Sex = sex
				u.Year = year
				u.JoinTime = join
				u.Email = email
				u.PhoneNumber = phone
				u.Description = descrip
			}
			s.User = *u
			s.SessionPlus = sessionPlus
			if s.JoinTime == "" {
				ctx.Set("makeUserIsUser",false)
				return
			}
			SessionMap[sessionPlus] = s
			ctx.Set("makeUserIsUser",true)
			fmt.Println(u)
			//s.SetSessionToRedis()
		}
	} else {
		ctx.Set("makeUserIsUser",true)
	}
}
