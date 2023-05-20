package commons

const (
	TypeNotificationStatus    = "status"
	TypeNotificationNews      = "news"
	TypeNotificationMarketing = "marketing"
	TypeNotificationDaily     = "daily"
)

type MetaData struct {
	LangCode string `json:"lang_code"`
	Template string `json:"template"`
}

type Notification struct {
	To               string   `json:"to"`
	From             string   `json:"from"`
	Subject          string   `json:"subject"`
	Body             string   `json:"body"`
	WayToNotify      string   `json:"way_to_notify"`
	TypeNotification string   `json:"type_notification"`
	Meta             MetaData `json:"meta"`
}

type RateLimitCounter struct {
	Key            string `json:"key"`
	AttemptCounter int    `json:"attempt_counter"`
	PerSeconds     int    `json:"per_seconds"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
