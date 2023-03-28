package aiposter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
)

type Dictionary struct {
	ThemeList []struct {
		Key    string   `json:"key"`
		Themes []string `json:"themes"`
	} `json:"list"`
}

func GetThemeByKey(source io.Reader, key string) (theme string, err error) {
	var dictionary Dictionary
	if err = json.NewDecoder(source).Decode(&dictionary); err != nil {
		return theme, fmt.Errorf("failed to unmarshal dictionary: %w", err)
	}

	var list []string
	for _, themeList := range dictionary.ThemeList {
		if themeList.Key == key {
			if len(themeList.Themes) == 0 {
				return theme, errors.New("empty theme list")
			}

			list = themeList.Themes
			break
		}
	}

	if len(list) == 0 {
		return theme, errors.New("key not found")
	}

	theme = list[rand.Intn(len(list))]
	return
}
