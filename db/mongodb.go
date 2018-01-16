package db

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/arha/Ertebot/configuration"
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
	session, err = mgo.Dial(configuration.ErtebotConfig.GetString("mongodb-address"))
	if err != nil {
		log.WithError(err).Fatalln("MongoDB session can't be created")
	}

	PeopleCollection = session.DB(dbName).C(peopleCollectionName)
	MessagesCollection = session.DB(dbName).C(messagesCollectionName)
}

func Close() {
	session.Close()
}
