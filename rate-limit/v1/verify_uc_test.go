package main

import (
	"errors"
	"reflect"
	"technical-exersive/commons"
	"testing"
	"time"
)

func TestNewLimitRateStatusStrategy(t *testing.T) {
	type args struct {
		mtime  *commons.MDTimeMock
		config *ConfigRepositoryMock
	}

	a := args{
		mtime:  &commons.MDTimeMock{},
		config: &ConfigRepositoryMock{},
	}

	tests := []struct {
		name string
		args args
		want *VerifyUC
	}{
		{
			name: "should_construct_success",
			args: args{
				mtime:  a.mtime,
				config: a.config,
			},
			want: &VerifyUC{mtime: a.mtime, config: a.config},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVerifyUC(tt.args.mtime, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVerifyUC() = %v, want %v", got, tt.want)
			}
		})
	}
}

const (
	TestLimitAttempts int = 2
	TestPerSeconds    int = 60
)

func TestVerifyUC_handle(t *testing.T) {
	type fields struct {
		mtime  *commons.MDTimeMock
		config *ConfigRepositoryMock
	}
	type args struct {
		typeNotification string
		rateLimitCounter *RateLimitCounter
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		mock               func(fs fields)
		want               bool
		wantAttemptCounter int
		wantErr            bool
	}{
		// counter is zero and time is less
		{
			name: "should_verify_counter_is_zero_and_time_is_less",
			fields: fields{
				mtime:  &commons.MDTimeMock{},
				config: &ConfigRepositoryMock{},
			},
			args: args{
				typeNotification: TypeNotificationStatus,
				rateLimitCounter: &RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#project_invitations",
					AttemptCounter: 0,
					PerSeconds:     60,
					CreatedAt:      "2023-05-18T12:43:50Z",
				},
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T12:43:59Z")
				fs.mtime.On("GetNowUTC").Return(nowMock).Once()

				fs.config.On("getNotificationLimitsByType", TypeNotificationStatus).
					Return(TestLimitAttempts, TestPerSeconds, nil).Once()
			},
			want:               true,
			wantAttemptCounter: 1,
			wantErr:            false,
		},
		{
			// counter and time is less
			name: "should_verify_counter_and_time_is_less",
			fields: fields{
				mtime:  &commons.MDTimeMock{},
				config: &ConfigRepositoryMock{},
			},
			args: args{
				typeNotification: TypeNotificationStatus,
				rateLimitCounter: &RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#project_invitations",
					AttemptCounter: 1,
					PerSeconds:     60,
					CreatedAt:      "2023-05-18T12:43:50Z",
				},
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T12:43:59Z")
				fs.mtime.On("GetNowUTC").Return(nowMock).Once()

				fs.config.On("getNotificationLimitsByType", TypeNotificationStatus).
					Return(TestLimitAttempts, TestPerSeconds, nil).Once()
			},
			want:               true,
			wantAttemptCounter: 2,
			wantErr:            false,
		},
		{
			// counter is greater and time is less
			name: "should_verify_counter_is_greater_and_time_is_less",
			fields: fields{
				mtime:  &commons.MDTimeMock{},
				config: &ConfigRepositoryMock{},
			},
			args: args{
				typeNotification: TypeNotificationStatus,
				rateLimitCounter: &RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#project_invitations",
					AttemptCounter: 2,
					PerSeconds:     60,
					CreatedAt:      "2023-05-18T12:43:50Z",
				},
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T12:43:59Z")
				fs.mtime.On("GetNowUTC").Return(nowMock).Once()

				fs.config.On("getNotificationLimitsByType", TypeNotificationStatus).
					Return(TestLimitAttempts, TestPerSeconds, nil).Once()
			},
			want:               false,
			wantAttemptCounter: 2,
			wantErr:            false,
		},
		{
			// counter is less and time is greater
			name: "should_verify_counter_is_less_and_time_is_greater",
			fields: fields{
				mtime:  &commons.MDTimeMock{},
				config: &ConfigRepositoryMock{},
			},
			args: args{
				typeNotification: TypeNotificationStatus,
				rateLimitCounter: &RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#project_invitations",
					AttemptCounter: 1,
					PerSeconds:     60,
					CreatedAt:      "2023-05-18T12:43:50Z",
				},
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T12:48:59Z")
				fs.mtime.On("GetNowUTC").Return(nowMock).Once()

				fs.config.On("getNotificationLimitsByType", TypeNotificationStatus).
					Return(TestLimitAttempts, TestPerSeconds, nil).Once()
			},
			want:               true,
			wantAttemptCounter: 1,
			wantErr:            false,
		},
		{
			// counter is greater and time is greater
			name: "should_verify_counter_is_greater_and_time_is_greater",
			fields: fields{
				mtime:  &commons.MDTimeMock{},
				config: &ConfigRepositoryMock{},
			},
			args: args{
				typeNotification: TypeNotificationStatus,
				rateLimitCounter: &RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#project_invitations",
					AttemptCounter: 4,
					PerSeconds:     60,
					CreatedAt:      "2023-05-18T12:43:50Z",
				},
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T12:49:59Z")
				fs.mtime.On("GetNowUTC").Return(nowMock).Once()

				fs.config.On("getNotificationLimitsByType", TypeNotificationStatus).
					Return(TestLimitAttempts, TestPerSeconds, nil).Once()
			},
			want:               true,
			wantAttemptCounter: 1,
			wantErr:            false,
		},
		{
			// error when try did get config
			name: "should_verify_error_when_try_did_get_config",
			fields: fields{
				mtime:  &commons.MDTimeMock{},
				config: &ConfigRepositoryMock{},
			},
			args: args{
				typeNotification: TypeNotificationStatus,
				rateLimitCounter: &RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#project_invitations",
					AttemptCounter: 4,
					PerSeconds:     60,
					CreatedAt:      "2023-05-18T12:43:50Z",
				},
			},
			mock: func(fs fields) {
				nowMock, _ := time.Parse(time.RFC3339, "2023-05-18T12:49:59Z")
				fs.mtime.On("GetNowUTC").Return(nowMock).Once()

				fs.config.On("getNotificationLimitsByType", TypeNotificationStatus).
					Return(0, 0, errors.New("there no exit type")).Once()
			},
			want:               false,
			wantAttemptCounter: 0,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.fields)
			c := VerifyUC{
				mtime:  tt.fields.mtime,
				config: tt.fields.config,
			}
			got, err := c.handle(tt.args.typeNotification, tt.args.rateLimitCounter)
			if (err != nil) != tt.wantErr {
				t.Errorf("handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handle() got = %v, want %v", got, tt.want)
			}
		})
	}
}
