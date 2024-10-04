package cmd

import (
	"context"
	"github.com/engpetarmarinov/redis-msg-processor-demo/internal/consumer"
	"github.com/engpetarmarinov/redis-msg-processor-demo/internal/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func NewStreamCMD(redisClient *redis.Client) *cobra.Command {
	var streamCMD = &cobra.Command{
		Use:     "stream <channel> <stream>",
		Aliases: []string{"s"},
		Short:   "Stream a pubsub channel to a stream",
		Example: `
$ consumer-cli stream messages:published messages:stream`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			channel := args[0]
			stream := args[1]

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			s := consumer.NewStreamer(redisClient)
			go s.Start(ctx, channel, stream)
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			<-sigCh
			logger.Info("Shutting down")
			//TODO: graceful shutdown
		},
	}

	return streamCMD
}
