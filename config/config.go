package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Config struct {
	Cookies []Cookie `json:"cookies"`
	Token   string   `json:"token"`
}

func cacheFileName() (fileName string, err error) {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	fileName = filepath.Join(cachedir, "bento", "cache")

	return fileName, nil
}

func LoadCache() (conf Config, exist bool, err error) {
	cacheFile, err := cacheFileName()
	if err != nil {
		return Config{}, false, err
	}

	f, err := os.Open(cacheFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return Config{}, false, err
		}
		return Config{}, false, nil
	}
	defer f.Close()
	fi, err := os.Stat(cacheFile)
	if err != nil {
		return Config{}, false, err
	}
	now := time.Now()
	if !fi.ModTime().Add(30 * time.Minute).After(now) {
		return Config{}, false, nil
	}

	conf = Config{}
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		return Config{}, false, err
	}

	return conf, true, nil
}

func DumpCache(conf Config) error {
	cacheFile, err := cacheFileName()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(cacheFile), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(conf)
	if err != nil {
		return err
	}

	return nil
}

func RemoveCache() error {
	cacheFile, err := cacheFileName()
	if err != nil {
		return err
	}
	err = os.Remove(cacheFile)
	if err != nil {
		return err
	}
	return nil
}
