package controller

import (
	"douyin/web/db"
	"douyin/web/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Account1  转移所有子账号的余额到主账号上
func Account1(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.Account1(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

//////////////////////////////////////////////////////////////
// RenwuStep5  任务操作：5 主播提前下播
func RenwuStep5(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.RenwuStep5(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// RenwuStep4  任务操作：4任务失败
func RenwuStep4(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.RenwuStep4(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1, Code: 1, Data: res}
}

// RenwuStep3  任务操作：3 任务提交
func RenwuStep3(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.RenwuStep3(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return res
}

// RenwuStep2  任务操作：2 礼物送出
func RenwuStep2(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.RenwuStep2(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Code: 1, Msg: 1, Data: res}
}

// RenwuStep1  任务操作：1 进入任务
func RenwuStep1(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.RenwuStep1(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Code: 1, Msg: 1, Data: res}
}

// YonghuGetRenwu  获取任务
func YonghuGetRenwu(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}
	res, err := service.YonghuGetRenwu(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return res
}

// AddRenwu  添加任务
func AddRenwu(c *gin.Context) interface{} {
	req := db.RenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}

	uid, _ := c.Get("user_id")
	t, _ := uid.(int)
	req.Userid = t

	_, err = service.YonghuAddRenwu(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1}
}

// GetRenwuByID  get xxx by id
func GetRenwuByID(c *gin.Context) interface{} {
	_id := c.Param("oid")
	id, err := strconv.Atoi(_id)
	if err != nil {
		return Response{Code: -1}
	}
	_, err = service.GetRenwuByID(id)
	if err != nil {
		return Response{Code: -1}
	}
	return Response{Code: 0}
}

// ListRenwu // list by page condition
func ListRenwu(c *gin.Context) interface{} {
	var req = db.Renwu{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -1}
	}
	_, err = service.ListRenwu(&req)
	if err != nil {
		return Response{Code: -1}
	}
	return Response{Code: 0}
}
