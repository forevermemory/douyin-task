package global

const (
	REDIS_PREFIX_USER       string = "user_"
	REDIS_PREFIX_RENWU      string = "renwu_"
	REDIS_PREFIX_USER_TOKEN string = ""

	// renwulog_${userid}_${renwuid}
	REDIS_PREFIX_RENWU_LOG string = "renwulog_"

	// renwulogip_${ip}_${renwuid}
	REDIS_PREFIX_RENWU_LOG_IP string = "renwulogip_"
)
