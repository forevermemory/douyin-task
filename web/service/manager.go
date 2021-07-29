package service

import (
	"douyin/global"
	"douyin/web/db"
	"encoding/json"
	"fmt"
	"time"
)

var manager = &RedisSyncToMysqlManager{
	update:   make(chan interface{}, 1024),
	create:   make(chan interface{}, 1024),
	renwuSet: map[int]int{},
}

type RedisSyncToMysqlManager struct {
	// 同步数据到mysql 任务和用户
	update chan interface{}

	// 新增 chan
	create chan interface{}

	// 任务编号是否被锁
	renwuSet map[int]int
}

func RunRedisSyncToMysqlManager() {
	manager.Run()
}

func (m *RedisSyncToMysqlManager) initRenwu() {
	conn := global.REDIS.Get()
	defer conn.Close()

	// shengyusl
	renwus, err := db.RenuwuShenyuGreaterZero()
	if err != nil {

	}

	for _, renwu := range renwus {

		res, err := conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid))
		if err != nil {
			continue
		}
		if res != nil {
			continue
		}

		//
		rb, err := json.Marshal(renwu)
		if err != nil {
			continue
		}
		_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid), string(rb))
		if err != nil {
			continue
		}
	}

}
func (m *RedisSyncToMysqlManager) update_method() {
	for {
		select {
		case data, ok := <-m.update:
			if !ok {
				continue
			}
			fmt.Println("update::", data)

			// 获取cpu负载情况  负载高sleep久一点
			// 不会并发对mysql进行写入
			time.Sleep(time.Millisecond * 100)

			if r, ok := data.(*db.Renwu); ok {
				// 同步任务
				db.UpdateRenwu(r)
			} else if u, ok := data.(*db.Yonghu); ok {
				// 同步数据到用户
				db.UpdateYonghu(u)
			}

		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
func (m *RedisSyncToMysqlManager) create_method() {

	for {
		select {
		case data, ok := <-m.create:
			if !ok {
				continue
			}
			fmt.Println("update::", data)

			if r, ok := data.(*db.Rwlogs); ok {
				// 新增任务日志
				db.AddRwlogs(r)
			}
			// else if u, ok := data.(*db.Yonghu); ok {
			// 	// 同步数据到用户
			// 	db.UpdateYonghu(u)
			// }

		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
func (m *RedisSyncToMysqlManager) Run() {

	m.initRenwu()

	// 怕挂掉丢数据 可以先存到redis 在慢慢从redis取出来

	go m.create_method()
	go m.update_method()
}
