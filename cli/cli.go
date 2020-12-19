package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/catatsuy/bento/config"
	"github.com/catatsuy/bento/mirait"
	"github.com/catatsuy/bento/util"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

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

type CLI struct {
	appVersion           string
	outStream, errStream io.Writer
}

func NewCLI(outStream, errStream io.Writer) *CLI {
	log.SetOutput(errStream)
	return &CLI{appVersion: version(), outStream: outStream, errStream: errStream}
}

func version() string {
	if Version != "" {
		return Version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	return info.Main.Version
}

func (c *CLI) Run(args []string) int {
	if len(args) <= 1 {
		fmt.Fprintln(c.errStream, "must provide a input value")
		return ExitCodeFail
	}

	var (
		version  bool
		refresh  bool
		trim     bool
		filename string
		from     string
		to       string
	)

	flags := flag.NewFlagSet("bento", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.StringVar(&filename, "file", "", "translate a file")
	flags.BoolVar(&trim, "trim", false, "print text which remove the unnecessary characters")
	flags.BoolVar(&version, "version", false, "print version information and quit")
	flags.BoolVar(&refresh, "refresh", false, "refresh cache file")
	flags.StringVar(&from, "from", "", "from language")
	flags.StringVar(&to, "to", "", "to language")

	err := flags.Parse(args[1:])
	if err != nil {
		log.Print(err)
		return ExitCodeParseFlagError
	}
	if version {
		fmt.Fprintf(c.errStream, "bento version %s; %s\n", c.appVersion, runtime.Version())
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

	input := ""

	argv := flags.Args()
	if len(argv) == 1 {
		input = argv[0]
	} else if len(argv) > 1 {
		input = argv[0]
		err := flags.Parse(argv[1:])
		if err != nil {
			return ExitCodeParseFlagError
		}
	}

	if filename != "" {
		bb, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Print(err)
			return ExitCodeFail
		}
		input = string(bb)
	}

	if trim {
		return c.trim(input)
	}

	isJP := false
	if from == "" && to == "" {
		isJP = util.AutoDetectJP(input)
	} else if from == "ja" || to == "en" {
		isJP = true
	}

	return c.translate(input, isJP)
}

func (c *CLI) trim(input string) int {
	words, err := config.LoadWords()
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	r := strings.NewReplacer(". \n", ".\n\n", ".\n", ".\n\n")
	input = util.TrimUnnecessary(r.Replace(input))

	oldnew, _ := config.Replacer(words)
	replacerNoun := strings.NewReplacer(oldnew...)

	output := replacerNoun.Replace(input)

	fmt.Fprintln(c.outStream, output)

	return ExitCodeOK
}

func (c *CLI) translate(input string, isJP bool) int {
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
	input = util.TrimUnnecessary(r.Replace(input))

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

	characters := utf8.RuneCountInString(input)
	if characters < maxCharacters {
		output, err := sess.PostTranslate(input, isJP)
		if err != nil {
			log.Print(err)
			return ExitCodeFail
		}
		fmt.Fprintln(c.outStream, revertNoun.Replace(output))
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
				fmt.Fprintln(c.errStream, "you must split input")
				return ExitCodeFail
			}
		}

		for _, sinput := range inputs {
			output, err := sess.PostTranslate(sinput, isJP)
			if err != nil {
				log.Print(err)
				return ExitCodeFail
			}
			fmt.Fprintln(c.outStream, revertNoun.Replace(output))
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
