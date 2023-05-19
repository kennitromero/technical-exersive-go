package main

import "errors"

type IConfigRepository interface {
	getNotificationLimitsByType(typeNotification string) (int, int, error)
}

type ConfigRepository struct {
}

func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{}
}

const (
	TypeNotificationStatus    = "status"
	TypeNotificationNews      = "news"
	TypeNotificationMarketing = "marketing"
	TypeNotificationDaily     = "daily"
)

func (c *ConfigRepository) getNotificationLimitsByType(typeNotification string) (int, int, error) {
	config := map[string]map[string]int{
		TypeNotificationStatus: {
			"limit_attempts": 2,
			"per_seconds":    60,
		},
		TypeNotificationNews: {
			"limit_attempts": 1,
			"per_seconds":    86400,
		},
		TypeNotificationMarketing: {
			"limit_attempts": 3,
			"per_seconds":    1,
		},
		TypeNotificationDaily: {
			"limit_attempts": 3,
			"per_seconds":    3600,
		},
	}

	if len(config[typeNotification]) == 0 {
		return 0, 0, errors.New("notification type does not exist")
	}

	return config[typeNotification]["limit_attempts"], config[typeNotification]["per_seconds"], nil
}
