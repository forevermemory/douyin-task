package service

import (
	"douyin/global"
	"douyin/web/db"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Down6 更新账户抖币余额
func Down6(req *db.RenwuRequest) (interface{}, error) {
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
	//////////////////////////
	user.Money += req.Money
	//////////////////////////
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

	manager.addUpdate(&user)

	return nil, nil
}

// Down5 查询提现记录
func Down5(req *db.RenwuRequest) (interface{}, error) {
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

	//////////////////////////
	qu := &db.Txlogs{
		Uid: user.Uid,
		Page: db.Page{
			PageSize: 99999999,
		},
	}
	los, err := db.ListTxlogs(qu)
	if err != nil {
		return nil, err
	}

	/////////////////////////

	return los, nil

}

// Down4 更新账户信息
func Down4(req *db.RenwuRequest) (interface{}, error) {
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

	//////////////////////////
	user.Dyid = req.Dyid
	user.Dbye = req.Dbye
	user.Ksyz = req.Bdksyz
	user.Dyyz = req.Bddyyz
	user.Xtbbh = req.Bbh
	user.Xtbbh = req.Xtbbh
	user.Cfdj = req.Cfdj

	//////////////////////////
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

	manager.addUpdate(&user)

	return nil, nil
}

// Down3 查询用户总余额
func Down3(req *db.RenwuRequest) (interface{}, error) {
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

	///////////////////////
	// 走mysql 查询子用户
	sons, err := db.ListYonghuBySuper(user.Account)
	if err != nil {
		return nil, err
	}
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
		if son.Money == -1 {
			continue
		}
		user.Money += son.Money
	}

	//////////////////////

	res := make(map[string]interface{})
	res["msg"] = 1
	res["money"] = user.Money

	return res, nil
}

// Down2 提现
func Down2(req *db.RenwuRequest) (interface{}, error) {
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

	///////////////////////
	user.Money -= req.Money
	if user.Money < 0 {
		return nil, errors.New("提现金额不足")
	}
	// 保存提现记录
	txl := db.Txlogs{
		Uid:   user.Uid,
		Money: req.Money,
		Zfb:   user.Zfb,
		Isadd: 0,
		Day:   time.Now(),
	}
	_, err = db.AddTxlogs(&txl)
	if err != nil {
		return nil, err
	}
	//////////////////////////
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

	manager.addUpdate(&user)

	return nil, nil
}

// Down1
// 转移所有子账号的余额到主账号上
func Down1(req *db.RenwuRequest) (interface{}, error) {
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
		if stmp.Money > 0 {
			sonMoney += stmp.Money
			stmp.Money = -1 // qwq
			tmpSons = append(tmpSons, son)
		}
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

	manager.addUpdate(&user)
	for _, son := range sons {
		manager.addUpdate(son)
	}

	return nil, nil
}
