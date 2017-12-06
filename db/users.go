package db

import (
	"gopkg.in/mgo.v2/bson"
)

const UsersColl = "users"

// User schema for users collection
type User struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Rank        int    `json:"rank"`
	AccessLevel string `json:"accessLevel"`
	AccessToken string `json:"accessToken"`
}

// Model for new user insert query
type InsertUserQuery struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Rank        int    `json:"rank"`
	AccessLevel string `json:"accessLevel"`
	AccessToken string `json:"accessToken"`
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
