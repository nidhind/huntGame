package db

import (
	"log"

	"gopkg.in/mgo.v2"
)

const DB = "mmdb_dev"

var Session *mgo.Session

func InitMongo() bool {
	url := "mongodb://localhost"
	log.Println("Establishing MongoDB connection...")
	var err error
	Session, err = mgo.Dial(url)
	if err != nil {
		log.Fatal("Cannot connect to MongoDB!")
		return true
	} else {
		log.Println("Connected to ", url)
		return false
	}
}

func GetNewSession() mgo.Session {
	return *Session.Copy()
}
