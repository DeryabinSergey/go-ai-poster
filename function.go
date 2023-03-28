package aiposter

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	functions.CloudEvent("EventHandler", EventHandler)
}

const (
	messageTypeKey    = "key"
	messageTypeMsg    = "msg"
	dialogueLength    = 7
	chatSystemMessage = "You are the content manager of a telegram group with the theme of IT developers."
)

type TaskType string
type Task struct {
	Theme string
	Type  TaskType
}

type ChatMessages struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
}

func EventHandler(ctx context.Context, e event.Event) error {
	if err := InitConfig(); err != nil {
		return fmt.Errorf("init config: %v", err)
	}

	task, err := getTaskFromEvent(e)
	if err != nil {
		return fmt.Errorf("get task from event: %v", err)
	}

	stor, err := NewStorage(ctx, Conf.BucketName)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}

	if task.Type == messageTypeKey {
		reader, err := stor.GetReader(ctx, Conf.DictionaryFile)
		if err != nil {
			return fmt.Errorf("dictonary read error: %w", err)
		}
		defer func(reader *storage.Reader) {
			err := reader.Close()
			if err != nil {
				log.Printf("Failed to close dictionary reader: %v", err)
			}
		}(reader)
		if task.Theme, err = GetThemeByKey(reader, task.Theme); err != nil {
			return fmt.Errorf("failed to get theme: %w", err)
		}
	}

	var dialog ChatMessages
	{
		reader, err := stor.GetReader(ctx, Conf.DialogueFile)
		if err != nil {
			if err == storage.ErrObjectNotExist {
				dialog.Messages = make([]openai.ChatCompletionMessage, 0)
				dialog.Messages = append(dialog.Messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleSystem,
					Content: chatSystemMessage,
				})
			} else {
				return fmt.Errorf("fetching dialogue: %w", err)
			}
		} else {
			defer func(reader *storage.Reader) {
				err := reader.Close()
				if err != nil {
					log.Printf("Failed to close dialogue reader: %v", err)
				}
			}(reader)
			if err = json.NewDecoder(reader).Decode(&dialog); err != nil {
				return fmt.Errorf("failed to unmarshal dialogue: %w", err)
			}
			err = reader.Close()
			if err != nil {
				log.Printf("Failed to close dialogue reader: %v", err)
			}

			if len(dialog.Messages) > dialogueLength {
				lastMessages := dialog.Messages[len(dialog.Messages)-dialogueLength+1:]
				for i, m := range lastMessages {
					dialog.Messages[i+1] = m
				}
				dialog.Messages = dialog.Messages[:dialogueLength]
			}
		}
	}

	dialog.Messages = append(dialog.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: task.Theme,
	})

	{
		client := openai.NewClient(Conf.OpenAiToken)
		resp, err := client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: dialog.Messages,
			},
		)
		if err != nil {
			return fmt.Errorf("ChatCompletion error: %w", err)
		}

		content := cleanString(resp.Choices[0].Message.Content)
		dialog.Messages = append(dialog.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})

		if err = postMessage(content); err != nil {
			return fmt.Errorf("failed to post message: %w", err)
		}
	}

	rawDialogue, err := json.Marshal(dialog)
	if err != nil {
		return fmt.Errorf("failed to marshal dialogue: %w", err)
	}
	writer := stor.GetWriter(ctx, Conf.DialogueFile)
	defer func(writer *storage.Writer) {
		err := writer.Close()
		if err != nil {
			log.Printf("Failed to close dialogue writer: %v", err)
		}
	}(writer)
	writer.ContentType = "application/json"
	if _, err = writer.Write(rawDialogue); err != nil {
		return fmt.Errorf("failed to write dialogue: %w", err)
	}

	return nil
}

// cleanString removes the bot's phrases from the response, if exist
func cleanString(str string) string {
	// a list of bot phrases that are known at the moment
	listBotPhrases := []string{
		"sure,",
		"here's",
		"here is",
	}

	list := strings.SplitN(str, "\n", 2)
	if len(list) < 2 {
		return str
	}

	firstLine := strings.ToLower(strings.TrimSpace(list[0]))
	for _, botPhrase := range listBotPhrases {
		if strings.HasPrefix(firstLine, botPhrase) {
			return strings.TrimSpace(list[1])
		}
	}

	return str
}

// postMessage sends a message to the telegram channel.
func postMessage(message string) error {
	response, err := http.PostForm("https://api.telegram.org/bot"+Conf.BotToken+"/sendMessage", url.Values{
		"chat_id":    {Conf.ChannelId},
		"text":       {message},
		"parse_mode": {"markdown"},
	})
	if err != nil {
		return err
	}

	if err := response.Body.Close(); err != nil {
		log.Printf("failed to close telegram response: %v\n", err)
	}

	return nil
}
