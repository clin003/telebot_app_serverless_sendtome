package sendtome

import (
	"encoding/json"
	"fmt"
	"os"
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

	if c.Message().Private() {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到私聊消息：", c.Message())
		} else {
			fmt.Println("收到私聊消息：", c.Message())
		}

		sendToMeID = os.Getenv("SENDTOME_ID")
		if len(sendToMeID) <= 0 {
			return nil
		}
		reciverId, err := strconv.ParseInt(sendToMeID, 10, 64)
		reciver := &tele.User{
			ID: reciverId, //int64(reciverId),
		}

		return c.ForwardTo(reciver)
	}
	if c.Message().IsReply() {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到回复消息：", c.Message())
		} else {
			fmt.Println("收到回复消息：", c.Message())
		}

	}
	return nil
}
