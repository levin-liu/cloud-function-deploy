package scan_key_v2

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-redis/redis/v8"
)

var _redisClient *redis.Client

func init() {
	functions.HTTP("ScanKeyV2", scanKey)
}

func initRedis() (redisClient *redis.Client, err error) {
	redisHost := os.Getenv("REDISHOST")
	if redisHost == "" {
		return nil, errors.New("REDISHOST must be set")
	}
	redisPort := os.Getenv("REDISPORT")
	if redisPort == "" {
		return nil, errors.New("REDISPORT must be set")
	}

	if _redisClient != nil {
		return _redisClient, nil
	}

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	options := &redis.Options{
		Addr: addr,
	}

	_redisClient = redis.NewClient(options)
	return
}

func checkRedisClientConnection(ctx context.Context, redisClient *redis.Client) (string, error) {
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return "", err
	}
	return pong, nil
}

func scanKey(w http.ResponseWriter, r *http.Request) {
	// Initialize connection pool on first invocation
	ctx := context.Background()
	redisClient, err := initRedis()
	if err != nil {
		log.Printf("connect to Redis failed: %v", err)
		return
	}
	if _, err := checkRedisClientConnection(ctx, redisClient); err != nil {
		log.Printf("connect to Redis failed: %v", err)
		return
	}

	results, err := redisClient.Keys(ctx, "*").Result()
	if err != nil {
		log.Printf("redis.Keys: %v", err)
		http.Error(w, "Error incrementing visit count", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Exist keys: %v", results)
}
