package db

import (
	"douyin/global"
	"time"

	"github.com/jinzhu/gorm"
)

type Iplogs struct {
	Id     int    `gorm:"column:id;primary_key;auto_increment;" json:"id" form:"id"`
	Uid    int    `gorm:"column:uid" json:"uid" form:"uid"`
	Rid    int    `gorm:"column:rid" json:"rid" form:"rid"`
	Userid int    `gorm:"column:userid" json:"userid" form:"userid"`
	IP     string `gorm:"column:ip" json:"ip" form:"ip"`

	Times int `gorm:"column:times" json:"times" form:"times"`

	Day time.Time `gorm:"column:day" json:"day" form:"day"`
	Page
}

// TableName 表名
func (o *Iplogs) TableName() string {
	return "iplogs"
}

// DeleteIplogs 根据id删除
func DeleteIplogs(id int) error {
	db := global.MYSQL
	return db.Table("iplogs").Where("id = ?", id).Error
}

// GetIplogsByID 根据id查询一个
func GetIplogsByID(id int) (*Iplogs, error) {
	db := global.MYSQL
	o := &Iplogs{}
	err := db.Table("iplogs").Where("id = ?", id).First(o).Error
	return o, err
}

// AddIplogs 新增
func AddIplogs(o *Iplogs, tx ...*gorm.DB) (*Iplogs, error) {
	db := global.MYSQL
	if len(tx) > 0 {
		db = tx[0]
	}
	err := db.Create(o).Error
	return o, err
}

// UpdateIplogs 修改
func UpdateIplogs(o *Iplogs, tx ...*gorm.DB) (*Iplogs, error) {
	db := global.MYSQL
	if len(tx) > 0 {
		db = tx[0]
	}
	err := db.Table("iplogs").Where("id=?", o.Id).Update(o).Error
	return o, err
}

// ListIplogs 分页条件查询
func ListIplogs(o *Iplogs) ([]*Iplogs, error) {
	db := global.MYSQL
	res := make([]*Iplogs, 0)
	err := db.Table("iplogs").Where(o).Offset((o.PageNo - 1) * o.PageSize).Limit(o.PageSize).Find(&res).Error
	return res, err
}

// CountIplogs 条件数量
func CountIplogs(o *Iplogs) (int, error) {
	db := global.MYSQL
	var count int
	err := db.Table("iplogs").Where(o).Count(&count).Error
	return count, err
}
