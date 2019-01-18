package main

import (
	"log"
	"time"
	"net/http"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"github.com/BTUS/util"
	"github.com/BTUS/model"
	"github.com/BTUS/mdb"
	"github.com/BTUS/static"
)


//**********************************************************************************//
// http requests
func testConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, your connection is fine. %s!", r.URL.Path[1:])
}


func Register(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, user := model.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user.MongoID = bson.NewObjectId()
	user.IsVerify = false

	err = user.Insert(user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	uu := mdb.QueryById(static.TBL_USERS, user.MongoID)
	json.NewEncoder(w).Encode(&uu)
}


func Login(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, user := model.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}


	var u model.User

	// check if has email -> login by email
	if len(user.Email) > 0 {
		err, u = user.QueryByEmailPassword(user.Email, user.Password)
	}

	// check if has phone -> login by phone
	if len(user.Phone) > 0 {
		err, u = user.QueryByPhonePassword(user.Phone, user.Password)
	}

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// if Token expired, create new token
	// hash the username and password and timestamp
	// then insert the hash into user.Token, and user.TokenExpiryTime
	if u.TokenExpiryTime.Before(time.Now().UTC()) || u.TokenExpiryTime.IsZero() {

		u.LoginNonce++
		u.TokenExpiryTime = time.Now().UTC().Add(time.Hour * 24 * 10)

		err := u.UpdateByMongoId(u)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}

	// create token and return to client
	// because login, so that we have to return Token, cannot hide it
	u.Token = util.Hash(fmt.Sprintf("%s%s%s%d", u.Phone, u.Password, u.Email, u.LoginNonce))

	model.HideSensitiveUser(&u)

	json.NewEncoder(w).Encode(&u)
}


func Logout(w http.ResponseWriter, r *http.Request) {
	// disable push notification

	// disable login token


}


func ReportUser(w http.ResponseWriter, r *http.Request) {
	// add a warning message into this user profile -> so that other users will see the reported messages

}



func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, rUser := model.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	isAuth := model.CheckAuth(rUser.Email, rUser.Phone, rUser.Token)
	if isAuth == false {
		http.Error(w, "Authentication fail.", 400)
		return
	}

	if len(rUser.MongoID) > 0 {
		err = rUser.UpdateByMongoId(rUser)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}

	model.HideSensitiveUser(&rUser)

	json.NewEncoder(w).Encode(&rUser)
}


func LoadProfile(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, u := model.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	isAuth := model.CheckAuth(u.Email, u.Phone, u.Token)
	if isAuth == false {
		http.Error(w, "Authentication fail.", 400)
		return
	}

	var mongoId = ""

	err, user := u.QueryByMongoID(mongoId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	model.HideSensitiveUser(&user)

	json.NewEncoder(w).Encode(&user)
}





func getAllBlogWithPaging(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	pageNum := params["no"]
	log.Printf("page %s", pageNum)

	i, e := strconv.Atoi(pageNum)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}

	rBlog := model.Blog{}
	err, blogs := rBlog.QueryAllPaging(i)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&blogs)
}


func getBlogById(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	blogId := params["id"]
	//log.Printf("blog id %s", blogId)

	rBlog := model.NewBlog(bson.ObjectIdHex(blogId))

	err, blog := rBlog.QueryById()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&blog)
}

func findUser(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, rUser := model.DecodeRequestIntoUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err, user := chat.FindUser(rUser.Username)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&user)
}


func createChannel11(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, rRoom := model.DecodeRequestIntoRoom(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err, room11 := chat.CreateChannel11(rRoom)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// go routine push notification
	go fcm.NotifyAllDevicesOfUser(rRoom.Username2, room11.ChannelId)

	json.NewEncoder(w).Encode(&room11)
}

func addToken(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, p := model.DecodeRequestIntoPhone(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	p.MongoID = bson.NewObjectId()

	err = p.Insert()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	phone := mdb.QueryById(static.TBL_DEVICES, p.MongoID)

	json.NewEncoder(w).Encode(&phone)
}


func removeToken(w http.ResponseWriter, r *http.Request) {

	util.CheckBodyNil(w, r)

	err, p := model.DecodeRequestIntoPhone(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = p.HardDelete()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&p)
}


func pushToDevice(w http.ResponseWriter, r *http.Request) {
	// device token
	// data to push

	util.CheckBodyNil(w, r)

	err, mess := model.DecodeRequestIntoFcmMessage(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if len(mess.DeviceToken) <= 0 {
		http.Error(w, "Device token cannot be nil", 400)
		return
	}

	err = fcm.Notify1Device(mess.DeviceToken, mess.Content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&mess)
}


func pushToUser(w http.ResponseWriter, r *http.Request) {
	// device token
	// data to push

	util.CheckBodyNil(w, r)

	err, mess := model.DecodeRequestIntoFcmMessage(w, r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if len(mess.Username) <= 0 {
		http.Error(w, "Username cannot be nil", 400)
		return
	}

	err = fcm.NotifyAllDevicesOfUser(mess.Username, mess.Content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(&mess)
}

// http requests
//**********************************************************************************//
