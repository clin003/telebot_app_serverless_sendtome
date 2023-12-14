package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/clin003/tgbot_app_dev/common"
	"github.com/clin003/tgbot_app_dev/features"
	_ "github.com/clin003/tgbot_app_dev/main/distro/all"
	"github.com/clin003/tgbot_app_dev/utils"

	tele "gopkg.in/telebot.v3"
)

var (
	bot *tele.Bot
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// fmt.Fprint(w, "Hi!")
		return
	}
	body, err := io.ReadAll(r.Body)
	common.Must(err)
	log.Println(string(body))

	var u tele.Update
	common.Must(json.Unmarshal(body, &u))

	bot.ProcessUpdate(u)
}

func init() {
	var err error
	botToken := os.Getenv("BAICAI_BOT_TELEGRAM_TOKEN")
	bot, err = tele.NewBot(tele.Settings{
		Token:       botToken,
		Synchronous: true,
	})
	common.Must(err)

	commands := []tele.Command{
		{
			Text:        "/id",
			Description: "查询自己的用户id信息",
		},
		{
			Text:        "/info",
			Description: "查询公开群组频道信息",
		},
		{
			Text:        "/ping",
			Description: "Ping",
		},
		{
			Text:        "/about",
			Description: "About",
		},
		{
			Text:        "/start",
			Description: "Start",
		},
		// {
		// 	Text:        "/sendCrypto",
		// 	Description: "Send crypto (发送加密货币)",
		// },
		// {
		// 	Text:        "/sendCryptoUSDT",
		// 	Description: "Send crypto USDT (发送加密货币 USDT)",
		// },
	}

	if len(os.Getenv("SEND_CRYPTO_MSG")) > 0 {
		commands = append(commands, tele.Command{
			Text:        "/sendcrypto",
			Description: "Send crypto (发送加密货币)",
		})
	}
	if len(os.Getenv("SEND_CRYPTO_USDT_MSG")) > 0 {
		commands = append(commands, tele.Command{
			Text:        "/sendUSDT",
			Description: "Send crypto USDT (发送加密货币 USDT)",
		})
	}

	bot.SetCommands(commands)
	webhookURL := os.Getenv("BAICAI_BOT_TELEGRAM_WEBHOOK_URL")
	if len(webhookURL) > 0 && strings.HasPrefix(webhookURL, "https") {
		utils.SetTelegramWebhook(botToken, webhookURL)
	}
	features.Handle(bot)
}
