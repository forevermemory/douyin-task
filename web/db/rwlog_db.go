package db

import (
	"douyin/global"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	Rwlogs_isadd_GET_TASK                      int = 1
	Rwlogs_isadd_FINISH_GUAJI                  int = 2
	Rwlogs_isadd_FINISH_SEND_PRESENT           int = 3
	Rwlogs_isadd_ABADON_TASK_NOT_IN            int = 4
	Rwlogs_isadd_ABADON_TASK_SEND_PRESENT_FAIL int = 5
	Rwlogs_isadd_ABADON_TASK_ZB_OFFLINE        int = 6
	Rwlogs_isadd_ABADON_TASK_OFF_LINE          int = 7
	Rwlogs_isadd_ABADON_TASK_LEAVE             int = 8
	Rwlogs_isadd_ALLREADY_DONE                 int = 9
)

type Rwlogs struct {
	Id     int    `gorm:"column:id;primary_key;auto_increment;" json:"id" form:"id"`
	Uid    int    `gorm:"column:uid" json:"uid" form:"uid"`
	Rid    int    `gorm:"column:rid" json:"rid" form:"rid"`
	Userid int    `gorm:"column:userid" json:"userid" form:"userid"`
	Zbid   string `gorm:"column:zbid" json:"zbid" form:"zbid"`

	// 1: 领取任务
	// 2: 完成挂机任务
	// 3: 完成送礼任务
	// 4: 放弃任务，进不去
	// 5: 放弃任务，送礼失败
	// 6: 放弃任务，没开播
	// 7: 放弃任务，不在线
	// 8: 放弃任务，中途离开直播间
	// 9: 已经做过
	Isadd int `gorm:"column:isadd" json:"isadd" form:"isadd"`

	Day time.Time `gorm:"column:day" json:"day" form:"day"`
	Page
}

// TableName 表名
func (o *Rwlogs) TableName() string {
	return "rwlogs"
}

// DeleteRwlogs 根据id删除
func DeleteRwlogs(id int) error {
	db := global.MYSQL
	return db.Table("rwlogs").Where("id = ?", id).Error
}

// GetRwlogsByID 根据id查询一个
func GetRwlogsByID(id int) (*Rwlogs, error) {
	db := global.MYSQL
	o := &Rwlogs{}
	err := db.Table("rwlogs").Where("id = ?", id).First(o).Error
	return o, err
}

// GetRwlogsByruandyonghuid 根据wrid wid查询一个
func GetRwlogsByruandyonghuid(uid, rid int) (*Rwlogs, error) {
	db := global.MYSQL
	o := &Rwlogs{}
	err := db.Table("rwlogs").Where("uid = ? and rid = ? ", uid, rid).First(o).Error
	return o, err
}

// AddRwlogs 新增
func AddRwlogs(o *Rwlogs, tx ...*gorm.DB) (*Rwlogs, error) {
	db := global.MYSQL
	if len(tx) > 0 {
		db = tx[0]
	}
	err := db.Create(o).Error
	return o, err
}

// UpdateRwlogs 修改
func UpdateRwlogs(o *Rwlogs, tx ...*gorm.DB) (*Rwlogs, error) {
	db := global.MYSQL
	if len(tx) > 0 {
		db = tx[0]
	}
	err := db.Table("rwlogs").Where("id=?", o.Id).Update(o).Error
	return o, err
}

// ListRwlogs 分页条件查询
func ListRwlogs(o *Rwlogs) ([]*Rwlogs, error) {
	db := global.MYSQL
	res := make([]*Rwlogs, 0)
	err := db.Table("rwlogs").Where(o).Offset((o.PageNo - 1) * o.PageSize).Limit(o.PageSize).Find(&res).Error
	return res, err
}

// CountRwlogs 条件数量
func CountRwlogs(o *Rwlogs) (int, error) {
	db := global.MYSQL
	var count int
	err := db.Table("rwlogs").Where(o).Count(&count).Error
	return count, err
}
