package service

import (
	"douyin/global"
	"douyin/web/db"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Account1 任务操作：5 主播提前下播
// 转移所有子账号的余额到主账号上
func Account1(req *db.RenwuRequest) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	// user
	userStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER_TOKEN, req.Token)))
	if err != nil {
		return nil, err
	}
	user := db.Yonghu{}
	err = json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		return nil, err
	}

	// 走mysql 查询子用户
	sons, err := db.ListYonghuBySuper(user.Account)
	if err != nil {
		return nil, err
	}
	// 转移所有子账号的余额到主账号上
	var sonMoney int
	tmpSons := make([]*db.Yonghu, 0)
	for _, son := range sons {
		// 这里只取son的id 再从redis取
		stmp := db.Yonghu{}

		userStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, son.Uid)))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(userStr), &stmp)
		if err != nil {
			return nil, err
		}
		//////////////////
		sonMoney += stmp.Money
		stmp.Money = -1 // qwq
		tmpSons = append(tmpSons, son)
		//////////////////

	}

	for _, son := range tmpSons {
		// update
		uy, err := json.Marshal(son)
		if err != nil {
			return nil, err
		}
		_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, son.Uid), string(uy))
		if err != nil {
			return nil, err
		}
		_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER_TOKEN, son.Token), string(uy))
		if err != nil {
			return nil, err
		}
	}

	user.Money += sonMoney
	// update
	uy, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid), string(uy))
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER_TOKEN, req.Token), string(uy))
	if err != nil {
		return nil, err
	}

	manager.addUpdate(user)
	for _, son := range sons {
		manager.addUpdate(son)
	}

	return nil, nil
}
