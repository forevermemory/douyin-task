package main

import (
	"douyin/config"
	"douyin/global"
	"douyin/web/router"
	"douyin/web/service"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// init config
	config.InitConfig()

	//
	global.InitConnect()

	conn := global.REDIS.Get()
	_, err := conn.Do("hset", "user", "Birthday", time.Now().Format("2006-01-02 15:04:05"))
	_, err = conn.Do("hset", "user", "Age", 111)
	fmt.Println(err)
	fmt.Println(err)
	fmt.Println(err)

	type user struct {
		Birthday string `redis:"Birthday"`
		Age      int    `redis:"Age"`
	}

	// res2 := []interface{}{"user", "Birthday", "age"}
	// res, err := redis.Values(conn.Do("hmget", res2...))
	res, err := redis.Values(conn.Do("hmget", "user", "Birthday", "Age"))
	u := new(user)
	err = redis.ScanStruct(res, u) // 只能用hgetall
	fmt.Println(err)
	fmt.Println(res)
	fmt.Println(u)
	fmt.Println(u.Age)
	return

	//
	go service.RunRedisSyncToMysqlManager()
	go service.RunCronTasks()

	// router
	r := router.InitRouter(nil, "")

	fmt.Println("Listening and serving HTTP on 0.0.0.0:", config.CONFIG.HttpConfig.Port)
	r.Run("0.0.0.0:" + fmt.Sprintf("%d", config.CONFIG.HttpConfig.Port))

}
