package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/catatsuy/bento/config"
	"github.com/catatsuy/bento/mirait"
)

const (
	ExitCodeOK   = 0
	ExitCodeFail = 1
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	if len(args) <= 1 {
		log.Println("must provide a input value")
		return ExitCodeFail
	}

	input := trimUnnecessary(args[1])

	conf, exist, err := config.LoadCache()
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	sess, err := mirait.NewSession()
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	token := ""
	if !exist {
		token, err = sess.GetToken()
		if err != nil {
			log.Print(err)
			return ExitCodeFail
		}
	} else {
		sess.SetCacheCookie(conf.Cookies)
		token = conf.Token
	}
	sess.SetToken(token)

	output, err := sess.PostTranslate(input, isJP(input))
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	fmt.Println(output)

	ccs := sess.DumpCookies()
	err = config.DumpCache(config.Config{
		Cookies: ccs,
		Token:   token,
	})
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	return ExitCodeOK
}

func isJP(input string) bool {
	ratio := float64(utf8.RuneCountInString(input)) / float64(len(input))

	return ratio < 0.5
}

func trimUnnecessary(input string) string {
	strs := strings.Split(input, "\n")

	newStrs := make([]string, 0, len(strs))
	for _, s := range strs {
		tmp := strings.TrimLeft(s, " /\t")
		if tmp == "" {
			tmp = "\n"
		}
		newStrs = append(newStrs, tmp)
	}

	return strings.Join(newStrs, " ")
}
