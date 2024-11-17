package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID    string `json:"user_id"`
	Fullname  string `json:"fullname" binding:"required"`
	Email     string `json:"email" binding:"required,email" gorm:"unique"`
	Password  string `json:"password" binding:"required,min=4"`
	CreatedAt time.Time
}

func ping_func(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "PONG",
	})
}

func registerUser(c *gin.Context) {
	// Get and validate user data ...
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return

	}
	newUser.Password = string(hashedPassword)
	AWS_ACC_KEY := os.Getenv("AWS_ACC_KEY")
	AWS_SECRETE_KEY := os.Getenv("AWS_SECRETE_KEY")
	AWS_REGION := os.Getenv("AWS_REGION")

	// Create an AWS session
	awsCreds := credentials.NewStaticCredentials(AWS_ACC_KEY, AWS_SECRETE_KEY, "")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: awsCreds,
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return
	}

	newUser.UserID = uuid.New().String()

	db := dynamodb.New(sess)
	// Create the input for DynamoDB PutItem
	input := &dynamodb.PutItemInput{
		TableName: aws.String("track_money_user"),
		Item: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(newUser.UserID),
			},
			"email": {
				S: aws.String(newUser.Email),
			},
			"fullname": {
				S: aws.String(newUser.Fullname),
			},
			"password": {
				S: aws.String(newUser.Password),
			},
		},
	}

	// Put the item in DynamoDB
	_, err = db.PutItem(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user in DynamoDB"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
