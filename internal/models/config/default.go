package config

var defaultBase = Base{
	Server: ServerConfig{
		Port: 8080,
		Host: "0.0.0.0",
	},
	Database: DatabaseConfig{
		Type: "sqlite",
	},
	Log: LogConfig{
		Level:  "debug",
		Output: "console",
	},
}

var system = GroupAdvance{
	GroupName:   "system",
	Description: "系统配置",
	Data: []Advance{
		{
			Type:        "bool",
			Key:         "proxy.enable",
			Value:       "false",
			Description: "是否启用代理",
		},
		{
			Type:        "string",
			Key:         "proxy.url",
			Value:       "socks5://user:pass@127.0.0.1:1080",
			Description: "代理地址",
		},
		{
			Type:        "number",
			Key:         "task.max_timeout",
			Value:       "60",
			Description: "任务最大超时时间（秒）",
		},
		{
			Type:        "number",
			Key:         "task.max_retry",
			Value:       "3",
			Description: "任务最大重试次数",
		},
		{
			Type:        "number",
			Key:         "log.retention_days",
			Value:       "7",
			Description: "日志保留天数",
		},
		{
			Type:        "number",
			Key:         "notify.operation",
			Value:       "0",
			Description: "需要通知的操作类型",
		},
		{
			Type:        "number",
			Key:         "notify.id",
			Value:       "0",
			Description: "系统默认通知渠道",
		},
	},
}

var defaultAdvance = []GroupAdvance{
	system,
}

func DefaultAdvance() []GroupAdvance {
	return defaultAdvance
}

func DefaultBase() Base {
	return defaultBase
}
