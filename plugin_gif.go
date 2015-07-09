package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
)

// Giphy API
type Giphy struct {
	Data []struct {
		Type             string `json:"type"`
		ID               string `json:"id"`
		URL              string `json:"url"`
		BitlyGifURL      string `json:"bitly_gif_url"`
		BitlyURL         string `json:"bitly_url"`
		EmbedURL         string `json:"embed_url"`
		Username         string `json:"username"`
		Source           string `json:"source"`
		Rating           string `json:"rating"`
		Caption          string `json:"caption"`
		ContentURL       string `json:"content_url"`
		ImportDatetime   string `json:"import_datetime"`
		TrendingDatetime string `json:"trending_datetime"`
		Images           struct {
			FixedHeight struct {
				URL      string `json:"url"`
				Width    string `json:"width"`
				Height   string `json:"height"`
				Size     string `json:"size"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
			} `json:"fixed_height"`
			FixedHeightStill struct {
				URL    string `json:"url"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"fixed_height_still"`
			FixedHeightDownsampled struct {
				URL      string `json:"url"`
				Width    string `json:"width"`
				Height   string `json:"height"`
				Size     string `json:"size"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
			} `json:"fixed_height_downsampled"`
			FixedWidth struct {
				URL      string `json:"url"`
				Width    string `json:"width"`
				Height   string `json:"height"`
				Size     string `json:"size"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
			} `json:"fixed_width"`
			FixedWidthStill struct {
				URL    string `json:"url"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"fixed_width_still"`
			FixedWidthDownsampled struct {
				URL      string `json:"url"`
				Width    string `json:"width"`
				Height   string `json:"height"`
				Size     string `json:"size"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
			} `json:"fixed_width_downsampled"`
			FixedHeightSmall struct {
				URL      string `json:"url"`
				Width    string `json:"width"`
				Height   string `json:"height"`
				Size     string `json:"size"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
			} `json:"fixed_height_small"`
			FixedHeightSmallStill struct {
				URL    string `json:"url"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"fixed_height_small_still"`
			FixedWidthSmall struct {
				URL      string `json:"url"`
				Width    string `json:"width"`
				Height   string `json:"height"`
				Size     string `json:"size"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
			} `json:"fixed_width_small"`
			FixedWidthSmallStill struct {
				URL    string `json:"url"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"fixed_width_small_still"`
			Downsized struct {
				URL    string `json:"url"`
				Width  string `json:"width"`
				Height string `json:"height"`
				Size   string `json:"size"`
			} `json:"downsized"`
			DownsizedStill struct {
				URL    string `json:"url"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"downsized_still"`
			DownsizedLarge struct {
				URL    string `json:"url"`
				Width  string `json:"width"`
				Height string `json:"height"`
				Size   string `json:"size"`
			} `json:"downsized_large"`
			Original struct {
				URL      string `json:"url"`
				Width    string `json:"width"`
				Height   string `json:"height"`
				Size     string `json:"size"`
				Frames   string `json:"frames"`
				Mp4      string `json:"mp4"`
				Mp4Size  string `json:"mp4_size"`
				Webp     string `json:"webp"`
				WebpSize string `json:"webp_size"`
			} `json:"original"`
			OriginalStill struct {
				URL    string `json:"url"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"original_still"`
		} `json:"images"`
	} `json:"data"`
	Meta struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	} `json:"meta"`
	Pagination struct {
		TotalCount int `json:"total_count"`
		Count      int `json:"count"`
		Offset     int `json:"offset"`
	} `json:"pagination"`
}

func pluginGif(result Result) {
	logrus.Infof("pluginGif: request %d", result.UpdateID)

	err := sendChatAction(result.Message.Chat.ID, result.Message.MessageID, "upload_document")
	if err != nil {
		logrus.Error("pluginGif:", err)
		return
	}

	key := url.QueryEscape(strings.TrimPrefix(result.Message.Text, "/gif "))

	content, err := getRemoteURL(fmt.Sprintf("http://api.giphy.com/v1/gifs/search?q=%s&api_key=dc6zaTOxFJmzC", key))
	if err != nil {
		logrus.Error("pluginGif:", err)
		return
	}

	var giphy Giphy
	err = json.Unmarshal(content, &giphy)
	if err != nil {
		logrus.Error("pluginGif:", err)
		return
	}

	if len(giphy.Data) > 0 {
		index := rand.Intn(len(giphy.Data))
		object, err := getRemoteObject(giphy.Data[index].Images.Downsized.URL)
		if err != nil {
			logrus.Error("pluginGif:", err)
			return
		}

		err = sendDocument(result.Message.Chat.ID, result.Message.MessageID, object, fmt.Sprintf("%s.gif", key))
		if err != nil {
			logrus.Error("pluginGif:", err)
			return
		}
	} else {
		err = sendMessage(result.Message.Chat.ID, result.Message.MessageID, "No matches", false)
		if err != nil {
			logrus.Error("pluginGif:", err)
			return
		}
	}
}
