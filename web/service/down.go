package service

import (
	"douyin/global"
	"douyin/web/db"
	"errors"
	"time"
)

// Down6 更新账户抖币余额
func Down6(req *db.RenwuRequest) (interface{}, error) {
	user, err := manager.getUserByToken(req.Token)
	if err != nil {
		return nil, err
	}

	//////////////////////////
	user.Money += req.Money
	//////////////////////////
	// update
	manager.setUser(user)

	return nil, nil
}

// Down5 查询提现记录
func Down5(req *db.RenwuRequest) (interface{}, error) {
	user, err := manager.getUserByToken(req.Token)
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

	user, err := manager.getUserByToken(req.Token)
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
	manager.setUser(user)

	return nil, nil
}

// Down3 查询用户总余额
func Down3(req *db.RenwuRequest) (interface{}, error) {
	user, err := manager.getUserByToken(req.Token)
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
		tmps, err := manager.getUser(son.Uid)
		if err != nil {
			return nil, err
		}
		if tmps.Money == -1 {
			continue
		}
		user.Money += tmps.Money
	}

	//////////////////////

	res := make(map[string]interface{})
	res["msg"] = 1
	res["money"] = user.Money

	return res, nil
}

// Down2 提现
func Down2(req *db.RenwuRequest) (interface{}, error) {
	user, err := manager.getUserByToken(req.Token)
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
	manager.setUser(user)

	return nil, nil
}

// Down1
// 转移所有子账号的余额到主账号上
func Down1(req *db.RenwuRequest) (interface{}, error) {
	user, err := manager.getUserByToken(req.Token)
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

	ts := make(map[int]int) // sonid money

	for _, son := range sons {
		// 这里只取son的id 再从redis取

		stmp, err := manager.getUser(son.Uid)
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
		ts[stmp.Uid] = stmp.Money

	}
	user.Money += sonMoney

	/////////////////////
	// update
	manager.setUser(user)
	for _, son := range tmpSons {
		// 主要防止转移过程中子账户有任务完成 增加金币
		// 方案: 再获取一次子账号 防止再增加的金币被刷没了
		stmp, err := manager.getUser(son.Uid)
		if err != nil {
			return nil, err
		}
		if ts[son.Uid] != son.Money {
			// 金额发生变动了 这里肯定是增加了
		}
		// 减去旧的
		stmp.Money = stmp.Money - ts[son.Uid]
		manager.setUser(stmp)
	}
	return nil, nil
}
