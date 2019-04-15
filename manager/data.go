package main

import (
	"database/sql"
	"math/rand"
	"time"
)

var (
	dbHere  *sql.DB
	userMap = make(map[string]ManagerSession)
	r       = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type ManagerSession struct {
	sessionID   int64
	sessionPlus string
	ManagerUser
}
type ManagerUser struct {
	isRoot   int // 1 root 2 commonUser
	userID   int64
	userName string
	level    int // 职务的级别 1 2 3 4 5 6 7 root:1
	db       string
	salt     string
}

type ManagerCookie struct {
	name   string
	maxAge int64
	value  string
}
