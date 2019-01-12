package cachedb

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

var (
	log       = logrus.New()
	redisPool *redis.Pool
)

func init() {
	// Log as JSON instead of the default ASCII formatter
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Output to a file instead of stderr
	file, err := os.OpenFile("logs/db.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr", err)
	}
}

// ConnectPostgresql creates a postgresql db connection, and sends back a pointer to the object along with any errors
func ConnectPostgresql() (*pgx.Conn, error) {
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		return nil, err
	}
	log.Info("Creating Postgresql connection")
	pgdb, err := pgx.Connect(config)
	if err != nil {
		fmt.Println("Postgresql connection failed: ", err)
		log.Panic("Postgresql connection failed: ", err)
	}
	log.Info("Postgresql successfully connected")
	return pgdb, nil
}

// ConnectPostgresqlConfig creates a postgresql db connection with a manual config,
// and sends back a pointer to the object along with any errors
func ConnectPostgresqlConfig(config pgx.ConnConfig) (*pgx.Conn, error) {
	log.Info("Creating Postgresql connection")
	pgdb, err := pgx.Connect(config)
	if err != nil {
		fmt.Println("Postgresql connection failed: ", err)
		log.Panic("Postgresql connection failed: ", err)
	}
	log.Info("Postgresql successfully connected")
	return pgdb, nil
}

// ConnectRedis creates a redis db connection, and sends back a pointer to the object
func ConnectRedis(database int) *redis.Pool {
	log.Info("Creating Redis connection")
	// then the Redis Database, rddb, for caching and sessions.
	rddb := newPool()
	log.Info("Redis successfully connected")
	return rddb
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10000,
		MaxActive:   100000,
		IdleTimeout: time.Duration(15) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				fmt.Println("Redis failed to connect: ", err.Error())
				log.Error("Redis failed to connect: ", err.Error())
			}
			if _, err := c.Do("AUTH", os.Getenv("RDPASSWORD")); err != nil {
				fmt.Println("Redis failed to authenticate: ", err.Error())
				log.Error("Redis failed to authenticate: ", err.Error())
				c.Close()
				return nil, err
			}
			return c, err
		},
	}
}
