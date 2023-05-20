package main

import (
	"os"
	"reflect"
	"technical-exersive/commons"
	"testing"
)

func TestDetermineTemplateUC_handle(t *testing.T) {
	type fields struct {
		util       *commons.MdUtil
		mail       *SendMailMock
		currentDir string
	}
	type args struct {
		notification commons.Notification
	}
	currentDir, _ := os.Getwd()

	tests := []struct {
		name    string
		fields  fields
		mock    func(fs fields)
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should_get_html_template_status_rendered_success",
			fields: fields{
				util:       &commons.MdUtil{},
				mail:       &SendMailMock{},
				currentDir: currentDir,
			},
			mock: func(fs fields) {
				fs.mail.On(
					"SendMail",
					"henry@modak.live",
					"kennitromero@gmail.com",
					"Hi",
					"<!DOCTYPE html><html lang=\"es\"><head>  "+
						"  <title>Status Template Strategy</title></head><body><h1>Hi, "+
						"this is email status</h1><p>New Status: delivered</p></body></html>",
				).Return(nil).Once()
			},
			args: args{
				notification: commons.Notification{
					To:               "henry@modak.live",
					From:             "kennitromero@gmail.com",
					Subject:          "Hi",
					Body:             "{\"status\":\"delivered\"}",
					TypeNotification: commons.TypeNotificationStatus,
				},
			},
			wantErr: false,
		},
		{
			name: "should_get_html_template_news_rendered_success",
			fields: fields{
				util: &commons.MdUtil{},
				mail: &SendMailMock{},
			},
			mock: func(fs fields) {
				fs.mail.On(
					"SendMail",
					"henry@modak.live",
					"kennitromero@gmail.com",
					"Hi",
					"<!DOCTYPE html><html lang=\"es\"><head>    <title>News Template Strategy</title></head><body><h1>Hi, this is email news</h1><p>We have a good new: you&#39;re hired Kennit</p></body></html>",
				).Return(nil).Once()
			},
			args: args{
				notification: commons.Notification{
					To:               "henry@modak.live",
					From:             "kennitromero@gmail.com",
					Subject:          "Hi",
					Body:             "{\"new\":\"you're hired Kennit\"}",
					TypeNotification: commons.TypeNotificationNews,
				},
			},
			wantErr: false,
		},
		{
			name: "should_get_html_template_daily_rendered_success",
			fields: fields{
				util: &commons.MdUtil{},
				mail: &SendMailMock{},
			},
			mock: func(fs fields) {
				fs.mail.On(
					"SendMail",
					"henry@modak.live",
					"kennitromero@gmail.com",
					"Hi",
					"<!DOCTYPE html><html><head>    <title>Daily Template Strategy</title></head><body><h1>Hi, this is email daily</h1><p>Hey, the daily is I&#39;m was working Hard 😅</p></body></html>",
				).Return(nil).Once()
			},
			args: args{
				notification: commons.Notification{
					To:               "henry@modak.live",
					From:             "kennitromero@gmail.com",
					Subject:          "Hi",
					Body:             "{\"daily\":\"I'm was working Hard 😅\"}",
					TypeNotification: commons.TypeNotificationDaily,
				},
			},
			wantErr: false,
		},
		{
			name: "should_get_html_template_marketing_rendered_success",
			fields: fields{
				util: &commons.MdUtil{},
				mail: &SendMailMock{},
			},
			mock: func(fs fields) {
				fs.mail.On(
					"SendMail",
					"henry@modak.live",
					"kennitromero@gmail.com",
					"Hi",
					"<!DOCTYPE html><html lang=\"es\"><head>    <title>Marketing Template Strategy</title></head><body><h1>Hi, this is email marketing</h1><p>We have a best price 10.500</p></body></html>",
				).Return(nil).Once()
			},
			args: args{
				notification: commons.Notification{
					To:               "henry@modak.live",
					From:             "kennitromero@gmail.com",
					Subject:          "Hi",
					Body:             "{\"price\":\"10.500\"}",
					TypeNotification: commons.TypeNotificationMarketing,
				},
			},
			want:    "<!DOCTYPE html><html lang=\"es\"><head>    <title>Marketing Template Strategy</title></head><body><h1>Hi, this is email marketing</h1><p>We have a best price 10.500</p></body></html>",
			wantErr: false,
		},
		{
			name: "should_cannot_get_html_template_because_body_is_bad",
			fields: fields{
				util: &commons.MdUtil{},
				mail: &SendMailMock{},
			},
			mock: func(fs fields) {},
			args: args{
				notification: commons.Notification{
					Body:             "{",
					TypeNotification: commons.TypeNotificationMarketing,
				},
			},
			wantErr: true,
		},
		{
			name: "should_cannot_get_html_template_because_template_not_exits",
			fields: fields{
				util: &commons.MdUtil{},
				mail: &SendMailMock{},
			},
			mock: func(fs fields) {},
			args: args{
				notification: commons.Notification{
					Body:             "{\"price\":\"10.500\"}",
					TypeNotification: "TypeNotificationNoExist",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.fields)
			oss := &UseCase{
				util:       tt.fields.util,
				mail:       tt.fields.mail,
				currentDir: tt.fields.currentDir,
			}
			err := oss.handle(tt.args.notification)
			if (err != nil) != tt.wantErr {
				t.Errorf("handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.fields.mail.AssertExpectations(t)
		})
	}
}

func TestNewUseCase(t *testing.T) {
	type args struct {
		util       commons.ImdUtil
		mail       *SendMailMock
		currentDir string
	}

	a := args{
		mail: &SendMailMock{},
		util: &commons.MdUtil{},
	}

	tests := []struct {
		name string
		args args
		want *UseCase
	}{
		{
			name: "should_construct_use_case_success",
			args: args{mail: a.mail, util: a.util},
			want: &UseCase{mail: a.mail, util: a.util},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUseCase(tt.args.util, tt.args.mail, tt.args.currentDir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
