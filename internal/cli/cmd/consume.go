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

func NewConsumeCMD(redisClient *redis.Client) *cobra.Command {
	var consumeCMD = &cobra.Command{
		Use:     "consume <stream> [flags]",
		Aliases: []string{"c"},
		Short:   "Consume a stream of messages",
		Long: `
	Consume a stream with messages.

	The --size flag can be used to set number of consumers.`,
		Example: `
$ consumer-cli consume messages:streamed --size=4 --limit=10`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stream := args[0]
			size, err := cmd.Flags().GetInt("size")
			if err != nil {
				logger.Error("Error", "error", err)
				os.Exit(1)
			}

			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				logger.Error("Error", "error", err)
				os.Exit(1)
			}

			//TODO: size and limit validation
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := consumer.NewConsumer(redisClient)
			go c.Start(ctx, stream, size, limit)

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			<-sigCh
			logger.Info("Shutting down")
			c.Stop(ctx, size)
		},
	}

	consumeCMD.Flags().Int("size", 1, "number of consumers")
	consumeCMD.Flags().Int("limit", 10, "limit of how many messages will be fetched at once")
	return consumeCMD
}
