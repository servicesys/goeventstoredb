package store

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/servicesys/goeventstoredb/core"
	"github.com/servicesys/jsonschema/schema"
)

type EventStorePostgresql struct {
	PoolDB        *pgxpool.Pool
	JsonValidator schema.JsonSchemaValidator
}

func NewEventStore(poolDB *pgxpool.Pool) *EventStorePostgresql {

	if poolDB == nil {
		err := errors.New("goeventstoredb:DATABASE NIL")
		log.Println(err)
		os.Exit(1)
	}

	return &EventStorePostgresql{PoolDB: poolDB, JsonValidator: schema.JsonSchemaValidatorQri{}}
}

func (eventSore *EventStorePostgresql) SaveEvent(event core.Event) error {

	strSQL := ` INSERT INTO eventstore.event(id, event_type, 
             event_version,aggregate_id, event_data, user_id,aggregate_type,domain_tenant, app_name , transaction_id, created_at) 
             VALUES($1, $2 ,  $3 , $4 , $5 , $6 , $7 , $8 , $9 , $10 ,  CURRENT_TIMESTAMP  at time zone 'utc' )`

	_, err := eventSore.PoolDB.Exec(context.Background(), strSQL,
		event.EventID,
		event.EventType.ID,
		event.EventVersion,
		event.AggregateID,
		event.Payload,
		event.UserID,
		event.AggregateType,
		event.DomainTenant,
		event.AppName,
		event.TransactionID)

	return err
}

func (eventSore *EventStorePostgresql) SaveEventType(eventType core.EventType) error {

	strSQL := `INSERT INTO eventstore.event_type (id, meta_data, created_at)  VALUES($1, $2 ,CURRENT_TIMESTAMP  at time zone 'utc' );`

	_, err := eventSore.PoolDB.Exec(context.Background(), strSQL,
		eventType.ID,
		eventType.MetaData)

	return err
}

func (eventSore *EventStorePostgresql) Load(aggregateID int64, aggregateType string, domain string) ([]core.Event, error) {

	strQuery :=
		`SELECT 
         e.id,
		 e.event_version, 
		 e.aggregate_id, 
		 e.event_data ,
		 e.created_at,
		 e.user_id, 
		 e.aggregate_type,
         e.domain_tenant,
		 e.app_name,
		 e.transaction_id,
		 e.event_type,
		 et.meta_data 
		 FROM eventstore.event e INNER JOIN  eventstore.event_type et ON (e.event_type=et.id)
		 WHERE aggregate_id = $1 AND aggregate_type = $2 AND  domain_tenant = $3
       	 ORDER BY e.created_at`

	rows, err := eventSore.PoolDB.Query(context.Background(), strQuery, aggregateID, aggregateType, domain)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return rowsEvents(rows)
}

func (eventSore *EventStorePostgresql) GetByTransactionID(transactionID string) (core.Event, error) {

	strQuery :=
		`SELECT 
         e.id,
		 e.event_version, 
		 e.aggregate_id, 
		 e.event_data ,
		 e.created_at,
		 e.user_id, 
		 e.aggregate_type,
         e.domain_tenant,
		 e.app_name,
		 e.transaction_id,
		 e.event_type,
		 et.meta_data 
		 FROM eventstore.event e INNER JOIN  eventstore.event_type et ON (e.event_type=et.id)
		 WHERE  e.transaction_id = $1 
         ORDER BY e.created_at`

	rows, err := eventSore.PoolDB.Query(context.Background(), strQuery, transactionID)
	defer rows.Close()
	if err != nil {
		return core.Event{}, err
	}

	event := core.Event{}

	if rows.Next() {

		err := rows.Scan(
			&event.EventID,
			&event.EventVersion,
			&event.AggregateID,
			&event.Payload,
			&event.CreatedAt,
			&event.UserID,
			&event.AggregateType,
			&event.DomainTenant,
			&event.AppName,
			&event.TransactionID,
			&event.EventType.ID,
			&event.EventType.MetaData)
		if err != nil {
			return core.Event{}, err
		}
	}

	return event, err
}

func (eventSore *EventStorePostgresql) GetByPeriod(start, finish string) ([]core.Event, error) {

	layoutISO := "2006-01-02"
	date_start, err_start := time.Parse(layoutISO, start)
	date_finish, err_finish := time.Parse(layoutISO, finish)

	if err_start != nil || err_finish != nil {
		return nil, err_start
	}

	strQuery :=
		`SELECT 
         e.id,
		 e.event_version, 
		 e.aggregate_id, 
		 e.event_data ,
		 e.created_at,
		 e.user_id, 
		 e.aggregate_type,
         e.domain_tenant,
		 e.app_name,
		 e.transaction_id,
		 e.event_type,
		 et.meta_data 
		 FROM eventstore.event e INNER JOIN  eventstore.event_type et ON (e.event_type=et.id)
         WHERE  DATE(e.created_at) >= $1 AND DATE( e.created_at) <=  $2
       	 ORDER BY e.created_at`

	rows, err := eventSore.PoolDB.Query(context.Background(), strQuery, date_start, date_finish)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return rowsEvents(rows)
}

func (eventSore *EventStorePostgresql) Validate(event core.Event) (bool, []string) {

	validJson, errorStr := eventSore.JsonValidator.ValidatorBytes(event.EventType.MetaData, event.Payload)
	if !validJson {
		return validJson, errorStr
	} else {
		return true, nil
	}

}

func rowsEvents(rows pgx.Rows) ([]core.Event, error) {

	events := make([]core.Event, 0)
	for rows.Next() {
		event := core.Event{}
		err := rows.Scan(
			&event.EventID,
			&event.EventVersion,
			&event.AggregateID,
			&event.Payload,
			&event.CreatedAt,
			&event.UserID,
			&event.AggregateType,
			&event.DomainTenant,
			&event.AppName,
			&event.TransactionID,
			&event.EventType.ID,
			&event.EventType.MetaData)

		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
