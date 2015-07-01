package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/Sirupsen/logrus"
)

// Series from OMDB API
type Series struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	Rated      string `json:"Rated"`
	Released   string `json:"Released"`
	Runtime    string `json:"Runtime"`
	Genre      string `json:"Genre"`
	Director   string `json:"Director"`
	Writer     string `json:"Writer"`
	Actors     string `json:"Actors"`
	Plot       string `json:"Plot"`
	Language   string `json:"Language"`
	Country    string `json:"Country"`
	Awards     string `json:"Awards"`
	Poster     string `json:"Poster"`
	Metascore  string `json:"Metascore"`
	Imdbrating string `json:"imdbRating"`
	Imdbvotes  string `json:"imdbVotes"`
	Imdbid     string `json:"imdbID"`
	Type       string `json:"Type"`
	Response   string `json:"Response"`
}

func pluginSeries(result Result) {
	logrus.Infof("pluginSeries: request %d", result.UpdateID)

	err := sendChatAction(result.Message.Chat.ID, result.Message.MessageID, "upload_photo")
	if err != nil {
		logrus.Error("pluginSeries:", err)
		return
	}

	key := url.QueryEscape(strings.TrimPrefix(result.Message.Text, "/series "))

	res, err := http.Get(fmt.Sprintf("http://www.omdbapi.com/?t=%s&plot=short&type=series&r=json", key))
	if err != nil {
		logrus.Error("pluginSeries:", err)
		return
	}
	defer res.Body.Close()

	var series Series
	err = json.NewDecoder(res.Body).Decode(&series)
	if err != nil {
		logrus.Error("pluginSeries:", err)
		return
	}

	if series.Response == "True" {
		var message bytes.Buffer
		w := bufio.NewWriter(&message)
		t, err := template.New("message").Parse("{{if .Title}}{{.Title}}{{end}}{{if .Year}} ({{.Year}}){{end}}{{if .Imdbrating}} - {{.Imdbrating}}{{end}}{{if .Imdbid}}\nhttp://imdb.com/title/{{.Imdbid}}{{end}}{{if .Plot}}\n\n{{.Plot}}{{end}}")
		if err != nil {
			logrus.Error("pluginSeries:", err)
			return
		}
		err = t.Execute(w, series)
		if err != nil {
			logrus.Error("pluginSeries:", err)
			return
		}
		w.Flush()

		if len(series.Poster) > 0 {
			object, err := getRemoteObject(series.Poster)
			if err != nil {
				logrus.Error("pluginSeries:", err)
				return
			}
			err = sendPhoto(result.Message.Chat.ID, result.Message.MessageID, object, fmt.Sprintf("%s.jpg", key), message.String())
			if err != nil {
				logrus.Error("pluginSeries:", err)
				return
			}
		} else {
			err = sendMessage(result.Message.Chat.ID, result.Message.MessageID, message.String(), true)
			if err != nil {
				logrus.Error("pluginSeries:", err)
				return
			}
		}
	} else {
		err = sendMessage(result.Message.Chat.ID, result.Message.MessageID, "No matches", false)
		if err != nil {
			logrus.Error("pluginSeries:", err)
			return
		}
	}
}
