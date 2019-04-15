// 后期将session储存在redis中，前期先储存在map中即可。处理session的地方
package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gomodule/redigo/redis"
)

// 得到这个 Session
func GetSessionFromSessionMap(sessionID string) *Session {
	return SessionMap[sessionID]
}

//因为sessionID被删除，将这个session数据，从这个mysql数据表中删除掉。
//将sessionID 加入到这个 map中
func (s *Session) AddMap() {
	s.setSession()
	SessionMap[s.SessionPlus] = s
}

func (s *Session) Delete() {
	delete(SessionMap, s.SessionPlus)
	//todo: 将这个session从数据库中删除。
}

// NewSession() 返回一个空的Session的指针。
func NewSession() *Session {
	return &Session{}
}

// 设置 session
func (s *Session) setSessionPlus() {
	s.SessionPlus, _ = Encryption(s.SessionID, s.UserPlus)
}

func (s *Session) setSession() {
	// todo 从数据库中取到数据，赋值到session中。
	// 将数据给导入到这个Session中。
	s.setSessionPlus()
}

// 将数据写入到缓存redis中
func (s *Session) SetSessionToRedis() {
	s.lock.Lock()
	defer s.lock.Unlock()
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		glog.Error(err)
		fmt.Println("debug", err)
	}
	defer c.Close()
	if _, err = c.Do("HMSET", redis.Args{}.Add(s.SessionPlus).AddFlat(s)...); err != nil {
		glog.Error(err)
		fmt.Println(err)
	}

}

// 将数据从缓存redis中取出来,这个s必须现有 SessionPlus对象才能取，然后取是将这个s给充满即可。
func (s *Session) GetSessionFromRedis() {
	s.lock.Lock()
	defer s.lock.Unlock()
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		glog.Error(err)
		return
	}
	defer c.Close()
	v, err := redis.Values(c.Do("HGETALL", s.SessionPlus))
	if err != nil {
		glog.Error(err)
		fmt.Println(err)
		return
	}

	if err := redis.ScanStruct(v, s); err != nil {
		glog.Error(err)
		fmt.Println(err)
		return
	}
}

// 将缓存中的数据删除
func (s *Session) DeleteSessionFromRedis() {
	s.lock.Lock()
	defer s.lock.Unlock()
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		glog.Error(err)
	}
	_,err = c.Do("HDEL",s.SessionPlus)
	if err != nil {
		fmt.Println(err)
		glog.Error(err)
	}

}
