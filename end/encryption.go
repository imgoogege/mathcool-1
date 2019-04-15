package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
)

// 加密的方式是这样的，使用 参数(int64) + 参数(string) + 一个随机的数字（盐） 然后 返回值位是这样的 第一个是 加密的结果，第二个是那个盐 都是 string。
func Encryption(id int64, st string) (result string, salt string) {
	rand := rand.New(source)
	exp := rand.ExpFloat64()
	idValue := strconv.FormatInt(id, 10)
	expValue := strconv.FormatFloat(exp, 'f', 8, 64)
	mdValue := md5.Sum([]byte(st + idValue + expValue))
	return fmt.Sprintf("%x", mdValue), expValue
}
