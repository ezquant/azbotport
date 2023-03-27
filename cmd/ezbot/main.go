package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ezquant/ezbot/internal/backtesting"
	"github.com/ezquant/ezbot/internal/exchanges/binance"
	"github.com/ezquant/ezbot/internal/exchanges/paperwallet"
	"github.com/ezquant/ezbot/internal/models"
)

func main() {
	tradeCmd := flag.NewFlagSet("trade", flag.ExitOnError)
	dryRun := tradeCmd.Bool("dry-run", false, "Trade in dry-run mode")
	tradeConfigPath := tradeCmd.String("config", "user_data/config.yml", "Configuration file path")
	databasePath := tradeCmd.String("database", "user_data/db", "Database path")

	testCmd := flag.NewFlagSet("test", flag.ExitOnError)
	testConfigPath := testCmd.String("config", "user_data/config.yml", "Configuration file path")

	if len(os.Args) < 2 {
		log.Fatalln("expected 'trade' or 'test' subcommands")
	}

	switch os.Args[1] {

	case "trade":
		tradeCmd.Parse(os.Args[2:])
		config, err := readConfig(tradeConfigPath)
		if err != nil {
			log.Fatalf("cannot read config file: %v", err)
		}

		if *dryRun {
			paperwallet.Run(config, databasePath)
			return
		}

		binance.Run(config, databasePath)
	case "test":
		testCmd.Parse(os.Args[2:])
		config, err := readConfig(testConfigPath)
		if err != nil {
			log.Fatalf("cannot read config file: %v", err)
		}

		backtesting.Run(config, databasePath)
	default:
		fmt.Println("expected 'trade' or 'test' subcommands")
		os.Exit(1)
	}
}

func readConfig(path *string) (config *models.Config, err error) {
	data, err := os.ReadFile(*path)
	if err != nil {
		return nil, err
	}

	config = &models.Config{}
	err = yaml.Unmarshal(data, config)

	return
}
