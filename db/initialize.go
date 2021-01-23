package db

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v7"
	"xorm.io/xorm"

	// postgres driver to be used by xorm
	_ "github.com/lib/pq"

	"github.com/boof/umg/settings"
)

var (
	Engine *xorm.Engine

	redisClient *redis.Client
	onlineUsers *redis.Client
)

func init() {
	createEngine()
	createRedisClient()
	createOnlineUsers()
}

// createEngine creates xorm postgres Engine
func createEngine() {
	eng, err := xorm.NewEngine(settings.DriverName, GetDataSourceName())
	if err != nil {
		log.Fatalf("error while creating new database Engine: %v", err)
	}

	err = eng.Ping()
	if err != nil {
		log.Fatalf("error while connecting database: %v", err)
	} else {
		log.Println("postgres Engine created successfully")
	}

	Engine = eng
}

func createRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Error while creating redis client: %v \n", err)
	} else {
		log.Println("Redis client created successfully!")
	}
}

func createOnlineUsers() {
	onlineUsers = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	_, err := onlineUsers.Ping().Result()
	if err != nil {
		log.Fatalf("Error while creating redis client for online users: %v \n", err)
	} else {
		log.Println("Redis client for online users created successfully!")
	}
}

// GetDataSourceName returns data source name
func GetDataSourceName() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv(settings.DBHost), os.Getenv(settings.DBPort),
		os.Getenv(settings.PostgresUser), os.Getenv(settings.PostgresPass), os.Getenv(settings.PostgresDB))
}

func Sync(table interface{}) {
	err := Engine.Sync(table)
	if err != nil {
		log.Fatalf("error while syncing tables: %v", err)
	}
}
