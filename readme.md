### api

-   https://easydoc.net/s/20973294/X3ZT7UO2/Fy7Pp11S

txlogs 提现记录
jblogs 是记录 money 变动的
rwlogs 是记录任务变动

ALTER TABLE `renwu`
ADD COLUMN `is_only_one_time` int(1) NOT NULL DEFAULT 0 COMMENT '是否一个用户只能领取一次 0 否 1 是' AFTER `xtbbh`;
ADD COLUMN `lqzbyc` int(1) NOT NULL DEFAULT 0 COMMENT '一天只能领取那个主播任务一次 0 否 1 是' AFTER `is_only_one_time`;
ADD COLUMN `ipsync` int(11) NOT NULL DEFAULT 0 COMMENT '同 ip 只能进多少台' AFTER `lqzbyc`;

CREATE TABLE `iplogs` (
`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
`uid` int(10) NOT NULL,
`rid` int(10) unsigned NOT NULL,
`userid` bigint(20) unsigned DEFAULT NULL,
`ip` varchar(32) DEFAULT NULL,
`times` int(1) NOT NULL DEFAULT '0',
`day` datetime DEFAULT NULL,
PRIMARY KEY (`id`) USING BTREE,
) ENGINE=MyISAM AUTO_INCREMENT=115369 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC

1. 用户添加任务从 redis 过滤出符合的任务
2. 部分任务一个用户只能领取一次 (mysql 加了字段 is_only_one_time)
3. 一天只能领取那个主播任务一次 (mysql 加了字段 lqzbyc)
4. 用户 5 分钟内做这个任务失败了 下次就不让他领取这个任务
5. 还有个条件是限制一个任务 同 ip 只能进多少台 (mysql 加了字段 ipsync)
