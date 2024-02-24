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

// OnText              = "\atext"
// OnPhoto             = "\aphoto"
// OnAudio             = "\aaudio"
// OnAnimation         = "\aanimation"
// OnDocument          = "\adocument"
// OnSticker           = "\asticker"
// OnVideo             = "\avideo"
// OnVoice             = "\avoice"
// OnVideoNote         = "\avideo_note"
func init() {
	features.RegisterFeature(tele.OnText, OnPrivateSendToMe)
	features.RegisterFeature(tele.OnPhoto, OnPrivateSendToMeByPhoto)
	features.RegisterFeature(tele.OnAudio, OnPrivateSendToMeByAudio)
	features.RegisterFeature(tele.OnDocument, OnPrivateSendToMeByDocument)
	features.RegisterFeature(tele.OnVideo, OnPrivateSendToMeByVideo)
	features.RegisterFeature(tele.OnVoice, OnPrivateSendToMeByVoice)
}

// OnText
func OnPrivateSendToMe(c tele.Context) error {
	if c.Message().FromChannel() || c.Message().FromGroup() {
		return nil
	}
	if !c.Message().Private() && !c.Message().IsReply() {
		return nil
	}
	msgId := ""
	if len(c.Message().AlbumID) > 0 {
		msgId = fmt.Sprintf("%d_%s", c.Message().Chat.ID, c.Message().AlbumID)
	} else { // c.Message().ID > 0
		msgId = msgId + fmt.Sprintf("%d_%d", c.Message().Chat.ID, c.Message().ID)
	}
	if _, ok := syncMap.LoadOrStore(msgId, ""); ok {
		return nil
	}

	adminID := os.Getenv("SENDTOME_ID")
	if len(adminID) == 0 {
		return nil
	}
	senderID := fmt.Sprintf("%d", c.Message().Sender.ID)
	// 管理员回复信息
	if c.Message().IsReply() && strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到回复消息(err)：", c.Message())
		} else {
			fmt.Println("收到回复消息：", string(jsonText))
		}
		prefixLine, _, isFound := strings.Cut(c.Message().ReplyTo.Text, "\n")
		if !isFound {
			return c.Reply("回复消息格式异常：" + c.Message().ReplyTo.Text)
		}
		_, sendToID, isFound := strings.Cut(prefixLine, "#id")
		if !isFound {
			return c.Reply("回复消息格式异常：" + c.Message().ReplyTo.Text)
		}

		reciverId, err := strconv.ParseInt(sendToID, 10, 64)
		if err != nil {
			fmt.Print("回复消息格式异常：待回复id %s\n%s", sendToID, c.Message().ReplyTo.Text)
		}
		reciver := &tele.User{
			ID: reciverId, //int64(reciverId),
		}

		if _, err := c.Bot().Send(reciver, c.Message().Text); err != nil {
			// return err
			return c.Reply("⚠️回复内容转投失败，请重试。" + err.Error())
		}
		// return nil
		return c.Reply("✅回复内容转投成功。")
	}
	// 收到私聊消息
	if c.Message().Private() && !strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到私聊消息(err)：", c.Message())
		} else {
			fmt.Println("收到私聊消息：", string(jsonText))
		}

		reciverId, err := strconv.ParseInt(adminID, 10, 64)
		if err != nil {
			fmt.Println("设置有误：环境变量(SENDTOME_ID)：", adminID)
		}
		reciver := &tele.User{
			ID: reciverId, //int64(reciverId),
		}

		newMsg := fmt.Sprintf("@%s #id%d\n%s",
			c.Message().Sender.Username,
			c.Message().Sender.ID,
			c.Message().Text)

		if _, err := c.Bot().Send(reciver, newMsg); err != nil {
			return err
			// return c.Reply("感谢！⚠️私聊内容转投失败。" + err.Error())
		}
		return nil
		// return c.Reply("感谢！✅私聊内容转投成功。")
	}
	what := fmt.Sprintf("@%s #id%d\n%s %s\nHi,Admin! 别逗了，宝！\n%s",
		c.Message().Sender.Username,
		c.Message().Sender.ID,
		c.Message().Sender.FirstName,
		c.Message().Sender.LastName,
		c.Message().Text)
	return c.Reply(what)
}

// OnPhoto
func OnPrivateSendToMeByPhoto(c tele.Context) error {
	if c.Message().FromChannel() || c.Message().FromGroup() {
		return nil
	}
	if !c.Message().Private() && !c.Message().IsReply() {
		return nil
	}
	msgId := ""
	if len(c.Message().AlbumID) > 0 {
		msgId = fmt.Sprintf("%d_%s", c.Message().Chat.ID, c.Message().AlbumID)
	} else { // c.Message().ID > 0
		msgId = msgId + fmt.Sprintf("%d_%d", c.Message().Chat.ID, c.Message().ID)
	}
	if _, ok := syncMap.LoadOrStore(msgId, ""); ok {
		return nil
	}

	adminID := os.Getenv("SENDTOME_ID")
	if len(adminID) == 0 {
		return nil
	}
	senderID := fmt.Sprintf("%d", c.Message().Sender.ID)
	// 管理员回复信息
	if c.Message().IsReply() && strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到回复消息(err)：", c.Message())
		} else {
			fmt.Println("收到回复消息：", string(jsonText))
		}
		prefixLine, _, isFound := strings.Cut(c.Message().ReplyTo.Text, "\n")
		if !isFound {
			return c.Reply("回复消息格式异常(prefixLine): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}
		_, sendToID, isFound := strings.Cut(prefixLine, "#id")
		if !isFound {
			return c.Reply("回复消息格式异常(sendToID): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}

		reciverId, err := strconv.ParseInt(sendToID, 10, 64)
		if err != nil {
			fmt.Print("回复消息格式异常：待回复id %s\n%+v", sendToID, c.Message().ReplyTo)
		}
		reciver := &tele.User{
			ID: reciverId,
		}

		if c.Message().Photo != nil {
			newMsg := c.Message().Photo
			// newMsg.Caption = c.Message().Photo.Caption
			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return c.Reply("⚠️回复内容转投失败，请重试。" + err.Error())
			}
			return c.Reply("✅回复内容转投成功。")
		}
		return c.Reply("⚠️回复内容转投失败，请重试。" + fmt.Sprintf("获取图片信息失败: %+v", c.Message()))
	}

	// 收到私聊消息
	if c.Message().Private() && !strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到私聊消息(err)：", c.Message())
		} else {
			fmt.Println("收到私聊消息：", string(jsonText))
		}

		reciverId, err := strconv.ParseInt(adminID, 10, 64)
		if err != nil {
			fmt.Println("设置有误：环境变量(SENDTOME_ID)：", adminID)
		}
		reciver := &tele.User{
			ID: reciverId,
		}
		if c.Message().Photo != nil {
			newMsg := c.Message().Photo
			msgText := fmt.Sprintf("@%s #id%d\n%s",
				c.Message().Sender.Username,
				c.Message().Sender.ID,
				c.Message().Caption)
			newMsg.Caption = msgText

			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return err
				// return c.Reply("感谢！⚠️私聊内容转投失败。" + err.Error())
			}
			return nil
			// return c.Reply("感谢！✅私聊内容转投成功。")
		}
		return fmt.Errorf("获取图片信息失败")
	}

	what := fmt.Sprintf("@%s #id%d\n%s %s\nHi,Admin! 别逗了，宝！\n%s",
		c.Message().Sender.Username,
		c.Message().Sender.ID,
		c.Message().Sender.FirstName,
		c.Message().Sender.LastName,
		c.Message().Caption)
	return c.Reply(what)
}

// OnAudio
func OnPrivateSendToMeByAudio(c tele.Context) error {
	if c.Message().FromChannel() || c.Message().FromGroup() {
		return nil
	}
	if !c.Message().Private() && !c.Message().IsReply() {
		return nil
	}
	msgId := ""
	if len(c.Message().AlbumID) > 0 {
		msgId = fmt.Sprintf("%d_%s", c.Message().Chat.ID, c.Message().AlbumID)
	} else { // c.Message().ID > 0
		msgId = msgId + fmt.Sprintf("%d_%d", c.Message().Chat.ID, c.Message().ID)
	}
	if _, ok := syncMap.LoadOrStore(msgId, ""); ok {
		return nil
	}

	adminID := os.Getenv("SENDTOME_ID")
	if len(adminID) == 0 {
		return nil
	}
	senderID := fmt.Sprintf("%d", c.Message().Sender.ID)
	// 管理员回复信息
	if c.Message().IsReply() && strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到回复消息(err)：", c.Message())
		} else {
			fmt.Println("收到回复消息：", string(jsonText))
		}
		prefixLine, _, isFound := strings.Cut(c.Message().ReplyTo.Text, "\n")
		if !isFound {
			return c.Reply("回复消息格式异常(prefixLine): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}
		_, sendToID, isFound := strings.Cut(prefixLine, "#id")
		if !isFound {
			return c.Reply("回复消息格式异常(sendToID): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}

		reciverId, err := strconv.ParseInt(sendToID, 10, 64)
		if err != nil {
			fmt.Print("回复消息格式异常：待回复id %s\n%+v", sendToID, c.Message().ReplyTo)
		}
		reciver := &tele.User{
			ID: reciverId,
		}

		if c.Message().Audio != nil {
			newMsg := c.Message().Audio
			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return c.Reply("⚠️回复内容转投失败，请重试。" + err.Error())
			}
			return c.Reply("✅回复内容转投成功。")
		}
		return c.Reply("⚠️回复内容转投失败，请重试。" + fmt.Sprintf("获取图片信息失败: %+v", c.Message()))
	}

	// 收到私聊消息
	if c.Message().Private() && !strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到私聊消息(err)：", c.Message())
		} else {
			fmt.Println("收到私聊消息：", string(jsonText))
		}

		reciverId, err := strconv.ParseInt(adminID, 10, 64)
		if err != nil {
			fmt.Println("设置有误：环境变量(SENDTOME_ID)：", adminID)
		}
		reciver := &tele.User{
			ID: reciverId,
		}
		if c.Message().Audio != nil {
			newMsg := c.Message().Audio
			msgText := fmt.Sprintf("@%s #id%d\n%s",
				c.Message().Sender.Username,
				c.Message().Sender.ID,
				c.Message().Caption)
			newMsg.Caption = msgText

			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return err
				// return c.Reply("感谢！⚠️私聊内容转投失败。" + err.Error())
			}
			return nil
			// return c.Reply("感谢！✅私聊内容转投成功。")
		}
		return fmt.Errorf("获取图片信息失败")
	}

	what := fmt.Sprintf("@%s #id%d\n%s %s\nHi,Admin! 别逗了，宝！\n%s",
		c.Message().Sender.Username,
		c.Message().Sender.ID,
		c.Message().Sender.FirstName,
		c.Message().Sender.LastName,
		c.Message().Caption)
	return c.Reply(what)
}

// OnDocument
func OnPrivateSendToMeByDocument(c tele.Context) error {
	if c.Message().FromChannel() || c.Message().FromGroup() {
		return nil
	}
	if !c.Message().Private() && !c.Message().IsReply() {
		return nil
	}
	msgId := ""
	if len(c.Message().AlbumID) > 0 {
		msgId = fmt.Sprintf("%d_%s", c.Message().Chat.ID, c.Message().AlbumID)
	} else { // c.Message().ID > 0
		msgId = msgId + fmt.Sprintf("%d_%d", c.Message().Chat.ID, c.Message().ID)
	}
	if _, ok := syncMap.LoadOrStore(msgId, ""); ok {
		return nil
	}

	adminID := os.Getenv("SENDTOME_ID")
	if len(adminID) == 0 {
		return nil
	}
	senderID := fmt.Sprintf("%d", c.Message().Sender.ID)
	// 管理员回复信息
	if c.Message().IsReply() && strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到回复消息(err)：", c.Message())
		} else {
			fmt.Println("收到回复消息：", string(jsonText))
		}
		prefixLine, _, isFound := strings.Cut(c.Message().ReplyTo.Text, "\n")
		if !isFound {
			return c.Reply("回复消息格式异常(prefixLine): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}
		_, sendToID, isFound := strings.Cut(prefixLine, "#id")
		if !isFound {
			return c.Reply("回复消息格式异常(sendToID): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}

		reciverId, err := strconv.ParseInt(sendToID, 10, 64)
		if err != nil {
			fmt.Print("回复消息格式异常：待回复id %s\n%+v", sendToID, c.Message().ReplyTo)
		}
		reciver := &tele.User{
			ID: reciverId,
		}

		if c.Message().Document != nil {
			newMsg := c.Message().Document
			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return c.Reply("⚠️回复内容转投失败，请重试。" + err.Error())
			}
			return c.Reply("✅回复内容转投成功。")
		}
		return c.Reply("⚠️回复内容转投失败，请重试。" + fmt.Sprintf("获取图片信息失败: %+v", c.Message()))
	}

	// 收到私聊消息
	if c.Message().Private() && !strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到私聊消息(err)：", c.Message())
		} else {
			fmt.Println("收到私聊消息：", string(jsonText))
		}

		reciverId, err := strconv.ParseInt(adminID, 10, 64)
		if err != nil {
			fmt.Println("设置有误：环境变量(SENDTOME_ID)：", adminID)
		}
		reciver := &tele.User{
			ID: reciverId,
		}
		if c.Message().Document != nil {
			newMsg := c.Message().Document
			msgText := fmt.Sprintf("@%s #id%d\n%s",
				c.Message().Sender.Username,
				c.Message().Sender.ID,
				c.Message().Caption)
			newMsg.Caption = msgText

			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return err
				// return c.Reply("感谢！⚠️私聊内容转投失败。" + err.Error())
			}
			return nil
			// return c.Reply("感谢！✅私聊内容转投成功。")
		}
		return fmt.Errorf("获取图片信息失败")
	}

	what := fmt.Sprintf("@%s #id%d\n%s %s\nHi,Admin! 别逗了，宝！\n%s",
		c.Message().Sender.Username,
		c.Message().Sender.ID,
		c.Message().Sender.FirstName,
		c.Message().Sender.LastName,
		c.Message().Caption)
	return c.Reply(what)
}

// OnVideo
func OnPrivateSendToMeByVideo(c tele.Context) error {
	if c.Message().FromChannel() || c.Message().FromGroup() {
		return nil
	}
	if !c.Message().Private() && !c.Message().IsReply() {
		return nil
	}
	msgId := ""
	if len(c.Message().AlbumID) > 0 {
		msgId = fmt.Sprintf("%d_%s", c.Message().Chat.ID, c.Message().AlbumID)
	} else { // c.Message().ID > 0
		msgId = msgId + fmt.Sprintf("%d_%d", c.Message().Chat.ID, c.Message().ID)
	}
	if _, ok := syncMap.LoadOrStore(msgId, ""); ok {
		return nil
	}

	adminID := os.Getenv("SENDTOME_ID")
	if len(adminID) == 0 {
		return nil
	}
	senderID := fmt.Sprintf("%d", c.Message().Sender.ID)
	// 管理员回复信息
	if c.Message().IsReply() && strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到回复消息(err)：", c.Message())
		} else {
			fmt.Println("收到回复消息：", string(jsonText))
		}
		prefixLine, _, isFound := strings.Cut(c.Message().ReplyTo.Text, "\n")
		if !isFound {
			return c.Reply("回复消息格式异常(prefixLine): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}
		_, sendToID, isFound := strings.Cut(prefixLine, "#id")
		if !isFound {
			return c.Reply("回复消息格式异常(sendToID): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}

		reciverId, err := strconv.ParseInt(sendToID, 10, 64)
		if err != nil {
			fmt.Print("回复消息格式异常：待回复id %s\n%+v", sendToID, c.Message().ReplyTo)
		}
		reciver := &tele.User{
			ID: reciverId,
		}

		if c.Message().Video != nil {
			newMsg := c.Message().Video
			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return c.Reply("⚠️回复内容转投失败，请重试。" + err.Error())
			}
			return c.Reply("✅回复内容转投成功。")
		}
		return c.Reply("⚠️回复内容转投失败，请重试。" + fmt.Sprintf("获取图片信息失败: %+v", c.Message()))
	}

	// 收到私聊消息
	if c.Message().Private() && !strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到私聊消息(err)：", c.Message())
		} else {
			fmt.Println("收到私聊消息：", string(jsonText))
		}

		reciverId, err := strconv.ParseInt(adminID, 10, 64)
		if err != nil {
			fmt.Println("设置有误：环境变量(SENDTOME_ID)：", adminID)
		}
		reciver := &tele.User{
			ID: reciverId,
		}
		if c.Message().Video != nil {
			newMsg := c.Message().Video
			msgText := fmt.Sprintf("@%s #id%d\n%s",
				c.Message().Sender.Username,
				c.Message().Sender.ID,
				c.Message().Caption)
			newMsg.Caption = msgText

			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return err
				// return c.Reply("感谢！⚠️私聊内容转投失败。" + err.Error())
			}
			return nil
			// return c.Reply("感谢！✅私聊内容转投成功。")
		}
		return fmt.Errorf("获取图片信息失败")
	}

	what := fmt.Sprintf("@%s #id%d\n%s %s\nHi,Admin! 别逗了，宝！\n%s",
		c.Message().Sender.Username,
		c.Message().Sender.ID,
		c.Message().Sender.FirstName,
		c.Message().Sender.LastName,
		c.Message().Caption)
	return c.Reply(what)
}

// OnVoice
func OnPrivateSendToMeByVoice(c tele.Context) error {
	if c.Message().FromChannel() || c.Message().FromGroup() {
		return nil
	}
	if !c.Message().Private() && !c.Message().IsReply() {
		return nil
	}
	msgId := ""
	if len(c.Message().AlbumID) > 0 {
		msgId = fmt.Sprintf("%d_%s", c.Message().Chat.ID, c.Message().AlbumID)
	} else { // c.Message().ID > 0
		msgId = msgId + fmt.Sprintf("%d_%d", c.Message().Chat.ID, c.Message().ID)
	}
	if _, ok := syncMap.LoadOrStore(msgId, ""); ok {
		return nil
	}

	adminID := os.Getenv("SENDTOME_ID")
	if len(adminID) == 0 {
		return nil
	}
	senderID := fmt.Sprintf("%d", c.Message().Sender.ID)
	// 管理员回复信息
	if c.Message().IsReply() && strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到回复消息(err)：", c.Message())
		} else {
			fmt.Println("收到回复消息：", string(jsonText))
		}
		prefixLine, _, isFound := strings.Cut(c.Message().ReplyTo.Text, "\n")
		if !isFound {
			return c.Reply("回复消息格式异常(prefixLine): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}
		_, sendToID, isFound := strings.Cut(prefixLine, "#id")
		if !isFound {
			return c.Reply("回复消息格式异常(sendToID): " + fmt.Sprintf("%+v", c.Message().ReplyTo))
		}

		reciverId, err := strconv.ParseInt(sendToID, 10, 64)
		if err != nil {
			fmt.Print("回复消息格式异常：待回复id %s\n%+v", sendToID, c.Message().ReplyTo)
		}
		reciver := &tele.User{
			ID: reciverId,
		}

		if c.Message().Voice != nil {
			newMsg := c.Message().Voice
			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return c.Reply("⚠️回复内容转投失败，请重试。" + err.Error())
			}
			return c.Reply("✅回复内容转投成功。")
		}
		return c.Reply("⚠️回复内容转投失败，请重试。" + fmt.Sprintf("获取图片信息失败: %+v", c.Message()))
	}

	// 收到私聊消息
	if c.Message().Private() && !strings.EqualFold(senderID, adminID) {
		if jsonText, err := json.Marshal(c.Message()); err != nil {
			fmt.Println("收到私聊消息(err)：", c.Message())
		} else {
			fmt.Println("收到私聊消息：", string(jsonText))
		}

		reciverId, err := strconv.ParseInt(adminID, 10, 64)
		if err != nil {
			fmt.Println("设置有误：环境变量(SENDTOME_ID)：", adminID)
		}
		reciver := &tele.User{
			ID: reciverId,
		}
		if c.Message().Voice != nil {
			newMsg := c.Message().Voice
			msgText := fmt.Sprintf("@%s #id%d\n%s",
				c.Message().Sender.Username,
				c.Message().Sender.ID,
				c.Message().Caption)
			newMsg.Caption = msgText

			if _, err := c.Bot().Send(reciver, newMsg); err != nil {
				return err
				// return c.Reply("感谢！⚠️私聊内容转投失败。" + err.Error())
			}
			return nil
			// return c.Reply("感谢！✅私聊内容转投成功。")
		}
		return fmt.Errorf("获取图片信息失败")
	}

	what := fmt.Sprintf("@%s #id%d\n%s %s\nHi,Admin! 别逗了，宝！\n%s",
		c.Message().Sender.Username,
		c.Message().Sender.ID,
		c.Message().Sender.FirstName,
		c.Message().Sender.LastName,
		c.Message().Caption)
	return c.Reply(what)
}
