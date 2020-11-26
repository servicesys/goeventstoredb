package main

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/servicesys/goeventstoredb/core"
	"github.com/servicesys/goeventstoredb/store"
	"math/rand"
	"os"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "agenda"
)

func main() {

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

	var eventStore = store.NewEventStore(db)

	uuid, _ := uuid.DefaultGenerator.NewV1()
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

	event := core.Event{
		EventID:       uuid.String(),
		EventType:     "eventFired",
		EventVersion:  "1",
		AggregateID:   rand.Int63n(1000),
		AggregateType: "mytype",
		Payload:       []byte(jsonRow),
		MetaData:      []byte(schemaRaw),
		UserID:        "merlin" ,
	    DomainTenant: "castelo-brasil"}


	valid , errStr := eventStore.Validate(event)
	if valid {

		errSave := eventStore.Save(event)
		if err != nil {
			fmt.Println(errSave)
		}


		events, errLoad := eventStore.Load(410, "mytype" , "castelo-brasil")
		fmt.Println(errLoad)
		fmt.Println(string(events[0].Payload))

	}else {

		fmt.Println(errStr)
	}



}
