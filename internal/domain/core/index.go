package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type St struct {
	slackWebhookUrl string
	channel         string
	glLink          string
}

type MsgSt struct {
	Backlog []MsgBacklogSt         `json:"backlog"`
	Fields  map[string]interface{} `json:"fields"`
}

type MsgBacklogSt struct {
	Message string             `json:"message"`
	Fields  MsgBacklogFieldsSt `json:"fields"`
}

type MsgBacklogFieldsSt struct {
	ContainerName string `json:"container_name"`
	GlMessageId   string `json:"gl2_message_id"`
}

type SlackMsgSt struct {
	Channel  *string           `json:"channel,omitempty"`
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

func NewSt(slackWebhookUrl, channel, glLink string) *St {
	return &St{
		slackWebhookUrl: slackWebhookUrl,
		channel:         channel,
		glLink:          glLink,
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
		if bl.Fields.GlMessageId != "" {
			rows = append(rows, "       message_id: *"+bl.Fields.GlMessageId+"*")
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

		if len(msg.Fields) > 0 {
			rows = append(rows, fmt.Sprintf("    fields:"))
			for k, v := range msg.Fields {
				rows = append(rows, fmt.Sprintf("       %s: *%v*", k, v))
			}
		}

		if c.glLink != "" {
			rows = append(rows, "<"+c.glLink+"|GrayLog>")
		}
		texts = append(texts, strings.Join(rows, "\n"))
	}

	log.Println(string(msgBytes))

	if len(texts) > 0 {
		err = c.sendSlackText(texts)
		if err != nil {
			return err
		}
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

	if c.channel != "" {
		slackMsg.Channel = &c.channel
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
