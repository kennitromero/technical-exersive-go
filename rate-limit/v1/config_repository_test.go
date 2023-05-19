package main

import (
	"reflect"
	"testing"
)

func TestConfigRepositoryGetNotificationLimitsByType(t *testing.T) {
	type args struct {
		typeNotification string
	}
	tests := []struct {
		name          string
		args          args
		limitAttempts int
		perSeconds    int
		wantErr       bool
	}{
		{
			name: "should_get_notification_limits_by_type_success",
			args: args{
				typeNotification: TypeNotificationStatus,
			},
			limitAttempts: 2,
			perSeconds:    60,
			wantErr:       false,
		},
		{
			name: "should_get_notification_limits_by_type_fail_type_not_exist",
			args: args{
				typeNotification: "NoExistTypeNotificationStatus",
			},
			limitAttempts: 0,
			perSeconds:    0,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigRepository{}
			limitAttempts, perSeconds, err := c.getNotificationLimitsByType(tt.args.typeNotification)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if limitAttempts != tt.limitAttempts {
				t.Errorf("get() limitAttempts = %v, want %v", limitAttempts, tt.limitAttempts)
			}
			if perSeconds != tt.perSeconds {
				t.Errorf("get() perSeconds = %v, want %v", perSeconds, tt.perSeconds)
			}
		})
	}
}

func TestNewConfigRepository(t *testing.T) {
	tests := []struct {
		name string
		want *ConfigRepository
	}{
		{
			name: "should_construct_success",
			want: &ConfigRepository{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfigRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
