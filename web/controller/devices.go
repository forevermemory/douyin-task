package controller

import (
	"douyin/web/db"
	"douyin/web/service"

	"github.com/gin-gonic/gin"
)

// DevicesList DevicesList
func DevicesList(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1, Code: -1}
	}

	// // 获取ip
	// req.Registerip = c.Request.Host

	res, err := service.DevicesList(&req)
	if err != nil {
		return Response{Msg: -1, Code: -1}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}
