package main

import (
	"errors"
	"technical-exersive/commons"
)

type IConfigRepository interface {
	getNotificationLimitsByType(typeNotification string) (int, int, error)
}

type ConfigRepository struct {
}

func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{}
}

func (c *ConfigRepository) getNotificationLimitsByType(typeNotification string) (int, int, error) {
	config := map[string]map[string]int{
		commons.TypeNotificationStatus: {
			"limit_attempts": 2,
			"per_seconds":    60,
		},
		commons.TypeNotificationNews: {
			"limit_attempts": 1,
			"per_seconds":    86400,
		},
		commons.TypeNotificationMarketing: {
			"limit_attempts": 3,
			"per_seconds":    1,
		},
		commons.TypeNotificationDaily: {
			"limit_attempts": 3,
			"per_seconds":    3600,
		},
	}

	if len(config[typeNotification]) == 0 {
		return 0, 0, errors.New("notification type does not exist")
	}

	return config[typeNotification]["limit_attempts"], config[typeNotification]["per_seconds"], nil
}
