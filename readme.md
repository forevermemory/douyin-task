### api

-   https://easydoc.net/s/20973294/X3ZT7UO2/Fy7Pp11S

txlogs 提现记录
jblogs 是记录 money 变动的
rwlogs 是记录任务变动

-   table change
    ALTER TABLE `renwu`
    ADD COLUMN `is_only_one_time` int(1) NOT NULL DEFAULT 0 COMMENT '是否一个用户只能领取一次 0 否 1 是' AFTER `xtbbh`;
    ADD COLUMN `lqzbyc` int(1) NOT NULL DEFAULT 0 COMMENT '一天只能领取那个主播任务一次 0 否 1 是' AFTER `is_only_one_time`;
    ADD COLUMN `ipsync` int(11) NOT NULL DEFAULT 0 COMMENT '同 ip 只能进多少台' AFTER `lqzbyc`;

### 7.31 上午

1. 用户添加任务从 redis 过滤出符合的任务
2. 部分任务一个用户只能领取一次 (mysql 加了字段 is_only_one_time)
3. 一天只能领取那个主播任务一次 (mysql 加了字段 lqzbyc)
4. 用户 5 分钟内做这个任务失败了 下次就不让他领取这个任务
5. 还有个条件是限制一个任务 同 ip 只能进多少台 (mysql 加了字段 ipsync)

### 7.31 待优化

1、任务架构问题
建议用一个线程定时获取 mysql 有效任务(shengyusl>0 and stop=0)，到 redis 里面，redis 可以用集合或者其他方式

2、我看目前用户表是用的 token 做的主键，我这边 token 是用的 aes 加密 解密后可以得到 uid，我待会把加密解密函数发给你，用 uid 做 key

3、top6：可以直接用语句 SELECT `account` FROM `yonghu` WHERE `dyid`= 返回用户提交的 dyid 绑定的账户名

4、top101：添加任务的时候 tqjs 为数量的 0.5 ，stop 默认为 0，biaoshi 是唯一的，防止重复放单

5、用户获取任务不能把任务表的所有信息都返回给用户，这个接口里面有。可以考虑把需要的信息写入 redis 其他的不写

6、限制 ip 可以写进 redis 用 rid+ip 地址 存储计数，达到限制比如 20 返回失败

### 8.1 优化 7.31 问题

-   1. 7.31 的修改已完成
-   2. table change 的几个字段变更 必须要有
-   3. 代码优化
-   4. Top2 生成 token `utils.GetToken`你来吧

### 8.3

同步任务表只获取下单时间 5 小时内的

rwlogs 不写进内存 读 mysql

Top1、用户注册的时候 如果带有上级需要把上级 uid 查询出来写进 mysql，上级在用户增加余额的时候会分 10% -->done

Top2、Top5 登录的时候会获取用户 onlie 如果大于 0 表示用户在线 不允许重复登录 --> done

Middle4、 取消任务的时候记得进锁 防止该任务有其他人在更改 --> done

Down1、 我目前服务器是加锁 然后使用 mysql sum 语句统计子账户金币并且清零 一次性加到主账号上，看看有什么高效方法，主要防止转移过程中子账户有任务完成 增加金币 -->done

### 8.4

-   1. 把 yonghu 的 rwksTime 字段 无符号去掉

### 8.6
- 优化更新指定字段
