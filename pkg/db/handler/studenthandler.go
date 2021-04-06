package handler

import (
	"errors"
	"github.com/goonode/mogo"
	"labix.org/v2/mgo/bson"
	"studentDB.go/config"
	"studentDB.go/pkg/db/model"
)

//Studentservice is to handle user relation db query
type Studentservice struct{}

//Create is to register new user
func (studentservice Studentservice) Create(student *model.Student) error {
	conn := config.GetConnection()
	defer conn.Session.Close()

	doc := mogo.NewDoc(model.Student{}).(*model.Student)
	err := doc.FindOne(bson.M{"email": student.Email}, doc)
	if err == nil {
		return errors.New("AlreadyExist")
	}
	studentModel := mogo.NewDoc(student).(*model.Student)
	err = mogo.Save(studentModel)
	if vErr, ok := err.(*mogo.ValidationError); ok {
		return vErr
	}
	return err
}

// Delete a user from DB
func (studentservice Studentservice) Delete(email string) error {
	student, _ := studentservice.FindByEmail(email)
	conn := config.GetConnection()
	defer conn.Session.Close()
	err := student.Remove()
	return err
}

//Find user
func (studentservice Studentservice) Find(student *model.Student) (*model.Student, error) {
	conn := config.GetConnection()
	defer conn.Session.Close()

	doc := mogo.NewDoc(model.Student{}).(*model.Student)
	err := doc.FindOne(bson.M{"email": student.Email}, doc)

	if err != nil {
		return nil, err
	}
	return doc, nil
}

//Find user from email
func (studentservice Studentservice) FindByEmail(email string) (*model.Student, error) {
	conn := config.GetConnection()
	defer conn.Session.Close()

	student := new(model.Student)
	student.Email = email
	return studentservice.Find(student)
}
