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

type Cache struct {
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

func LoadCache() (cache Cache, exist bool, err error) {
	cacheFile, err := cacheFileName()
	if err != nil {
		return Cache{}, false, err
	}

	f, err := os.Open(cacheFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return Cache{}, false, err
		}
		return Cache{}, false, nil
	}
	defer f.Close()
	fi, err := os.Stat(cacheFile)
	if err != nil {
		return Cache{}, false, err
	}
	now := time.Now()
	if !fi.ModTime().Add(30 * time.Minute).After(now) {
		return Cache{}, false, nil
	}

	cache = Cache{}
	err = json.NewDecoder(f).Decode(&cache)
	if err != nil {
		return Cache{}, false, err
	}

	return cache, true, nil
}

func DumpCache(cache Cache) error {
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

	err = json.NewEncoder(f).Encode(cache)
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
