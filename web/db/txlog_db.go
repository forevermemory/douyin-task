package db

import (
	"douyin/global"
	"time"

	"github.com/jinzhu/gorm"
)

type Txlogs struct {
	Id    int       `gorm:"column:id;primary_key;auto_increment;" json:"id" form:"id"`
	Uid   int       `gorm:"column:uid" json:"uid" form:"uid"`
	Money int       `gorm:"column:money" json:"money" form:"money"`
	Zfb   string    `gorm:"column:zfb" json:"zfb" form:"zfb"`
	Isadd int       `gorm:"column:isadd" json:"isadd" form:"isadd"` //0：未支付 1：已支付
	Day   time.Time `gorm:"column:day" json:"day" form:"day"`

	Page
}

// TableName 表名
func (o *Txlogs) TableName() string {
	return "txlogs"
}

// DeleteTxlogs 根据id删除
func DeleteTxlogs(id int) error {
	db := global.MYSQL
	return db.Table("txlogs").Where("id = ?", id).Error
}

// GetTxlogsByID 根据id查询一个
func GetTxlogsByID(id int) (*Txlogs, error) {
	db := global.MYSQL
	o := &Txlogs{}
	err := db.Table("txlogs").Where("id = ?", id).First(o).Error
	return o, err
}

// AddTxlogs 新增
func AddTxlogs(o *Txlogs, tx ...*gorm.DB) (*Txlogs, error) {
	db := global.MYSQL
	if len(tx) > 0 {
		db = tx[0]
	}
	err := db.Create(o).Error
	return o, err
}

// UpdateTxlogs 修改
func UpdateTxlogs(o *Txlogs, tx ...*gorm.DB) (*Txlogs, error) {
	db := global.MYSQL
	if len(tx) > 0 {
		db = tx[0]
	}
	err := db.Table("txlogs").Where("id=?", o.Id).Update(o).Error
	return o, err
}

// ListTxlogs 分页条件查询
func ListTxlogs(o *Txlogs) ([]*Txlogs, error) {
	db := global.MYSQL
	res := make([]*Txlogs, 0)
	err := db.Table("txlogs").Where(o).Offset((o.PageNo - 1) * o.PageSize).Limit(o.PageSize).Find(&res).Error
	return res, err
}

// CountTxlogs 条件数量
func CountTxlogs(o *Txlogs) (int, error) {
	db := global.MYSQL
	var count int
	err := db.Table("txlogs").Where(o).Count(&count).Error
	return count, err
}
