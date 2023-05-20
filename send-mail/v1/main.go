package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"strconv"
	"technical-exersive/commons"
)

func handler(_ context.Context, sqsEvent events.SQSEvent) error {

	// ToDo, dependency injection
	logger := log.New(os.Stdout, "[technical-exercise-modak][send-mail] ", log.LstdFlags)

	util := commons.NewMdUtil()

	mailHost := os.Getenv("MAIL_HOST")
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		return err
	}
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	mail := NewSendMail(mailHost, mailPort, mailUsername, mailPassword)
	uc := NewUseCase(util, mail, currentDir+"/send-mail/v1")

	for _, message := range sqsEvent.Records {
		body := message.Body
		logger.Println("message received", body)

		var notification commons.Notification
		err := json.Unmarshal([]byte(body), &notification)
		if err != nil {
			logger.Fatalln("failed conversion of message to JSON", err)
			return nil
		}

		err = uc.handle(notification)
		if err != nil {
			logger.Fatalln("failed to send email", err)
			return nil
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)

	/*
		currentDir, err := os.Getwd()

		logger := log.New(os.Stdout, "modak", log.LstdFlags)
		util := commons.NewMdUtil()
		mail := NewSendMail("host", 1234, "kennitromero@gmail.com", "p123")

		uc := NewUseCase(util, mail, currentDir+"/send-mail/v1")

		notification := commons.Notification{
			To:               "kennitromero@gmail.com",
			From:             "Henry de Modak",
			Subject:          "Welcome to the Modak Challenge",
			Body:             "{\"status\":\"delivered\"}",
			WayToNotify:      "email",
			TypeNotification: "status",
			Meta: commons.MetaData{
				LangCode: "en_US",
				Template: "invitation",
			},
		}

		err = uc.handle(notification)
		if err != nil {
			logger.Fatalln("failure", err.Error())
		}
	*/
}
