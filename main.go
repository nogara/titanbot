package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	initializeConfig()
	initializeDB()
}

func main() {
	logrus.Infof("Starting TitanBot")
	lastUpdateID := 0
	for {
		res, err := http.Get(fmt.Sprintf("%s/getUpdates?offset=%d", "https://api.telegram.org/bot"+config.TelegramAPIKey, lastUpdateID+1))
		if err != nil {
			logrus.Error(err)
			time.Sleep(10 * time.Second)
			continue
		}
		defer res.Body.Close()

		var updates Updates
		err = json.NewDecoder(res.Body).Decode(&updates)
		if err != nil {
			logrus.Error(err)
			time.Sleep(10 * time.Second)
			continue
		}

		if updates.Ok == true {
			for _, result := range updates.Result {
				if result.UpdateID > lastUpdateID {
					lastUpdateID = result.UpdateID
					go processUpdate(result)
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func processUpdate(result Result) {
	switch {
	case strings.HasPrefix(result.Message.Text, "/start"):
		pluginHelp(result)
	case strings.HasPrefix(result.Message.Text, "/help"):
		pluginHelp(result)
	case strings.HasPrefix(result.Message.Text, "/g "):
		pluginSearch(result)
	case strings.HasPrefix(result.Message.Text, "/gif "):
		pluginGif(result)
	case strings.HasPrefix(result.Message.Text, "/movie "):
		pluginMovie(result)
	case strings.HasPrefix(result.Message.Text, "/series "):
		pluginSeries(result)
	case strings.HasPrefix(result.Message.Text, "/artist "):
		pluginArtist(result)
	case strings.HasPrefix(result.Message.Text, "/img "):
		pluginImg(result)
	}
}
