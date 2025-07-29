package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	dsn := "postgres://postgres:postgres@localhost:5432/main?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	ln := pgdriver.NewListener(db)
	if err := ln.Listen(context.Background(), "events:created"); err != nil {
		panic(err)
	}

	for notif := range ln.Channel() {
		fmt.Println(notif.Channel, notif.Payload)
	}
}

// docker run --name pg_test -e "POSTGRES_USER=postgres" -e "POSTGRES_PASSWORD=postgres" -e "POSTGRES_DB=main" -d -p 5432:5432 postgres

// CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

// CREATE TABLE IF NOT EXISTS events(
//   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
//   body TEXT NULL,
//   status VARCHAR (50) NOT NULL DEFAULT 'PENDING'
// );

// INSERT INTO events("body") VALUES ('example body 4');

// CREATE FUNCTION events_created_trigger()
// RETURNS TRIGGER AS $$
// BEGIN
//   PERFORM pg_notify('events:created', NEW.id::text); -- NOTIFY events:created, NEW.id::text;
//   RETURN NULL;
// END;
// $$
// LANGUAGE plpgsql;

// CREATE TRIGGER events_created_trigger
// AFTER INSERT ON events
// FOR EACH ROW EXECUTE PROCEDURE events_created_trigger();

// DROP TRIGGER IF EXISTS events_created_trigger ON events;
// DROP FUNCTION IF EXISTS events_created_trigger;

// NOTIFY "events:created", 'optional payload';
