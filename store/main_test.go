package store

import (
	"context"
	"fmt"
	"github.com/jackc/tern/migrate"
	"log"
	"os"
	"testing"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/servicesys/goeventstoredb/core"
)

const (
	host     = "localhost"
	port     = 5438
	user     = "valter"
	password = "valter"
	dbname   = "app_sistema"
)

var eventStore *EventStorePostgresql

func TestMain(t *testing.T) {

	log.Println("START")
	pool, err := dockertest.NewPool("")

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12",
		Env: []string{
			"POSTGRES_PASSWORD=valter",
			"POSTGRES_USER=valter",
			"POSTGRES_DB=app_sistema",
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: "5433"},
			},
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Tell docker to hard kill the container in 120 seconds
	resource.Expire(120)

	//hostAndPort := resource.GetHostPort("5432/tcp")
	//databaseUrl := fmt.Sprintf("postgres://valter:valter@%s/dbname?sslmode=disable", hostAndPort)

	poolDB := createPGXPoolDB(user, password, host, "5433", dbname, "rascunho")

	//migratex(ctx context.Context, poolConn *pgxpool.Conn,  path string) (err error)
	ctx := context.Background()
	con, _ := poolDB.Acquire(ctx)
	errorMigrate := migratex(ctx, con, "../migrations/")
	if errorMigrate != nil {
		panic(errorMigrate)
	}

	eventStore = NewEventStore(poolDB)

	//fmt.Println(databaseUrl)
	//log.Println(databaseUrl)
	pool.MaxWait = 120 * time.Second
	// if err = pool.Retry(func() error {
	//	db, err = sql.Open("postgres", databaseUrl)
	//	if err != nil {
	// ...
}

func TestSomething(t *testing.T) {

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

	eventType := core.EventType{
		ID:       "person5",
		MetaData: []byte(schemaRaw),
	}

	erro := eventStore.SaveEventType(eventType)

	if erro != nil {
		t.Error("Not insert event type", erro.Error())
	}
}

func createPGXPoolDB(user string, pass string, host string, port string, database string, app string) *pgxpool.Pool {

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

func migratex(ctx context.Context, poolConn *pgxpool.Conn, path string) (err error) {

	migrator, err := migrate.NewMigrator(ctx, poolConn.Conn(), "eventstore_schema_version")
	if err != nil {
		return fmt.Errorf("cannot run migration: %w", err)
	}

	migrator.OnStart = func(sequence int32, name, direction, sql string) {
		log.Printf("executing %s %s\n", name, direction)
	}

	// Test the migration scripts and prepare database for integration tests.
	if err := migrator.LoadMigrations(path); err != nil {
		return fmt.Errorf("cannot load migrations: %w", err)
	}

	// Undo database migrations.
	if err := migrator.MigrateTo(ctx, 0); err != nil {
		return fmt.Errorf("cannot undo database migrations: %v", err)
	}

	// Migrate to latest version of the database
	if err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("cannot apply migrations: %v", err)
	}
	return nil
}