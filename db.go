package main

import (
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"
)

var session *gorethink.Session

func initializeDB() {
	var err error
	session, err = gorethink.Connect(gorethink.ConnectOpts{Address: config.DBHost, Database: config.DBName})
	if err != nil {
		logrus.Fatal("db:", err)
	}
}

func getRemoteObject(url string) (Object, error) {
	rows, err := gorethink.Table("objects").GetAllByIndex("url", url).Run(session)

	var row Object
	err = rows.One(&row)
	if err != nil && err != gorethink.ErrEmptyResult {
		return Object{}, err
	}

	if err == gorethink.ErrEmptyResult {
		res, err := http.Get(url)
		if err != nil {
			return Object{}, err
		}
		defer res.Body.Close()

		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return Object{}, err
		}

		result, err := gorethink.Table("objects").Insert(Object{URL: url, Content: content}).RunWrite(session)
		if err != nil {
			return Object{}, err
		}
		return Object{ID: result.GeneratedKeys[0], URL: url, Content: content}, nil
	}
	return row, nil
}

func updateRemoteObject(objectID string, fileID string) error {
	_, err := gorethink.Table("objects").Get(objectID).Update(Object{FileID: fileID}).RunWrite(session)
	if err != nil {
		return err
	}
	return nil
}
