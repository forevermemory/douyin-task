package service

import (
	"douyin/global"
	"douyin/web/db"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// DevicesList DevicesList
func DevicesList(req *db.YonghuRequest) (interface{}, error) {
	// 直接从redis取
	conn := global.REDIS.Get()
	defer conn.Close()

	userStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER_TOKEN, req.Token)))
	if err != nil {
		return nil, err
	}

	user := db.Yonghu{}
	err = json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		return nil, err
	}

	// 判断用户状态
	if user.Guishu != "" {
		// 非主账号   主账号guishu是空guishuuid是0
		return nil, errors.New("非主账号,操作失败")
	}

	// 根据主账号查询关联子账号
	qu := db.Yonghu{
		Guishu: user.Account,
		Page: db.Page{
			PageSize: 99999999,
		},
	}

	secondUsers, err := db.ListYonghu(&qu)
	if err != nil {
		return nil, err
	}
	// TODO 是否需要加入yonghu到redis

	return secondUsers, nil

}
