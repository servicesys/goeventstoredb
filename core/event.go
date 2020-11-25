package core

import "time"

type Event struct {
	EventID       string
	EventType     string
	DomainTenant  string
	EventVersion  string
	TimeStamp     string
	AggregateID   int64
	AggregateType string
	Payload       []byte
	MetaData      []byte
	CreatedAt     time.Time
	UserID        string
}
