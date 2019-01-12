package cachedb

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx"
)

type dataManager interface {
	New()
	Get(name, lovetoken string)
}

// DataManager contains the redis pool, the pgdb connection
type DataManager struct {
	rddb *redis.Pool
	pgdb *pgx.Conn
}
