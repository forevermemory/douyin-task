package db

import (
	"douyin/global"
	"time"

	"github.com/jinzhu/gorm"
)

type RenwuResponse struct {
	Code  int `json:"code"`
	On    int `json:"on"`
	Money int `json:"money"`

	Lx    int    `json:"lx"`               //  //任务类型
	Dzcs  int    `json:"dzcs" form:"dzcs"` //点赞次数
	Xhc   int    `json:"xhc"`
	Rwxh  string `json:"rwxh"` // 订单号
	Time  int    `json:"time"` // 时间
	Name  string `json:"name"`
	Id    int    `json:"id"`
	Url   string `json:"url"`
	Jrfs  int    `json:"jrfs"`
	Sfsl  int    `json:"sfsl"`
	Sfgj  int    `json:"sfgj"`
	Rwjd  int    `json:"rwjd"`
	Gsfsp int    `json:"gsfsp"`
}

type RenwuRequest struct {
	ID    int    `json:"ID" form:"ID"`
	Code  int    `json:"code" form:"code"` //
	Token string `json:"token" form:"token"`

	Dzcs int `gorm:"column:dzcs" json:"dzcs" form:"dzcs"` //点赞次数

	Rwmoney int `json:"rwmoney" form:"rwmoney"` //任务单价
	Sfsl    int `json:"sfsl" form:"sfsl"`       //任务权重 1=送礼物的任务 0=不送礼物的任务
	Sfgj    int `json:"sfgj" form:"sfgj"`       //0=不挂机任务 1=挂机任务

	Zbid   int `json:"zbid" form:"zbid"`     //直播id
	Userid int `json:"userid" form:"userid"` //直播id

	Dyid  string ` json:"dyid" form:"dyid"`  //抖音id
	Dbye  int    ` json:"dbye" form:"dbye"`  //抖币余额
	Dyyz  int    ` json:"dyyz" form:"dyyz"`  //抖音验证状态0-1
	Ksyz  int    ` json:"ksyz" form:"ksyz"`  //可去掉
	Rjbbh int    `json:"rjbbh" form:"rjbbh"` //任务限制-软件版本号

	Money int `json:"money" form:"money"` //

	Bddyyz      int    ` json:"bddyyz" form:"bddyyz"` //抖音验证状态0-1
	Bdksyz      int    ` json:"bdksyz" form:"bdksyz"` //可去掉
	Bdzhqz      int    ` json:"bdzhqz" form:"bdzhqz"` //对应任务sfgj
	Dysign      string ` json:"dysign" form:"dysign"`
	UserNameAdd int    ` json:"user_name_add" form:"user_name_add"`

	Bbh   int `json:"bbh" form:"bbh"`     //任务限制-软件版本号
	Xtbbh int `json:"xtbbh" form:"xtbbh"` //任务限制-系统版本号
	Cfdj  int ` json:"cfdj" form:"cfdj"`  // 抖音等级
	Pb80  int ` json:"pb80" form:"pb80"`

	Isadd int `json:"isadd" form:"isadd"` // 失败原因

	//
	Ipaddr string
}

type AddRenwuRequest struct {
	ID     int    `json:"ID" form:"ID"`
	Lx     int    `json:"lx" form:"lx"`
	Sl     int    `json:"sl" form:"sl"`
	Sj     int    `json:"sj" form:"sj"`
	Url    string `json:"url" form:"url"`
	In     int    `json:"in" form:"in"`
	Name   string `json:"name" form:"name"`
	Gsfsp  int    `json:"gsfsp" form:"gsfsp"`   //光送粉丝牌-任务条件
	Dzcs   int    `json:"dzcs" form:"dzcs"`     //点赞次数
	Ddh    string ` json:"ddh" form:"ddh"`      //订单号
	Money  int    ` json:"money" form:"money"`  //订单号
	Sfsl   int    `json:"sfsl" form:"sfsl"`     //任务权重 1=送礼物的任务 0=不送礼物的任务
	Sfgj   int    `gjson:"sfgj" form:"sfgj"`    //0=不挂机任务 1=挂机任务
	Rjbbh  int    `json:"rjbbh" form:"rjbbh"`   //任务限制-软件版本号
	Xtbbh  int    ` json:"xtbbh" form:"xtbbh"`  //任务限制-系统版本号
	Zbid   int    `json:"zbid" form:"zbid"`     //直播id
	Userid string `json:"userid" form:"userid"` //直播id

	// IsOnlyOneTime int `json:"is_only_one_time" form:"is_only_one_time"` //是否一个用户只能领取一次 0否 1是
	// Lqzbyc        int `json:"lqzbyc" form:"lqzbyc"`                     //一天只能领取那个主播任务一次 0 否 1是
	// Ipsync        int ` json:"ipsync" form:"ipsync"`                    //  同ip只能进多少台

}

type Renwu struct {
	Rid int `gorm:"column:Rid;primary_key;auto_increment;" json:"Rid" form:"Rid"`

	Douyinid     string `gorm:"column:douyinID" json:"douyinID" form:"douyinID"`             //主页ID
	Zbid         int    `gorm:"column:zbid" json:"zbid" form:"zbid"`                         //直播id
	Shengyusl    int    `gorm:"column:shengyusl" json:"shengyusl" form:"shengyusl"`          //剩余数量
	Shichang     int    `gorm:"column:shichang" json:"shichang" form:"shichang"`             //时长 小时
	Leixing      int    `gorm:"column:leixing" json:"leixing" form:"leixing"`                //任务类型
	Zhuanshuuid  int    `gorm:"column:zhuanshuUID" json:"zhuanshuUID" form:"zhuanshuUID"`    //指定用户UID放单-可以去掉
	Zongshuliang int    `gorm:"column:zongshuliang" json:"zongshuliang" form:"zongshuliang"` //总数量
	Stop         int    `gorm:"column:stop" json:"stop" form:"stop"`                         //0=正常 1=停止
	Sbcs         int    `gorm:"column:sbcs" json:"sbcs" form:"sbcs"`                         //失败次数
	Url          string `gorm:"column:url" json:"url" form:"url"`
	Fangdantime  int    `gorm:"column:fangdantime" json:"fangdantime" form:"fangdantime"` //放单时间
	Fangdanren   string `gorm:"column:fangdanren" json:"fangdanren" form:"fangdanren"`    //放单人可去掉
	In           int    `gorm:"column:in" json:"in" form:"in"`                            //进入直播模式

	Tiqianjieshu int `gorm:"column:tiqianjieshu" json:"tiqianjieshu" form:"tiqianjieshu"` //提前结束数量 初始为放单数量一半，当用户反馈提前结束数量-1 数量为负数时用户可提前结束任务

	Name          string `gorm:"column:name" json:"name" form:"name"`                            //主播名字
	Xianghuangche int    `gorm:"column:xianghuangche" json:"xianghuangche" form:"xianghuangche"` //可去掉
	Biaoshi       string `gorm:"column:biaoshi" json:"biaoshi" form:"biaoshi"`                   //订单号
	Gsfsp         int    `gorm:"column:gsfsp" json:"gsfsp" form:"gsfsp"`                         //光送粉丝牌-任务条件
	Dzcs          int    `gorm:"column:dzcs" json:"dzcs" form:"dzcs"`                            //点赞次数
	Sfsl          int    `gorm:"column:sfsl" json:"sfsl" form:"sfsl"`                            //任务权重 2就是送礼物任务 1就是不送礼物任务
	Sfgj          int    `gorm:"column:sfgj" json:"sfgj" form:"sfgj"`                            //0=不挂机任务 1=挂机任务
	Ksqz          int    `gorm:"column:ksqz" json:"ksqz" form:"ksqz"`                            //快手权重-可去掉
	Rwmoney       int    `gorm:"column:rwmoney" json:"rwmoney" form:"rwmoney"`                   //任务单价
	Rjbbh         int    `gorm:"column:rjbbh" json:"rjbbh" form:"rjbbh"`                         //任务限制-软件版本号
	Xtbbh         int    `gorm:"column:xtbbh" json:"xtbbh" form:"xtbbh"`                         //任务限制-系统版本号

	// IsOnlyOneTime int `gorm:"column:is_only_one_time" json:"is_only_one_time" form:"is_only_one_time"` //是否一个用户只能领取一次 0否 1是
	// Lqzbyc        int `gorm:"column:lqzbyc" json:"lqzbyc" form:"lqzbyc"`                               //一天只能领取那个主播任务一次 0 否 1是
	// Ipsync        int `gorm:"column:ipsync" json:"ipsync" form:"ipsync"`                               //  同ip只能进多少台

	UpdateType int `gorm:"-" json:"-" form:"-"`

	Page
}

// TableName 表名
func (o *Renwu) TableName() string {
	return "renwu"
}

// DeleteRenwu 根据id删除
func DeleteRenwu(id int) error {
	db := global.MYSQL
	return db.Table("renwu").Where("Rid = ?", id).Error
}

// GetRenwuByID 根据id查询一个
func GetRenwuByID(id int) (*Renwu, error) {
	db := global.MYSQL
	o := &Renwu{}
	err := db.Table("renwu").Where("Rid = ?", id).First(o).Error
	return o, err
}

// AddRenwu 新增
func AddRenwu(o *Renwu, tx ...*gorm.DB) (*Renwu, error) {
	db := global.MYSQL
	if len(tx) > 0 {
		db = tx[0]
	}
	err := db.Create(o).Error
	return o, err
}

const (
	RENWU_UPDATE_ALL               int = 1
	RENWU_UPDATE_Shengyusl         int = 2
	RENWU_UPDATE_STOP_Tiqianjieshu int = 3
	RENWU_UPDATE_MIDDLE_4          int = 4
)

// UpdateRenwu 修改
func UpdateRenwu(o *Renwu, _type int) (*Renwu, error) {
	db := global.MYSQL

	u2 := new(Renwu)

	switch _type {
	case RENWU_UPDATE_ALL:
		u2 = o
		u2.Rid = 0
	case RENWU_UPDATE_Shengyusl:
		u2.Shengyusl = o.Shengyusl
	case RENWU_UPDATE_STOP_Tiqianjieshu:
		u2.Tiqianjieshu = o.Tiqianjieshu
		u2.Stop = o.Stop
	case RENWU_UPDATE_MIDDLE_4:

	}

	err := db.Table("renwu").Where("Rid=?", o.Rid).Updates(u2).Error
	return o, err
}

// ListRenwu 分页条件查询
func ListRenwu(o *Renwu) ([]*Renwu, error) {
	db := global.MYSQL
	res := make([]*Renwu, 0)
	err := db.Table("renwu").Where(o).Offset((o.PageNo - 1) * o.PageSize).Limit(o.PageSize).Find(&res).Error
	return res, err
}

// CountRenwu 条件数量
func CountRenwu(o *Renwu) (int, error) {
	db := global.MYSQL
	var count int
	err := db.Table("renwu").Where(o).Count(&count).Error
	return count, err
}

/////////////////////

//
// RenuwuShenyuGreaterZero >0
// 同步任务表只获取下单时间 5 小时内的
func RenuwuShenyuGreaterZero() ([]*Renwu, error) {
	db := global.MYSQL
	res := make([]*Renwu, 0)

	_time := time.Now().Unix() - 3600*5
	err := db.Table("renwu").Where("shengyusl > 0 and fangdantime >= ?", _time).Find(&res).Error
	return res, err
}
