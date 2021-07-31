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

	case 1:
		return controller.Top1(c)
	case 2:
		return controller.Top2(c)
	case 3:
		return controller.Top3(c)
	case 5:
		return controller.Top5(c)
	case 6:
		return controller.Top6(c)
	case 101:
		return controller.Top101(c)
	case 1001:
		// middle
		switch req.Code {

		case 110:
			return controller.Top1001_110(c) // 获取任务
		case 1:
			return controller.Middle1(c)
		case 2:
			return controller.Middle2(c)
		case 3:
			return controller.Middle3(c)
		case 4:
			return controller.Middle4(c)
		case 5:
			return controller.Middle5(c)
		// case 501:
		// 	return controller.Middle501(c)
		// down
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
