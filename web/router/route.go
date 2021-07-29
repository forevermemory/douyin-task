package router

import (
	"douyin/web/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, prefix string) *gin.Engine {
	if r == nil {
		r = gin.Default()
	}

	r.Use(Cors())

	// r.Static("/moban", "./static")
	// 获取登录验证码
	// r.POST(prefix+"/login", route(controller.Login))

	// ///////////////////////////////////////////////////////////
	r.GET("/", func(c *gin.Context) {

	})
	// user := r.Group(prefix + "/user")
	// {
	// 	user.POST("/add", route(controller.AddUser))

	// }

	//  注册
	r.POST("/register", route(controller.AddYonghu))
	//  登录
	r.POST("/login", route(controller.Login))
	//  token登录
	r.POST("/tokenLogin", route(controller.TokenLogin))

	// 获取所有设备信息
	r.POST("/devices", route(controller.DevicesList))

	// 查询dyid重复
	r.POST("/dyidrepeat", route(controller.CheckDouyinIDRepeat))

	// 添加任务 todo
	r.POST("/taskadd", route(controller.AddRenwu))

	// 获取任务
	r.POST("/getrenwu", route(controller.YonghuGetRenwu))

	// 任务操作 ---------------------------
	//  1 进入任务
	r.POST("/renwu_1", route(controller.RenwuStep1))
	//  2 礼物送出
	r.POST("/renwu_2", route(controller.RenwuStep2))
	//  3 任务提交
	r.POST("/renwu_3", route(controller.RenwuStep3))
	//  4  任务失败
	r.POST("/renwu_4", route(controller.RenwuStep4))
	//  5 主播提前下播
	r.POST("/renwu_5", route(controller.RenwuStep5))

	//////////////////
	// 转移所有子账号的余额到主账号上
	r.POST("/down_1", route(controller.Account1))

	// r.Use(JWTAuth())
	return r
}
