package controller

import (
	"douyin/web/db"
	"douyin/web/service"

	"github.com/gin-gonic/gin"
)

// Middle501  任务操作任务暂停 暂停任务，即使有剩余数量设备也不能获取该任务。
func Middle501(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Middle501(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// Middle5  任务操作：5 主播提前下播
func Middle5(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Middle5(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// Middle4  任务操作：4任务失败
func Middle4(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Middle4(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// Middle3  任务操作：3 任务提交
func Middle3(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Middle3(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return res
}

// Middle2  任务操作：2 礼物送出
func Middle2(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Middle2(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Code: 1, Msg: 1, Data: res}
}

// Middle1  任务操作：1 进入任务
func Middle1(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Middle1(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Code: 1, Msg: 1, Data: res}
}
