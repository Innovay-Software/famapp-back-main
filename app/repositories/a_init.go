package repositories

import (
	"context"

	"github.com/innovay-software/famapp-main/app/services"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// DB connection
var mainDBCon *gorm.DB
var readDBCon *gorm.DB

// Mongo Collections
var mongoDB *mongo.Database
var mongoUsersColl *mongo.Collection

// Redis connection
var redisCon *redis.Client
var redisCtx context.Context

// Repo instances
var ConfigRepoIns *configRepo
var UserRepoIns *userRepo
var LockerNoteRepoIns *lockerNoteRepo
var FolderRepoIns *folderRepo
var CloudUploadRepoIns *cloudUploadRepo
var JobRepoIns *jobRepo
var UtilsRepoIns *utilsRepo
var MongoRepoIns *mongoRepo

func RepoInit() {

	// Init database and redis
	dbInit()
	redisInit()
	mongoInit()

	ConfigRepoIns = &configRepo{
		mainDBCon: mainDBCon,
		readDBCon: readDBCon,
		rd:        &redisRepo{redisClient: redisCon, redisCtx: redisCtx},
	}
	UserRepoIns = &userRepo{
		mainDBCon: mainDBCon,
		readDBCon: readDBCon,
		rd:        &redisRepo{redisClient: redisCon, redisCtx: redisCtx},
	}
	LockerNoteRepoIns = &lockerNoteRepo{
		mainDBCon: mainDBCon,
		readDBCon: readDBCon,
		rd:        &redisRepo{redisClient: redisCon, redisCtx: redisCtx},
	}
	FolderRepoIns = &folderRepo{
		mainDBCon: mainDBCon,
		readDBCon: readDBCon,
		rd:        &redisRepo{redisClient: redisCon, redisCtx: redisCtx},
	}
	CloudUploadRepoIns = &cloudUploadRepo{
		mainDBCon: mainDBCon,
		readDBCon: readDBCon,
	}
	JobRepoIns = &jobRepo{
		mainDBCon: mainDBCon,
	}
	UtilsRepoIns = &utilsRepo{
		mainDBCon: mainDBCon,
		readDBCon: readDBCon,
	}
	MongoRepoIns = &mongoRepo{
		mongoDB:        mongoDB,
		mongoUsersColl: mongoUsersColl,
	}
}

func dbInit() {
	mainDBCon = services.GetMainDBConnection()
	readDBCon = services.GetReadDBConnection()
}

func redisInit() {
	redisCon, redisCtx = services.GetRedisConnection()
}

func mongoInit() {
	mongoDB = services.GetMongoDBConnection()
	mongoUsersColl = mongoDB.Collection("users")
}
