package main

import (
	"fmt"
	"github.com/go-mail/mail"
	"github.com/stretchr/testify/mock"
)

type ISendMail interface {
	SendMail(to string, from string, subject string, body string) error
}

type SendMailMock struct {
	mock.Mock
}

func (sm *SendMailMock) SendMail(to string, from string, subject string, body string) error {
	args := sm.Called(to, from, subject, body)
	return args.Error(0)
}

type SendMail struct {
	host     string
	port     int
	username string
	password string
}

func NewSendMail(host string, port int, username string, password string) *SendMail {
	return &SendMail{host: host, port: port, username: username, password: password}
}

func (sm *SendMail) SendMail(to string, from string, subject string, body string) error {

	fmt.Println("To", to)
	fmt.Println("From", from)

	fmt.Println("Subject", subject)
	fmt.Println("text/html", body)

	fmt.Println(sm.host, sm.port, sm.username, sm.password)
	return nil

	m := mail.NewMessage()

	m.SetHeader("To", to)
	m.SetHeader("From", from)

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := mail.NewDialer(sm.host, sm.port, sm.username, sm.password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
