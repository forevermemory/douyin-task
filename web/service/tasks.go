package service

import (
	"douyin/global"
	"douyin/web/db"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/robfig/cron"
)

func RunCronTasks() {

	c := cron.New()
	// every 5s
	c.AddFunc("*/5 * * * * *", func() {
		fmt.Println("task run...")
		go run()
	})
	c.Start()
	select {}
}

func run() {
	t := &Task{}
	t.Run()
}

type Task struct {
}

func (t *Task) Run() (interface{}, error) {
	// 服务器会定时把在任务进度1也就是刚领了任务超过5分钟的
	// 和任务进度2 也就是领了任务在送礼物  超过8分钟的
	// 用户取消任务  因为超时了

	conn := global.REDIS.Get()
	defer conn.Close()

	for userid := range manager.yonghuSet {

		// user
		userStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, userid)))
		if err != nil {
			return nil, err
		}
		user := db.Yonghu{}
		err = json.Unmarshal([]byte(userStr), &user)
		if err != nil {
			return nil, err
		}

		// 领了任务超过5分钟的
		// 还有就是用户5分钟内做这个任务失败了 下次就不让他领取这个任务
		if user.Rwjd == 2 {
			if int(time.Now().Unix())-user.Rwkstime > 60*5 {
				////////////////////
				user.Rwjd = -1
				user.Rid = -1
				// 更新到任务log

				// renwulog
				renwulogStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_LOG, user.Uid, user.Rid)))
				if err != nil {
					return nil, err
				}
				rwlog := db.Rwlogs{}
				err = json.Unmarshal([]byte(renwulogStr), &rwlog)
				if err != nil {
					return nil, err
				}
				/////////////
				rwlog.Isadd = db.Rwlogs_isadd_ABADON_TASK_EXCEPT_FIVE_MIN
				/////////////

				rb, err := json.Marshal(rwlog)
				if err != nil {
					return nil, err
				}
				_, err = conn.Do("set", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_LOG, rwlog.Userid, rwlog.Rid), string(rb))
				if err != nil {
					return nil, err
				}
				///////////////////
				// update
				uy, err := json.Marshal(user)
				if err != nil {
					return nil, err
				}
				_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid), string(uy))
				if err != nil {
					return nil, err
				}
				_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER_TOKEN, user.Token), string(uy))
				if err != nil {
					return nil, err
				}

				manager.addUpdate(&user)
				manager.addUpdate(&rwlog)
			}

			continue
		}

		// 任务进度2 也就是领了任务在送礼物  超过8分钟的
		if user.Rwjd == 3 {
			if int(time.Now().Unix())-user.Rwkstime > 60*8 {
				////////////////////
				user.Rwjd = -1
				user.Rid = -1
				///////////////////
				// renwulog
				renwulogStr, err := redis.String(conn.Do("get", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_LOG, user.Uid, user.Rid)))
				if err != nil {
					return nil, err
				}
				rwlog := db.Rwlogs{}
				err = json.Unmarshal([]byte(renwulogStr), &rwlog)
				if err != nil {
					return nil, err
				}
				/////////////
				rwlog.Isadd = db.Rwlogs_isadd_ABADON_TASK_EXCEPT_EIGHT_MIN
				/////////////

				rb, err := json.Marshal(rwlog)
				if err != nil {
					return nil, err
				}
				_, err = conn.Do("set", fmt.Sprintf("%v_%v_%v", global.REDIS_PREFIX_RENWU_LOG, rwlog.Userid, rwlog.Rid), string(rb))
				if err != nil {
					return nil, err
				}
				///////////////////

				///////////
				// update
				uy, err := json.Marshal(user)
				if err != nil {
					return nil, err
				}
				_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER, user.Uid), string(uy))
				if err != nil {
					return nil, err
				}
				_, err = conn.Do("set", fmt.Sprintf("%v%v", global.REDIS_PREFIX_USER_TOKEN, user.Token), string(uy))
				if err != nil {
					return nil, err
				}

				manager.addUpdate(&user)
				manager.addUpdate(&rwlog)
			}

			continue
		}

	}

	return nil, nil

}