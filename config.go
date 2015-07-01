package main

import (
	"encoding/json"
	"os"

	"github.com/Sirupsen/logrus"
)

var config Config

func initializeConfig() {
	f, err := os.Open("config.json")
	if err != nil {
		logrus.Fatal("config:", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		logrus.Fatal("config:", err)
	}
}
