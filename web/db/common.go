package db

import "database/sql"

// DataStore list的结构体
type DataStore struct {
	Total     int         `json:"total"`
	TotalPage int         `json:"totalPage"`
	Data      interface{} `json:"data"`
}

// Page 分页参数
type Page struct {
	PageNo   int `gorm:"-" json:"page,default=1,omitempty" form:"page,default=1"`
	PageSize int `gorm:"-" json:"page_size,default=10,omitempty" form:"page_size,default=10"`
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
