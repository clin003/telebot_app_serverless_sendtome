package sendtome

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/clin003/tgbot_app_dev/features"
	tele "gopkg.in/telebot.v3"
)

// var sendToMeID string
var syncMap sync.Map

func init() {
	features.RegisterFeature(tele.OnText, OnPrivateSendToMe)
	// features.RegisterFeature(tele.OnPhoto, OnChannelLinkGroup)
	// features.RegisterFeature(tele.OnVideo, OnChannelLinkGroup)
	// features.RegisterFeature(tele.OnMedia, OnChannelLinkGroup)

	// sendToMeID = os.Getenv("SENDTOME_ID")
}

func OnPrivateSendToMe(c tele.Context) error {
	// if len(sendToMeID) <= 0 {
	// 	return nil
	// }
	if c.Message().FromChannel() || c.Message().FromGroup() {
		return nil
	}

	if !c.Message().Private() &&
		!c.Message().IsReply() {
		return nil
	}
	// // fmt.Println("OnPrivateSendToMe", 1)
	// if !(c.Message().OriginalChat != nil) || !(c.Message().SenderChat != nil) {
	// 	return nil
	// }
	// // fmt.Println("OnPrivateSendToMe", 2)
	// if c.Message().OriginalChat.Type != tele.ChatPrivate ||
	// 	c.Message().SenderChat.Type != tele.ChatPrivate {
	// 	return nil
	// }
	msgId := ""
	if len(c.Message().AlbumID) > 0 {
		msgId = fmt.Sprintf("%d_%s", c.Message().Chat.ID, c.Message().AlbumID)
	} else { // c.Message().ID > 0
		msgId = msgId + fmt.Sprintf("%d_%d", c.Message().Chat.ID, c.Message().ID)
	}

	if _, ok := syncMap.LoadOrStore(msgId, ""); ok {
		return nil
	}
	sendToMeID := os.Getenv("SENDTOME_ID")
	if len(sendToMeID) == 0 {
		return nil
	}
	// 管理员回复信息
	fmt.Println("c.Message().Sender.ID", c.Message().Sender.ID, string(c.Message().Sender.ID), c.Message().IsReply())
	fmt.Println("sendToMeID", sendToMeID, c.Message().IsReply())
	if c.Message().IsReply() && strings.EqualFold(fmt.Sprintf("%d", c.Message().Sender.ID), sendToMeID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到回复消息(err)：", c.Message())
		} else {
			fmt.Println("收到回复消息：", string(jsonText))
		}

		return nil
	}
	// 收到私聊消息
	if c.Message().Private() && !strings.EqualFold(string(c.Message().Sender.ID), sendToMeID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到私聊消息(err)：", c.Message())
		} else {
			fmt.Println("收到私聊消息：", string(jsonText))
		}

		reciverId, err := strconv.ParseInt(sendToMeID, 10, 64)
		if err != nil {
			fmt.Println("设置有误：环境变量(SENDTOME_ID)：", sendToMeID)
		}
		reciver := &tele.User{
			ID: reciverId, //int64(reciverId),
		}

		newMsg := fmt.Sprintf("@%s:%d\n%s",
			c.Message().Sender.Username,
			c.Message().Sender.ID,
			c.Message().Text)

		if _, err := c.Bot().Send(reciver, newMsg); err != nil {
			return err
		}
		return nil
	}
	what := fmt.Sprintf("%s %s\n@%s:%d\nIs Owner\n%s",
		c.Message().Sender.FirstName,
		c.Message().Sender.LastName,
		c.Message().Sender.Username,
		c.Message().Sender.ID,
		c.Message().Text)
	return c.Reply(what)
	// if c.Message().IsReply() {
	// 	if jsonText, err := json.Marshal(c.Message()); err != nil {
	// 		fmt.Println("收到回复消息(err)：", c.Message())
	// 	} else {
	// 		fmt.Println("收到回复消息：", string(jsonText))
	// 	}
	// }
	// return nil
}
