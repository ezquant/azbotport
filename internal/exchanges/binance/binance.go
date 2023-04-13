package binance

import (
	"context"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/ezquant/azbot/azbot"
	"github.com/ezquant/azbot/azbot/exchange"
	"github.com/ezquant/azbot/azbot/storage"
	"github.com/ezquant/azbotport/internal/localkv"
	"github.com/ezquant/azbotport/internal/models"
	"github.com/ezquant/azbotport/internal/strategies"
)

func Run(config *models.Config, databasePath *string) {
	var (
		ctx              = context.Background()
		binanceAPIKey    = os.Getenv("BINANCE_API_KEY")
		binanceSecretKey = os.Getenv("BINANCE_SECRET_KEY")
		telegramToken    = os.Getenv("TELEGRAM_TOKEN")
		telegramUser, _  = strconv.Atoi(os.Getenv("TELEGRAM_USER"))
		pairs            = make([]string, 0, len(config.AssetWeights))
	)

	for pair := range config.AssetWeights {
		pairs = append(pairs, pair)
	}

	settings := azbot.Settings{
		Pairs: pairs,
		Telegram: azbot.TelegramSettings{
			Enabled: true,
			Token:   telegramToken,
			Users:   []int{telegramUser},
		},
	}

	// creating a storage to save trades
	storage, err := storage.FromFile(path.Join(*databasePath, "trades.db"))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize exchange
	binanceCredential := exchange.WithBinanceCredentials(binanceAPIKey, binanceSecretKey)
	binance, err := exchange.NewBinance(ctx, binanceCredential)
	if err != nil {
		log.Fatalln(err)
	}

	// initialize local KV store for strategies
	kv, err := localkv.NewLocalKV(databasePath)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize strategy and bot
	switch config.Strategy {
	case "DiamondHands":
		strat, err := strategies.NewDiamondHands(config, kv)
		if err != nil {
			log.Fatal(err)
		}

		bot, err := azbot.NewBot(ctx, settings, binance, strat, azbot.WithStorage(storage))
		if err != nil {
			log.Fatalln(err)
		}

		// Run azbot
		err = bot.Run(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		log.Fatal("Invalid strategy")
	}
}
