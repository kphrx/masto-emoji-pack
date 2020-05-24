package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Result struct {
	Server string
	Emojis Emojis
	Error  error
}

func main() {
	opts := parseOptions()

	c := make(chan Result, 5)
	defer close(c)

	for _, d := range opts.Servers {
		go saveEmojiList(d, opts, c)
	}

	for _, _ = range opts.Servers {
		if r := <-c; r.Error != nil {
			fmt.Println(r.Error)
		} else {
			fmt.Printf("Success: %s, emoji count: %d\n", r.Server, len(r.Emojis))
			//for _, e := range r.Emojis {
			//    fmt.Printf("%s, %s, %s\n", e.Shortcode, e.Url, e.Category)
			//}
		}
	}
}

func saveEmojiList(domain string, opts Options, c chan Result) {
	r := Result{
		Server: domain,
	}
	r.Emojis, r.Error = NewEmojiList(domain)
	if r.Error != nil {
		c <- r
		return
	}

	p := NewEmojiPack()
	p.SetFiles(r.Emojis)

	var b []byte
	b, r.Error = p.Json()
	if r.Error != nil {
		c <- r
		return
	}

	// opts.OutputDir
	out := filepath.Join(opts.OutputDir, domain)

	if _, err := os.Stat(out); os.IsNotExist(err) {
		os.Mkdir(out, 0777)
	}
	var f *os.File
	f, r.Error = os.Create(filepath.Join(out, "pack.json"))
	if r.Error != nil {
		c <- r
		return
	}

	_, r.Error = f.Write(b)
	if r.Error != nil {
		c <- r
		return
	}

	c <- r
}
