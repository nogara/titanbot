package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"math/rand"
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
			URLZ     string `xml:"url_z,attr"`
			HeightZ  string `xml:"height_z,attr"`
			WidthZ   string `xml:"width_z,attr"`
			URLN     string `xml:"url_n,attr"`
			HeightN  string `xml:"height_n,attr"`
			WidthN   string `xml:"width_n,attr"`
			URLO     string `xml:"url_o,attr"`
			HeightO  string `xml:"height_o,attr"`
			WidthO   string `xml:"width_o,attr"`
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

	content, err := getRemoteURL(fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=%s&safe_search=3&sort=relevance&per_page=10&extras=url_z,url_n,url_o&text=%s", config.FlickrAPIKey, key))
	if err != nil {
		logrus.Error("pluginImg:", err)
		return
	}

	var image Image
	err = xml.Unmarshal(content, &image)
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

		var imageURL string
		if len(image.Photos.Photo[index].URLZ) > 0 {
			imageURL = image.Photos.Photo[index].URLZ
		} else if len(image.Photos.Photo[index].URLN) > 0 {
			imageURL = image.Photos.Photo[index].URLN
		} else if len(image.Photos.Photo[index].URLO) > 0 {
			imageURL = image.Photos.Photo[index].URLO
		}

		object, err := getRemoteObject(imageURL)
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
