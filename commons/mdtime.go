package commons

import "time"

type ImdTime interface {
	GetNowUTC() time.Time
}

type MDTime struct {
}

func NewMTime() *MDTime {
	return &MDTime{}
}

func (mTime *MDTime) GetNowUTC() time.Time {
	return time.Now().UTC()
}
