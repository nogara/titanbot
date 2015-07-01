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

// Person from TMDB API
type Person struct {
	Page    int `json:"page"`
	Results []struct {
		Adult    bool `json:"adult"`
		ID       int  `json:"id"`
		KnownFor []struct {
			Adult            bool    `json:"adult"`
			BackdropPath     string  `json:"backdrop_path"`
			GenreIds         []int   `json:"genre_ids"`
			ID               int     `json:"id"`
			OriginalLanguage string  `json:"original_language"`
			OriginalTitle    string  `json:"original_title"`
			Overview         string  `json:"overview"`
			ReleaseDate      string  `json:"release_date"`
			PosterPath       string  `json:"poster_path"`
			Popularity       float64 `json:"popularity"`
			Title            string  `json:"title"`
			Video            bool    `json:"video"`
			VoteAverage      float64 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
			MediaType        string  `json:"media_type"`
		} `json:"known_for"`
		Name        string  `json:"name"`
		Popularity  float64 `json:"popularity"`
		ProfilePath string  `json:"profile_path"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// PersonDetail from TMDB API
type PersonDetail struct {
	Adult        bool          `json:"adult"`
	AlsoKnownAs  []interface{} `json:"also_known_as"`
	Biography    string        `json:"biography"`
	Birthday     string        `json:"birthday"`
	Deathday     string        `json:"deathday"`
	Homepage     string        `json:"homepage"`
	ID           int           `json:"id"`
	ImdbID       string        `json:"imdb_id"`
	Name         string        `json:"name"`
	PlaceOfBirth string        `json:"place_of_birth"`
	Popularity   float64       `json:"popularity"`
	ProfilePath  string        `json:"profile_path"`
}

func pluginArtist(result Result) {
	logrus.Infof("pluginArtist: request %d", result.UpdateID)

	err := sendChatAction(result.Message.Chat.ID, result.Message.MessageID, "upload_photo")
	if err != nil {
		logrus.Error("pluginArtist:", err)
		return
	}

	key := url.QueryEscape(strings.TrimPrefix(result.Message.Text, "/artist "))

	res, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/search/person?api_key=%s&query=%s", config.TMDBAPIKey, key))
	if err != nil {
		logrus.Error("pluginArtist:", err)
		return
	}
	defer res.Body.Close()

	var person Person
	err = json.NewDecoder(res.Body).Decode(&person)
	if err != nil {
		logrus.Error("pluginArtist:", err)
		return
	}

	if person.TotalResults > 0 {
		res, err := http.Get(fmt.Sprintf("https://api.themoviedb.org/3/person/%d?api_key=%s", person.Results[0].ID, config.TMDBAPIKey))
		if err != nil {
			logrus.Error("pluginArtist:", err)
			return
		}
		defer res.Body.Close()

		var personDetail PersonDetail
		err = json.NewDecoder(res.Body).Decode(&personDetail)
		if err != nil {
			logrus.Error("pluginArtist:", err)
			return
		}

		var message bytes.Buffer
		w := bufio.NewWriter(&message)
		t, err := template.New("message").Parse("{{if .Name}}{{.Name}}{{end}}{{if .Birthday}} ({{.Birthday}}){{end}}{{if .PlaceOfBirth}}\n{{.PlaceOfBirth}}{{end}}{{if or .ImdbID .Homepage}}\n{{end}}{{if .ImdbID}}\nhttp://imdb.com/name/{{.ImdbID}}{{end}}{{if .Homepage}}\n{{.Homepage}}{{end}}")
		if err != nil {
			logrus.Error("pluginArtist:", err)
			return
		}
		err = t.Execute(w, personDetail)
		if err != nil {
			logrus.Error("pluginArtist:", err)
			return
		}
		w.Flush()

		if strings.HasPrefix(personDetail.ProfilePath, "/") {
			object, err := getRemoteObject("http://image.tmdb.org/t/p/w300" + personDetail.ProfilePath)
			if err != nil {
				logrus.Error("pluginArtist:", err)
				return
			}
			err = sendPhoto(result.Message.Chat.ID, result.Message.MessageID, object, fmt.Sprintf("%s.jpg", key), message.String())
			if err != nil {
				logrus.Error("pluginArtist:", err)
				return
			}
		} else {
			sendMessage(result.Message.Chat.ID, result.Message.MessageID, message.String(), true)
		}
	} else {
		sendMessage(result.Message.Chat.ID, result.Message.MessageID, "No matches", false)
	}
}
