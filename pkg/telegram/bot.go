package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"telegram-bot-pocket/pkg/config"
	"telegram-bot-pocket/pkg/repository"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
	messages        config.Messages
}

func NewBot(bot *tgbotapi.BotAPI,
	pocketClient *pocket.Client,
	tr repository.TokenRepository,
	redirectURL string,
	messages config.Messages,
) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, redirectURL: redirectURL, tokenRepository: tr, messages: messages}
}

func (b *Bot) Start() error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.initUpdatesChannel()

	b.handlerUpdates(updates)

	return nil
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handlerUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				if err := b.handleCommand(update.Message); err != nil {
					b.handleError(update.Message.Chat.ID, err)
				}
				continue
			}
			if err := b.handleMessage(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
		}
	}
}
