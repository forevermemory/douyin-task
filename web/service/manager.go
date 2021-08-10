package service

import (
	"context"
	"douyin/global"
	"douyin/utils"
	"douyin/web/db"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var manager = &RedisSyncToMysqlManager{
	update: make(chan interface{}, 1024),
	create: make(chan interface{}, 1024),

	rwlock: &sync.Mutex{},
}

type RedisSyncToMysqlManager struct {
	rwlock *sync.Mutex
	// 同步数据到mysql 任务和用户
	update chan interface{}

	// 新增 chan
	create chan interface{}
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

// getRenwuLockWithTimeout 循环获取 直到超时
func (m *RedisSyncToMysqlManager) getRenwuLockWithTimeout(rid int, timeout time.Duration) (bool, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU_LOCK, rid)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

TIMEOUT:
	for {
		select {
		case <-ctx.Done():
			break TIMEOUT
		default:
			break
		}
		res, err := redis.Int64(conn.Do("setnx", _key, 1))
		if err != nil {
			return false, err
		}

		if res == 1 {
			// get lock
			// expire
			_, err = conn.Do("expire", _key, global.MAX_GET_LOCK_TIMEOUT)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}

	// 超时了还没有获得锁
	return false, nil
}

// getRenwuLock 尝试一次获取锁
func (m *RedisSyncToMysqlManager) getRenwuLock(rid int) (bool, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU_LOCK, rid)

	res, err := redis.Int64(conn.Do("setnx", _key, 1))
	if err != nil {
		return false, err
	}

	if res == 1 {
		// expire
		_, err = conn.Do("expire", _key, global.MAX_GET_LOCK_TIMEOUT)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (m *RedisSyncToMysqlManager) delRenwuLock(rid int) error {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU_LOCK, rid)
	_, err := conn.Do("del", _key)
	return err
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

////////////// 用户集合
func (m *RedisSyncToMysqlManager) getYonghuSet() ([]string, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	ids, err := redis.Strings(conn.Do("SMEMBERS", global.REDIS_PREFIX_USERS))
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (m *RedisSyncToMysqlManager) setYonghuSet(id int) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_, err := conn.Do("sadd", global.REDIS_PREFIX_USERS, id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *RedisSyncToMysqlManager) delYonghuSet(id int) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	res, err := redis.Int64(conn.Do("srem", global.REDIS_PREFIX_USERS, id))
	if err != nil {
		return nil, err
	}
	if res == 0 {
		// 删除失败 或者重复删除 无所谓
	}
	return nil, nil
}

////////////// 任务集合
func (m *RedisSyncToMysqlManager) getRenwuSet() ([]string, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	renwuIDs, err := redis.Strings(conn.Do("SMEMBERS", global.REDIS_PREFIX_RENWUS))
	if err != nil {
		return nil, err
	}

	return renwuIDs, nil
}

func (m *RedisSyncToMysqlManager) setRenwuSet(id int) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_, err := conn.Do("sadd", global.REDIS_PREFIX_RENWUS, id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *RedisSyncToMysqlManager) delRenwuSet(id int) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	res, err := redis.Int64(conn.Do("srem", global.REDIS_PREFIX_RENWUS, id))
	if err != nil {
		return nil, err
	}
	if res == 0 {
		// 删除失败 或者重复删除 无所谓
	}
	return nil, nil
}

//////////////ip

func (m *RedisSyncToMysqlManager) getIpLimit(ip string, rid int) (int, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	limit, err := redis.Int(conn.Do("get", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_IP, ip, rid)))
	if err != nil {
		return 0, err
	}
	return limit, nil
}

func (m *RedisSyncToMysqlManager) setIpLimit(ip string, rid int) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_IP, ip, rid)
	limit, err := redis.Int(conn.Do("get", _key))
	if err != nil {
		return 0, err
	}
	limit += 1
	_, err = conn.Do("set", _key, limit)
	if err != nil {
		return nil, err
	}
	return limit, nil
}

////////////////renwu
func (m *RedisSyncToMysqlManager) delRenwu(renwu *db.Renwu) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_, err := conn.Do("del", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid))
	if err != nil {
		return nil, err
	}

	_, err = m.delRenwuSet(renwu.Rid)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
func (m *RedisSyncToMysqlManager) setRenwu(renwu *db.Renwu) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid)

	switch renwu.UpdateType {
	case db.RENWU_UPDATE_ALL:
		m.setRenwu_hsetall(renwu)
	case db.RENWU_UPDATE_Shengyusl:
		conn.Do("hset", _key, "Shengyusl", renwu.Shengyusl)
	case db.RENWU_UPDATE_STOP_Tiqianjieshu:
		conn.Do("hset", _key, "Tiqianjieshu", renwu.Tiqianjieshu)
		conn.Do("hset", _key, "Stop", renwu.Stop)
	default:
		return nil, nil
	}

	// 更新到redis
	m.setRenwuSet(renwu.Rid)

	m.addUpdate(renwu)
	return nil, nil
}
func (m *RedisSyncToMysqlManager) getRenwu(renwuid int) (*db.Renwu, error) {

	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwuid)
	// user
	renwu := &db.Renwu{}

	_, err := conn.Do("hget", _key, "Rid")
	if err != nil {
		if !errors.Is(err, redis.ErrNil) {
			return nil, err
		}
		// 走mysql查询
		renwu, err = db.GetRenwuByID(renwuid)
		if err != nil {
			return nil, err
		}
		m.setRenwu_hsetall(renwu)

	} else {
		renwu, err = m.getRenwu_hgetall(renwuid)
		if err != nil {
			return nil, err
		}

	}
	// 更新到redis
	m.setRenwuSet(renwu.Rid)

	return renwu, nil
}

/////////////user

func (m *RedisSyncToMysqlManager) setUser(user *db.Yonghu) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid)

	// add lock

	switch user.UpdateType {
	case db.USER_UPDATE_MONEY:
		conn.Do("hset", _key, "Money", user.Money)
	case db.USER_UPDATE_DOWN_4:
		conn.Do("hset", _key, "Dyid", user.Dyid)
		conn.Do("hset", _key, "Dbye", user.Dbye)
		conn.Do("hset", _key, "Ksyz", user.Ksyz)
		conn.Do("hset", _key, "Dyyz", user.Dyyz)
		conn.Do("hset", _key, "Xtbbh", user.Xtbbh)
		conn.Do("hset", _key, "Cfdj", user.Cfdj)
	case db.USER_UPDATE_TOP2:
		conn.Do("hset", _key, "Token", user.Token)
		conn.Do("hset", _key, "Lastloginip", user.Lastloginip)
		conn.Do("hset", _key, "Lastlogintime", m.stringfyTime(user.Lastlogintime))
	case db.USER_UPDATE_TOP5:
		conn.Do("hset", _key, "Lastlogintime", m.stringfyTime(user.Lastlogintime))
		conn.Do("hset", _key, "Lastloginip", user.Lastloginip)
	case db.USER_UPDATE_ONLY_RID:
		conn.Do("hset", _key, "Rid", user.Rid)
	case db.USER_UPDATE_ONLY_RWID:
		conn.Do("hset", _key, "Rwjd", user.Rwjd)
	case db.USER_UPDATE_ONLY_RWID_RWKSSJ:
		conn.Do("hset", _key, "Rwjd", user.Rwjd)
		conn.Do("hset", _key, "Rwkstime", user.Rwkstime)
	case db.USER_UPDATE_ONLY_RIDAND_RWID:
		conn.Do("hset", _key, "Rwjd", user.Rwjd)
		conn.Do("hset", _key, "Rid", user.Rid)

	default:
		return nil, nil
	}

	// 更新到redis
	m.setYonghuSet(user.Rid)

	m.addUpdate(user)
	return nil, nil
}

func (m *RedisSyncToMysqlManager) getUserByToken(token string) (*db.Yonghu, error) {
	userid := utils.TokenDecrypt(token)
	return m.getUser(userid)
}

func (m *RedisSyncToMysqlManager) getUser(userid int) (*db.Yonghu, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, userid)

	// user
	user := &db.Yonghu{}

	_, err := conn.Do("hget", _key, "Uid")
	if err != nil {
		if !errors.Is(err, redis.ErrNil) {
			return nil, err
		}

		// set
		// 走mysql查询
		user, err = db.GetYonghuByID(userid)
		if err != nil {
			return nil, err
		}
		m.setUser_hsetall(user)
	} else {
		user, err = m.getUser_hgetall(userid)
	}

	return user, nil

}

///////////////rwlog

func (m *RedisSyncToMysqlManager) setRenwulog(rwlog *db.Rwlogs) (interface{}, error) {
	m.addUpdate(rwlog)
	return nil, nil
}

func (m *RedisSyncToMysqlManager) getRenwulog(userid int, renwuid int) (*db.Rwlogs, error) {
	// 走mysql查询
	rwlog, err := db.GetRwlogsByruandyonghuid(userid, renwuid)
	if err != nil {
		return nil, err
	}

	return rwlog, nil

}
