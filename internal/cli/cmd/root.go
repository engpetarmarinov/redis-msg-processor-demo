package cmd

import (
	"github.com/engpetarmarinov/redis-msg-processor-demo/internal/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"os"
)

const Version = "0.0.1"

type CLI struct {
	rootCmd *cobra.Command
}

func NewCLI(redisClient *redis.Client) *CLI {
	var rootCMD = &cobra.Command{
		Use:           "consumer-cli <command> <subcommand> [flags]",
		Short:         "consumer-cli",
		Long:          `Command line tool to demo redis pubsub and streams`,
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	consumeCMD := NewConsumeCMD(redisClient)
	rootCMD.AddCommand(consumeCMD)
	streamerCMD := NewStreamCMD(redisClient)
	rootCMD.AddCommand(streamerCMD)
	monitorCMD := NewMonitorCMD(redisClient)
	rootCMD.AddCommand(monitorCMD)

	return &CLI{
		rootCmd: rootCMD,
	}
}

func (cli *CLI) Run() {
	if err := cli.rootCmd.Execute(); err != nil {
		logger.Error("Error", "error", err)
		os.Exit(1)
	}
}
