package main

import (
	"technical-exersive/commons"
	"time"
)

type IVerifyUC interface {
	handle(typeNotification string, rateLimitCounter *commons.RateLimitCounter) (bool, error)
}

type VerifyUC struct {
	mtime  commons.ImdTime
	config IConfigRepository
}

func NewVerifyUC(mtime commons.ImdTime, config IConfigRepository) *VerifyUC {
	return &VerifyUC{mtime: mtime, config: config}
}

func (c VerifyUC) handle(typeNotification string, rateLimitCounter *commons.RateLimitCounter) (bool, error) {
	limitAttempts, perSeconds, err := c.config.getNotificationLimitsByType(typeNotification)
	if err != nil {
		return false, err
	}

	now := c.mtime.GetNowUTC()

	rateLimitCounter.PerSeconds = perSeconds
	rateLimitCounter.UpdatedAt = now.Format(time.RFC3339)

	previousDate, _ := time.Parse(time.RFC3339, rateLimitCounter.CreatedAt)
	diff := now.Sub(previousDate)
	if diff.Seconds() >= float64(perSeconds) {
		rateLimitCounter.AttemptCounter = 0
		rateLimitCounter.CreatedAt = now.Format(time.RFC3339)
	}

	if rateLimitCounter.AttemptCounter < limitAttempts {
		rateLimitCounter.AttemptCounter++
		return true, nil
	}

	return false, nil
}
