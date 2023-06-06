package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	db "github.com/mohamedmuhsinJ/loginAssignment/Db"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName" `
	DateOfBirth time.Time `json:"dateofBirth"`
	Email       string    `json:"email" gorm:"unique" validate:"email,required" `
	PhoneNumber string    `json:"phone"`
	Cv          string    `json:"cv"`
}

func Register(c *gin.Context) {

	fName := c.PostForm("firstName")
	lName := c.PostForm("lastName")
	dateOfBirth := c.PostForm("DateOfBirth")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	cvPath, err := c.FormFile("cv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	allowedexe := []string{".pdf", ".doc", ".docx"}
	ext := filepath.Ext(cvPath.Filename)
	isValid := false
	for _, allowedExt := range allowedexe {
		if ext == allowedExt {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	c.SaveUploadedFile(cvPath, "./public/"+cvPath.Filename)

	dOB, err := time.Parse(`"2006-01-02"`, dateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	validPhone := validatePhone(phone)
	if !validPhone {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Phone number",
		})
		return
	}

	user := User{
		FirstName:   fName,
		LastName:    lName,
		Email:       email,
		PhoneNumber: phone,
		DateOfBirth: dOB,
		Cv:          cvPath.Filename,
	}
	rec := db.Db.Create(&user)

	if rec.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store data",
		})
		return
	}
	sendEmailConfirmation(user.Email)
	c.JSON(http.StatusOK, user)
}

func validatePhone(number string) bool {
	regPattern := `^\d{10}$`
	valid, err := regexp.MatchString(regPattern, number)
	if err != nil {
		panic(err)

	}
	return valid
}

func sendEmailConfirmation(email string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env")
	}
	// smtpServer := os.Getenv("SMTP_SERVER")
	// smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	subject := "Please confirm your email"

	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{email}, []byte(subject))
	if err != nil {

		fmt.Printf("failed to connect smtp")
	}

}

func Home(c *gin.Context) {
	email := c.Param("email")
	var user User
	db.Db.First(&user, email)
	if user.ID == 0 {
		c.JSON(400, gin.H{
			"error": "user doesnot exists",
		})

		return

	}
	cv := user.Cv
	ext := filepath.Ext(cv)
	cvcontent, err := ioutil.ReadFile("./public/" + cv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctype := getFileContentType(ext)
	c.Data(http.StatusOK, ctype, cvcontent)

}

func getFileContentType(Ext string) string {
	switch Ext {
	case ".pdf":
		return "application/pdf"
	case ".doc", ".docx":
		return "application/msword"
	default:
		return "application/octet-stream"
	}
}
