package main

import (
	"encoding/json"
)

type EmojiPack struct {
	Files    map[string]string `json:"files"`
	Metadata EmojiPackMetadata `json:"pack"`
}

type EmojiPackMetadata struct {
	License          string `json:"license,omitempty"`
	Homepage         string `json:"homepage,omitempty"`
	Description      string `json:"description,omitempty"`
	Fallback         string `json:"fallback-src,omitempty"`
	FallbackChecksum string `json:"fallback-src-sha256,omitempty"`
	AllowSharing     bool   `json:"share-files,omitempty"`
}

func NewEmojiPack() EmojiPack {
	e := EmojiPack{}
	e.Files = map[string]string{}
	e.Metadata = EmojiPackMetadata{}
	return e
}

func (p *EmojiPack) SetFiles(es Emojis) {
	for _, e := range es {
		p.Files[e.Shortcode] = e.Url
	}
}

func (p *EmojiPack) Json() ([]byte, error) {
	return json.Marshal(p)
}
