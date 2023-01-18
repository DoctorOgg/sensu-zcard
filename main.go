package main

// https://redis.uptrace.dev
// github.com/spf13/viper@v1.7.0
import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	RedisHost     string
	RedisKey      string
	RedisPort     int
	RedisDB       int
	RedisPassword string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-zcard",
			Short:    "Sensu Check to return zcard metrics",
			Keyspace: "sensu.io/plugins/sensu-zcard/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "redis-host",
			Env:       "REDIS_HOST",
			Argument:  "host",
			Shorthand: "i",
			Default:   "localhost",
			Usage:     "The Redis host to connect to",
			Value:     &plugin.RedisHost,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "redis-port",
			Env:       "REDIS_PORT",
			Argument:  "port",
			Shorthand: "p",
			Default:   6379,
			Usage:     "The Redis port to connect to",
			Value:     &plugin.RedisPort,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "redis-key",
			Env:       "REDIS_KEY",
			Argument:  "key",
			Shorthand: "k",
			Default:   "fubar:farts*",
			Usage:     "The Redis key to report on",
			Value:     &plugin.RedisKey,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "redis-password",
			Env:       "REDIS_PASSWORD",
			Argument:  "password",
			Shorthand: "w",
			Default:   "",
			Usage:     "The Redis password",
			Value:     &plugin.RedisPassword,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "redis-db",
			Env:       "REDIS_DB",
			Argument:  "database",
			Shorthand: "d",
			Default:   0,
			Usage:     "The Redis db to connect to",
			Value:     &plugin.RedisDB,
		},
	}
)

func main() {
	useStdin := false
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error check stdin: %v\n", err)
		panic(err)
	}
	//Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		log.Println("using stdin")
		useStdin = true
	}

	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func getRedisKeys(rdb *redis.Client, ctx context.Context, prefix string) []string {
	var keys []string
	iter := rdb.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	return keys
}

func getzcard(rdb *redis.Client, ctx context.Context, key string) int64 {
	zcard, err := rdb.ZCard(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return zcard
}

func checkArgs(event *corev2.Event) (int, error) {
	// if len(plugin.Example) == 0 {
	// 	return sensu.CheckStateWarning, fmt.Errorf("--example or CHECK_EXAMPLE environment variable is required")
	// }
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	now := time.Now()
	ctx := context.Background()

	redis_connection_string := fmt.Sprintf("%s:%d", plugin.RedisHost, plugin.RedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_connection_string,
		Password: plugin.RedisPassword, // no password set
		DB:       plugin.RedisDB,       // use default DB
	})

	var accumulator int64 = 0

	keys := getRedisKeys(rdb, ctx, plugin.RedisKey)
	//log.Println("keys", keys)
	for _, key := range keys {
		zcard := getzcard(rdb, ctx, key)
		//fmt.Printf("%s %d\n", key, zcard)
		metrics_ouput := fmt.Sprintf("%s %d %d", key, zcard, now.Unix()) // type Graphite format
		accumulator += zcard
		fmt.Println(metrics_ouput)
		time.Sleep(time.Second)
	}
	metrics_ouput := fmt.Sprintf("total" %d %d", accumulator, now.Unix()) // type Graphite format

	// log.Println("executing check with --example", plugin.Example)
	return sensu.CheckStateOK, nil
}
