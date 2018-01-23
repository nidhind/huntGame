package db

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const UsersColl = "users"

// User schema for users collection
type User struct {
	FirstName               string    `bson:"firstName"`
	LastName                string    `bson:"lastName"`
	Email                   string    `bson:"email"`
	Password                string    `bson:"password"`
	Level                   int       `bson:"level"`
	AccessLevel             string    `bson:"accessLevel"`
	AccessToken             string    `bson:"accessToken"`
	PreviousLevelFinishTime time.Time `bson:"previousLevelFinishTime"`
}

// Model for new user insert query
type InsertUserQuery struct {
	FirstName               string    `bson:"firstName"`
	LastName                string    `bson:"lastName"`
	Email                   string    `bson:"email"`
	Password                string    `bson:"password"`
	Level                   int       `bson:"level"`
	AccessLevel             string    `bson:"accessLevel"`
	AccessToken             string    `bson:"accessToken"`
	PreviousLevelFinishTime time.Time `bson:"previousLevelFinishTime"`
}

func GetUserByEmail(emailId string) (User, error) {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(UsersColl)

	var user User
	err := c.Find(bson.M{"email": emailId}).One(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserByAccessToken(t string) (User, error) {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(UsersColl)

	var user User
	err := c.Find(bson.M{"accessToken": t}).One(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func InsertNewUser(u *InsertUserQuery) error {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(UsersColl)
	err := c.Insert(u)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAccessTokenByEmailId(e string, t string) error {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(UsersColl)
	q := bson.M{"email": e}
	u := bson.M{"$set": bson.M{"accessToken": t}}
	err := c.Update(&q, &u)
	if err != nil {
		return err
	}
	return nil
}

func UpdateLevelByEmailId(e string, l int) error {
	s := GetSession()
	defer s.Close()
	c := s.DB(DB).C(UsersColl)
	q := bson.M{"email": e}
	u := bson.M{"$set": bson.M{"level": l, "previousLevelFinishTime": time.Now()}}
	err := c.Update(&q, &u)
	if err != nil {
		return err
	}
	return nil
}
