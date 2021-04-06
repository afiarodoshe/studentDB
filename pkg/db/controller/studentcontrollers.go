package controller

import (
	"log"
	service "studentDB.go/pkg/db/handler"
	"studentDB.go/pkg/db/model"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

//AuthController is for auth logic
type AuthController struct{}

//Login is to process login request
func (auth *AuthController) Login(c *gin.Context) {

	var loginInfo model.Student
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	//TODO
	studentService := service.Studentservice{}
	student, errf := studentService.Find(&loginInfo)
	if errf != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "Not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(loginInfo.Password))
	if err != nil {
		c.AbortWithStatusJSON(402, gin.H{"error": "Email or password is invalid."})
		return
	}

	token, err := student.GetJwtToken()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	//-------
	c.JSON(200, gin.H{
		"token": token,
	})
}

//Profile is to provide current user info
func (auth *AuthController) Profile(c *gin.Context) {
	user := c.MustGet("student").(*model.Student)

	c.JSON(200, gin.H{
		"user_name": user.Name,
		"email":     user.Email,
	})
}

//Signup is for student signup
func (auth *AuthController) Signup(c *gin.Context) {

	type signupInfo struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"`
	}
	var info signupInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "Please input all fields"})
		return
	}
	student := model.Student{}
	student.Email = info.Email
	hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
		return
	}

	student.Password = string(hash)
	student.Name = info.Name
	userservice := service.Studentservice{}
	err = userservice.Create(&student)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"result": "ok"})
	}
	return
}