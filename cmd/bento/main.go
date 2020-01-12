package main

import (
	"log"

	"github.com/catatsuy/bento/mirait"
)

func main() {
	sess, err := mirait.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	err = sess.SetToken()
	if err != nil {
		log.Fatal(err)
	}

	output, err := sess.PostTranslate("Hello")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(output)
}
