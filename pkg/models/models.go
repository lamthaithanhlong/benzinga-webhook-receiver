package models

import "time"

type LogEntry struct {
	UserID    int      `json:"user_id"`
	Total     float64  `json:"total"`
	Title     string   `json:"title"`
	Meta      MetaData `json:"meta"`
	Completed bool     `json:"completed"`
}

type MetaData struct {
	Logins       []LoginRecord `json:"logins"`
	PhoneNumbers PhoneNumbers  `json:"phone_numbers"`
}

type LoginRecord struct {
	Time time.Time `json:"time"`
	IP   string    `json:"ip"`
}

type PhoneNumbers struct {
	Home   string `json:"home"`
	Mobile string `json:"mobile"`
}
