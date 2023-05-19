package main

import (
	"encoding/json"
	"technical-exersive/commons"
)

const FormatCompositeKey string = "TO#WAY_TO_NOTIFY#TYPE_NOTIFICATION"
const FormatCompositeKeyTo string = "TO"
const FormatCompositeKeyWayToNotify string = "WAY_TO_NOTIFY"
const FormatCompositeKeyTypeNotification string = "TYPE_NOTIFICATION"

type ICacheRepository interface {
	get(compositeKey string) (RateLimitCounter, error)
	set(rlc RateLimitCounter) error
}

type CacheRepository struct {
	cache commons.ImdCache
}

func NewCacheRepository(cache commons.ImdCache) *CacheRepository {
	return &CacheRepository{cache: cache}
}

func (rl *CacheRepository) get(compositeKey string) (RateLimitCounter, error) {
	valueAsString, err := rl.cache.Get(compositeKey)

	if err == nil && valueAsString == "" {
		return RateLimitCounter{}, nil
	}

	if err != nil {
		return RateLimitCounter{}, err
	}

	var rateLimitCounter RateLimitCounter
	err = json.Unmarshal([]byte(valueAsString), &rateLimitCounter)

	if err != nil {
		return RateLimitCounter{}, err
	}

	return rateLimitCounter, nil
}

func (rl *CacheRepository) set(rlc RateLimitCounter) error {
	valueAsString, err := json.Marshal(rlc)
	if err != nil {
		return err
	}

	return rl.cache.Set(rlc.Key, string(valueAsString), rlc.PerSeconds)
}
