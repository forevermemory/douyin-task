package service

import (
	"douyin/global"
	"douyin/web/db"
	"fmt"
	"sync"
	"time"
)

var manager = &RedisSyncToMysqlManager{
	update:     make(chan interface{}, 1024),
	create:     make(chan interface{}, 1024),
	renwuLock:  make(map[int]int),
	renwuIDSet: make(map[int]*db.Renwu),
	yonghuSet:  make(map[int]int),

	lock: &sync.Mutex{},
}

type RedisSyncToMysqlManager struct {
	lock *sync.Mutex
	// 同步数据到mysql 任务和用户
	update chan interface{}

	// 新增 chan
	create chan interface{}

	// 任务编号是否被锁
	renwuLock map[int]int

	// 任务编号集合
	renwuIDSet map[int]*db.Renwu

	// 用户id集合
	yonghuSet map[int]int
}

func RunRedisSyncToMysqlManager() {

	fmt.Println("加载任务到redis...")
	manager.initRenwu()

	// 注册或者登陆之后再加入内存

	fmt.Println("加载用户到redis...")
	manager.initYonghu()

	// fmt.Println("加载任务日志到redis...")
	// manager.initRenwuLog()

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

	// 刚启动加载用户列表 只需要把onlie>0的加载就行了
	users, err := db.ListYonghuV3()
	if err != nil {
		return
	}
	for _, u := range users {
		m.yonghuSet[u.Uid] = 1
		m.setUser(u, 1)
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
		manager.setRenwulog(lo, 1)
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
		manager.setRenwu(renwu, 1)
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
				db.UpdateRenwu(r, r.UpdateType)
			} else if u, ok := data.(*db.Yonghu); ok {
				// 同步数据到用户
				db.UpdateYonghu(u, u.UpdateType)
			} else if u, ok := data.(*db.Rwlogs); ok {
				// 同步数据到任务日志
				db.UpdateRwlogs(u, u.UpdateType)
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
			}

		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
