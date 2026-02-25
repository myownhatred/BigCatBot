package mlog

import "time"

// mlog represents a message from the Telegram chat.
type Mlog struct {
    Timestamp   time.Time `json:"timestamp"`
    MessageId   int       `json:"message_id"`
    Sender      string    `json:"sender"`
    Content     string    `json:"content"`
}

// msum holds a collection of messages for summary.
type Msum struct {
    Messages []Mlog `json:"messages"`
}
