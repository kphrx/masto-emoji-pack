package main

import (
	"fmt"
	"log"
)

func main() {
	es, err := NewEmojiList("pl.kpherox.dev")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range es {
		fmt.Printf("%s, %s, %s\n", e.Shortcode, e.Url, e.Category)
	}
}
