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

// Middle5 任务操作：5 主播提前下播
// 任务 tqjs+1 到达一定数量检测 并且暂停任务stop=1
func Middle5(req *db.RenwuRequest) (interface{}, error) {
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

	// renwu
	renwuStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, user.Rid)))
	if err != nil {
		return nil, err
	}
	renwu := db.Renwu{}
	err = json.Unmarshal([]byte(renwuStr), &renwu)
	if err != nil {
		return nil, err
	}

	/////////////////////////////
	// 提前结束数量 初始为放单数量一半，当用户反馈提前结束数量-1 数量为负数时用户可提前结束任务
	renwu.Tiqianjieshu += 1
	renwu.Stop = 1
	/////////////////////////////
	// update
	rb, err := json.Marshal(renwu)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid), string(rb))
	if err != nil {
		return nil, err
	}
	manager.addUpdate(renwu)

	return nil, nil
}

// Middle4 任务操作：4任务失败
// 取消用户任务，任务数量+1
func Middle4(req *db.RenwuRequest) (interface{}, error) {
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

	// renwu
	renwuStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, user.Rid)))
	if err != nil {
		return nil, err
	}
	renwu := db.Renwu{}
	err = json.Unmarshal([]byte(renwuStr), &renwu)
	if err != nil {
		return nil, err
	}

	/////////////////////////////
	rwlog, err := db.GetRwlogsByruandyonghuid(user.Uid, user.Rid)
	if err != nil {
		return nil, err
	}
	rwlog.Isadd = req.Isadd

	renwu.Shengyusl += 1

	user.Rid = -1 // if 0 gorm will ignore it
	user.Rwjd = -1
	/////////////////////////////
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
	rb, err := json.Marshal(renwu)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid), string(rb))
	if err != nil {
		return nil, err
	}

	manager.addUpdate(renwu)
	manager.addUpdate(user)
	manager.addUpdate(rwlog)

	return nil, nil
}

// Middle3 任务操作：3 任务提交
// 设备完成任务后提交服务器，服务器判断数据库用户信息rid是否正常大于0，如果正常读取任务信息，判断任务放单+时长是否小于当前时间，
// 防止设备提前结束任务，如果正常则给用户余额增加任务佣金。
func Middle3(req *db.RenwuRequest) (interface{}, error) {
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

	// renwu
	renwuStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, user.Rid)))
	if err != nil {
		return nil, err
	}
	renwu := db.Renwu{}
	err = json.Unmarshal([]byte(renwuStr), &renwu)
	if err != nil {
		return nil, err
	}

	/////////////////////////////
	// rid是否正常大于0
	if user.Rid <= 0 {
		return nil, errors.New("rid异常")
	}
	// 判断任务放单+时长是否小于当前时间
	if (renwu.Fangdantime + renwu.Shichang*3600) > int(time.Now().Unix()) {
		return nil, errors.New("设备提前结束任务异常")
	}
	// 用户余额增加任务佣金。
	user.Money += renwu.Rwmoney
	/////////////////////////////
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

	////// response

	res := make(map[string]interface{})
	res["code"] = 1
	res["money"] = renwu.Rwmoney

	return res, nil
}

// Middle2 任务操作：2 礼物送出
// 任务如果不需要送礼物或者设备提交送礼物完成后，设置用户信息rwjd=3，rwkstime为当前时间
func Middle2(req *db.RenwuRequest) (interface{}, error) {
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

	// renwu
	renwuStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, user.Rid)))
	if err != nil {
		return nil, err
	}
	renwu := db.Renwu{}
	err = json.Unmarshal([]byte(renwuStr), &renwu)
	if err != nil {
		return nil, err
	}

	/////////////////////////////
	// sfsl 任务权重 1=送礼物的任务 0=不送礼物的任务
	user.Rwjd = 3
	user.Rwkstime = int(time.Now().Unix())
	/////////////////////////////
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

	return nil, nil
}

// Middle1 任务操作：1 进入任务
// 设备进入指定直播间后会提交服务器，服务器先判断该任务需不需要送礼物，如果需要则设置rwjd=2并且rwkstime=当前时间
func Middle1(req *db.RenwuRequest) (interface{}, error) {
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

	// // 查询rwlog
	// rwlog, err := db.GetRwlogsByruandyonghuid(user.Uid, user.Rid)
	// if err != nil {
	// 	return nil, err
	// }

	// renwu
	renwuStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, user.Rid)))
	if err != nil {
		return nil, err
	}
	renwu := db.Renwu{}
	err = json.Unmarshal([]byte(renwuStr), &renwu)
	if err != nil {
		return nil, err
	}

	/////////////////////////////
	// 看看sfsl 是否送礼 是不是=1如果是1就设置 用户 rwjd=2 否则设置rwjd=3
	// 设备获取到任务后设置用户信息rid rwjd=1 rwkstime=任务领取时间

	if renwu.Sfsl == 0 {
		user.Rwjd = 3
	} else if renwu.Sfsl == 1 {
		user.Rwjd = 2
	}
	// user.Rwjd = 1
	// user.Rwkstime = int(time.Now().Unix())
	user.Rwkstime = int(time.Now().Unix())
	/////////////////////////////
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

	/////

	return nil, nil
}
