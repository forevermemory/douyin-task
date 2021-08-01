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

	}
	user.Money += sonMoney

	/////////////////////
	// update
	manager.setUser(user)
	for _, son := range tmpSons {
		manager.setUser(son)
	}
	return nil, nil
}
