package cachedb

import (
	"testing"

	"github.com/jackc/pgx"
)

func TestConnectPostgresql(t *testing.T) {
	// see if we can connect to the database and query it
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		t.Errorf("Couldn't read environment config for PGDB: %v", err)
	}
	_, err = ConnectPostgresqlConfig(config)
	if err != nil {
		t.Errorf("Oh noes! Couldn't connect to PGDB: %v", err)
	}
}

func TestConnectRedis(t *testing.T) {

}
