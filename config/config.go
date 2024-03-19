package config

import (
	"bufio"
	mrand "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	toml "github.com/pelletier/go-toml/v2"
)

const (
	defaultUserAgent = "waiwai client"
)

var (
	random = mrand.New(mrand.NewSource(time.Now().UnixNano()))
)

type Config struct {
	UserAgent string
}

func configFileName() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homedir, ".bento", "config.toml"), nil
}

func LoadConfig() (Config, error) {
	cfn, err := configFileName()
	if err != nil {
		return Config{}, err
	}
	f, err := os.Open(cfn)
	if err != nil {
		if !os.IsNotExist(err) {
			return Config{}, err
		}
		return Config{}, nil
	}
	defer f.Close()

	cfg := Config{}
	tree, err := toml.LoadReader(f)
	if err != nil {
		return Config{}, err
	}

	cfg.UserAgent = tree.GetDefault("bento.user_agent", defaultUserAgent).(string)

	return cfg, nil
}

func configWordsFileName() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homedir, ".bento", "words.txt"), nil
}

func LoadWords() ([]string, error) {
	cfn, err := configWordsFileName()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(cfn)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return nil, nil
	}
	defer f.Close()

	words := make([]string, 0, 10)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.Trim(text, " ")
		if text == "" {
			continue
		}
		words = append(words, text)
	}

	return words, nil
}

func Replacer(words []string) ([]string, []string) {
	l := len(words)
	oldnew := make([]string, 2*l)
	newold := make([]string, 2*l)

	for i := 0; i < l; i++ {
		rstr := randomStr(4, "ABCD")

		oldnew[2*i] = words[i]
		oldnew[2*i+1] = rstr
		newold[2*i] = rstr
		newold[2*i+1] = words[i]
	}

	return oldnew, newold
}

func randomStr(n int, s string) string {
	buf := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		buf = append(buf, byte(s[random.Int()%len(s)]))
	}
	return string(buf)
}
