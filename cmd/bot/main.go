package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"telegram-bot-pocket/pkg/config"
	"telegram-bot-pocket/pkg/repository"
	"telegram-bot-pocket/pkg/repository/boltdb"
	"telegram-bot-pocket/pkg/server"
	"telegram-bot-pocket/pkg/telegram"
)

func main() {

	cnf, err := config.Init()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cnf.TelegramToken)

	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cnf.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cnf)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cnf.AuthServerURL, cnf.Messages)

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cnf.TelegramBotUrl)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}

}

func initDB(cfg *config.Config) (*bolt.DB, error) {

	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil

}
