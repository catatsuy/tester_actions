package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/catatsuy/bento/config"
	"github.com/catatsuy/bento/mirait"
)

const (
	ExitCodeOK             = 0
	ExitCodeFail           = 1
	ExitCodeParseFlagError = 1

	splitCharacters = 1500
	maxCharacters   = 2000
)

var (
	Version string
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	if len(args) <= 1 {
		log.Println("must provide a input value")
		return ExitCodeFail
	}

	input := args[1]

	var (
		version  bool
		refresh  bool
		filename string
	)

	flags := flag.NewFlagSet("bento", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	flags.StringVar(&filename, "file", "", "translate a file")
	flags.BoolVar(&version, "version", false, "print version information and quit")
	flags.BoolVar(&refresh, "refresh", false, "refresh cache file")

	err := flags.Parse(args[1:])
	if err != nil {
		log.Print(err)
		return ExitCodeParseFlagError
	}
	if version {
		fmt.Fprintf(os.Stderr, "bento version %s; %s\n", Version, runtime.Version())
		return ExitCodeOK
	}

	if refresh {
		err := config.RemoveCache()
		if err != nil {
			log.Print(err)
			return ExitCodeFail
		}
		return ExitCodeOK
	}

	if filename != "" {
		bb, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Print(err)
			return ExitCodeFail
		}
		input = string(bb)
	}

	return translate(input)
}

func translate(input string) int {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	words, err := config.LoadWords()
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	r := strings.NewReplacer(". \n", ".\n\n", ".\n", ".\n\n")
	input = trimUnnecessary(r.Replace(input))

	oldnew, newold := config.Replacer(words)
	replacerNoun := strings.NewReplacer(oldnew...)
	revertNoun := strings.NewReplacer(newold...)

	input = replacerNoun.Replace(input)

	conf, exist, err := config.LoadCache()
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	sess, err := mirait.NewSession(cfg)
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

	isJP := isJP(input)

	characters := utf8.RuneCountInString(input)
	if characters < maxCharacters {
		output, err := sess.PostTranslate(input, isJP)
		if err != nil {
			log.Print(err)
			return ExitCodeFail
		}
		fmt.Println(revertNoun.Replace(output))
	} else {
		inputSplits := strings.Split(input, "\n")

		chs := make([]int, len(inputSplits))
		for i, in := range inputSplits {
			chs[i] = utf8.RuneCountInString(in)
		}

		inputs := make([]string, 0, len(chs))
		index := 0
		count := 0
		for i := range chs {
			count += chs[i]
			if count < splitCharacters {
				continue
			}
			if count < maxCharacters {
				inputs = append(inputs, strings.Join(inputSplits[index:i+1], "\n"))
				index = i + 1
				count = 0
			} else if i > index {
				inputs = append(inputs, strings.Join(inputSplits[index:i], "\n"))
				index = i
				count = chs[i]
			} else {
				log.Print("you must split input")
				return ExitCodeFail
			}
		}

		for _, sinput := range inputs {
			output, err := sess.PostTranslate(sinput, isJP)
			if err != nil {
				log.Print(err)
				return ExitCodeFail
			}
			fmt.Println(revertNoun.Replace(output))
			time.Sleep(4 * time.Second)
		}
	}

	if !exist {
		ccs := sess.DumpCookies()
		err = config.DumpCache(config.Cache{
			Cookies: ccs,
			Token:   token,
		})
		if err != nil {
			log.Print(err)
			return ExitCodeFail
		}
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
