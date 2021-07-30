package controller

import (
	"douyin/web/db"
	"douyin/web/service"

	"github.com/gin-gonic/gin"
)

// Top1 注册
func Top1(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1}
	}

	// 获取ip
	req.Registerip = c.Request.Host

	_, err = service.Top1(&req)
	if err != nil {
		return Response{Msg: -1}
	}
	return Response{Msg: 1}
}

// Top2 Top2
func Top2(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1}
	}

	// 获取ip
	req.Registerip = c.Request.Host

	user, err := service.Top2(&req)
	if err != nil {
		return Response{Msg: -1}
	}
	//
	return user
}

// Top3 Top3
func Top3(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1, Code: -1}
	}

	res, err := service.Top3(&req)
	if err != nil {
		return Response{Msg: -1, Code: -1}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// Top5 Top5
func Top5(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1}
	}

	// // 获取ip
	req.Registerip = c.Request.Host

	res, err := service.Top5(&req)
	if err != nil {
		return Response{Msg: -1}
	}
	return res
}

// Top6 Top6
func Top6(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1, Code: -1}
	}

	res, err := service.Top6(&req)
	if err != nil {
		return Response{Msg: -1, Code: -1}
	}
	return Response{Data: res, Msg: 1, Code: 1}
}

// Top101  添加任务
func Top101(c *gin.Context) interface{} {
	req := db.AddRenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}

	res, err := service.Top101(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// Top1001_110  获取任务
func Top1001_110(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Top1001_110(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}
