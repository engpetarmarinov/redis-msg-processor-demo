package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/engpetarmarinov/redis-msg-processor-demo/internal/logger"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type Consumer struct {
	redisClient *redis.Client
}

// NewConsumer inits *Consumer
func NewConsumer(client *redis.Client) *Consumer {
	return &Consumer{
		redisClient: client,
	}
}

// Start creates a consumer group and starts {size} consumers for a stream, fetching {limit} messages in batches
func (c *Consumer) Start(ctx context.Context, stream string, size int, limit int) {
	err := c.initConsumerGroup(ctx, stream)
	if err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists" {
		logger.Info("Consumer group already exists", "group", getGroupName(stream))
	} else if err != nil && !errors.Is(err, redis.Nil) {
		logger.Error("Failed to create consumer group", "error", err, "group", getGroupName(stream))
		return
	}

	logger.Info("subscribing to stream", "stream", stream, "size", size)
	// Start the configured number of consumers
	for i := 0; i < size; i++ {
		consumerID := getConsumerID(i)
		err := c.registerConsumer(ctx, consumerID)
		if err != nil {
			logger.Error("Failed to register consumer", "consumer", consumerID, "error", err)
		}
		go c.startConsumer(ctx, consumerID, stream, limit)
	}
}

func (c *Consumer) Stop(ctx context.Context, size int) {
	for i := 0; i < size; i++ {
		consumerID := getConsumerID(i)
		err := c.deregisterConsumer(ctx, consumerID)
		if err != nil {
			logger.Error("Failed to deregister consumer", "error", err)
		}
	}
}

func (c *Consumer) initConsumerGroup(ctx context.Context, stream string) error {
	return c.redisClient.XGroupCreateMkStream(ctx, stream, getGroupName(stream), "$").Err()
}

// registerConsumer Add consumer ID to the Redis List "consumer:ids"
func (c *Consumer) registerConsumer(ctx context.Context, consumerID string) error {
	return c.redisClient.LPush(ctx, "consumer:ids", consumerID).Err()
}

// deregisterConsumer Remove consumer ID from the Redis List "consumer:ids"
func (c *Consumer) deregisterConsumer(ctx context.Context, consumerID string) error {
	return c.redisClient.LRem(ctx, "consumer:ids", 0, consumerID).Err()
}

func (c *Consumer) processMessage(ctx context.Context, consumerID string, msgValues map[string]interface{}) {
	processCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	done := make(chan struct{})
	go func() {
		// Simulate processing by adding a random property
		messageJSON, err := json.Marshal(msgValues)
		if err != nil {
			logger.Error("Error marshalling message", "error", err)
			return
		}

		var message Message
		err = json.Unmarshal(messageJSON, &message)
		if err != nil {
			logger.Error("Error unmarshalling message", "error", err)
			return
		}

		_, err = c.redisClient.XAdd(ctx, &redis.XAddArgs{
			//TODO: param or config?
			Stream: "messages:processed",
			Values: map[string]interface{}{
				"consumer_id":  consumerID,
				"message_id":   message.MessageID,
				"processed_at": time.Now().Format(time.RFC3339),
			},
		}).Result()

		if err != nil {
			logger.Error("Error adding to stream", "consumer", consumerID, "message", messageJSON, "error", err)
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
		return
	case <-processCtx.Done():
		logger.Warn("Processing message timed out", "consumer", consumerID, "message", msgValues)
		return
	}
}

func (c *Consumer) startConsumer(ctx context.Context, consumerID string, streamName string, limit int) {
	for {
		stream, err := c.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    getGroupName(streamName),
			Consumer: consumerID,
			Streams:  []string{streamName, ">"},
			Count:    int64(limit),
			Block:    0,
		}).Result()

		if err != nil {
			logger.Error("Error reading from stream", "error", err)
			continue
		}

		for _, message := range stream[0].Messages {
			//TODO: take out processing of messages out of the consumer
			//TODO: return error from processMessage
			c.processMessage(ctx, consumerID, message.Values)

			// Acknowledge the message after processing
			c.redisClient.XAck(ctx, streamName, getGroupName(streamName), message.ID)
		}
	}
}

func getGroupName(streamName string) string {
	return fmt.Sprintf("consumer-group-%s", streamName)
}

func getConsumerID(seq int) string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s-consumer-%d", hostname, seq+1)
}
