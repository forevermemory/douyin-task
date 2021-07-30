package controller

import (
	"douyin/web/db"
	"douyin/web/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

//////////////////////////////////////////////////////////////

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
	req := db.AddRenwuRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		return Response{Code: -101}
	}

	_, err = service.YonghuAddRenwuZhen(&req)
	if err != nil {
		return Response{Code: -101}
	}
	return Response{Msg: 1}
}

// // AddRenwu  添加任务
// func AddRenwu(c *gin.Context) interface{} {
// 	req := db.AddRenwuRequest{}
// 	err := c.ShouldBind(&req)
// 	if err != nil {
// 		return Response{Code: -101}
// 	}

// 	uid, _ := c.Get("user_id")
// 	t, _ := uid.(int)
// 	req.Userid = t

// 	_, err = service.YonghuAddRenwuZhen(&req)
// 	if err != nil {
// 		return Response{Code: -101}
// 	}
// 	return Response{Msg: 1}
// }

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
