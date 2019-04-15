package main

import "github.com/gin-gonic/gin"

// 一个 前端的判断 是否登陆的装置 ，不然每次等在fuc中判断是复杂的行为,只需要判断 有没有cookie的值就行，不需要判断这个sessionPlus是否是真是的
//因为你每次的请求都是带有了判断。每次请求都判断了。这里就不用赘写了。
func USER(ctx *gin.Context) {
	ctx.Set("makeSureUser", false)
	s, err := ctx.Cookie("mathcool")
	if err != nil || s == "" { // 如果没能从cookie中获取信息，那么就是确定这个并没有登陆，
		ctx.Set("makeSureUser", false)
		return
	}
	ctx.Set("makeSureUser", true)
}
