package service

import (
	"douyin/global"
	"douyin/web/db"
	"fmt"
	"time"

	"github.com/robfig/cron"
)

func RunCronTasks() {

	c := cron.New()
	// every 5s
	c.AddFunc("*/5 * * * * ?", func() {
		fmt.Println("task run...")
		go run()
	})
	c.Start()
	select {}
}

func run() {
	t := &Task{}
	t.task1()
	t.task2()
}

type Task struct {
}

func (t *Task) task2() (interface{}, error) {
	// 建议用一个线程定时获取 mysql 有效任务(shengyusl>0 and stop=0)，到 redis 里面，redis 可以用集合或者其他方式
	// 有效任务(shengyusl>0 and stop=0)在程序执行之后已经加载
	// 是定时清理redis中的数量为0的任务

	for rid := range manager.renwuIDSet {
		renwu, err := manager.getRenwu(rid)
		if err != nil {
			continue
		}
		if renwu.Shengyusl <= 0 || renwu.Stop <= 0 {
			// 从redis移除
			manager.delRenwu(renwu)
		}
	}

	return nil, nil
}
func (t *Task) task1() (interface{}, error) {
	// 服务器会定时把在任务进度1也就是刚领了任务超过5分钟的
	// 和任务进度2 也就是领了任务在送礼物  超过8分钟的
	// 用户取消任务  因为超时了

	conn := global.REDIS.Get()
	defer conn.Close()

	for userid := range manager.yonghuSet {

		// user
		user, err := manager.getUser(userid)
		if err != nil {
			continue
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
				rwlog, err := manager.getRenwulog(user.Uid, user.Rid)
				if err != nil {
					continue
				}
				/////////////
				rwlog.Isadd = db.Rwlogs_isadd_ABADON_TASK_NOT_IN
				/////////////

				manager.setUser(user)
				manager.setRenwulog(rwlog)
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
				rwlog, err := manager.getRenwulog(user.Uid, user.Rid)
				if err != nil {
					continue
				}
				/////////////
				rwlog.Isadd = db.Rwlogs_isadd_ABADON_TASK_EXCEPT_EIGHT_MIN
				/////////////

				manager.setRenwulog(rwlog)
				manager.setUser(user)
			}

			continue
		}

	}

	return nil, nil

}
