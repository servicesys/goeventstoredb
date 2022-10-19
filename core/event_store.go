package core

type EventStore interface {
	SaveEvent(event Event) error
	SaveEventType(event EventType) error
	Load(aggregateID int64, aggregateType string, domain string) ([]Event, error)
	GetByTransactionID(transactionID string) (Event, error)
	GetByPeriod(date_start, date_finish string) ([]Event, error)
	Validate(event Event) (bool, []string)
}
