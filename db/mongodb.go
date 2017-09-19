package db

import (
	mgo "gopkg.in/mgo.v2"
)

var (
	PeopleCollection   *mgo.Collection
	MessagesCollection *mgo.Collection
	session            *mgo.Session
)

const (
	dbName                 = "ertebot"
	peopleCollectionName   = "people"
	messagesCollectionName = "messages"
)

func NewMongoDB() {
	var err error
	session, err = mgo.Dial("aryaha.com:27017")
	if err != nil {
		panic(err)
	}

	PeopleCollection = session.DB(dbName).C(peopleCollectionName)
	MessagesCollection = session.DB(dbName).C(messagesCollectionName)
}

func Close() {
	session.Close()
}
