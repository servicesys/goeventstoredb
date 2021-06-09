package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/servicesys/goeventstoredb/core"
	"github.com/servicesys/jsonschema/schema"
)

type EventStorePostgresql struct {
	//DBConnection  *pgx.Conn
	PoolDB  *pgxpool.Pool
	JsonValidator schema.JsonSchemaValidator
}

func NewEventStore(poolDB *pgxpool.Pool) *EventStorePostgresql {

	if poolDB == nil {
		panic("goeventstoredb:DATABASE NIL")
	}

	return &EventStorePostgresql{PoolDB: poolDB, JsonValidator: schema.JsonSchemaValidatorQri{}}
}

func (eventSore *EventStorePostgresql) Save(event core.Event) error {

	strSQL := ` INSERT INTO eventstore.event(event_id, event_type, 
             event_version,aggregate_id, payload,
             meta_data,user_id,aggregate_type,domain_tenant, app_name , transaction_id , created_at) 
             VALUES($1, $2 ,  $3 , $4 , $5 , $6 , $7 , $8 , $9 , $10 , $11 , CURRENT_TIMESTAMP  at time zone 'utc' )`

	_, err := eventSore.PoolDB.Exec(context.Background(), strSQL,
		event.EventID,
		event.EventType,
		event.EventVersion,
		event.AggregateID,
		event.Payload,
		event.MetaData,
		event.UserID,
		event.AggregateType,
		event.DomainTenant,
		event.AppName,
		event.TransactionID)

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
         domain_tenant,
		 app_name,
		 transaction_id
		 FROM eventstore.event  
		WHERE aggregate_id = $1 AND aggregate_type = $2 AND  domain_tenant = $3`

	rows, err := eventSore.PoolDB.Query(context.Background(), strQuery, aggregateID, aggregateType, domain)

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
			&event.DomainTenant,
			&event.AppName,
			&event.TransactionID)

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
