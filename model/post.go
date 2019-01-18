package model

import (
	"github.com/globalsign/mgo/bson"
	"time"
	"github.com/BTUS/static"
	"github.com/BTUS/mdb"
	"log"
)

type Post struct {
	MongoID bson.ObjectId `bson:"_id,omitempty"`

	UserMongoID bson.ObjectId `bson:"_userid,omitempty"`

	Title string
	Content string
	PostedDate time.Time
	UpdatedDate time.Time
	Like int
	IsReport bool
}


func (p *Post) Insert(post Post ) error {
	err := mdb.Insert(static.TBL_POSTS, post)
	if err  != nil {
		return err
	}
	return nil
}


func (p *Post) QueryByMongoID(mongoId string) (error, Post) {

	db := mdb.InitDB()
	c := db.C(static.TBL_POSTS)

	var result Post
	err := c.Find(bson.M{"mongoid": mongoId}).One(&result)
	if err != nil {
		return err, result
	}

	return nil, result
}


func (p *Post) UpdateByMongoID(post Post) error {
	updateObject := bson.M{"$set":
	bson.M{
		"title" : post.Title,
		"content": post.Content,
		"posteddate" : post.PostedDate,
		"updateddate" : post.UpdatedDate,
		"like": post.Like,
		"isreport" : post.IsReport,
	}}

	err := mdb.UpdateById(static.TBL_USERS, post.MongoID, updateObject)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) HardDelete(post Post) error {
	err := mdb.Delete(static.TBL_POSTS, post.MongoID)
	if err != nil {
		return err
	}
	log.Println("Hard Deleted")
	return nil
}

func (p *Post) AddLike(post Post) error {
	// query like
	err, po := p.QueryByMongoID(post.MongoID.String())
	if err != nil {
		return err
	}

	po.Like ++

	err = p.UpdateByMongoID(po)
	if err != nil {
		return err
	}




}


