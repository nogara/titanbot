package main

import "github.com/Sirupsen/logrus"

func pluginComando(result Result) {
	logrus.Infof("pluginComando: request %d", result.UpdateID)
	err := sendMessage(result.Message.Chat.ID, 0, "/comando", true)
	if err != nil {
		logrus.Error("pluginComando:", err)
		return
	}
}
