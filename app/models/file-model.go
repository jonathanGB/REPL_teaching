package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sort"
)

type (
	FileModel struct {
		db *mgo.Database
	}

	File struct {
		Id        bson.ObjectId `bson:"_id"`
		Name      string        `bson:"name"`
		Owner     bson.ObjectId `bson:"owner"`
		OwnerName string        `bson:"ownerName"`
		Extension string        `bson:"extension"`
		Content   []byte        `bson:"content"`
		IsPrivate bool          `bson:"isPrivate"`
	}

	RenderedFile struct {
		Id        string // ObjectIdHex
		Name      string
		Owner     string // ObjectIdHex
		OwnerName string
		Extension string
		IsPrivate bool
	}

	StructuredRenderedFiles struct{
		TeacherFiles []RenderedFile
		StudentKeys []string
		StudentFiles map[string][]RenderedFile
	}
)

func NewFileModel(s *mgo.Session) *FileModel {
	return &FileModel{s.DB("repl")}
}

func (fm *FileModel) GetGroupFiles(tId, gId, uId bson.ObjectId, role string) *StructuredRenderedFiles {
	result := struct {
		Files []File `bson:"files"`
	}{}
	sResult := StructuredRenderedFiles{
		[]RenderedFile{},
		[]string{},
		map[string][]RenderedFile{},
	}
	studentsName := []string{}
	studentsNameToKey := map[string]string{}

	fm.db.C("groups").Find(bson.M{"_id": gId}).Select(bson.M{"files.content": 0, "_id": 0}).One(&result)

	for _, file := range result.Files {
		// visible file
		if file.Owner == uId || (role == "teacher" || file.Owner == tId) && !file.IsPrivate {
			rFile := RenderedFile{
				file.Id.Hex(),
				file.Name,
				file.Owner.Hex(),
				file.OwnerName,
				file.Extension,
				file.IsPrivate,
			}

			if file.Owner == tId { // teacher's file
				sResult.TeacherFiles = append(sResult.TeacherFiles, rFile)
			} else { // student file
				// new student name, initialize some structures
				if _, ok := sResult.StudentFiles[rFile.Owner]; !ok {
					studentsName = append(studentsName, rFile.OwnerName)
					studentsNameToKey[rFile.OwnerName] = rFile.Owner
					sResult.StudentFiles[rFile.Owner] = []RenderedFile{}
				}

				sResult.StudentFiles[rFile.Owner] = append(sResult.StudentFiles[rFile.Owner], rFile)
			}
		}
	}

	// sort slices
	sort.Strings(studentsName)
	for _, studentName := range studentsName {
		sResult.StudentKeys = append(sResult.StudentKeys, studentsNameToKey[studentName])
	}

	return &sResult
}

func (fm *FileModel) IsThereUserFile(fileName string, gId, uId bson.ObjectId) bool {
	result := struct {
		Id bson.ObjectId `bson:"_id"`
	}{}

	fm.db.C("groups").Find(bson.M{
		"_id":         gId,
		"files.owner": uId,
		"files.name":  fileName,
	}).Select(bson.M{"_id": 1}).One(&result)

	return result.Id != ""
}

func (fm *FileModel) AddFile(file *File, gId bson.ObjectId) error {
	return fm.db.C("groups").Update(
		bson.M{"_id": gId},
		bson.M{
			"$push": bson.M{"files": file},
		},
	)
}

func (fm *FileModel) GetFile(fId, gId bson.ObjectId) (*File, error) {
	result := struct {
		Files []File `bson:"files"`
	}{}

	fm.db.C("groups").Find(bson.M{"_id": gId, "files._id": fId}).Select(bson.M{"files.$": 1, "_id": 0}).One(&result)

	if len(result.Files) != 1 || result.Files[0].Id == "" {
		return nil, fmt.Errorf("no file found")
	}

	return &result.Files[0], nil
}
