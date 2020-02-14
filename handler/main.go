package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	tag     string
	webhook string
	channel string
)

func init() {
	tag = os.Getenv("NameTag")
	webhook = os.Getenv("Webhook")
	channel = os.Getenv("Channel")
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	if len(tag) < 1 {
		return fmt.Errorf("NameTag does not exist")
	}

	if len(webhook) < 1 {
		return fmt.Errorf("Webhook does not exist")
	}

	if len(channel) < 1 {
		return fmt.Errorf("Channel does not exist")
	}

	fmt.Printf("NameTag: %v, WebhookUrl: %v, Channel: %v", tag, webhook, channel)

	c := &Client{
		EC2API: ec2.New(session.Must(session.NewSession())),
	}

	instances, err := c.getInstances(tag)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Println(instances)

	// ToDo: goroutineにしたい
	var body string
	for _, i := range instances {
		id := i.InstanceId

		if i.State == "running" {
			if err := c.stopEc2(id); err != nil {
				return fmt.Errorf("%v", err)
			}

			body = id

			continue
		}

		if i.State == "stopped" {
			if err := c.startEc2(id); err != nil {
				return fmt.Errorf("%v", err)
			}

			body = id

			continue
		}

		body = body + ","
	}

	p := Payload{
		WebhookUrl: webhook,
		Channel:    channel,
		Body:       body,
	}

	if err = p.postSlack(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
