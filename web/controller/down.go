package controller

import (
	"douyin/web/db"
	"douyin/web/service"

	"github.com/gin-gonic/gin"
)

// Down6  更新账户抖币余额
func Down6(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Down6(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// Down5  查询提现记录
func Down5(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Down5(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// Down4  更新账户信息
func Down4(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Down4(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// Down3  查询用户总余额
func Down3(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Down3(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return res
}

// Down2  提现
func Down2(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Down2(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// Down1  转移所有子账号的余额到主账号上
func Down1(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Down1(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}
