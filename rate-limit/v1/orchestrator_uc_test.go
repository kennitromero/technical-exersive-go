package main

import (
	"errors"
	"reflect"
	"technical-exersive/commons"
	"testing"
	"time"
)

func TestOrchestratorUC_run(t *testing.T) {
	type fields struct {
		CacheRepositoryMock *CacheRepositoryMock
		VerifyUCMock        *VerifyUCMock
		MTimeMock           *commons.MDTimeMock
		MDSQSRepositoryMock *commons.MDSQSRepositoryMock
	}

	commonArgs := Notification{
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

	type args struct {
		notification Notification
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fs fields)
		wantErr bool
	}{
		{
			name: "should_when_rate_limit_exist_success",
			fields: fields{
				CacheRepositoryMock: &CacheRepositoryMock{},
				VerifyUCMock:        &VerifyUCMock{},
				MTimeMock:           &commons.MDTimeMock{},
				MDSQSRepositoryMock: &commons.MDSQSRepositoryMock{},
			},
			args: args{
				notification: commonArgs,
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T07:57:51Z")

				fs.MTimeMock.On("GetNowUTC").Return(nowMock).Once()

				rl := RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#status",
					AttemptCounter: 0,
					CreatedAt:      "2023-05-18T07:57:51Z",
				}

				fs.CacheRepositoryMock.On(
					"get",
					"kennitromero@gmail.com#email#status",
				).Return(RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#status",
					AttemptCounter: 0,
					CreatedAt:      "2023-05-18T07:57:51Z",
				}, nil).Once()

				fs.VerifyUCMock.On(
					"handle",
					"status",
					&rl,
				).Return(true, nil).Once()

				fs.CacheRepositoryMock.On("set", rl).Return(nil).Once()

				fs.MDSQSRepositoryMock.On("SendMessage", "{\"to\":\"kennitromero@"+
					"gmail.com\",\"from\":\"Henry de Modak\",\"subject\":\"Welcome to the Modak "+
					"Challenge\",\"body\":\"You start on June 5\",\"way_to_notify\":\"email\",\"type_"+
					"notification\":\"status\",\"meta\":{\"lang_code\":\"en_US\",\"tem"+
					"plate\":\"invitation\"}}").Return("message_id_x", nil).Once()
			},
			wantErr: false,
		},
		{
			name: "should_when_rate_limit_is_new_success",
			fields: fields{
				CacheRepositoryMock: &CacheRepositoryMock{},
				VerifyUCMock:        &VerifyUCMock{},
				MTimeMock:           &commons.MDTimeMock{},
				MDSQSRepositoryMock: &commons.MDSQSRepositoryMock{},
			},
			args: args{
				notification: commonArgs,
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T07:57:51Z")
				fs.MTimeMock.On("GetNowUTC").Return(nowMock).Once()

				initRl := RateLimitCounter{}
				fs.CacheRepositoryMock.On(
					"get",
					"kennitromero@gmail.com#email#status",
				).Return(initRl, nil).Once()

				filledRl := RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#status",
					AttemptCounter: 0,
					CreatedAt:      "2023-05-18T07:57:51Z",
				}

				fs.VerifyUCMock.On(
					"handle",
					"status",
					&filledRl,
				).Return(true, nil).Once()

				fs.CacheRepositoryMock.On("set", filledRl).Return(nil).Once()

				fs.MDSQSRepositoryMock.On("SendMessage", "{\"to\":\"kennitromero@"+
					"gmail.com\",\"from\":\"Henry de Modak\",\"subject\":\"Welcome to the Modak "+
					"Challenge\",\"body\":\"You start on June 5\",\"way_to_notify\":\"email\",\"type_"+
					"notification\":\"status\",\"meta\":{\"lang_code\":\"en_US\",\"tem"+
					"plate\":\"invitation\"}}").Return("message_id_x", nil).Once()
			},
			wantErr: false,
		},
		{
			name: "should_when_rate_error_get_rate_limit",
			fields: fields{
				CacheRepositoryMock: &CacheRepositoryMock{},
				VerifyUCMock:        &VerifyUCMock{},
				MTimeMock:           &commons.MDTimeMock{},
				MDSQSRepositoryMock: &commons.MDSQSRepositoryMock{},
			},
			args: args{
				notification: commonArgs,
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T07:57:51Z")
				fs.MTimeMock.On("GetNowUTC").Return(nowMock).Once()

				initRl := RateLimitCounter{}
				fs.CacheRepositoryMock.On(
					"get",
					"kennitromero@gmail.com#email#status",
				).Return(initRl, errors.New("any")).Once()
			},
			wantErr: true,
		},
		{
			name: "should_when_rate_limit_error_handle",
			fields: fields{
				CacheRepositoryMock: &CacheRepositoryMock{},
				VerifyUCMock:        &VerifyUCMock{},
				MTimeMock:           &commons.MDTimeMock{},
				MDSQSRepositoryMock: &commons.MDSQSRepositoryMock{},
			},
			args: args{
				notification: commonArgs,
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T07:57:51Z")
				fs.MTimeMock.On("GetNowUTC").Return(nowMock).Once()

				rl := RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#status",
					AttemptCounter: 0,
					CreatedAt:      "2023-05-18T07:57:51Z",
				}

				fs.CacheRepositoryMock.On(
					"get",
					"kennitromero@gmail.com#email#status",
				).Return(RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#status",
					AttemptCounter: 0,
					CreatedAt:      "2023-05-18T07:57:51Z",
				}, nil).Once()

				fs.VerifyUCMock.On(
					"handle",
					"status",
					&rl,
				).Return(true, errors.New("any")).Once()
			},
			wantErr: true,
		},
		{
			name: "should_when_rate_error_set_rate_limit",
			fields: fields{
				CacheRepositoryMock: &CacheRepositoryMock{},
				VerifyUCMock:        &VerifyUCMock{},
				MTimeMock:           &commons.MDTimeMock{},
				MDSQSRepositoryMock: &commons.MDSQSRepositoryMock{},
			},
			args: args{
				notification: commonArgs,
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T07:57:51Z")
				fs.MTimeMock.On("GetNowUTC").Return(nowMock).Once()

				rl := RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#status",
					AttemptCounter: 0,
					CreatedAt:      "2023-05-18T07:57:51Z",
				}

				fs.CacheRepositoryMock.On(
					"get",
					"kennitromero@gmail.com#email#status",
				).Return(RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#status",
					AttemptCounter: 0,
					CreatedAt:      "2023-05-18T07:57:51Z",
				}, nil).Once()

				fs.VerifyUCMock.On(
					"handle",
					"status",
					&rl,
				).Return(true, nil).Once()

				fs.CacheRepositoryMock.On("set", rl).Return(errors.New("any")).Once()
			},
			wantErr: true,
		},
		{
			name: "should_when_rate_error_send_message_to_sqs",
			fields: fields{
				CacheRepositoryMock: &CacheRepositoryMock{},
				VerifyUCMock:        &VerifyUCMock{},
				MTimeMock:           &commons.MDTimeMock{},
				MDSQSRepositoryMock: &commons.MDSQSRepositoryMock{},
			},
			args: args{
				notification: commonArgs,
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T07:57:51Z")
				fs.MTimeMock.On("GetNowUTC").Return(nowMock).Once()

				rl := RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#status",
					AttemptCounter: 0,
					CreatedAt:      "2023-05-18T07:57:51Z",
				}

				fs.CacheRepositoryMock.On(
					"get",
					"kennitromero@gmail.com#email#status",
				).Return(RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#status",
					AttemptCounter: 0,
					CreatedAt:      "2023-05-18T07:57:51Z",
				}, nil).Once()

				fs.VerifyUCMock.On(
					"handle",
					"status",
					&rl,
				).Return(true, nil).Once()

				fs.CacheRepositoryMock.On("set", rl).Return(nil).Once()

				fs.MDSQSRepositoryMock.On("SendMessage", "{\"to\":\"kennitromero@"+
					"gmail.com\",\"from\":\"Henry de Modak\",\"subject\":\"Welcome to the Modak "+
					"Challenge\",\"body\":\"You start on June 5\",\"way_to_notify\":\"email\",\"type_"+
					"notification\":\"status\",\"meta\":{\"lang_code\":\"en_US\",\"tem"+
					"plate\":\"invitation\"}}").Return("", errors.New("any")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.fields)
			l := &OrchestratorUC{
				repo:      tt.fields.CacheRepositoryMock,
				verify:    tt.fields.VerifyUCMock,
				mtime:     tt.fields.MTimeMock,
				queueRepo: tt.fields.MDSQSRepositoryMock,
			}
			err := l.run(tt.args.notification)

			tt.fields.CacheRepositoryMock.AssertExpectations(t)
			tt.fields.VerifyUCMock.AssertExpectations(t)
			tt.fields.MTimeMock.AssertExpectations(t)
			tt.fields.MDSQSRepositoryMock.AssertExpectations(t)

			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewOrchestratorUC(t *testing.T) {
	type args struct {
		repo      ICacheRepository
		verify    IVerifyUC
		mtime     commons.ImdTime
		queueRepo commons.ImdQueue
	}

	a := args{
		repo:      &CacheRepositoryMock{},
		verify:    &VerifyUCMock{},
		mtime:     &commons.MDTimeMock{},
		queueRepo: &commons.MDSQSRepositoryMock{},
	}

	tests := []struct {
		name string
		args args
		want *OrchestratorUC
	}{
		{
			name: "should_construct_success",
			args: args{
				repo:   a.repo,
				verify: a.verify,
				mtime:  a.mtime,
			},
			want: &OrchestratorUC{repo: a.repo, verify: a.verify, mtime: a.mtime},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrchestratorUC(tt.args.repo, tt.args.verify, tt.args.mtime, tt.args.queueRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrchestratorUC() = %v, want %v", got, tt.want)
			}
		})
	}
}
