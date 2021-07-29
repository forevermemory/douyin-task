package controller

import (
	"douyin/web/db"
	"douyin/web/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CheckDouyinIDRepeat CheckDouyinIDRepeat
func CheckDouyinIDRepeat(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1, Code: -1}
	}

	// // 获取ip
	// req.Registerip = c.Request.Host

	res, err := service.CheckDouyinIDRepeat(&req)
	if err != nil {
		return Response{Msg: -1, Code: -1}
	}
	return Response{Data: res, Msg: 1, Code: 1}
}

// TokenLogin TokenLogin
func TokenLogin(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1}
	}

	// // 获取ip
	// req.Registerip = c.Request.Host

	res, err := service.TokenLogin(&req)
	if err != nil {
		return Response{Msg: -1}
	}
	return res
}

// Login Login
func Login(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1}
	}

	// 获取ip
	req.Registerip = c.Request.Host

	user, err := service.LoginUser(&req)
	if err != nil {
		return Response{Msg: -1}
	}
	//
	return user
}

// AddYonghu 注册
func AddYonghu(c *gin.Context) interface{} {
	var req = db.YonghuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Msg: -1}
	}

	// 获取ip
	req.Registerip = c.Request.Host

	_, err = service.AddYonghu(&req)
	if err != nil {
		return Response{Msg: -1}
	}
	return Response{Msg: 1}
}

/////////////////////////////////

// GetYonghuByID  get xxx by id
func GetYonghuByID(c *gin.Context) interface{} {
	_id := c.Param("oid")
	id, err := strconv.Atoi(_id)
	if err != nil {
		return Response{Code: -1}
	}
	_, err = service.GetYonghuByID(id)
	if err != nil {
		return Response{Code: -1}
	}
	return Response{Code: 0}
}

// ListYonghu // list by page condition
func ListYonghu(c *gin.Context) interface{} {
	var req = db.Yonghu{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -1}
	}
	_, err = service.ListYonghu(&req)
	if err != nil {
		return Response{Code: -1}
	}
	return Response{Code: 0}
}
