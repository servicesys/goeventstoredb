package store

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/servicesys/goeventstoredb/core"
	"github.com/servicesys/jsonschema/schema"
)

type EventStorePostgresql struct {
	DBConnection  *pgx.Conn
	JsonValidator schema.JsonSchemaValidator
}

func NewEventStore(db *pgx.Conn) *EventStorePostgresql {

	if db == nil {
		panic("goeventstoredb:DATABASE NIL")
	}

	return &EventStorePostgresql{DBConnection: db, JsonValidator: schema.JsonSchemaValidatorQri{}}
}

func (eventSore *EventStorePostgresql) Save(event core.Event) error {

	strSQL := ` INSERT INTO eventstore.event(event_id, event_type, 
             event_version,aggregate_id, payload,
             meta_data,user_id,aggregate_type,domain_tenant,created_at) 
             VALUES($1, $2 ,  $3 , $4 , $5 , $6 , $7 , $8 , $9 ,CURRENT_TIMESTAMP  at time zone 'utc' )`

	_, err := eventSore.DBConnection.Exec(context.Background(), strSQL,
		event.EventID,
		event.EventType,
		event.EventVersion,
		event.AggregateID,
		event.Payload,
		event.MetaData,
		event.UserID,
		event.AggregateType,
		event.DomainTenant)

	return err
}

func (eventSore *EventStorePostgresql) Load(aggregateID int64, aggregateType string, domain string) ([]core.Event, error) {

	strQuery :=
		`SELECT event_id, 
         event_type, 
		 event_version, 
		 aggregate_id, 
		 payload ,
		 meta_data ,
		 created_at,
		 user_id, 
		 aggregate_type,
         domain_tenant
		 FROM eventstore.event  
		WHERE aggregate_id = $1 AND aggregate_type = $2 AND  domain_tenant = $3`

	rows, err := eventSore.DBConnection.Query(context.Background(), strQuery, aggregateID, aggregateType, domain)

	var events []core.Event

	for rows.Next() {

		event := core.Event{}
		err := rows.Scan(
			&event.EventID,
			&event.EventType,
			&event.EventVersion,
			&event.AggregateID,
			&event.Payload,
			&event.MetaData,
			&event.CreatedAt,
			&event.UserID,
			&event.AggregateType,
			&event.DomainTenant)

		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	defer rows.Close()

	return events, err
}

func (eventSore *EventStorePostgresql) Validate(event core.Event) (bool, []string) {

	validJson, errorStr := eventSore.JsonValidator.ValidatorBytes(event.MetaData, event.Payload)
	if !validJson {
		return validJson, errorStr
	} else {
		return true, nil
	}

}
