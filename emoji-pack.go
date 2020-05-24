package main

import (
	"encoding/json"

	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	p := EmojiPack{}
	p.Files = map[string]string{}
	p.Metadata = EmojiPackMetadata{}
	return p
}

func (p *EmojiPack) SetFiles(es Emojis) {
	for _, e := range es {
		p.Files[e.Shortcode] = e.Url
	}
}

func (p *EmojiPack) Json() ([]byte, error) {
	return json.Marshal(p)
}

type EmojiResult struct {
	Shortcode string
	File      string
	Error     error
}

func (p *EmojiPack) GenerateEmojiPack(outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = mkdir(outputDir); err != nil {
			return err
		}
	}

	c := make(chan EmojiResult, 30)
	limitCh := make(chan struct{}, 20)
	defer close(c)
	defer close(limitCh)

	for code, ru := range p.Files {
		go downloadEmojiFile(outputDir, code, ru, c, limitCh)
	}

	files := map[string]string{}
	for i := 0; i < len(p.Files); i++ {
		result := <-c
		if result.Error != nil {
			continue
		}
		files[result.Shortcode] = result.File
	}
	p.Files = files

	// generate pack.json
	f, err := os.Create(filepath.Join(outputDir, "pack.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := p.Json()
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}

func mkdir(name string) error {
	err := os.Mkdir(name, 0777)
	if os.IsNotExist(err) {
		mkdir(filepath.Dir(name))
		return os.Mkdir(name, 0777)
	}
	return err
}

func downloadEmojiFile(outputDir string, shortcode string, fileUrl string, result chan EmojiResult, limit chan struct{}) {
	limit <- struct{}{}

	r := EmojiResult{
		Shortcode: shortcode,
	}

	var u *url.URL
	u, r.Error = url.Parse(fileUrl)
	if r.Error != nil {
		<-limit
		result <- r
		return
	}

	var resp *http.Response
	resp, r.Error = http.Get(u.String())
	if r.Error != nil {
		<-limit
		result <- r
		return
	}
	defer resp.Body.Close()

	r.File = filepath.Base(u.String())

	var f *os.File
	f, r.Error = os.Create(filepath.Join(outputDir, r.File))
	if r.Error != nil {
		<-limit
		result <- r
		return
	}
	defer f.Close()

	_, r.Error = io.Copy(f, resp.Body)
	<-limit
	result <- r
}
