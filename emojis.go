package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Emojis []Emoji

type Emoji struct {
	Shortcode       string `json:"shortcode"`
	Url             string `json:"url"`
	StaticUrl       string `json:"static_url"`
	VisibleInPicker bool   `json:"visible_in_picker"`
	Category        string `json:"category"`
}

func NewEmojiList(domain string) (Emojis, error) {
	emojis := Emojis{}

	bytes, err := fetchCustomEmojis(domain)
	if err != nil {
		return emojis, err
	}

	if err := json.Unmarshal(bytes, &emojis); err != nil {
		return emojis, err
	}

	return emojis, nil
}

func fetchCustomEmojis(domain string) ([]byte, error) {
	u, err := url.Parse("https://example.com/api/v1/custom_emojis")
	if err != nil {
		return nil, err
	}
	u.Host = domain

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
