package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	// ToDo, dependency injection
	logger := log.New(os.Stdout, "[technical-exersive-modak][send-mail] ", log.LstdFlags)

	for _, message := range sqsEvent.Records {
		body := message.Body
		logger.Println("message received", body)

		var notification Notification
		err = json.Unmarshal([]byte(body), &notification)

		err = uc.run(notification)
		if err != nil {
			logger.Fatalln("failure in the notification rate limit system")
			return nil
		}
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
