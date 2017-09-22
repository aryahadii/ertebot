package util

import (
	"errors"
	"strings"

	"gitlab.com/arha/Ertebot/db"
	"gopkg.in/mgo.v2/bson"
)

func GetUserID(username string) (string, error) {
	var id string
	err := db.PeopleCollection.Find(bson.M{"username": strings.ToLower(username)}).One(&id)
	if err != nil {
		return id, errors.New("Not found")
	}
	return id, nil
}
