package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/goonode/mogo"
	"studentDB.go/utils"
	"time"
)

type Student struct {
	mogo.DocumentModel `bson:",inline" coll:"users"`
	Email              string `idx:"{email},unique" json:"email" binding:"required"`
	Password           string `json:"password" binding:"required"`
	Name               string `json:"name"`
	// CreatedAt          *time.Time
	// UpdatedAt          *time.Time
	VerifiedAt *time.Time
}

//GetJwtToken returns jwt token with student email claims
func (student *Student) GetJwtToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": string(student.Email),
	})
	secretKey := utils.EnvVar("TOKEN_KEY", "")
	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}

func init() {
	mogo.ModelRegistry.Register(Student{})
}
