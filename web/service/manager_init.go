package service

import (
	"douyin/global"
	"douyin/web/db"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func (m *RedisSyncToMysqlManager) initRenwu() {
	conn := global.REDIS.Get()
	defer conn.Close()

	// shengyusl
	renwus, err := db.RenuwuShenyuGreaterZero()
	if err != nil {

	}
	for _, renwu := range renwus {
		manager.setRenwu_hsetall(renwu)
	}
}

func (m *RedisSyncToMysqlManager) getRenwu_hgetall(renwuid int) (*db.Renwu, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwuid)

	res, err := redis.Values(conn.Do("hgetall", _key))
	if err != nil {
		return nil, err
	}
	u := new(db.Renwu)
	err = redis.ScanStruct(res, u)
	return u, nil
}
func (m *RedisSyncToMysqlManager) getRenwu_hmget(renwuid int, fields []string) (*db.Renwu, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwuid)

	res, err := redis.Values(conn.Do("hmget", _key, fields))
	if err != nil {
		return nil, err
	}
	u := new(db.Renwu)
	err = redis.ScanStruct(res, u)
	return u, nil
}

func (m *RedisSyncToMysqlManager) setRenwu_hsetall(renwu *db.Renwu) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid)

	conn.Do("hset", _key, "Rid", renwu.Rid)
	conn.Do("hset", _key, "Douyinid", renwu.Douyinid)
	conn.Do("hset", _key, "Zbid", renwu.Zbid)
	conn.Do("hset", _key, "Shengyusl", renwu.Shengyusl)
	conn.Do("hset", _key, "Shichang", renwu.Shichang)
	conn.Do("hset", _key, "Leixing", renwu.Leixing)
	conn.Do("hset", _key, "Zhuanshuuid", renwu.Zhuanshuuid)
	conn.Do("hset", _key, "Zongshuliang", renwu.Zongshuliang)
	conn.Do("hset", _key, "Stop", renwu.Stop)
	conn.Do("hset", _key, "Sbcs", renwu.Sbcs)
	conn.Do("hset", _key, "Fangdantime", renwu.Fangdantime)
	conn.Do("hset", _key, "Fangdanren", renwu.Fangdanren)
	conn.Do("hset", _key, "In", renwu.In)
	conn.Do("hset", _key, "Tiqianjieshu", renwu.Tiqianjieshu)
	conn.Do("hset", _key, "Name", renwu.Name)
	conn.Do("hset", _key, "Xianghuangche", renwu.Xianghuangche)
	conn.Do("hset", _key, "Biaoshi", renwu.Biaoshi)
	conn.Do("hset", _key, "Gsfsp", renwu.Gsfsp)
	conn.Do("hset", _key, "Dzcs", renwu.Dzcs)
	conn.Do("hset", _key, "Sfsl", renwu.Sfsl)
	conn.Do("hset", _key, "Ksqz", renwu.Ksqz)
	conn.Do("hset", _key, "Rwmoney", renwu.Rwmoney)
	conn.Do("hset", _key, "Rjbbh", renwu.Rjbbh)
	conn.Do("hset", _key, "Xtbbh", renwu.Xtbbh)

	// Rid int `Rid;primary_key;auto_increment;" json:"Rid" form:"Rid"`
	// Douyinid     string `douyinID" json:"douyinID" form:"douyinID"`             //主页ID
	// Zbid         int    `zbid" json:"zbid" form:"zbid"`                         //直播id
	// Shengyusl    int    `shengyusl" json:"shengyusl" form:"shengyusl"`          //剩余数量
	// Shichang     int    `shichang" json:"shichang" form:"shichang"`             //时长 小时
	// Leixing      int    `leixing" json:"leixing" form:"leixing"`                //任务类型
	// Zhuanshuuid  int    `zhuanshuUID" json:"zhuanshuUID" form:"zhuanshuUID"`    //指定用户UID放单-可以去掉
	// Zongshuliang int    `zongshuliang" json:"zongshuliang" form:"zongshuliang"` //总数量
	// Stop         int    `stop" json:"stop" form:"stop"`                         //0=正常 1=停止
	// Sbcs         int    `sbcs" json:"sbcs" form:"sbcs"`                         //失败次数
	// Url          string `url" json:"url" form:"url"`
	// Fangdantime  int    `fangdantime" json:"fangdantime" form:"fangdantime"` //放单时间
	// Fangdanren   string `fangdanren" json:"fangdanren" form:"fangdanren"`    //放单人可去掉
	// In           int    `in" json:"in" form:"in"`                            //进入直播模式
	// Tiqianjieshu int `tiqianjieshu" json:"tiqianjieshu" form:"tiqianjieshu"` //提前结束数量 初始为放单数量一半，当用户反馈提前结束数量-1 数量为负数时用户可提前结束任务
	// Name          string `name" json:"name" form:"name"`                            //主播名字
	// Xianghuangche int    `xianghuangche" json:"xianghuangche" form:"xianghuangche"` //可去掉
	// Biaoshi       string `biaoshi" json:"biaoshi" form:"biaoshi"`                   //订单号
	// Gsfsp         int    `gsfsp" json:"gsfsp" form:"gsfsp"`                         //光送粉丝牌-任务条件
	// Dzcs          int    `dzcs" json:"dzcs" form:"dzcs"`                            //点赞次数
	// Sfsl          int    `sfsl" json:"sfsl" form:"sfsl"`                            //任务权重 2就是送礼物任务 1就是不送礼物任务
	// Sfgj          int    `sfgj" json:"sfgj" form:"sfgj"`                            //0=不挂机任务 1=挂机任务
	// Ksqz          int    `ksqz" json:"ksqz" form:"ksqz"`                            //快手权重-可去掉
	// Rwmoney       int    `rwmoney" json:"rwmoney" form:"rwmoney"`                   //任务单价
	// Rjbbh         int    `rjbbh" json:"rjbbh" form:"rjbbh"`                         //任务限制-软件版本号
	// Xtbbh         int    `xtbbh" json:"xtbbh" form:"xtbbh"`                         //任务限制-系统版本号

}

func (m *RedisSyncToMysqlManager) initYonghu() {

	// 刚启动加载用户列表 只需要把onlie>0的加载就行了
	users, err := db.ListYonghuV3()
	if err != nil {
		return
	}
	for _, u := range users {
		m.yonghuSet[u.Uid] = 1
		m.setUser_hsetall(u)
	}
}

func (m *RedisSyncToMysqlManager) setUser_hsetall(user *db.Yonghu) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid)

	conn.Do("hset", _key, "UID", user.Uid)
	conn.Do("hset", _key, "Account", user.Account)
	conn.Do("hset", _key, "Accountmd5", user.Accountmd5)
	conn.Do("hset", _key, "Password", user.Password)
	conn.Do("hset", _key, "Passwordmd5", user.Passwordmd5)
	conn.Do("hset", _key, "Onlie", user.Onlie)
	conn.Do("hset", _key, "Lastlogintime", m.stringfyTime(user.Lastlogintime))
	conn.Do("hset", _key, "Lastloginip", user.Lastloginip)
	conn.Do("hset", _key, "Registertime", m.stringfyTime(user.Registertime))
	conn.Do("hset", _key, "Registerip", user.Registerip)
	conn.Do("hset", _key, "State", user.State)
	conn.Do("hset", _key, "Stateinformation", user.Stateinformation)
	conn.Do("hset", _key, "Token", user.Token)
	conn.Do("hset", _key, "Tokentime", user.Tokentime)
	conn.Do("hset", _key, "Guishu", user.Guishu)
	conn.Do("hset", _key, "Guishuuid", user.Guishuuid)
	conn.Do("hset", _key, "Shangjiuid", user.Shangjiuid)
	conn.Do("hset", _key, "Rid", user.Rid)
	conn.Do("hset", _key, "Rwkstime", user.Rwkstime)
	conn.Do("hset", _key, "Rwjd", user.Rwjd)
	conn.Do("hset", _key, "Dymz", user.Dymz)
	conn.Do("hset", _key, "Dyid", user.Dyid)
	conn.Do("hset", _key, "Dbye", user.Dbye)
	conn.Do("hset", _key, "Dyyz", user.Dyyz)
	conn.Do("hset", _key, "Ksyz", user.Ksyz)
	conn.Do("hset", _key, "Money", user.Money)
	conn.Do("hset", _key, "Zfb", user.Zfb)
	conn.Do("hset", _key, "Zfbname", user.Zfbname)
	conn.Do("hset", _key, "Xtbbh", user.Xtbbh)
	conn.Do("hset", _key, "Rjbbh", user.Rjbbh)
	conn.Do("hset", _key, "Cfdj", user.Cfdj)

	// Uid int `UID;primary_key;auto_increment;" json:"UID" form:"UID"`
	// Account          string    `Account" json:"Account" form:"Account"` //账号
	// Accountmd5       string    `AccountMD5" json:"AccountMD5" form:"AccountMD5"`
	// Password         string    `Password" json:"Password" form:"Password"`
	// Passwordmd5      string    `PasswordMD5" json:"PasswordMD5" form:"PasswordMD5"`
	// Onlie            int       `Onlie" json:"Onlie" form:"Onlie"`                   //websoket句柄
	// Lastlogintime    time.Time `LastLoginTime" json:"LastLoginTime" form:"LastLoginTime"`    //最后登录时间
	// Lastloginip      string    `LastLoginIP" json:"LastLoginIP" form:"LastLoginIP"`                //最后登录IP
	// Registertime     time.Time `RegisterTime" json:"RegisterTime" form:"RegisterTime"`             //注册时间
	// Registerip       string    `RegisterIP" json:"RegisterIP" form:"RegisterIP"`                   //注册IP
	// State            int       `State" json:"State" form:"State"`      //0=正常 1=冻结
	// Stateinformation string    `StateInformation" json:"StateInformation" form:"StateInformation"` //状态原因
	// Token            string    `Token" json:"Token" form:"Token"`
	// Tokentime        int       `TokenTime" json:"TokenTime" form:"TokenTime"`
	// Guishu           string    `guishu" json:"guishu" form:"guishu"`             //归属哪个用户
	// Guishuuid        int       `guishuUID" json:"guishuUID" form:"guishuUID"`    //归属用户uid
	// Shangjiuid       int       `shangjiUID" json:"shangjiUID" form:"shangjiUID"` //上级uid
	// Rid              int       `RID" json:"RID" form:"RID"`                      //任务ID
	// Rwkstime         int       `rwksTime" json:"rwksTime" form:"rwksTime"`       //任务开始时间
	// Rwjd             int       `rwjd" json:"rwjd" form:"rwjd"`                   //1=正在进直播间 2=正在送礼物 3=正在挂机
	// Dymz             string    `dymz" json:"dymz" form:"dymz"`                   //可去掉
	// Dyid             string    `dyid" json:"dyid" form:"dyid"`                   //抖音id
	// Dbye             int       `dbye" json:"dbye" form:"dbye"`                   //抖币余额
	// Dyyz             int       `dyyz" json:"dyyz" form:"dyyz"`                   //抖音验证状态0-1
	// Ksyz             int       `ksyz" json:"ksyz" form:"ksyz"`                   //可去掉
	// Money            int       `money" json:"money" form:"money"`
	// Zfb              string    `zfb" json:"zfb" form:"zfb"`             //支付宝
	// Zfbname          string    `zfbname" json:"zfbname" form:"zfbname"` //支付宝姓名
	// Xtbbh            int       `xtbbh" json:"xtbbh" form:"xtbbh"`       //系统版本号
	// Rjbbh            int       `rjbbh" json:"rjbbh" form:"rjbbh"`       //软件版本号
	// Cfdj int `cfdj" json:"cfdj" form:"cfdj"` // 抖音等级

}

func (m *RedisSyncToMysqlManager) getUser_hgetall(id int) (*db.Yonghu, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, id)

	res, err := redis.Values(conn.Do("hgetall", _key))
	if err != nil {
		return nil, err
	}
	u := new(db.Yonghu)
	err = redis.ScanStruct(res, u)
	return u, nil
}
