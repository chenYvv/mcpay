package tg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mcpay/pkg/logger"
)

func TgBot() {
	token := viper.GetString("G_MConfig.tg_bot_token")
	//token := "7411793460:AAH-Miz5KOqYMqCGCyD6UYMgBybed_Ze-J4"
	if token == "" {
		return
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Info("机器人 初始化", zap.String("err", err.Error()))
		return
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	//{"ok":true,"result":[{"update_id":754793215,"pre_checkout_query":{"id":"1994504975831482675","from":{"id":6906832821,"is_bot":false,"first_name":"Vicki","last_name":"Cooper","username":"nvS1mple","language_code":"zh-hans"},"currency":"CNY","total_amount":2000,"invoic e_payload":"payload"}}]}
	//{"ok":true,"result":[{"update_id":754793216, "message":{"message_id":19,"from":{"id":6906832821,"is_bot":false,"first_name":"Vicki","last_name":"Cooper","username":"nvS1mple","language_code":"zh-hans"},"chat":{"id":6906832821,"first_name":"Vicki","last_name":"Cooper", "username":"nvS1mple","type":"private"},"date":1722498276,"successful_payment":{"currency":"CNY","total_amount":2000,"invoice_payload":"payload","telegram_payment_charge_id":"7263151821_6906832821_1432483489_322_73980731374 25856208","provider_payment_charge_id":"ch_3PitgcLsgrYi154l08slpmgf"}}}]}
	for update := range updates {
		if update.Message != nil { // If we got a message
			//log.Debug("获取到消息 [%s] %s", update.Message.From.UserName, update.Message.Text)
			logger.Info("机器人", zap.String("获取到消息", fmt.Sprintf("获取到消息 [%s] %s", update.Message.From.UserName, update.Message.Text)))
			if update.Message.Text != "" {
				msg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath("./upload/ad.jpg"))
				msg.Caption = "Welcome to Hipoker, the most realistic Texas Hold game！👇"
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Play", "https://t.me/Hi_Poker_bot/TexasHold")),
					tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Join Player Group", "https://t.me/+ApOD71xvfew2YjEy")),
					tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Join Agent Group", "https://t.me/+HQrLcVrgz8YyNjYy")),
				)
				_, errSent := bot.Send(msg)
				if errSent != nil {
					logger.Info("机器人", zap.String("发送失败", errSent.Error()))
					//log.Error("tgbot send err: %s", errSent.Error())
				}
			}
			//else {
			//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//	msg.ReplyToMessageID = update.Message.MessageID
			//	bot.Send(msg)
			//}
		}
	}
}
