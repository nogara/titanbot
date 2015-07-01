package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/edgard/goutil"
)

func sendMessage(chatID int, messageID int, text string, preview bool) error {
	var content bytes.Buffer
	w := multipart.NewWriter(&content)
	contentType := w.FormDataContentType()

	w.WriteField("chat_id", strconv.Itoa(chatID))
	w.WriteField("reply_to_message_id", strconv.Itoa(messageID))
	w.WriteField("text", text)
	w.WriteField("disable_web_page_preview", strconv.FormatBool(preview))

	w.Close()

	res, err := http.Post(fmt.Sprintf("%s/sendMessage", "https://api.telegram.org/bot"+config.TelegramAPIKey), contentType, &content)
	if err != nil {
		return err
	}

	var result ResultMessage
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result.Ok != true {
		return errors.New("error sending message: " + strconv.Itoa(messageID))
	}

	return nil
}

func sendPhoto(chatID int, messageID int, object Object, filename string, caption string) error {
	var content bytes.Buffer
	w := multipart.NewWriter(&content)
	contentType := w.FormDataContentType()

	w.WriteField("chat_id", strconv.Itoa(chatID))
	w.WriteField("reply_to_message_id", strconv.Itoa(messageID))
	if len(caption) > 0 {
		w.WriteField("caption", goutil.StringCap(caption, 200))
	}
	if len(object.FileID) > 0 {
		w.WriteField("photo", object.FileID)
	} else {
		filePart, _ := w.CreateFormFile("photo", filename)
		filePart.Write(object.Content)
	}

	w.Close()

	res, err := http.Post(fmt.Sprintf("%s/sendPhoto", "https://api.telegram.org/bot"+config.TelegramAPIKey), contentType, &content)
	if err != nil {
		return err
	}

	var result ResultPhoto
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result.Ok != true {
		return errors.New("error sending photo: " + strconv.Itoa(messageID))
	}

	if len(object.FileID) == 0 {
		updateRemoteObject(object.ID, result.Result.Photo[0].FileID)
	}

	return nil
}

func sendDocument(chatID int, messageID int, object Object, filename string) error {
	var content bytes.Buffer
	w := multipart.NewWriter(&content)
	contentType := w.FormDataContentType()

	w.WriteField("chat_id", strconv.Itoa(chatID))
	w.WriteField("reply_to_message_id", strconv.Itoa(messageID))

	if len(object.FileID) > 0 {
		w.WriteField("document", object.FileID)
	} else {
		filePart, _ := w.CreateFormFile("document", filename)
		filePart.Write(object.Content)
	}

	w.Close()

	res, err := http.Post(fmt.Sprintf("%s/sendDocument", "https://api.telegram.org/bot"+config.TelegramAPIKey), contentType, &content)
	if err != nil {
		return err
	}

	var result ResultDocument
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result.Ok != true {
		return errors.New("error sending document: " + strconv.Itoa(messageID))
	}

	if len(object.FileID) == 0 {
		updateRemoteObject(object.ID, result.Result.Document.FileID)
	}

	return nil
}

func sendChatAction(chatID int, messageID int, action string) error {
	var content bytes.Buffer
	w := multipart.NewWriter(&content)
	contentType := w.FormDataContentType()

	w.WriteField("chat_id", strconv.Itoa(chatID))
	w.WriteField("action", action)

	w.Close()

	res, err := http.Post(fmt.Sprintf("%s/sendChatAction", "https://api.telegram.org/bot"+config.TelegramAPIKey), contentType, &content)
	if err != nil {
		return err
	}

	var result ResultAction
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result.Ok != true {
		return errors.New("error sending action: " + strconv.Itoa(messageID))
	}

	return nil
}
