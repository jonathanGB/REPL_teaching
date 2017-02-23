package models

import (
	//"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	GroupModel struct {
		db *mgo.Database
	}

	Group struct {
		Id          bson.ObjectId `bson:"_id"`
		Name        string        `bson:"name"`
		Teacher     bson.ObjectId `bson:"teacher"`
		TeacherName string        `bson:"teacher-name"`
		Files       []File        `bson:"files"`
		Password    []byte        `bson:"password"`
	}

	RenderedGroup struct {
		Id string
		Name string
		TeacherName string
		FilesLen int
	}
)


func NewGroupModel(s *mgo.Session) *GroupModel {
	return &GroupModel{s.DB("repl")}
}

func (gm *GroupModel) GetUserGroups(userId bson.ObjectId) []RenderedGroup {
	userGroups := struct{
		GroupIDs []bson.ObjectId `bson:"groups"`
	}{}
	groups := []Group{}

	gm.db.C("users").Find(bson.M{"_id": userId}).Select(bson.M{"groups": 1, "_id": 0}).One(&userGroups)
	gm.db.C("groups").Find(bson.M{
		"_id": bson.M{"$in": userGroups.GroupIDs},
	}).All(&groups)

	rGroups := []RenderedGroup{}
	for _, group := range groups {
		rGroups = append(rGroups, RenderedGroup{
			group.Id.Hex(),
			group.Name,
			group.TeacherName,
			len(group.Files),
		})
	}

	return rGroups
}