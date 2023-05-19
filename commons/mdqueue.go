package commons

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type ImdQueue interface {
	SendMessage(messageBody string) (string, error)
}

type MDSQSRepository struct {
	queueUrl string
}

func NewMDSQSRepository(queueUrl string) *MDSQSRepository {
	return &MDSQSRepository{queueUrl: queueUrl}
}

func (s *MDSQSRepository) SendMessage(messageBody string) (string, error) {
	// ToDo I should improve it (dependency injection)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsClient := sqs.New(sess)
	queueURL := s.queueUrl

	sendMessageOutput, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    aws.String(queueURL),
	})

	if err != nil {
		handleError(err)
		return "", err
	}

	messageID := aws.StringValue(sendMessageOutput.MessageId)
	return messageID, nil
}

func handleError(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		log.Println("error AWS:", awsErr.Code(), awsErr.Message())
	} else {
		log.Println("error:", err.Error())
	}
}
