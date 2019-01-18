package model

import (
	"time"
	"github.com/globalsign/mgo/bson"
	"log"
	"github.com/mitchellh/mapstructure"
	"github.com/BTUS/mdb"
	"github.com/BTUS/static"
)

type User struct {
	MongoID bson.ObjectId `bson:"_id,omitempty"`
	Email string
	Phone string
	Password string
	DisplayName string
	Avatar string

	Token string
	TokenExpiryTime time.Time
	LoginNonce uint64

	IsVerify bool

	// FCM
	FcmToken string
}


func (user *User) Insert(u User) error {
	err := mdb.Insert(static.TBL_USERS, u)
	if err  != nil {
		return err
	}
	return nil
}

func (user *User) GetAll() []User {
	results := mdb.QueryAll(static.TBL_USERS)
	var users []User

	for _, u := range results {
		var user User
		err := mapstructure.Decode(u, &user)
		if err  != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users
}


func (user *User) QueryByEmailPassword(email string, pass string) (error, User) {

	db := mdb.InitDB()
	c := db.C(static.TBL_USERS)

	var result User
	err := c.Find(bson.M{"email": email, "password": pass}).One(&result)
	if err != nil {
		log.Println(err)
		return err, result
	}

	return nil, result
}


func (user *User) QueryByPhonePassword(phone string, pass string) (error, User) {

	db := mdb.InitDB()
	c := db.C(static.TBL_USERS)

	var result User
	err := c.Find(bson.M{"phone": phone, "password": pass}).One(&result)
	if err != nil {
		log.Println(err)
		return err, result
	}

	return nil, result
}


func (user *User) QueryByEmail(email string) (error, User) {

	db := mdb.InitDB()
	c := db.C(static.TBL_USERS)

	var result User
	err := c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return err, result
	}

	return nil, result
}

func (user *User) QueryByPhone(phone string) (error, User) {

	db := mdb.InitDB()
	c := db.C(static.TBL_USERS)

	var result User
	err := c.Find(bson.M{"phone": phone}).One(&result)
	if err != nil {
		return err, result
	}

	return nil, result
}

func (user *User) QueryByMongoID(mongoId string) (error, User) {

	db := mdb.InitDB()
	c := db.C(static.TBL_USERS)

	var result User
	err := c.Find(bson.M{"mongoid": mongoId}).One(&result)
	if err != nil {
		return err, result
	}

	return nil, result
}


func (user *User) UpdateByMongoId(u User) error  {

	updateObject := bson.M{"$set":
		bson.M{
		"password" : u.Password,
		"email": u.Email,
		"displayname" : u.DisplayName,
		"phone" : u.Phone,
		"loginnonce": u.LoginNonce,
		"token" : u.Token,
		"tokenexpirytime" : u.TokenExpiryTime,
		"fcmtoken" : u.FcmToken }}

	err := mdb.UpdateById(static.TBL_USERS, u.MongoID, updateObject)
	if err != nil {
		return err
	}
	return nil
}


func (user *User) HardDelete(u User) error {
	err := mdb.Delete(static.TBL_USERS, u.MongoID)
	if err != nil {
		return err
	}
	log.Println("Hard Deleted")
	return nil
}


