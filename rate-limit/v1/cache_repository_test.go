package main

import (
	"errors"
	"reflect"
	"technical-exersive/commons"
	"testing"
)

func TestNewCacheRepository(t *testing.T) {
	type args struct {
		cache *commons.MDCacheMock
	}

	a := args{
		cache: &commons.MDCacheMock{},
	}

	tests := []struct {
		name string
		args args
		want *CacheRepository
	}{
		{
			name: "should_construct_success",
			args: args{
				cache: a.cache,
			},
			want: &CacheRepository{cache: a.cache},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCacheRepository(tt.args.cache); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCacheRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRateLimitCacheRepository_get(t *testing.T) {
	type fields struct {
		cache *commons.MDCacheMock
	}
	type args struct {
		compositeKey string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		mock    func(fs fields)
		want    commons.RateLimitCounter
		wantErr bool
	}{
		{
			name: "should_get_and_found_success",
			args: args{
				compositeKey: "kennitromero@gmail.com#email#news",
			},
			fields: fields{
				cache: &commons.MDCacheMock{},
			},
			mock: func(fs fields) {
				fs.cache.On(
					"Get",
					"kennitromero@gmail.com#email#news",
				).Return("{\"key\":\"kennitromero@gmail.com#email"+
					"#news\",\"attempt_counter\":1,\"per_seconds\":60,\"cre"+
					"ated_at\":\"2023-05-18T12:43:50Z\"}", nil,
				).Once()
			},
			want: commons.RateLimitCounter{
				Key:            "kennitromero@gmail.com#email#news",
				AttemptCounter: 1,
				PerSeconds:     60,
				CreatedAt:      "2023-05-18T12:43:50Z",
			},
			wantErr: false,
		},
		{
			name: "should_get_and_not_found_success",
			args: args{
				compositeKey: "kennitromero@gmail.com#email#news",
			},
			fields: fields{
				cache: &commons.MDCacheMock{},
			},
			mock: func(fs fields) {
				fs.cache.On(
					"Get",
					"kennitromero@gmail.com#email#news",
				).Return("", nil).Once()
			},
			want:    commons.RateLimitCounter{},
			wantErr: false,
		},
		{
			name: "should_get_fail_error_cache_library",
			args: args{
				compositeKey: "kennitromero@gmail.com#email#news",
			},
			fields: fields{
				cache: &commons.MDCacheMock{},
			},
			mock: func(fs fields) {
				fs.cache.On(
					"Get",
					"kennitromero@gmail.com#email#news",
				).Return("", errors.New("any")).Once()
			},
			want:    commons.RateLimitCounter{},
			wantErr: true,
		},
		{
			name: "should_get_fail_bad_json_unmarshal",
			args: args{
				compositeKey: "kennitromero@gmail.com#email#news",
			},
			fields: fields{
				cache: &commons.MDCacheMock{},
			},
			mock: func(fs fields) {
				fs.cache.On(
					"Get",
					"kennitromero@gmail.com#email#news",
				).Return("{{{{", nil).Once()
			},
			want:    commons.RateLimitCounter{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.fields)
			rl := &CacheRepository{
				cache: tt.fields.cache,
			}
			got, err := rl.get(tt.args.compositeKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() got = %v, want %v", got, tt.want)
			}

			tt.fields.cache.AssertExpectations(t)
		})
	}
}

func TestRateLimitCacheRepository_set(t *testing.T) {
	type fields struct {
		cache *commons.MDCacheMock
	}
	type args struct {
		rlc commons.RateLimitCounter
	}
	tests := []struct {
		name    string
		fields  fields
		mock    func(fs fields)
		args    args
		wantErr bool
	}{
		{
			name: "should_set_success",
			args: args{
				rlc: commons.RateLimitCounter{
					Key:            "kennitromero@gmail.com#email#news",
					AttemptCounter: 1,
					PerSeconds:     60,
					CreatedAt:      "2023-05-18T12:43:50Z",
					UpdatedAt:      "2023-05-18T12:55:50Z",
				},
			},
			fields: fields{
				cache: &commons.MDCacheMock{},
			},
			mock: func(fs fields) {
				fs.cache.On(
					"Set",
					"kennitromero@gmail.com#email#news",
					"{\"key\":\"kennitromero@gmail.com#email"+
						"#news\",\"attempt_counter\":1,\"per_seconds\":60,\"cre"+
						"ated_at\":\"2023-05-18T12:43:50Z\",\"updated_at\":\"2023-05-18T12:55:50Z\"}",
					60,
				).Return(nil).Once()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.fields)
			rl := &CacheRepository{
				cache: tt.fields.cache,
			}
			if err := rl.set(tt.args.rlc); (err != nil) != tt.wantErr {
				t.Errorf("set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
