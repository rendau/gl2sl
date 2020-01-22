package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type St struct {
	slackWebhookUrl string
}

type MsgSt struct {
	Backlog []MsgBacklogSt `json:"backlog"`
}

type MsgBacklogSt struct {
	Message string             `json:"message"`
	Fields  MsgBacklogFieldsSt `json:"fields"`
}

type MsgBacklogFieldsSt struct {
	ContainerName string `json:"container_name"`
}

type SlackMsgSt struct {
	Username string            `json:"username"`
	Text     string            `json:"text"`
	Blocks   []SlackMsgBlockSt `json:"blocks"`
}

type SlackMsgBlockSt struct {
	Type string              `json:"type"`
	Text SlackMsgBlockTextSt `json:"text"`
}

type SlackMsgBlockTextSt struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewSt(slackWebhookUrl string) *St {
	return &St{
		slackWebhookUrl: slackWebhookUrl,
	}
}

func (c *St) HandleMessage(msgBytes []byte) error {
	msg := MsgSt{}
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		return err
	}

	texts := make([]string, 0, len(msg.Backlog))
	var rows []string
	for _, bl := range msg.Backlog {
		rows = nil
		if bl.Fields.ContainerName != "" {
			rows = append(rows, "*"+bl.Fields.ContainerName+"*:")
		}
		blMsg := map[string]interface{}{}
		err = json.Unmarshal([]byte(bl.Message), &blMsg)
		if err == nil {
			for k, v := range blMsg {
				rows = append(rows, fmt.Sprintf("       %s: *%v*", k, v))
			}
		} else {
			rows = append(rows, fmt.Sprintf("       message: *%s*", bl.Message))
		}
		texts = append(texts, strings.Join(rows, "\n"))
	}

	err = c.sendSlackText(texts)
	if err != nil {
		return err
	}

	return nil
}

func (c *St) sendSlackText(texts []string) error {
	if len(texts) == 0 {
		return nil
	}

	slackMsg := SlackMsgSt{
		Username: "GrayLog",
		Text:     "Message from GrayLog",
	}

	for _, text := range texts {
		slackMsg.Blocks = append(slackMsg.Blocks, SlackMsgBlockSt{
			Type: "section",
			Text: SlackMsgBlockTextSt{
				Type: "mrkdwn",
				Text: text,
			},
		})
	}

	slackMsgBytes, err := json.Marshal(slackMsg)
	if err != nil {
		return err
	}

	httpClient := http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("POST", c.slackWebhookUrl, bytes.NewBuffer(slackMsgBytes))
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("bad status code from slack")
	}

	return nil
}
