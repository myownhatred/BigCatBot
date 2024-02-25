package entities

import "time"

type TimeWithOut struct {
	ID     int       `json:"id"`
	Name   string    `json:"name"`
	Time   time.Time `json:"time"`
	ChatId int64     `json:"chat_id"`
}
