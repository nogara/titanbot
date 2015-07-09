package main

import "github.com/Sirupsen/logrus"

func pluginHelp(result Result) {
	logrus.Infof("pluginHelp: request %d", result.UpdateID)

	err := sendMessage(result.Message.Chat.ID, result.Message.MessageID, `This bot have multiple uses, currently the following commands are available:

/g [keyword] - Search Google
/img [keyword] - Search for a random image
/gif [keyword] - Search for a random gif
/movie [keyword] - Search for movie info
/series [keyword] - Search for tv series info
/artist [keyword] - Search for movie/tv series artist info

More commands will be added soon, check back often for new stuff!

Results are retrieved from multiple sources, including: Google, Flickr, Giphy, OMDb and TMDb.`, true)
	if err != nil {
		logrus.Error("pluginHelp:", err)
		return
	}
}
