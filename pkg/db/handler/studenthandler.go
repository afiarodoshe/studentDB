package handler

import (
	"errors"
	"github.com/goonode/mogo"
	"labix.org/v2/mgo/bson"
	"regexp"
	"studentDB.go/config"
	"studentDB.go/pkg/db/model"
)

//Studentservice is to handle student relation db query
type Studentservice struct{}

//Create is to register new student
func (studentservice Studentservice) Create(student *model.Student) error {
	conn := config.GetConnection()
	defer conn.Session.Close()
	student.Status= "V"
	if !isEmailValid(student.Email) {
		return errors.New("InvalidEmail")
	}

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

// Delete a student from DB
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



var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}