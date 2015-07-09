package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/Sirupsen/logrus"
)

// Image from Flickr API
type Image struct {
	Photos struct {
		Page    int    `xml:"page,attr"`
		Pages   int    `xml:"pages,attr"`
		PerPage int    `xml:"perpage,attr"`
		Total   string `xml:"total,attr"`
		Photo   []struct {
			ID       string `xml:"id,attr"`
			Owner    string `xml:"owner,attr"`
			Secret   string `xml:"secret,attr"`
			Server   string `xml:"server,attr"`
			Farm     int    `xml:"farm,attr"`
			Title    string `xml:"title,attr"`
			IsPublic int    `xml:"ispublic,attr"`
			IsFriend int    `xml:"isfriend,attr"`
			IsFamily int    `xml:"isfamily,attr"`
			URLN     string `xml:"url_n,attr"`
			HeightN  string `xml:"height_n,attr"`
			WidthN   string `xml:"width_n,attr"`
		} `xml:"photo"`
	} `xml:"photos"`
	Stat string `xml:"stat,attr"`
}

func pluginImg(result Result) {
	logrus.Infof("pluginImg: request %d", result.UpdateID)

	err := sendChatAction(result.Message.Chat.ID, result.Message.MessageID, "upload_photo")
	if err != nil {
		logrus.Error("pluginImg:", err)
		return
	}

	key := url.QueryEscape(strings.TrimPrefix(result.Message.Text, "/img "))

	res, err := http.Get(fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=%s&extras=url_n&text=%s", config.FlickrAPIKey, key))
	if err != nil {
		logrus.Error("pluginImg:", err)
		return
	}
	defer res.Body.Close()

	var image Image
	err = xml.NewDecoder(res.Body).Decode(&image)
	if err != nil {
		logrus.Error("pluginImg:", err)
		return
	}

	if image.Stat == "ok" && len(image.Photos.Photo) > 0 {
		index := rand.Intn(len(image.Photos.Photo))

		var message bytes.Buffer
		w := bufio.NewWriter(&message)
		t, err := template.New("message").Parse("{{if .Title}}{{.Title}}{{end}}")
		if err != nil {
			logrus.Error("pluginImg:", err)
			return
		}
		err = t.Execute(w, image.Photos.Photo[index])
		if err != nil {
			logrus.Error("pluginImg:", err)
			return
		}
		w.Flush()

		object, err := getRemoteObject(image.Photos.Photo[index].URLN)
		if err != nil {
			logrus.Error("pluginImg:", err)
			return
		}
		err = sendPhoto(result.Message.Chat.ID, result.Message.MessageID, object, fmt.Sprintf("%s.jpg", key), message.String())
		if err != nil {
			logrus.Error("pluginImg:", err)
			return
		}
	} else {
		err = sendMessage(result.Message.Chat.ID, result.Message.MessageID, "No matches", false)
		if err != nil {
			logrus.Error("pluginImg:", err)
			return
		}
	}
}
