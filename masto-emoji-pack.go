package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Result struct {
	Server string
	Output []string
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
			fmt.Printf("Success: %s\n", r.Server)
			for _, name := range r.Output {
				fmt.Println(name)
			}
		}
	}
}

func saveEmojiList(domain string, opts Options, c chan Result) {
	r := Result{
		Server: domain,
	}

	var es Emojis
	es, r.Error = NewEmojiList(domain)
	if r.Error != nil {
		c <- r
		return
	}

	out := filepath.Join(opts.OutputDir, strings.Replace(domain, ".", "_", -1))

	if !opts.Split {
		p := NewEmojiPack()
		p.SetFiles(es)
		r.Error = p.GenerateEmojiPack(out)
		r.Output = append(r.Output, out)
		c <- r
		return
	}

	var ces = map[string]Emojis{}
	for _, e := range es {
		ces[e.Category] = append(ces[e.Category], e)
	}

	for c, es := range ces {
		p := NewEmojiPack()
		p.SetFiles(es)
		dir := filepath.Join(out, c)
		if err := p.GenerateEmojiPack(dir); err != nil {
			r.Error = err
		}
		r.Output = append(r.Output, dir)
	}

	c <- r
}
