package notifications

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/nlopes/slack"
)

type NotificationMessage struct {
	Message           string
	ResourceType      string
	ResourceNamespace string
}

func NotifyViolation(message NotificationMessage, botToken, channelString, color string) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	api := slack.New(botToken)
	attachment := slack.Attachment{
		Fallback: "A rule violation has occurred.",
		Pretext:  "The following rules violations has been occurred:",
		Text:     message.Message,
		Footer:   "Aegir",
		Ts:       json.Number(ts),
		Color:    color,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "ResourceType",
				Value: message.ResourceType,
			},
			slack.AttachmentField{
				Title: "Namespace",
				Value: message.ResourceNamespace,
			},
		},
	}

	_, timestamp, err := api.PostMessage(channelString, slack.MsgOptionText("", true), slack.MsgOptionAttachments(attachment))
	if channelString == "" || botToken == "" {
		return
	}
	if err != nil {
		log.Printf("Error sending slack notification: %s\n", err)
		return
	}
	log.Printf("Message successfully sent to channel %s at %s\n", channelString, timestamp)
}
