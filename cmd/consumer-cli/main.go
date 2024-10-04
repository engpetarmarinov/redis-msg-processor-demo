package main

import (
	"fmt"
	"github.com/engpetarmarinov/redis-msg-processor-demo/internal/cli/cmd"
	"github.com/engpetarmarinov/redis-msg-processor-demo/internal/config"
	"github.com/engpetarmarinov/redis-msg-processor-demo/internal/logger"
	"github.com/redis/go-redis/v9"
)

func main() {
	opts := config.NewConfigOpt().WithLogLevel().WithBrokerAddress().WithBrokerPort().WithBrokerPassword()
	conf := config.NewConfig(opts)
	logger.Init(logger.NewConfigOpt().WithLevel(conf.LogLevel))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.BrokerAddress, conf.BrokerPort),
		Password: conf.BrokerPassword,
	})

	cli := cmd.NewCLI(redisClient)
	cli.Run()
}
