package consumer

import (
	"context"
	"fmt"
	"github.com/engpetarmarinov/redis-msg-processor-demo/internal/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type Monitor struct {
	redisClient *redis.Client
}

func NewMonitor(redisClient *redis.Client) *Monitor {
	return &Monitor{redisClient: redisClient}
}

func (m *Monitor) Start(ctx context.Context, stream string, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	var previousLength int64 = 0

	for {
		select {
		case <-ticker.C:
			currentLength, err := m.redisClient.XLen(ctx, stream).Result()
			if err != nil {
				logger.Error("Error fetching stream length", "error", err)
				continue
			}

			newMessages := currentLength - previousLength
			if newMessages < 0 {
				newMessages = 0
			}

			processedPerSecond := fmt.Sprintf("%.2f", float64(newMessages)/float64(interval))

			logger.Info("Messages per interval", "count", processedPerSecond, "stream", stream, "interval", interval)

			previousLength = currentLength
		case <-ctx.Done():
			return
		}
	}
}
