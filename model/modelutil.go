package model

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"strings"
	"github.com/InclusION/util"
)

func DecodeRequestIntoFcmMessage(w http.ResponseWriter, r *http.Request) (error, FcmMessage) {

	decoder := json.NewDecoder(r.Body)

	var mess FcmMessage
	err := decoder.Decode(&mess)
	if err != nil {
		return err, mess
	}

	return nil, mess
}

func DecodeRequestIntoPhone(w http.ResponseWriter, r *http.Request) (error, Phone) {

	decoder := json.NewDecoder(r.Body)

	var phone Phone
	err := decoder.Decode(&phone)
	if err != nil {
		return err, phone
	}

	return nil, phone
}

func DecodeRequestIntoUser(w http.ResponseWriter, r *http.Request) (error, User) {

	decoder := json.NewDecoder(r.Body)

	var user User
	err := decoder.Decode(&user)
	if err != nil {
		return err, user
	}

	return nil, user
}

func DecodeRequestIntoRoom(w http.ResponseWriter, r *http.Request) (error, Room) {

	decoder := json.NewDecoder(r.Body)

	var room Room
	err := decoder.Decode(&room)
	if err != nil {
		return err, room
	}

	return nil, room
}

func DecodeRequestIntoHealth(w http.ResponseWriter, r *http.Request) (error, Health) {

	decoder := json.NewDecoder(r.Body)

	var health Health
	err := decoder.Decode(&health)
	if err != nil {
		return err, health
	}

	return nil, health
}

func DecodeRequestIntoCommand(cmd string) (error, Command) {
	r := strings.NewReader(cmd)
	decoder := json.NewDecoder(r)

	var c Command
	err := decoder.Decode(&c)

	if err != nil {
		return err, c
	}

	return nil, c
}

func CheckAuth(email string, phone string, token string) bool {


	var checkType = ""

	// query user information
	u := User{}

	if len(email) > 0 {
		checkType = "EMAIL"
	}

	if len(phone) > 0 {
		checkType = "PHONE"

	}

	var err error
	var user User

	if checkType == "PHONE" {
		err, user = u.QueryByPhone(phone)
	} else {
		err, user = u.QueryByEmail(email)
	}

	if err != nil {
		return false
	}

	log.Println(user.Phone)
	log.Println(user.Password)
	log.Println(user.Email)
	log.Println(user.LoginNonce)


	// hash then compare with current hash
	t := util.Hash(fmt.Sprintf("%s%s%s%d", user.Phone, user.Password, user.Email, user.LoginNonce))
	log.Printf("complied token: ", t)
	log.Printf("client token: ", token)
	isEqual := strings.Compare(t, token)

	if isEqual != 0 {
		return false
	}

	return true
}


func HideSensitiveUser(user *User) {
	user.Password = ""
	user.LoginNonce = 0
}

func HideSensitiveHealth(health *Health) {
	health.Token = ""
}


func HideSensitiveData(input interface{}) interface{} {

	// if model = User
	if user, ok := input.(User) ; ok {
		log.Println("hideSensitiveData: User")
		user.Token = ""
		user.Password = ""
		user.LoginNonce = 00

		return user
	}

	// if model = Health
	if health, ok := input.(Health) ; ok {
		log.Println("hideSensitiveData: Health")
		health.Token = ""

		return health
	}

	return nil
}
