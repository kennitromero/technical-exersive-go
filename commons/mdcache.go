package commons

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type ImdCache interface {
	Get(key string) (string, error)
	Set(key string, value string, seconds int) error
}

type MDRedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewMDRedisCache(client *redis.Client, ctx context.Context) *MDRedisCache {
	return &MDRedisCache{
		client: client,
		ctx:    ctx,
	}
}

func (mCache *MDRedisCache) Get(key string) (string, error) {
	val, err := mCache.client.Get(mCache.ctx, key).Result()

	//redis there not exist for key
	if err == redis.Nil {
		return "", nil
	}

	//redis general error
	if err != nil {
		return "", err
	}

	return val, nil
}

func (mCache *MDRedisCache) Set(key string, value string, seconds int) error {
	err := mCache.client.Set(mCache.ctx, key, value, time.Duration(seconds)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}
