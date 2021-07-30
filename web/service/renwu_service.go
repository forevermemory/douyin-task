package service

import (
	"douyin/web/db"
)

//////////////////////////////////////////////////

// AddRenwu add
func AddRenwu(req *db.Renwu) (*db.Renwu, error) {
	return db.AddRenwu(req)
}

// UpdateRenwu update
func UpdateRenwu(req *db.Renwu) (*db.Renwu, error) {
	return db.UpdateRenwu(req)
}

// GetRenwuByID get by id
func GetRenwuByID(id int) (*db.Renwu, error) {
	return db.GetRenwuByID(id)
}

// ListRenwu  page by condition
func ListRenwu(req *db.Renwu) (*db.DataStore, error) {
	list, err := db.ListRenwu(req)
	if err != nil {
		return nil, err
	}
	total, err := db.CountRenwu(req)
	if err != nil {
		return nil, err
	}
	return &db.DataStore{Total: total, Data: list, TotalPage: (int(total) + req.PageSize - 1) / req.PageSize}, nil
}

// DeleteRenwu delete
func DeleteRenwu(id int) error {
	return db.DeleteRenwu(id)
}
