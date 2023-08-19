package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/catatsuy/bento/openai"
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
		filename string
	)

	flags := flag.NewFlagSet("bento", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.StringVar(&filename, "file", "", "translate a file")
	flags.BoolVar(&version, "version", false, "print version information and quit")

	err := flags.Parse(args[1:])
	if err != nil {
		log.Print(err)
		return ExitCodeParseFlagError
	}
	if version {
		fmt.Fprintf(c.errStream, "bento version %s; %s\n", c.appVersion, runtime.Version())
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
		bb, err := os.ReadFile(filename)
		if err != nil {
			log.Print(err)
			return ExitCodeFail
		}
		input = string(bb)
	}

	client, err := openai.NewClient(openai.OpenAIAPIURL)
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	payload := &openai.Payload{
		Model: "gpt-3.5-turbo",
		Messages: []openai.Message{
			{
				Role:    "user",
				Content: "英語を日本語に翻訳してください。返事は翻訳された文章のみにしてください。" + input,
			},
		},
	}

	res, err := client.Chat(context.Background(), payload)
	if err != nil {
		log.Print(err)
		return ExitCodeFail
	}

	fmt.Fprintf(c.outStream, "%s\n", res.Choices[0].Message.Content)
	return 0
}
