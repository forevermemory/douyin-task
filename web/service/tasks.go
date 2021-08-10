package service

import (
	"douyin/web/db"
	"fmt"
	"strconv"

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

	renwuids, err := manager.getRenwuSet()
	if err != nil {
		return nil, err
	}

	for _, rwidStr := range renwuids {
		rid, _ := strconv.Atoi(rwidStr)
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

	users, err := db.ListYonghuV2()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		user.Rwjd = -1
		user.Rid = -1
		// renwulog
		rwlog, err := manager.getRenwulog(user.Uid, user.Rid)
		if err != nil {
			continue
		}
		/////////////
		rwlog.Isadd = db.Rwlogs_isadd_ABADON_TASK_NOT_IN
		/////////////
		user.UpdateType = db.USER_UPDATE_ONLY_RIDAND_RWID
		manager.setUser(user)
		manager.setRenwulog(rwlog)
	}

	return nil, nil

}
