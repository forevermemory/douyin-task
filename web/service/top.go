package service

import (
	"douyin/global"
	"douyin/utils"
	"douyin/web/db"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Top1 add
func Top1(req *db.YonghuRequest) (interface{}, error) {
	user := db.Yonghu{
		Account:     req.User,
		Accountmd5:  utils.Md5Encrypt(req.User),
		Password:    req.Password,
		Passwordmd5: utils.Md5Encrypt(req.Password),

		Registertime: time.Now(),
		Registerip:   req.Registerip,
		State:        0,
		Guishu:       req.Sj,

		Zfb:     req.Zfb,
		Zfbname: req.Zfbname,
	}

	u, err := db.AddYonghu(&user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Top2 Top2
func Top2(req *db.YonghuRequest) (interface{}, error) {

	conn := global.REDIS.Get()
	defer conn.Close()

	////////////////////////////
	//  user password
	user, err := db.LoginUser(req.User, req.Password)
	if err != nil {
		return nil, err
	}

	// token
	claims := db.CustomClaims{
		Uid: user.Uid,
	}
	token, err := db.NewJWT().CreateToken(claims)
	if err != nil {
		return nil, err
	}

	user.Token = token
	user.Lastloginip = req.Registerip
	user.Lastlogintime = time.Now()

	//////////////////////////
	// redis ...
	uy, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	// set 两个redis key
	// user_id {user}
	// token {user}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid), string(uy))
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER_TOKEN, token), string(uy))
	if err != nil {
		return nil, err
	}

	manager.addUpdate(user)

	// response
	res := db.YonghuResponse{
		Msg:           1,
		Code:          1,
		Guishu:        user.Guishu,
		Money:         user.Money,
		LastLoginTime: user.Lastlogintime.String(),
		LastLoginIP:   user.Lastloginip,
		Token:         token,
		Userid:        user.Uid,
	}
	return res, nil

}

// Top3 Top3
func Top3(req *db.YonghuRequest) (interface{}, error) {
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
	///////////////////////////
	// 判断用户状态
	if user.Guishu != "" || user.Guishuuid != 0 {
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

	return secondUsers, nil

}

// Top5 Top5
func Top5(req *db.YonghuRequest) (interface{}, error) {
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
	///////////////////////////
	// 更新用户登陆时间
	user.Lastloginip = req.Registerip
	user.Lastlogintime = time.Now()
	///////////////////////////

	manager.addUpdate(user)

	// response
	res := db.YonghuResponse{
		Msg:           1,
		Guishu:        user.Guishu,
		Money:         user.Money,
		LastLoginTime: user.Lastlogintime.String(),
		LastLoginIP:   user.Lastloginip,
		Token:         user.Token,
		Userid:        user.Uid,
	}
	return res, nil
}

// Top6 Top6 查询dyid重复
func Top6(req *db.YonghuRequest) (interface{}, error) {
	// 查询dyid重复
	data, err := db.CheckDouyinIDRepeat(req.Dyid)
	if err != nil {
		return nil, err
	}

	return data, nil

}

// Top101 添加任务
func Top101(req *db.AddRenwuRequest) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	var err error

	//////////////
	renwu := &db.Renwu{
		Leixing:      req.Lx,
		Shengyusl:    req.Lx,
		Zongshuliang: req.Lx,
		Shichang:     req.Sj,
		Url:          req.Url,
		In:           req.In,
		Name:         req.Name,
		Gsfsp:        req.Gsfsp,
		Dzcs:         req.Dzcs,
		Biaoshi:      req.Ddh,
		Rwmoney:      req.Money,
		Sfsl:         req.Sfsl,
		Sfgj:         req.Sfgj,
		Rjbbh:        req.Rjbbh,
		Xtbbh:        req.Xtbbh,
		Zbid:         req.Zbid,
		Douyinid:     req.Userid,
	}

	renwu, err = db.AddRenwu(renwu)
	if err != nil {
		return nil, err
	}
	////////////
	// update
	rb, err := json.Marshal(renwu)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid), string(rb))
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// Top1001_110 获取任务
// 如果任务库里面有任务，则通过提取用户提交的信息，看看哪个任务满足条件，返回给用户，任务数量-1，
// 这里需要考虑高并发问题，防止任务数量放出去的大于任务总数量。

// 设备访问服务器获取任务 先查询该设备有没有历史任务未完成 如果有就返回历史任务
func Top1001_110(req *db.RenwuRequest) (interface{}, error) {
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
	///////////////////////////

	// get renwu 用户存在任务
	if user.Rid > 0 {
		renwuStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, user.Rid)))
		if err != nil {
			return nil, err
		}
		renwu := db.Renwu{}
		err = json.Unmarshal([]byte(renwuStr), &renwu)
		if err != nil {
			return nil, err
		}
		return renwu, nil
	}

	// 查询是否存在满足条件的任务
	/*
		bdzhqz 任务表的是否送礼 sfsl 2就是送礼物任务 1就是不送礼物任务
	*/
	QU := &db.Renwu{
		Sfsl: req.Bdzhqz,
		Page: db.Page{
			PageSize: 99999999,
		},
	}
	okrenwus, err := db.ListRenwu(QU)
	if err != nil {
		return nil, err
	}
	if len(okrenwus) == 0 {
		return nil, errors.New("无满足的任务")
	}

	toGetRenwu := okrenwus[0]
	if toGetRenwu.Shengyusl == 0 {
		return nil, errors.New("任务数量为0")
	}

	// lock
	_, ok := manager.renwuSet[toGetRenwu.Rid]
	if ok {
		// lock
		return nil, errors.New("renwu is lock")
	}
	// unlock --> add lock
	manager.renwuSet[toGetRenwu.Rid] = 1
	defer func() {
		// unlock
		delete(manager.renwuSet, toGetRenwu.Rid)
	}()

	user.Rid = toGetRenwu.Rid
	// num--
	toGetRenwu.Shengyusl = toGetRenwu.Shengyusl - 1
	// 记录任务添加历史记录
	txlog := db.Rwlogs{
		Uid:    user.Uid,
		Rid:    toGetRenwu.Rid,
		Userid: user.Uid,
		Zbid:   strconv.Itoa(toGetRenwu.Zbid),
		Isadd:  db.Rwlogs_isadd_GET_TASK,
		Day:    time.Now(),
	}
	manager.addCreate(&txlog)
	/////////////////////////////
	// update
	rb, err := json.Marshal(toGetRenwu)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, toGetRenwu.Rid), string(rb))
	if err != nil {
		return nil, err
	}

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
	manager.addUpdate(toGetRenwu)

	// //// response
	// res.Code = 1
	// res.Lx = renwu.Leixing
	// res.Dzcs = renwu.Dzcs
	// res.Xhc = renwu.Xianghuangche
	// res.Rwxh = ""
	// res.Time = 0
	// res.Name = renwu.Name
	// res.Id = renwu.Rid
	// res.Url = renwu.Url
	// res.Jrfs = 0
	// res.Sfsl = renwu.Sfsl
	// res.Sfgj = renwu.Sfgj
	// res.Rwjd = 0
	// res.Gsfsp = renwu.Gsfsp

	return toGetRenwu, nil

}
