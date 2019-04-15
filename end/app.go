// 全局使用的某些函数。
package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
)

// 验证密码格式
func verificationPassword(value *string) {

}

// 验证邮箱
func verificationEmail(value *string) {

}

// 过滤用户的js代码和sql代码
func deleteUserJSAndSql(value *string){

}

// 过滤用户的非法信息，例如关于政治，色情的东西
func deleteColorAndParty(value *string){

}

func PreAction(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
}

//cors中间件
func cores()func( *gin.Context){
	return cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost","https://nudao.xyz"},
		AllowMethods:     []string{"PUT", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	})
}

// 前面是k后面是v
func sortMap(t map[int]int)([]int,[]int){
	aSlice := make([]int,0)
	aMap := make(map[int]int)
	bSlice := make([]int,0)
	cSlice := make([]int,0)
	for k,v := range t {
		aMap[v] = k
		aSlice = append(aSlice,v)
	}
	sort.Ints(aSlice)
	for _,v := range aSlice  {
		bSlice =append(bSlice,aMap[v])
	}
	for _,v := range bSlice  {
		cSlice = append(cSlice,t[v])
	}
	//这个时候bSlice是递增的，那么我们需要一个递减的[] 这个时候就需要排序算法了。
	sort.Sort(sort.Reverse(sort.IntSlice(bSlice)))
	return bSlice,cSlice
}
//return resultKye resultVALUE
// 传入一个t的map然后按照后面的值排好顺序然后返回这个map
func sortMapReturnMap(t map[int]int)[]map[int]interface{}  {
	data := make([]map[int]interface{},0)
	d,dd := sortMap(t)
	for k,v := range d{
		data = append(data,map[int]interface{}{
			v:dd[k],
		})
	}
	return data
}

func ifErrReturn(err error ,ctx *gin.Context,value interface{}){
	if err != nil {
		ctx.JSON(http.StatusOK,gin.H{
			"success":"error",
			"data":fmt.Sprint(value,err),
		})
		return
	}
}

