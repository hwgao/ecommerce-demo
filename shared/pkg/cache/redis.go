package cache

import (
    "encoding/json"
    "time"
    "github.com/go-redis/redis/v8"
    "context"
)

type Cache interface {
    Set(key string, value interface{}, expiration time.Duration) error
    Get(key string) (interface{}, bool)
    Delete(key string) error
}

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(redisURL string) *RedisCache {
    opt, _ := redis.ParseURL(redisURL)
    client := redis.NewClient(opt)
    
    return &RedisCache{client: client}
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return r.client.Set(context.Background(), key, data, expiration).Err()
}

func (r *RedisCache) Get(key string) (interface{}, bool) {
    val, err := r.client.Get(context.Background(), key).Result()
    if err != nil {
        return nil, false
    }
    
    var result interface{}
    if err := json.Unmarshal([]byte(val), &result); err != nil {
        return nil, false
    }
    
    return result, true
}

func (r *RedisCache) Delete(key string) error {
    return r.client.Del(context.Background(), key).Err()
}
