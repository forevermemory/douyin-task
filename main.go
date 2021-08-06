package main

import (
	"douyin/config"
	"douyin/global"
	"douyin/web/router"
	"douyin/web/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// init config
	config.InitConfig()

	//
	global.InitConnect()

	//
	go service.RunRedisSyncToMysqlManager()
	go service.RunCronTasks()

	// router
	r := router.InitRouter(nil, "")

	fmt.Println("Listening and serving HTTP on 0.0.0.0:", config.CONFIG.HttpConfig.Port)
	r.Run("0.0.0.0:" + fmt.Sprintf("%d", config.CONFIG.HttpConfig.Port))

}
