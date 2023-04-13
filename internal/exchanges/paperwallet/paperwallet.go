package paperwallet

import (
	"context"
	"os"
	"strconv"

	"github.com/ezquant/azbotport/internal/localkv"
	"github.com/ezquant/azbotport/internal/models"
	"github.com/ezquant/azbotport/internal/strategies"

	"github.com/ezquant/azbot/azbot"
	"github.com/ezquant/azbot/azbot/exchange"
	"github.com/ezquant/azbot/azbot/storage"

	log "github.com/sirupsen/logrus"
)

func Run(config *models.Config, databasePath *string) {
	var (
		ctx             = context.Background()
		telegramToken   = os.Getenv("TELEGRAM_TOKEN")
		telegramUser, _ = strconv.Atoi(os.Getenv("TELEGRAM_USER"))
		pairs           = make([]string, 0, len(config.AssetWeights))
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

	// Use binance for realtime data feed
	binance, err := exchange.NewBinance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// creating a storage to save trades
	storage, err := storage.FromMemory()
	if err != nil {
		log.Fatal(err)
	}

	// creating a paper wallet to simulate an exchange waller for fake operations
	paperWallet := exchange.NewPaperWallet(
		ctx,
		"BUSD",
		exchange.WithPaperFee(0.001, 0.001),
		exchange.WithPaperAsset("BUSD", 500),
		exchange.WithDataFeed(binance),
	)

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

		bot, err := azbot.NewBot(
			ctx,
			settings,
			paperWallet,
			strat,
			azbot.WithStorage(storage),
			azbot.WithPaperWallet(paperWallet),
		)
		if err != nil {
			log.Fatalln(err)
		}

		err = bot.Run(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		log.Fatal("Invalid strategy")
	}
}
