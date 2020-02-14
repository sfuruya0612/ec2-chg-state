package main

import (
	"fmt"

	slack "github.com/ashwanthkumar/slack-go-webhook"
)

// Payload is slack payload
type Payload struct {
	WebhookUrl string
	Channel    string
	Body       string
}

func (p *Payload) postSlack() error {
	field1 := slack.Field{Title: "Target Instances", Value: p.Body}

	attachment := slack.Attachment{}
	attachment.AddField(field1)

	color := "good"
	attachment.Color = &color

	payload := slack.Payload{
		Username:    "Lambda",
		IconEmoji:   ":lambda:",
		Channel:     p.Channel,
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.Send(p.WebhookUrl, "", payload)
	if err != nil {
		return fmt.Errorf("Send slack: %v", err)
	}

	return nil
}
