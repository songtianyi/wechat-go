package rrredis

import (
	"fmt"
	"gopkg.in/redis.v5"
	"sync"
	"time"
)

var (
	Nil = redis.Nil
)

type RedisOptions struct {
	dialTimeout  time.Duration
	db           int
	password     string
	connPoolSize int
}

var (
	cp              = &clientPool{pool: make(map[string]*redis.Client)}
	defaultPoolSize = 500
)

type clientPool struct {
	pool map[string]*redis.Client
	mu   sync.RWMutex
}

func (s *clientPool) add(addr string, rc *redis.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pool[addr] = rc
}

func (s *clientPool) get(addr string) *redis.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.pool[addr]; ok {
		return v
	}
	return nil
}

// if you wanna customize redis options, connect once when program start
func Connect(addr string, opt *RedisOptions) error {
	if addr == "" || opt == nil {
		return fmt.Errorf("Redis addr empty or options nil")
	}

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          opt.db,
		DialTimeout: opt.dialTimeout,
		PoolSize:    opt.connPoolSize,
		Password:    opt.password,
	})
	if client == nil {
		return fmt.Errorf(fmt.Sprintf("Connect to redis [%s] fail", addr))
	}
	if _, err := client.Ping().Result(); err != nil {
		return err
	}
	cp.add(addr, client)
	return nil
}

// get redis client before you want to do some operations
func GetRedisClient(addr string) (error, *RedisClient) {
	if c := cp.get(addr); c != nil {
		return nil, &RedisClient{c: c}
	}
	err := Connect(addr, &RedisOptions{
		db:           0,
		dialTimeout:  1 * time.Second,
		connPoolSize: defaultPoolSize,
		password:     "",
	})
	if err != nil {
		return fmt.Errorf("Couldn't find redis client for [%s], failed when try to get a new one, %s", addr, err), nil
	}
	return GetRedisClient(addr)
}

type RedisClient struct {
	c *redis.Client
}

// list of redis operations
// hash set
func (c *RedisClient) HMSet(key string, fields map[string]string) error {
	status := c.c.HMSet(key, fields)
	return status.Err()
}

func (c *RedisClient) HMGet(key string, fields ...string) ([]interface{}, error) {
	status := c.c.HMGet(key, fields...)
	return status.Result()
}

func (c *RedisClient) HMDelete(key string, fields ...string) error {
	status := c.c.HDel(key, fields...)
	return status.Err()
}

func (c *RedisClient) HMExists(key, field string) (bool, error) {
	status := c.c.HExists(key, field)
	return status.Result()
}

func (c *RedisClient) Exists(key string) bool {
	return c.c.Exists(key).Val()
}

func (c *RedisClient) Expire(key string, t time.Duration) (bool, error) {
	status := c.c.Expire(key, t)
	return status.Result()
}

func (c *RedisClient) Rename(okey, nkey string) error {
	status := c.c.Rename(okey, nkey)
	return status.Err()
}

func (c *RedisClient) Keys(pattern string) ([]string, error) {
	status := c.c.Keys(pattern)
	return status.Result()
}

func (c *RedisClient) KeyExists(key string) (bool, error) {
	status := c.c.Exists(key)
	return status.Result()
}

func (c *RedisClient) Del(key string) error {
	status := c.c.Del(key)
	return status.Err()
}

func (c *RedisClient) Set(key string, value string, expir time.Duration) error {
	status := c.c.Set(key, value, expir)
	return status.Err()
}

func (c *RedisClient) Get(key string) ([]byte, error) {
	status := c.c.Get(key)
	b, _ := status.Bytes()
	return b, status.Err()
}

// redis list operations
func (c *RedisClient) LPop(key string) (string, error) {
	status := c.c.LPop(key)
	return status.Result()
}

func (c *RedisClient) RPush(key string, values ...interface{}) (int64, error) {
	status := c.c.RPush(key, values...)
	return status.Result()
}

func (c *RedisClient) LLen(key string) (int64, error) {
	status := c.c.LLen(key)
	return status.Result()
}

// redis sorted set operations
func (c *RedisClient) ZInterStore(dest string, aggregate string, keys ...string) error {
	weights := make([]float64, len(keys))
	for i, n := 0, len(keys); i < n; i++ {
		weights[i] = 1
	}
	t := redis.ZStore{
		weights,
		aggregate,
	}
	status := c.c.ZInterStore(dest, t, keys...)
	return status.Err()
}

func (c *RedisClient) ZCard(key string) (int64, error) {
	status := c.c.ZCard(key)
	return status.Result()
}

func (c *RedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	status := c.c.ZRange(key, start, stop)
	return status.Result()
}

func (c *RedisClient) ZRevRange(key string, start, stop int64) ([]string, error) {
	status := c.c.ZRevRange(key, start, stop)
	return status.Result()
}

func (c *RedisClient) ZAddBatch(key string, score_lst []float64, data_lst []interface{}) int64 {
	mem_lst := make([]redis.Z, 0)
	for i := 0; i < len(score_lst); i++ {
		tmp := redis.Z{
			Score:  score_lst[i],
			Member: data_lst[i],
		}

		mem_lst = append(mem_lst, tmp)
	}
	status := c.c.ZAdd(key, mem_lst...)
	return status.Val()
}

func (c *RedisClient) ZRangeByScore(key, min, max string, offset, count int64) ([]string, error) {
	opt := &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}
	status := c.c.ZRangeByScore(key, *opt)
	return status.Result()
}

func (c *RedisClient) ZRangeByScoreWithScores(key, min, max string, offset, count int64) ([]redis.Z, error) {
	opt := &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}
	status := c.c.ZRangeByScoreWithScores(key, *opt)
	return status.Result()
}

func (c *RedisClient) Incr(key string) (int64, error) {
	status := c.c.Incr(key)
	return status.Result()
}
