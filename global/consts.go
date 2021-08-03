package global

const (
	REDIS_PREFIX_USER       string = "user_"
	REDIS_PREFIX_RENWU      string = "renwu_"
	REDIS_PREFIX_RENWU_LOCK string = "renwu_lock_"
	REDIS_PREFIX_USER_TOKEN string = ""

	// renwulog_${userid}_${renwuid}
	REDIS_PREFIX_RENWU_LOG string = "renwulog_"

	// renwulogip_${ip}_${renwuid}
	REDIS_PREFIX_RENWU_IP string = "renwulogip_"
)

const (
	// 同ip只能进多少台
	MAX_IP_TASK int = 20
)
