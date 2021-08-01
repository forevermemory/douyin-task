package service

import (
	"douyin/global"
	"douyin/utils"
	"douyin/web/db"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

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
	//
	delete(manager.renwuIDSet, renwu.Rid)
	return nil, nil
}
func (m *RedisSyncToMysqlManager) setRenwu(renwu *db.Renwu) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()
	// update
	uy, err := json.Marshal(renwu)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwu.Rid), string(uy))
	if err != nil {
		return nil, err
	}

	//
	if _, ok := m.renwuIDSet[renwu.Rid]; !ok {
		m.renwuIDSet[renwu.Rid] = 1
	}

	m.addUpdate(renwu)
	return nil, nil
}
func (m *RedisSyncToMysqlManager) getRenwu(renwuid int) (*db.Renwu, error) {
	conn := global.REDIS.Get()
	defer conn.Close()
	// renwu
	renwuStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwuid)))
	if err != nil {
		return nil, err
	}
	renwu := &db.Renwu{}
	if len(renwuStr) > 0 {
		err = json.Unmarshal([]byte(renwuStr), renwu)
		if err != nil {
			return nil, err
		}
		return renwu, nil
	}
	renwu, err = db.GetRenwuByID(renwuid)
	if err != nil {
		return nil, err
	}
	// update
	rb, err := json.Marshal(renwu)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_RENWU, renwuid), string(rb))
	if err != nil {
		return nil, err
	}

	return renwu, nil
}

/////////////user

func (m *RedisSyncToMysqlManager) setUser(user *db.Yonghu) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()
	// update
	uy, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid), string(uy))
	if err != nil {
		return nil, err
	}
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

	// user
	userStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, userid)))
	if err != nil {
		return nil, err
	}
	user := &db.Yonghu{}
	if len(userStr) > 0 {
		err = json.Unmarshal([]byte(userStr), user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	// 走mysql查询
	user, err = db.GetYonghuByID(userid)
	if err != nil {
		return nil, err
	}
	// update
	uy, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid), string(uy))
	if err != nil {
		return nil, err
	}
	return user, nil

}

///////////////rwlog

func (m *RedisSyncToMysqlManager) setRenwulog(rwlog *db.Rwlogs) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()
	// update
	renwlogb, err := json.Marshal(&rwlog)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_LOG, rwlog.Uid, rwlog.Rid), string(renwlogb))
	if err != nil {
		return nil, err
	}
	m.addUpdate(rwlog)
	return nil, nil
}
func (m *RedisSyncToMysqlManager) getRenwulog(userid int, renwuid int) (*db.Rwlogs, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	renwulogStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_LOG, userid, renwuid)))
	if err != nil {
		return nil, err
	}
	rwlog := &db.Rwlogs{}

	if len(renwulogStr) > 0 {
		err = json.Unmarshal([]byte(renwulogStr), &rwlog)
		if err != nil {
			return nil, err
		}
		return rwlog, nil
	}
	// 走mysql查询
	rwlog, err = db.GetRwlogsByruandyonghuid(userid, renwuid)
	if err != nil {
		return nil, err
	}
	// update
	renwlogb, err := json.Marshal(&rwlog)
	if err != nil {
		return nil, err
	}
	_, err = conn.Do("set", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_LOG, userid, renwuid), string(renwlogb))
	if err != nil {
		return nil, err
	}
	return rwlog, nil

}
