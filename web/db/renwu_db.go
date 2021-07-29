package db

import (
	"douyin/global"

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
	ID int

	Dzcs int `gorm:"column:dzcs" json:"dzcs" form:"dzcs"` //点赞次数

	Rwmoney int `json:"rwmoney" form:"rwmoney"` //任务单价
	Sfsl    int ` json:"sfsl" form:"sfsl"`      //任务权重 1=送礼物的任务 0=不送礼物的任务
	Sfgj    int ` json:"sfgj" form:"sfgj"`      //0=不挂机任务 1=挂机任务
	Rjbbh   int ` json:"rjbbh" form:"rjbbh"`    //任务限制-软件版本号
	Xtbbh   int ` json:"xtbbh" form:"xtbbh"`    //任务限制-系统版本号
	Zbid    int ` json:"zbid" form:"zbid"`      //直播id
	Userid  int ` json:"userid" form:"userid"`  //直播id

	Token string ` json:"token" form:"token"`
}

type Renwu struct {
	Rid int `gorm:"column:Rid;primary_key;auto_increment;" json:"Rid" form:"Rid"`

	Douyinid      string `gorm:"column:douyinID" json:"douyinID" form:"douyinID"`             //主页ID
	Zbid          int    `gorm:"column:zbid" json:"zbid" form:"zbid"`                         //直播id
	Shengyusl     int    `gorm:"column:shengyusl" json:"shengyusl" form:"shengyusl"`          //剩余数量
	Shichang      int    `gorm:"column:shichang" json:"shichang" form:"shichang"`             //时长
	Leixing       int    `gorm:"column:leixing" json:"leixing" form:"leixing"`                //任务类型
	Zhuanshuuid   int    `gorm:"column:zhuanshuUID" json:"zhuanshuUID" form:"zhuanshuUID"`    //指定用户UID放单-可以去掉
	Zongshuliang  int    `gorm:"column:zongshuliang" json:"zongshuliang" form:"zongshuliang"` //总数量
	Stop          int    `gorm:"column:stop" json:"stop" form:"stop"`                         //0=正常 1=停止
	Sbcs          int    `gorm:"column:sbcs" json:"sbcs" form:"sbcs"`                         //失败次数
	Url           string `gorm:"column:url" json:"url" form:"url"`
	Fangdantime   int    `gorm:"column:fangdantime" json:"fangdantime" form:"fangdantime"`       //放单时间
	Fangdanren    string `gorm:"column:fangdanren" json:"fangdanren" form:"fangdanren"`          //放单人可去掉
	In            int    `gorm:"column:in" json:"in" form:"in"`                                  //进入直播模式
	Tiqianjieshu  int    `gorm:"column:tiqianjieshu" json:"tiqianjieshu" form:"tiqianjieshu"`    //提前结束数量 初始为放单数量一半，当用户反馈提前结束数量-1 数量为负数时用户可提前结束任务
	Name          string `gorm:"column:name" json:"name" form:"name"`                            //主播名字
	Xianghuangche int    `gorm:"column:xianghuangche" json:"xianghuangche" form:"xianghuangche"` //可去掉
	Biaoshi       string `gorm:"column:biaoshi" json:"biaoshi" form:"biaoshi"`                   //订单号
	Gsfsp         int    `gorm:"column:gsfsp" json:"gsfsp" form:"gsfsp"`                         //光送粉丝牌-任务条件
	Dzcs          int    `gorm:"column:dzcs" json:"dzcs" form:"dzcs"`                            //点赞次数
	Sfsl          int    `gorm:"column:sfsl" json:"sfsl" form:"sfsl"`                            //任务权重 1=送礼物的任务 0=不送礼物的任务
	Sfgj          int    `gorm:"column:sfgj" json:"sfgj" form:"sfgj"`                            //0=不挂机任务 1=挂机任务
	Ksqz          int    `gorm:"column:ksqz" json:"ksqz" form:"ksqz"`                            //快手权重-可去掉
	Rwmoney       int    `gorm:"column:rwmoney" json:"rwmoney" form:"rwmoney"`                   //任务单价
	Rjbbh         int    `gorm:"column:rjbbh" json:"rjbbh" form:"rjbbh"`                         //任务限制-软件版本号
	Xtbbh         int    `gorm:"column:xtbbh" json:"xtbbh" form:"xtbbh"`                         //任务限制-系统版本号

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

// UpdateRenwu 修改
func UpdateRenwu(o *Renwu, tx ...*gorm.DB) (*Renwu, error) {
	db := global.MYSQL
	if len(tx) > 0 {
		db = tx[0]
	}
	err := db.Table("renwu").Where("Rid=?", o.Rid).Update(o).Error
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
func RenuwuShenyuGreaterZero() ([]*Renwu, error) {
	db := global.MYSQL
	res := make([]*Renwu, 0)
	err := db.Table("renwu").Where("shengyusl > 0 ").Find(&res).Error
	return res, err
}
