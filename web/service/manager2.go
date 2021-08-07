package service

import (
	"douyin/global"
	"douyin/utils"
	"douyin/web/db"
	"errors"
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
	delete(manager.renwuIDSet, renwu.Rid)
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

	// 更新到内存
	m.renwuIDSet[renwu.Rid] = renwu

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
		// set
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
	// 更新到内存
	m.renwuIDSet[renwu.Rid] = renwu

	return renwu, nil
}

/////////////user

func (m *RedisSyncToMysqlManager) setUser(user *db.Yonghu) (interface{}, error) {
	conn := global.REDIS.Get()
	defer conn.Close()

	_key := fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid)

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
