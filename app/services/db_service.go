package services

import (
	"fmt"
	"os"
	"time"

	"context"

	"github.com/innovay-software/famapp-main/app/utils"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mainDB *gorm.DB
var readDB *gorm.DB
// var mongoUsersColl *mongo.Collection

// Get MainDB Connection (Write)
func GetMainDBConnection() *gorm.DB {
	if mainDB != nil {
		return mainDB
	}
	mainDB = _getDBConnection("DB_MAIN_")
	return mainDB
}

// Get ReadDB Connection (Read)
func GetReadDBConnection() *gorm.DB {
	if readDB != nil {
		return readDB
	}
	readDB = _getDBConnection("DB_READ_")
	return readDB
}

// Get DB Connection Helper
func _getDBConnection(envPrefix string) *gorm.DB {
	utils.Log("Create DB connection with prefix: ", envPrefix)

	// DB_DSN :=="${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv(envPrefix+"USERNAME"),
		os.Getenv(envPrefix+"PASSWORD"),
		os.Getenv(envPrefix+"HOST"),
		os.Getenv(envPrefix+"PORT"),
		os.Getenv(envPrefix+"DATABASE"),
	)
	newDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := newDb.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return newDb
}

// Get MongoDB Connection
func GetMongoDBConnection() *mongo.Database {
	mongoUser := os.Getenv("MONGO_USERNAME")
	mongoPass := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	mongoConString := "mongodb://" + mongoUser + ":" + mongoPass + "@" + mongoHost + ":" + mongoPort

	clientOptions := options.Client().ApplyURI(mongoConString)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	database := client.Database(mongoDatabase)
	// mongoUsersColl = database.Collection("users")
	return database
}

// Get RedisDB Connection
func GetRedisConnection() (
	*redis.Client, context.Context,
) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DSN"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	return rdb, ctx
}
