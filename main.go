package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/sirupsen/logrus"
)

var plog = logrus.New()

//  Keep a global config, for now
var config Config
var rdb *redis.Client
var limiter *redis_rate.Limiter // TODO: Is this safe for concurrent use?

func main() {
	plog.Info("Starting rate-limiter-svc")
	// Load our config
	var err error
	config, err = LoadConfig()
	if err != nil {
		plog.Fatal(err)
	}
	plog.Info("Loaded config")

	// Create the redis Client
	rdb = redis.NewClient(&redis.Options{
		Addr: config.RedisAddress,
	})
	plog.Info("Created Redis Client")

	// Create the limiter
	limiter = redis_rate.NewLimiter(rdb)
	plog.Info("Created Limiter")

	// Start the listener
	handleRequests()
}
