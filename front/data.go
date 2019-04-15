// data.go 全局的data数据放这。不全局的不放。
package main

import (
	"math/rand"
	"time"
)

var (
	csrfMap = make(map[string]string) // 这个是防止csrf的这个东西的一个map
	source  = rand.NewSource(time.Now().UnixNano())
)

type ( // 在data里的写法，这样想比较好。
	SmallDonate struct { // 让人捐赠的一个struct
		Head_title string
		Value      []*Donate
	}
	Donate struct { // 捐赠
		Href            string
		SmallDonateData string
	}
)
