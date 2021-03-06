package db

import "database/sql"

/////////////////READMD

// 通过结构体变量更新字段值, gorm库会忽略零值字段。就是字段值等于0, nil, "", false这些值会被忽略掉，不会更新。如果想更新零值，可以使用map类型替代结构体

// DataStore list的结构体
type DataStore struct {
	Total     int         `json:"total"`
	TotalPage int         `json:"totalPage"`
	Data      interface{} `json:"data"`
}

// Page 分页参数
type Page struct {
	PageNo   int `gorm:"-" json:"page,default=1,omitempty" form:"page,default=1" redis:"-"`
	PageSize int `gorm:"-" json:"page_size,default=10,omitempty" form:"page_size,default=10" redis:"-"`
}

func SQLMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	rt := make([]map[string]interface{}, 0)
	columns, _ := rows.Columns()
	value := make([]interface{}, len(columns))
	valuePrt := make([]interface{}, len(columns))
	for rows.Next() {
		for i := range columns {
			valuePrt[i] = &value[i]
		}
		err := rows.Scan(valuePrt...)
		if err != nil {
			return nil, err
		}
		mp := map[string]interface{}{}
		for i, e := range columns {
			var v interface{}
			if co, ok := value[i].([]byte); ok {
				v = string(co)
			} else {
				v = value[i]
			}
			mp[e] = v
		}
		rt = append(rt, mp)
	}
	return rt, nil
}
