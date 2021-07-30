package router

import (
	"bytes"
	"douyin/web/controller"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Api struct {
}

var api = &Api{}

type BaseRequest struct {
	ID   int `json:"ID"`
	Code int `json:"code"`
}

type BaseResponse struct {
	Code int         `json:"code,omitempty"`
	Msg  int         `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func (a *Api) Handle(c *gin.Context) interface{} {
	req := BaseRequest{}

	buf, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	///////////
	err := c.ShouldBind(&req)
	if err != nil {
		return BaseResponse{Code: -101}
	}
	///////////
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	// 根据ID和code分发
	switch req.ID {
	case 1001:
		// //  1 进入任务
		// r.POST("/renwu_1", wrap(controller.RenwuStep1))
		// //  2 礼物送出
		// r.POST("/renwu_2", wrap(controller.RenwuStep2))
		// //  3 任务提交
		// r.POST("/renwu_3", wrap(controller.RenwuStep3))
		// //  4  任务失败
		// r.POST("/renwu_4", wrap(controller.RenwuStep4))
		// //  5 主播提前下播
		// r.POST("/renwu_5", wrap(controller.RenwuStep5))

		// middle
		switch req.Code {
		case 1:
			return controller.Down1(c)
		case 2:
			return controller.Down2(c)
		case 3:
			return controller.Down3(c)
		case 4:
			return controller.Down4(c)
		case 5:
			return controller.Down5(c)
		}
	case 101:
		// down
		switch req.Code {
		case 6:
			return controller.Down1(c)
		case 7:
			return controller.Down2(c)
		case 8:
			return controller.Down3(c)
		case 9:
			return controller.Down4(c)
		case 10:
			return controller.Down5(c)
		case 11:
			return controller.Down6(c)
		}

	}

	return nil
}

// //  注册
// r.POST("/register", wrap(controller.AddYonghu))
// //  登录
// r.POST("/login", wrap(controller.Login))
// //  token登录
// r.POST("/tokenLogin", wrap(controller.TokenLogin))

// // 获取所有设备信息
// r.POST("/devices", wrap(controller.DevicesList))

// // 查询dyid重复
// r.POST("/dyidrepeat", wrap(controller.CheckDouyinIDRepeat))

// // 添加任务 todo
// r.POST("/taskadd", wrap(controller.AddRenwu))

// // 获取任务
// r.POST("/getrenwu", wrap(controller.YonghuGetRenwu))

// // 任务操作 ---------------------------
