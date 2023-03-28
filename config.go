package aiposter

import (
	"fmt"
	"os"
)

const (
	openAiTokenVar    = "OPENAI_TOKEN"
	botTokenVar       = "TELEGRAM_BOT_TOKEN"
	channelIdVar      = "TELEGRAM_CHANNEL_ID"
	dictionaryFileVar = "DICTIONARY_FILE"
	dialogueFileVar   = "DIALOGUE_FILE"
	bucketNameVar     = "BUCKET_NAME"
)

var Conf Config

type Config struct {
	// Telegram bot token
	BotToken string
	// Telegram channel id
	ChannelId string
	// OpenAi API token
	OpenAiToken string
	// Dictionary file name in Google Cloud Storage
	DictionaryFile string
	// Dialogue file name in Google Cloud Storage
	DialogueFile string
	// Google Cloud Storage bucket name
	BucketName string
}

func InitConfig() error {
	Conf.BotToken = os.Getenv(botTokenVar)
	Conf.OpenAiToken = os.Getenv(openAiTokenVar)
	Conf.ChannelId = os.Getenv(channelIdVar)
	Conf.DictionaryFile = os.Getenv(dictionaryFileVar)
	Conf.DialogueFile = os.Getenv(dialogueFileVar)
	Conf.BucketName = os.Getenv(bucketNameVar)

	if Conf.BotToken == "" {
		return fmt.Errorf("bot token is not set")
	}

	if Conf.OpenAiToken == "" {
		return fmt.Errorf("open ai token is not set")
	}

	if Conf.ChannelId == "" {
		return fmt.Errorf("channel id is not set")
	}

	if Conf.DictionaryFile == "" {
		return fmt.Errorf("dictionary file is not set")
	}

	if Conf.DialogueFile == "" {
		return fmt.Errorf("dialogue file is not set")
	}

	if Conf.BucketName == "" {
		return fmt.Errorf("bucket name is not set")
	}

	return nil
}
