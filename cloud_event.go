package aiposter

import (
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
)

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

const (
	attributeTypeKey = "type"
)

func getTaskFromEvent(e event.Event) (t Task, err error) {
	var msg MessagePublishedData
	if err = e.DataAs(&msg); err != nil {
		return t, fmt.Errorf("event.DataAs: %v", err)
	}

	if t.Theme = string(msg.Message.Data); t.Theme == "" {
		return t, errors.New("empty message")
	}

	t.Type = messageTypeMsg
	if mType, ok := msg.Message.Attributes[attributeTypeKey]; ok {
		t.Type = TaskType(mType)
	}

	return
}
