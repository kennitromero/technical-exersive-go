package main

import (
	"encoding/json"
	"strings"
	"technical-exersive/commons"
	"time"
)

type OrchestratorUC struct {
	repo      ICacheRepository
	verify    IVerifyUC
	mtime     commons.ImdTime
	queueRepo commons.ImdQueue
}

func NewOrchestratorUC(repo ICacheRepository, verify IVerifyUC, mtime commons.ImdTime, queueRepo commons.ImdQueue) *OrchestratorUC {
	return &OrchestratorUC{repo: repo, verify: verify, mtime: mtime, queueRepo: queueRepo}
}

func (uc *OrchestratorUC) run(notification commons.Notification) error {
	now := uc.mtime.GetNowUTC()

	replacer := strings.NewReplacer(
		FormatCompositeKeyTo, notification.To,
		FormatCompositeKeyWayToNotify, notification.WayToNotify,
		FormatCompositeKeyTypeNotification, notification.TypeNotification,
	)
	compositeKey := replacer.Replace(FormatCompositeKey)

	rateLimitCounter, err := uc.repo.get(compositeKey)
	if err != nil {
		return err
	}

	if rateLimitCounter.Key == "" {
		rateLimitCounter = commons.RateLimitCounter{
			Key:            compositeKey,
			AttemptCounter: 0,
			CreatedAt:      now.Format(time.RFC3339),
		}
	}

	canBeSent, err := uc.verify.handle(notification.TypeNotification, &rateLimitCounter)
	if err != nil {
		return err
	}

	err = uc.repo.set(rateLimitCounter)
	if err != nil {
		return err
	}

	// If it should not be sent, the message is removed from the queue.
	// ToDo, define retry policy (or know if the message is really going to be discarded)
	if !canBeSent {
		return nil
	}

	bodyMessage, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	_, err = uc.queueRepo.SendMessage(string(bodyMessage))
	if err != nil {
		return err
	}

	rateLimitCounter.Key = ""

	return nil
}
