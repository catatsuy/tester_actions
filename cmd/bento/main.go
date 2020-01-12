package main

import (
	"context"
	"log"

	"github.com/catatsuy/bento/mirait"
)

func main() {
	sess, err := mirait.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = sess.SetToken(ctx)
	if err != nil {
		log.Fatal(err)
	}

	output, err := sess.PostTranslate(ctx, "Hello")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(output)
}
