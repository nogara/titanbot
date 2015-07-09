package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"github.com/Sirupsen/logrus"
)

// Movie from OMDB API
type Movie struct {
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

func pluginMovie(result Result) {
	logrus.Infof("pluginMovie: request %d", result.UpdateID)

	err := sendChatAction(result.Message.Chat.ID, result.Message.MessageID, "upload_photo")
	if err != nil {
		logrus.Error("pluginMovie:", err)
		return
	}

	key := url.QueryEscape(strings.TrimPrefix(result.Message.Text, "/movie "))

	content, err := getRemoteURL(fmt.Sprintf("http://www.omdbapi.com/?t=%s&plot=short&type=movie&r=json", key))
	if err != nil {
		logrus.Error("pluginMovie:", err)
		return
	}

	var movie Movie
	err = json.Unmarshal(content, &movie)
	if err != nil {
		logrus.Error("pluginMovie:", err)
		return
	}

	if movie.Response == "True" {
		var message bytes.Buffer
		w := bufio.NewWriter(&message)
		t, err := template.New("message").Parse("{{if .Title}}{{.Title}}{{end}}{{if .Year}} ({{.Year}}){{end}}{{if .Imdbrating}} - {{.Imdbrating}}{{end}}{{if .Imdbid}}\nhttp://imdb.com/title/{{.Imdbid}}{{end}}{{if .Plot}}\n\n{{.Plot}}{{end}}")
		if err != nil {
			logrus.Error("pluginMovie:", err)
			return
		}
		err = t.Execute(w, movie)
		if err != nil {
			logrus.Error("pluginMovie:", err)
			return
		}
		w.Flush()

		if len(movie.Poster) > 0 {
			object, err := getRemoteObject(movie.Poster)
			if err != nil {
				logrus.Error("pluginMovie:", err)
				return
			}
			err = sendPhoto(result.Message.Chat.ID, result.Message.MessageID, object, fmt.Sprintf("%s.jpg", key), message.String())
			if err != nil {
				logrus.Error("pluginMovie:", err)
				return
			}
		} else {
			err = sendMessage(result.Message.Chat.ID, result.Message.MessageID, message.String(), true)
			if err != nil {
				logrus.Error("pluginMovie:", err)
				return
			}
		}
	} else {
		err = sendMessage(result.Message.Chat.ID, result.Message.MessageID, "No matches", false)
		if err != nil {
			logrus.Error("pluginMovie:", err)
			return
		}
	}
}
