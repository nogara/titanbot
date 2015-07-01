package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

func pluginSearch(result Result) {
	logrus.Infof("pluginSearch: request %d", result.UpdateID)

	key := url.QueryEscape(strings.TrimPrefix(result.Message.Text, "/g "))

	doc, err := goquery.NewDocument(fmt.Sprintf("https://www.google.com/search?q=%s&hl=en", key))
	if err != nil {
		logrus.Error("pluginSearch:", err)
		return
	}

	var queryResult []string
	doc.Find("h3 a").Each(func(i int, s *goquery.Selection) {
		page, _ := s.Attr("href")
		if strings.HasPrefix(page, "/url?q=") {
			unescapedURL, err := url.QueryUnescape(strings.TrimSpace(strings.Split(strings.TrimPrefix(page, "/url?q="), "&")[0]))
			if err != nil {
				logrus.Error("pluginSearch:", err)
				return
			}
			queryResult = append(queryResult, unescapedURL)
		}
	})

	if len(queryResult) > 0 {
		err = sendMessage(result.Message.Chat.ID, result.Message.MessageID, queryResult[0], false)
		if err != nil {
			logrus.Error("pluginSearch:", err)
			return
		}
	} else {
		err = sendMessage(result.Message.Chat.ID, result.Message.MessageID, "No matches", false)
		if err != nil {
			logrus.Error("pluginSearch:", err)
			return
		}
	}
}
