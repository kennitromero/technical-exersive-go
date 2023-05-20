package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"strconv"
	"technical-exersive/commons"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-redis/redis/v8"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// ToDo, dependency injection
	logger := log.New(os.Stdout, "[technical-exersive-modak][rate-limit] ", log.LstdFlags)

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logger.Println("var REDIS_DB", os.Getenv("REDIS_DB"))
		panic("cannot convert REDIS_DB to int")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})

	mCache := commons.NewMDRedisCache(client, ctx)

	repo := NewCacheRepository(mCache)
	config := NewConfigRepository()
	mtime := commons.NewMTime()
	cs := NewVerifyUC(mtime, config)

	queueUrl := os.Getenv("QUEUE_URL")
	queueRepo := commons.NewMDSQSRepository(queueUrl)

	uc := NewOrchestratorUC(repo, cs, mtime, queueRepo)

	for _, message := range sqsEvent.Records {
		body := message.Body
		logger.Println("message received", body)

		var notification commons.Notification
		err = json.Unmarshal([]byte(body), &notification)
		if err != nil {
			logger.Fatalln("failed to convert message to struct")
			return nil
		}

		err = uc.run(notification)
	}

	return nil
}

func main() {
	lambda.Start(handler)

	// Manually Test
	// add at the beginning to not connect with AWS => commons/mdqueue.go:25
	// return "xxx", nil

	/*client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", "127.0.0.1", "6379"),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	mCache := commons.NewMDRedisCache(client, ctx)

	repo := NewCacheRepository(mCache)
	config := NewConfigRepository()
	mtime := commons.NewMTime()
	cs := NewVerifyUC(mtime, config)

	queueUrl := "QUEUE_URL"
	queueRepo := commons.NewMDSQSRepository(queueUrl)

	uc := NewOrchestratorUC(repo, cs, mtime, queueRepo)

	notification := Notification{
		To:               "kennitromero@gmail.com",
		From:             "Henry de Modak",
		Subject:          "Welcome to the Modak Challenge",
		Body:             "You start on June 5",
		WayToNotify:      "email",
		TypeNotification: "status",
		Meta: MetaData{
			LangCode: "en_US",
			Template: "invitation",
		},
	}

	err := uc.run(notification)
	if err != nil {
		fmt.Println("Parece que hay un error")
		panic(err)
	}*/
}
