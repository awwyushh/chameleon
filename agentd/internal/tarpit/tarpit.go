package tarpit

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

type Tarpit struct {
	rdb       *redis.Client
	threshold int
	baseMs    int
	growth    float64
	maxMs     int
	windowSec int
}

func New(rdb *redis.Client, threshold, baseMs int, growth float64, maxMs, windowSec int) *Tarpit {
	return &Tarpit{rdb: rdb, threshold: threshold, baseMs: baseMs, growth: growth, maxMs: maxMs, windowSec: windowSec}
}

// ShouldTarpit increments counter and returns delay_ms (0 if none)
func (t *Tarpit) ShouldTarpit(ctx context.Context, key string) (int, error) {
	redisKey := fmt.Sprintf("tarpit:%s", key)
	val, err := t.rdb.Incr(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	// set TTL
	t.rdb.Expire(ctx, redisKey, time.Duration(t.windowSec)*time.Second)
	if int(val) <= t.threshold {
		return 0, nil
	}
	count := float64(val - int64(t.threshold))
	delay := float64(t.baseMs) * math.Min(1.0+math.Log(1.0+count), t.growth)
	if delay > float64(t.maxMs) {
		delay = float64(t.maxMs)
	}
	return int(math.Round(delay)), nil
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
