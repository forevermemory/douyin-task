package db

import (
	"douyin/global"
	"time"

	"github.com/jinzhu/gorm"
)

type YonghuRequest struct {
	ID       int    `json:"ID"`
	User     string `json:"user"`
	Password string `json:"password"`
	Zfb      string `json:"zfb"`
	Zfbname  string `json:"zfbname"`
	Sj       string `json:"sj"`

	Token string `json:"token"`
	Dyid  string `json:"dyid"`

	Shangjiuid int ` json:"shangjiUID" form:"shangjiUID"` //上级uid

	Registerip string
}

type YonghuResponse struct {
	Msg           int    `json:"msg"`
	Code          int    `json:"code"`
	Guishu        string `json:"guishu"`
	Money         int    `json:"money"`
	LastLoginTime string `json:"LastLoginTime"`
	LastLoginIP   string `json:"LastLoginIP"`
	Token         string `json:"token"`
	Userid        int    `json:"userid"`
}

type Yonghu struct {
	Uid int `gorm:"column:UID;primary_key;auto_increment;" json:"UID" redis:"UID"`

	Account          string    `gorm:"column:Account" json:"Account" redis:"Account"` //账号
	Accountmd5       string    `gorm:"column:AccountMD5" json:"AccountMD5" redis:"AccountMD5"`
	Password         string    `gorm:"column:Password" json:"Password" redis:"Password"`
	Passwordmd5      string    `gorm:"column:PasswordMD5" json:"PasswordMD5" redis:"PasswordMD5"`
	Onlie            int       `gorm:"column:Onlie" json:"Onlie" redis:"Onlie"`                                  //websoket句柄
	Lastlogintime    time.Time `gorm:"column:LastLoginTime" json:"LastLoginTime" redis:"LastLoginTime"`          //最后登录时间
	Lastloginip      string    `gorm:"column:LastLoginIP" json:"LastLoginIP" redis:"LastLoginIP"`                //最后登录IP
	Registertime     time.Time `gorm:"column:RegisterTime" json:"RegisterTime" redis:"RegisterTime"`             //注册时间
	Registerip       string    `gorm:"column:RegisterIP" json:"RegisterIP" redis:"RegisterIP"`                   //注册IP
	State            int       `gorm:"column:State" json:"State" redis:"State"`                                  //0=正常 1=冻结
	Stateinformation string    `gorm:"column:Stateinformation" json:"Stateinformation" redis:"Stateinformation"` //状态原因
	Token            string    `gorm:"column:Token" json:"Token" redis:"Token"`
	Tokentime        int       `gorm:"column:TokenTime" json:"TokenTime" redis:"TokenTime"`
	Guishu           string    `gorm:"column:guishu" json:"guishu" redis:"guishu"`             //归属哪个用户
	Guishuuid        int       `gorm:"column:guishuUID" json:"guishuUID" redis:"guishuUID"`    //归属用户uid
	Shangjiuid       int       `gorm:"column:shangjiUID" json:"shangjiUID" redis:"shangjiUID"` //上级uid
	Rid              int       `gorm:"column:RID" json:"RID" redis:"RID"`                      //任务ID
	Rwkstime         int       `gorm:"column:rwksTime" json:"rwksTime" redis:"rwksTime"`       //任务开始时间
	Rwjd             int       `gorm:"column:rwjd" json:"rwjd" redis:"rwjd"`                   //1=正在进直播间 2=正在送礼物 3=正在挂机
	Dymz             string    `gorm:"column:dymz" json:"dymz" redis:"dymz"`                   //可去掉
	Dyid             string    `gorm:"column:dyid" json:"dyid" redis:"dyid"`                   //抖音id
	Dbye             int       `gorm:"column:dbye" json:"dbye" redis:"dbye"`                   //抖币余额
	Dyyz             int       `gorm:"column:dyyz" json:"dyyz" redis:"dyyz"`                   //抖音验证状态0-1
	Ksyz             int       `gorm:"column:ksyz" json:"ksyz" redis:"ksyz"`                   //可去掉
	Money            int       `gorm:"column:money" json:"money" redis:"money"`
	Zfb              string    `gorm:"column:zfb" json:"zfb" redis:"zfb"`             //支付宝
	Zfbname          string    `gorm:"column:zfbname" json:"zfbname" redis:"zfbname"` //支付宝姓名
	Xtbbh            int       `gorm:"column:xtbbh" json:"xtbbh" redis:"xtbbh"`       //系统版本号
	Rjbbh            int       `gorm:"column:rjbbh" json:"rjbbh" redis:"rjbbh"`       //软件版本号

	Cfdj int `gorm:"column:cfdj" json:"cfdj" redis:"cfdj"` // 抖音等级

	UpdateType int `gorm:"-" json:"-" redis:"-"`

	Page
}

// TableName 表名
func (o *Yonghu) TableName() string {
	return "yonghu"
}

// CheckDouyinIDRepeat CheckDouyinIDRepeat
func CheckDouyinIDRepeat(dyid string) ([]map[string]interface{}, error) {

	// sql := `
	// SELECT dyid,count(dyid) num ,GROUP_CONCAT(Account) account  from yonghu
	// where dyid = ?
	// GROUP BY dyid

	// `

	sql := "SELECT `account` FROM `yonghu` WHERE `dyid`= ?"
	rows, err := global.MYSQL.Raw(sql, dyid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := SQLMap(rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func LoginUser(user, password string) (*Yonghu, error) {
	o := &Yonghu{}
	err := global.MYSQL.Table("yonghu").Where("Account = ? and Password = ? ", user, password).First(o).Error
	if err != nil {
		return nil, err
	}
	return o, nil
}

// AddYonghu 新增
func AddYonghu(o *Yonghu, tx ...*gorm.DB) (*Yonghu, error) {
	db := global.MYSQL
	if len(tx) > 0 {
		db = tx[0]
	}
	err := db.Create(o).Error
	return o, err
}

const (
	USER_UPDATE_MONEY            int = 1
	USER_UPDATE_DOWN_4           int = 2
	USER_UPDATE_TOP2             int = 3
	USER_UPDATE_TOP5             int = 4
	USER_UPDATE_ONLY_RID         int = 5
	USER_UPDATE_ONLY_RWID        int = 6
	USER_UPDATE_ONLY_RWID_RWKSSJ int = 7
	USER_UPDATE_ONLY_RIDAND_RWID int = 8
)

// UpdateYonghu 修改
func UpdateYonghu(o *Yonghu, _type int) (*Yonghu, error) {
	db := global.MYSQL

	u2 := new(Yonghu)

	switch _type {
	case USER_UPDATE_MONEY:
		u2.Money = o.Money
	case USER_UPDATE_DOWN_4:
		u2.Dyid = o.Dyid
		u2.Dbye = o.Dbye
		u2.Ksyz = o.Ksyz
		u2.Dyyz = o.Dyyz
		u2.Xtbbh = o.Xtbbh
		u2.Cfdj = o.Cfdj
	case USER_UPDATE_TOP2:
		u2.Token = o.Token
		u2.Lastloginip = o.Lastloginip
		u2.Lastlogintime = o.Lastlogintime
	case USER_UPDATE_TOP5:
		u2.Lastloginip = o.Lastloginip
		u2.Lastlogintime = o.Lastlogintime
	case USER_UPDATE_ONLY_RID:
		u2.Rid = o.Rid
	case USER_UPDATE_ONLY_RWID:
		u2.Rwjd = o.Rwjd
	case USER_UPDATE_ONLY_RWID_RWKSSJ:
		u2.Rwjd = o.Rwjd
		u2.Rwkstime = o.Rwkstime
	case USER_UPDATE_ONLY_RIDAND_RWID:
		u2.Rwjd = o.Rwjd
		u2.Rid = o.Rid

	default:
		return nil, nil
	}

	err := db.Table("yonghu").Where("UID=?", o.Uid).Updates(u2).Error
	if err != nil {
		return nil, err
	}
	return o, err
}

///////////////////

// DeleteYonghu 根据id删除
func DeleteYonghu(id int) error {
	db := global.MYSQL
	return db.Table("yonghu").Where("UID = ?", id).Error
}

// GetYonghuByID 根据id查询一个
func GetYonghuByID(id int) (*Yonghu, error) {
	db := global.MYSQL
	o := &Yonghu{}
	err := db.Table("yonghu").Where("UID = ?", id).First(o).Error
	return o, err
}

// ListYonghu 分页条件查询
func ListYonghu(o *Yonghu) ([]*Yonghu, error) {
	db := global.MYSQL
	res := make([]*Yonghu, 0)
	err := db.Table("yonghu").Where(o).Offset((o.PageNo - 1) * o.PageSize).Limit(o.PageSize).Find(&res).Error
	return res, err
}

// ListYonghuV2
func ListYonghuV2() ([]*Yonghu, error) {
	db := global.MYSQL
	res := make([]*Yonghu, 0)

	// 	// `rid`>'0' and `rwjd`=1 and 'rwksTime`>'300'和
	// `rid`>'0' and `rwjd`=2 and '`rwksTime`>'480'

	now := time.Now().Unix()
	sql := `select * from  yonghu  where (rid>0 and rwjd=1 and (rwksTime - ?) > 300 ) or  (rid>0 and rwjd=2 and (rwksTime - ?) > 480) `
	err := db.Raw(sql, now, now).Scan(&res).Error
	return res, err
}

// ListYonghuV3
func ListYonghuV3() ([]*Yonghu, error) {
	db := global.MYSQL
	res := make([]*Yonghu, 0)

	// 	// `rid`>'0' and `rwjd`=1 and 'rwksTime`>'300'和
	// `rid`>'0' and `rwjd`=2 and '`rwksTime`>'480'

	sql := `select * from  yonghu  where  Onlie>0`
	err := db.Raw(sql).Scan(&res).Error
	return res, err
}

// CountYonghu 条件数量
func CountYonghu(o *Yonghu) (int, error) {
	db := global.MYSQL
	var count int
	err := db.Table("yonghu").Where(o).Count(&count).Error
	return count, err
}

///////
// ListYonghuBySuper ListYonghuBySuper
func ListYonghuBySuper(super string) ([]*Yonghu, error) {
	db := global.MYSQL
	res := make([]*Yonghu, 0)
	err := db.Table("yonghu").Where("guishu = ?", super).Find(&res).Error
	return res, err
}
