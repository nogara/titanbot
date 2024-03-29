package main

import (
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang/snappy/snappy"
)

var cacher *memcache.Client

func initializeCache() {
	var err error
	cacher = memcache.New(config.MCHost)
	if err != nil {
		logrus.Fatal("cache:", err)
	}
}

func getRemoteURL(url string) ([]byte, error) {
	item, err := cacher.Get(url)
	if err != nil && err != memcache.ErrCacheMiss {
		return nil, err
	}

	if err == memcache.ErrCacheMiss {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		content, err := snappy.Encode(nil, body)
		if err != nil {
			return nil, err
		}

		err = cacher.Set(&memcache.Item{Key: url, Value: content, Expiration: config.MCTTL})
		if err != nil {
			return nil, err
		}
		return body, nil
	}

	content, err := snappy.Decode(nil, item.Value)
	if err != nil {
		return nil, err
	}

	return content, nil
}
