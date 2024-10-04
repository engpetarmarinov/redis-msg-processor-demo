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

func NewMonitorCMD(redisClient *redis.Client) *cobra.Command {
	var monitorCMD = &cobra.Command{
		Use:     "monitor <stream>",
		Aliases: []string{"m"},
		Short:   "Monitors a stream and gives some stats",
		Example: `
$ consumer-cli monitor messages:processed --interval=3"`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stream := args[0]
			interval, err := cmd.Flags().GetInt("interval")
			if err != nil {
				logger.Error("Error", "error", err)
				os.Exit(1)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			m := consumer.NewMonitor(redisClient)
			go m.Start(ctx, stream, interval)
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			<-sigCh
			logger.Info("Shutting down")
			//TODO: graceful shutdown
		},
	}

	monitorCMD.Flags().Int("interval", 1, "interval in seconds")
	return monitorCMD
}
