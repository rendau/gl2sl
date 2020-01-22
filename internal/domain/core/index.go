package core

import "fmt"

type St struct {
	slackWebhookUrl string
}

func NewSt(slackWebhookUrl string) *St {
	return &St{
		slackWebhookUrl: slackWebhookUrl,
	}
}

func (c *St) HandleMessage(msg []byte) error {
	fmt.Println(string(msg))
	return nil
}
