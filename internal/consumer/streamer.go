package consumer

import (
	"context"
	"encoding/json"
	"github.com/engpetarmarinov/redis-msg-processor-demo/internal/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

// Streamer is a way to sub a channel and stream it to a stream
type Streamer struct {
	redisClient *redis.Client
}

// NewStreamer inits Streamer
func NewStreamer(redisClient *redis.Client) *Streamer {
	return &Streamer{redisClient: redisClient}
}

// Start starts listening to a pubsub channel and streams every message to a stream
func (s *Streamer) Start(ctx context.Context, channel string, stream string) {
	logger.Info("Streaming channel", "stream", stream, "channel", channel)
	pubsub := s.redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()
	ch := pubsub.Channel()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				logger.Info("PubSub channel closed, retrying subscription...")
				pubsub = s.redisClient.Subscribe(ctx, channel)
				ch = pubsub.Channel()
				continue
			}

			var decodedMsg Message
			err := json.Unmarshal([]byte(msg.Payload), &decodedMsg)
			if err != nil {
				logger.Error("Error unmarshalling message", "error", err, "payload", msg.Payload)
				continue
			}

			// Publish the received message to the Redis stream
			_, err = s.redisClient.XAdd(ctx, &redis.XAddArgs{
				Stream: stream,
				Values: map[string]interface{}{
					"message_id": decodedMsg.MessageID,
				},
			}).Result()
			if err != nil {
				logger.Error("Error publishing to stream", "error", err)
			}

		case <-time.After(10 * time.Second):
			//TODO: do I need to ping? go-redis pings internally
			logger.Info("No messages received, checking connection...")
			err := pubsub.Ping(ctx)
			if err != nil {
				logger.Error("Error checking connection", "error", err)
				return
			}
		case <-ctx.Done():
			logger.Info("Context done, exiting from streamer...")
			return
		}
	}
}
