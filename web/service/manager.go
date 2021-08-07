package service

import (
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

	rwlock: &sync.Mutex{},
}

type RedisSyncToMysqlManager struct {
	rwlock *sync.Mutex
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

func (m *RedisSyncToMysqlManager) getRenwuLock(rid int) bool {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	if _, ok := m.renwuLock[rid]; ok {
		return true
	}
	manager.renwuLock[rid] = 1
	return false
}
func (m *RedisSyncToMysqlManager) delRenwuLock(rid int) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	delete(m.renwuLock, rid)
}
func (m *RedisSyncToMysqlManager) parseTime(t string) time.Time {
	v, _ := time.Parse("2006-01-02 15:04:05", t)
	return v
}
func (m *RedisSyncToMysqlManager) stringfyTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
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
