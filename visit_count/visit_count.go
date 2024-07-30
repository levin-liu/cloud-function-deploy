package visit_count

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

func init() {
	// Register the HTTP handler with the Functions Framework
	functions.HTTP("VisitCount", visitCount)
}

// initializeRedis initializes and returns a connection pool
func initializeRedis() (*redis.Pool, error) {
	redisHost := os.Getenv("REDISHOST")
	if redisHost == "" {
		return nil, errors.New("REDISHOST must be set")
	}
	redisPort := os.Getenv("REDISPORT")
	if redisPort == "" {
		return nil, errors.New("REDISPORT must be set")
	}
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	const maxConnections = 10
	return &redis.Pool{
		MaxIdle: maxConnections,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, fmt.Errorf("redis.Dial: %w", err)
			}
			return c, err
		},
	}, nil
}

// visitCount increments the visit count on the Redis instance
// and prints the current count in the HTTP response.
func visitCount(w http.ResponseWriter, r *http.Request) {
	// Initialize connection pool on first invocation
	if redisPool == nil {
		// Pre-declare err to avoid shadowing redisPool
		var err error
		redisPool, err = initializeRedis()
		if err != nil {
			log.Printf("initializeRedis: %v", err)
			http.Error(w, "Error initializing connection pool", http.StatusInternalServerError)
			return
		}
	}

	conn := redisPool.Get()
	defer conn.Close()

	counter, err := redis.Int(conn.Do("INCR", "visits"))
	if err != nil {
		log.Printf("redis.Int: %v", err)
		http.Error(w, "Error incrementing visit count", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Visit count: %d", counter)
}
