package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"technical-exersive/commons"
)

type UseCase struct {
	util       commons.ImdUtil
	mail       ISendMail
	currentDir string
}

func NewUseCase(util commons.ImdUtil, mail ISendMail, currentDir string) *UseCase {
	return &UseCase{util: util, mail: mail, currentDir: currentDir}
}

func (uc *UseCase) handle(n commons.Notification) error {
	var payload map[string]interface{}
	err := json.Unmarshal([]byte(n.Body), &payload)
	if err != nil {
		return err
	}

	var filePath string

	templateDir := filepath.Join(uc.currentDir, "templates")

	// ToDo, replace it for strategy map

	switch n.TypeNotification {
	case commons.TypeNotificationStatus:
		filePath = filepath.Join(templateDir, "status.html")
	case commons.TypeNotificationNews:
		filePath = filepath.Join(templateDir, "news.html")
	case commons.TypeNotificationMarketing:
		filePath = filepath.Join(templateDir, "marketing.html")
	case commons.TypeNotificationDaily:
		filePath = filepath.Join(templateDir, "daily.html")
	default:
		return errors.New("template does not exist")
	}

	htmlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	htmlBody, err := uc.util.RenderedToHTML(string(htmlBytes), payload)
	if err != nil {
		return err
	}

	return uc.mail.SendMail(n.To, n.From, n.Subject, htmlBody)
}
