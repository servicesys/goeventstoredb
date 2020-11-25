package core

type EventStore interface {
	Save(event Event) error
	Load(aggregateID int64,  aggregateType string,domain string) ([]Event, error)
}
