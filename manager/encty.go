package main

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

// return result salt
func encty(a int, b string) (result string, salt string,pa string) {
	var as string
	p := make([]byte,4)
	salt = fmt.Sprint(r.ExpFloat64())
	if a == 0 {
		as = "0"
	} else {
		as = strconv.FormatInt(int64(a), 10)
	}
	r.Read([]byte(p))
	p1 := fmt.Sprintf("%x",p)
	return fmt.Sprintf("%x", md5.Sum([]byte(as+b+p1+salt))), salt,p1
}
