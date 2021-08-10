package service

import (
	"douyin/global"
	"douyin/web/db"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func (m *RedisSyncToMysqlManager) initRenwu() {
	conn := global.REDIS.Get()
	defer conn.Close()

	// shengyusl
	renwus, err := db.RenuwuShenyuGreaterZero()
	if err != nil {
		return
	}

	var rs = make([]interface{}, 0)
	rs = append(rs, global.REDIS_PREFIX_RENWUS)

	// 任务对象放到redis
	for _, renwu := range renwus {
		manager.setRenwu_hsetall(renwu)
		rs = append(rs, renwu.Rid)
	}
	conn.Do("sadd", rs...)
}

func (m *RedisSyncToMysqlManager) getRenwu_hgetall(renwuid int) (*db.Renwu, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwuid)

	res, err := redis.Values(conn.Do("hgetall", _key))
	if err != nil {
		return nil, err
	}
	u := new(db.Renwu)
	err = redis.ScanStruct(res, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *RedisSyncToMysqlManager) setRenwu_hsetall(renwu *db.Renwu) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid)

	_, err := conn.Do("hset", redis.Args{}.Add(_key).AddFlat(renwu)...)

	return nil, err
}

func (m *RedisSyncToMysqlManager) initYonghu() {
	conn := global.REDIS.Get()
	defer conn.Close()

	// 刚启动加载用户列表 只需要把onlie>0的加载就行了
	users, err := db.ListYonghuV3()
	if err != nil {
		return
	}

	var rs = make([]interface{}, 0)
	rs = append(rs, global.REDIS_PREFIX_USERS)

	for _, u := range users {

		m.setUser_hsetall(u)
		rs = append(rs, u.Uid)
	}

	conn.Do("sadd", rs...)
}

func (m *RedisSyncToMysqlManager) setUser_hsetall(user *db.Yonghu) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid)

	_, err := conn.Do("hset", redis.Args{}.Add(_key).AddFlat(user)...)
	if err != nil {
		return
	}
	conn.Do("hmset", _key, "Lastlogintime", m.stringfyTime(user.Lastlogintime), "Registertime", m.stringfyTime(user.Registertime))

}

func (m *RedisSyncToMysqlManager) getUser_hgetall(id int) (*db.Yonghu, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, id)

	res, err := redis.Values(conn.Do("hgetall", _key))
	if err != nil {
		return nil, err
	}
	u := new(db.Yonghu)
	// Integer, float, boolean, string and []byte fields are supported
	err = redis.ScanStruct(res, u)
	if err != nil {
		return nil, err
	}

	// time.Time
	_Lastlogintime, err := redis.String(conn.Do("hget", _key, "Lastlogintime"))
	_Registertime, err := redis.String(conn.Do("hget", _key, "Registertime"))
	if err != nil {
		return nil, err
	}
	u.Lastlogintime = m.parseTime(_Lastlogintime)
	u.Registertime = m.parseTime(_Registertime)

	return u, nil
}
