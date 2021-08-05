package service

import (
	"douyin/global"
	"douyin/utils"
	"douyin/web/db"
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
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

		Shangjiuid: req.Shangjiuid,
	}

	u, err := db.AddYonghu(&user)
	if err != nil {
		return nil, err
	}

	// add yonghu set
	manager.yonghuSet[user.Uid] = 1
	return u, nil
}

// Top2 Top2
func Top2(req *db.YonghuRequest) (interface{}, error) {

	////////////////////////////
	//  user password
	user, err := db.LoginUser(req.User, req.Password)
	if err != nil {
		return nil, err
	}
	if user.Onlie > 0 {
		return nil, errors.New("repeat login")
	}

	// token TODO
	token, err := utils.GetToken(user.Uid)
	if err != nil {
		return nil, err
	}

	user.Token = token
	user.Lastloginip = req.Registerip
	user.Lastlogintime = time.Now()

	//////////////////////////
	// redis ...
	user.UpdateType = db.USER_UPDATE_TOP2
	manager.setUser(user)

	// 登陆之后 再加入
	// add yonghu set
	manager.yonghuSet[user.Uid] = 1

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

	user, err := manager.getUserByToken(req.Token)
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

	user, err := manager.getUserByToken(req.Token)
	if err != nil {
		return nil, err
	}
	if user.Onlie > 0 {
		return nil, errors.New("repeat login")
	}

	///////////////////////////
	// 更新用户登陆时间
	user.Lastloginip = req.Registerip
	user.Lastlogintime = time.Now()
	///////////////////////////
	user.UpdateType = db.USER_UPDATE_TOP5
	manager.setUser(user)
	// 登陆之后 再加入
	// add yonghu set
	manager.yonghuSet[user.Uid] = 1

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

		// tqjs
		Tiqianjieshu: int(req.Lx / 2),
		Stop:         0,
	}

	renwu, err = db.AddRenwu(renwu)
	if err != nil {
		return nil, err
	}
	////////////
	renwu.UpdateType = db.RENWU_UPDATE_ALL
	manager.setRenwu(renwu)

	return nil, nil

}

// Top1001_110 获取任务
// 如果任务库里面有任务，则通过提取用户提交的信息，看看哪个任务满足条件，返回给用户，任务数量-1，
// 这里需要考虑高并发问题，防止任务数量放出去的大于任务总数量。

// 设备访问服务器获取任务 先查询该设备有没有历史任务未完成 如果有就返回历史任务
func Top1001_110(req *db.RenwuRequest) (interface{}, error) {

	user, err := manager.getUserByToken(req.Token)
	if err != nil {
		return nil, err
	}
	///////////////////////////
	// get renwu 用户存在任务
	if user.Rid > 0 {
		renwu, err := manager.getRenwu(user.Rid)
		if err != nil {
			return nil, err
		}
		return renwu, nil
	}

	// 查询是否存在满足条件的任务
	var toGetRenwu *db.Renwu

	// 遍历redis的所有任务
	for _, rw := range manager.renwuIDSet {
		// 1. 是否送礼
		if rw.Sfsl != req.Bdzhqz {
			continue
		}
		// 2. 如果任务类型是2 就说明一天只能领取一次
		renwulog, err := manager.getRenwulog(user.Uid, user.Rid)
		if !errors.Is(gorm.ErrRecordNotFound, err) {
			continue
		}
		if rw.Leixing == 2 {
			if renwulog != nil {
				if renwulog.Day.Day() == time.Now().Day() {
					// 说明是当天领的任务
					continue
				}
			}
		} else if rw.Leixing == 4 {
			// 3. 如果任务类型是4 点赞次数是777 就说明这个任务一个用户只能领取一次
			// 就从任务完成记录里面查找用户uid 如果有就跳过这个任务
			if rw.Dzcs == 777 {
				if renwulog != nil {
					continue
				}
			}
		}

		// 4. 用户5分钟内做这个任务失败了 下次就不让他领取这个任务
		if renwulog != nil {
			if renwulog.Isadd == db.Rwlogs_isadd_ABADON_TASK_NOT_IN {
				continue
			}
		}
		// 5.  还有个条件是限制一个任务  同ip只能进多少台
		_limit, err := manager.getIpLimit(req.Ipaddr, rw.Rid)
		if err != nil {
			continue
		}
		if _limit >= global.MAX_IP_TASK {
			continue
		}
		// 满足条件的任务////////////////
		// 获取任务 条件都满足后 进锁 进不去就换下一个任务判断条件
		// lock 任务
		_, ok := manager.renwuLock[toGetRenwu.Rid]
		if ok {
			// locked
			continue
		}
		// unlock --> add lock
		manager.renwuLock[toGetRenwu.Rid] = 1
		defer func() {
			// unlock
			delete(manager.renwuLock, toGetRenwu.Rid)
		}()
		toGetRenwu = rw
		break
	}

	if toGetRenwu == nil {
		return nil, errors.New("无满足的任务")
	}

	if toGetRenwu.Shengyusl == 0 {
		return nil, errors.New("当前任务数量为0")
	}
	//
	// 判断任务是否满足
	//

	////////////////////////// 添加任务了

	user.Rid = toGetRenwu.Rid
	// num--
	toGetRenwu.Shengyusl = toGetRenwu.Shengyusl - 1
	// 记录任务添加历史记录
	rwlog := &db.Rwlogs{
		Uid:    user.Uid,
		Rid:    toGetRenwu.Rid,
		Userid: user.Uid,
		Zbid:   strconv.Itoa(toGetRenwu.Zbid),
		Isadd:  db.Rwlogs_isadd_GET_TASK,
		Day:    time.Now(),
	}
	_, err = db.AddRwlogs(rwlog)
	if err != nil {
		return nil, err
	}

	/////////////////////////////
	// update
	manager.setIpLimit(req.Ipaddr, toGetRenwu.Rid)
	toGetRenwu.UpdateType = db.RENWU_UPDATE_Shengyusl
	manager.setRenwu(toGetRenwu)
	user.UpdateType = db.USER_UPDATE_ONLY_RID
	manager.setUser(user)

	// //// response  你根据需求添加或者减少
	res := &db.RenwuResponse{}
	res.Code = 1
	res.Lx = toGetRenwu.Leixing
	res.Dzcs = toGetRenwu.Dzcs
	res.Xhc = toGetRenwu.Xianghuangche
	res.Rwxh = ""
	res.Time = 0
	res.Name = toGetRenwu.Name
	res.Id = toGetRenwu.Rid
	res.Url = toGetRenwu.Url
	res.Jrfs = 0
	res.Sfsl = toGetRenwu.Sfsl
	res.Sfgj = toGetRenwu.Sfgj
	res.Rwjd = 0
	res.Gsfsp = toGetRenwu.Gsfsp

	// 5、用户获取任务不能把任务表的所有信息都返回给用户，这个接口里面有。可以考虑把需要的信息写入
	return res, nil

}
