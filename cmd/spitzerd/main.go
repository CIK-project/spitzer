package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/CIK-project/spitzer/types"
	"github.com/spf13/cobra"
	"github.com/CIK-project/spitzer/hub/gql"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	config = types.DBConfig{}
	endpoint = "0.0.0.0:8080"
	cors = "*"
)

func main() {
	rootCmd := &cobra.Command {
		Use: "spitzerd",
		Short: "Start spitzer's graphql endpoint deamon",
		Long: "Start spitzer's graphql endpoint deamon",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.ValidateBasic()
			if err != nil {
				return err
			}

			run(config, endpoint, cors)
			return nil
		},
	}

	rootCmd.Flags().StringVar(&config.Host, "db.host", "", "")
	rootCmd.Flags().StringVar(&config.User, "db.user", "", "")
	rootCmd.Flags().StringVar(&config.Password, "db.password", "", "")
	rootCmd.Flags().StringVar(&config.DBName, "db.name", "", "")
	rootCmd.Flags().StringVar(&endpoint, "endpoint", endpoint, "")
	rootCmd.Flags().StringVar(&cors, "cors", cors, "")

	if err := rootCmd.Execute(); err != nil {
		logger.Error(fmt.Sprintf("Failed to parse cli: %s", err.Error()))
		os.Exit(1)
	}
}

func run(config types.DBConfig, endpoint, cors string) {
	qls := gql.NewGqlService(logger, config, endpoint, cors)

	// Stop upon receiving SIGTERM or CTRL-C
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range exit {
			logger.Error(fmt.Sprintf("captured %v, exiting...", sig))
			if qls.IsRunning() {
				qls.Stop()
			}
			os.Exit(1)
		}
	}()

	if err := qls.Start(); err != nil {
		logger.Error(fmt.Sprintf("Failed to start: %v", err))
		os.Exit(1)
	}

	// Run forever
	select {}
}