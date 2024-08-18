package telegram

import "errors"
import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnauthorized = errors.New("user is not unauthorized")
	errUnableToSave = errors.New("unable to save")
)

func (b *Bot) handleError(chatId int64, err error) {

	msg := tgbotapi.NewMessage(chatId, b.messages.Default)

	switch err {
	case errInvalidURL:
		msg.Text = b.messages.InvalidUrl
		b.bot.Send(msg)
	case errUnauthorized:
		msg.Text = b.messages.Unauthorized
		b.bot.Send(msg)
	case errUnableToSave:
		msg.Text = b.messages.UnableToSave
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}
}
