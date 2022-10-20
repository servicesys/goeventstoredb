package main

import (
	"context"
	"fmt"
	"github.com/servicesys/goeventstoredb/core"
	"github.com/servicesys/goeventstoredb/store"

	//"github.com/gofrs/uuid"
	//"math/rand"

	//"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	//"math/rand"
	"os"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "valter"
	password = "valter"
	dbname   = "app_sistema"
)

func main() {

	//Colocar o PGPool
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := pgx.Connect(context.Background(), psqlInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		panic(err)
	}
	fmt.Println("Successfully connected!")
	fmt.Println(db.Config())

	poolDB := CreatePGXPool(user, password, host, "5432", dbname, "rascunho")
	var eventStore = store.NewEventStore(poolDB)

	jsonRow := `{
    "firstName" : "George" , 
    "lastName"  : "Lucas"
    }`

	schemaRaw := `{
    "$id": "https://qri.io/schema/",
    "$comment" : "sample comment",
    "title": "Person",
    "type": "object",
    "properties": {
        "firstName": {
            "type": "string"
        },
        "lastName": {
            "type": "string"
        }
    },
    "required": ["firstName", "lastName"]
  }`

	fmt.Println(jsonRow)

	eventType := core.EventType{
		ID:       "person",
		MetaData: []byte(schemaRaw),
	}
	fmt.Println(eventType)
	//errType := eventStore.SaveEventType(eventType)

	//fmt.Println(errType)

	/*
		for i := 0; i < 100; i++ {

			uuid, _ := uuid.DefaultGenerator.NewV1()
			event := core.Event{
				EventID:       uuid.String(),
				EventType:     eventType,
				EventVersion:  "1",
				AggregateID:   rand.Int63n(1000),
				AggregateType: "mytype",
				Payload:       []byte(jsonRow),
				UserID:        "merlin",
				AppName:       "cmd",
				TransactionID: core.GenerateTransactionID(),
				DomainTenant:  "castelo-brasil"}

			fmt.Println(event.TransactionID)

			valid, errStr := eventStore.Validate(event)
			if valid {

				errSave := eventStore.SaveEvent(event)
				if errSave != nil {
					fmt.Println(errSave)
					panic(errSave)
				}

				//events, errLoad := eventStore.Load(410, "mytype" , "castelo-brasil")
				//fmt.Println(errLoad)
				//fmt.Println(string(events[0].Payload))

			} else {

				fmt.Println(errStr)
			}

		}
	*/

	events, erro := eventStore.Load(410, "mytype", "castelo-brasil")

	fmt.Println(erro)

	if erro == nil {
		for _, evt := range events {
			fmt.Println(evt.TransactionID)
			fmt.Println(string(evt.EventType.MetaData[:]))
		}
	}

	evtTrnas, err := eventStore.GetByTransactionID("477905475755158499323137735008343377823")

	fmt.Println(err)
	fmt.Println(evtTrnas.EventID)
	fmt.Println(evtTrnas.TransactionID)

	eventsPeriod, erro := eventStore.GetByPeriod("2022-10-17", "2022-10-17")

	fmt.Println(erro)

	if erro == nil {
		for _, evt := range eventsPeriod {
			fmt.Println(evt.TransactionID)
			//fmt.Println(string(evt.EventType.MetaData[:]))
		}
	}

}

func CreatePGXPool(user string, pass string, host string, port string, database string, app string) *pgxpool.Pool {

	connString := "postgres://" + user + ":" + pass + "@" + host + "/" + database + "?sslmode=disable" + "&" + "application_name=" + app
	configPool, errConfPool := pgxpool.ParseConfig(connString)
	if errConfPool != nil {
		panic(errConfPool)
	}
	configPool.MinConns = 1
	configPool.MaxConns = 4

	poolConn, errPool := pgxpool.ConnectConfig(context.Background(), configPool)
	if errPool != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", errPool)
		os.Exit(1)
	}
	return poolConn
}
