package service

import (
	"douyin/global"
	"douyin/web/db"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var manager = &RedisSyncToMysqlManager{
	update:    make(chan interface{}, 1024),
	create:    make(chan interface{}, 1024),
	renwuSet:  map[int]int{},
	yonghuSet: map[int]int{},

	lock: &sync.Mutex{},
}

type RedisSyncToMysqlManager struct {
	lock *sync.Mutex
	// 同步数据到mysql 任务和用户
	update chan interface{}

	// 新增 chan
	create chan interface{}

	// 任务编号是否被锁
	renwuSet map[int]int
	// 用户id集合
	yonghuSet map[int]int
}

func RunRedisSyncToMysqlManager() {

	fmt.Println("加载任务到redis...")
	manager.initRenwu()

	fmt.Println("加载用户到redis...")
	manager.initYonghu()

	fmt.Println("加载任务日志到redis...")
	manager.initRenwuLog()

	fmt.Println("加载任务ip日志到redis...")
	manager.initRenwuIPLog()

	manager.Run()
}
func (m *RedisSyncToMysqlManager) Run() {
	go m.create_method()
	go m.update_method()
}
func (m *RedisSyncToMysqlManager) addUpdate(v interface{}) {
	m.update <- v
}
func (m *RedisSyncToMysqlManager) addCreate(v interface{}) {
	m.create <- v
}

func (m *RedisSyncToMysqlManager) initYonghu() {

	// get all
	qu := &db.Yonghu{
		Page: db.Page{
			PageSize: 99999999,
		},
	}
	users, err := db.ListYonghu(qu)
	if err != nil {

		return
	}

	for _, u := range users {
		m.yonghuSet[u.Uid] = 1
	}

}

func (m *RedisSyncToMysqlManager) initRenwuLog() {
	conn := global.REDIS.Get()
	defer conn.Close()

	qu := db.Rwlogs{
		Page: db.Page{
			PageSize: 99999999,
		},
	}
	renwulogss, err := db.ListRwlogs(&qu)
	if err != nil {
		return
	}

	for _, lo := range renwulogss {
		rb, err := json.Marshal(lo)
		if err != nil {
			continue
		}
		_, err = conn.Do("set", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_LOG, lo.Userid, lo.Rid), string(rb))
		if err != nil {
			continue
		}
	}

}

func (m *RedisSyncToMysqlManager) initRenwuIPLog() {
	conn := global.REDIS.Get()
	defer conn.Close()

	qu := db.Iplogs{
		Page: db.Page{
			PageSize: 99999999,
		},
	}
	renwuiplogs, err := db.ListIplogs(&qu)
	if err != nil {
		return
	}

	for _, rip := range renwuiplogs {
		rb, err := json.Marshal(rip)
		if err != nil {
			continue
		}
		_, err = conn.Do("set", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_LOG_IP, rip.IP, rip.Rid), string(rb))
		if err != nil {
			continue
		}
	}

}

func (m *RedisSyncToMysqlManager) initRenwu() {
	conn := global.REDIS.Get()
	defer conn.Close()

	// shengyusl
	renwus, err := db.RenuwuShenyuGreaterZero()
	if err != nil {

	}
	for _, renwu := range renwus {
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

			if r, ok := data.(*db.Renwu); ok {
				// 同步任务
				db.UpdateRenwu(r)
			} else if u, ok := data.(*db.Yonghu); ok {
				// 同步数据到用户
				db.UpdateYonghu(u)
			} else if u, ok := data.(*db.Rwlogs); ok {
				// 同步数据到任务日志
				db.UpdateRwlogs(u)
			} else if u, ok := data.(*db.Iplogs); ok {
				// 同步数据到任务ip日志
				db.UpdateIplogs(u)
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
			} else if u, ok := data.(*db.Renwu); ok {
				// 新增用户
				db.AddRenwu(u)
			} else if u, ok := data.(*db.Iplogs); ok {
				// 新增ip日志
				db.AddIplogs(u)
			}

		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
