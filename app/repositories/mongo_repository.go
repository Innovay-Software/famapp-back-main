package repositories

import (
	"context"

	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	mongoDB        *mongo.Database
	mongoUsersColl *mongo.Collection
}

func (repo *mongoRepo) CreateUserInMongo(user *models.User) {
	filter := bson.M{"mobile": (*user).Mobile}
	var mongoUser models.UserMongo
	err := repo.mongoUsersColl.FindOne(context.TODO(), filter).Decode(&mongoUser)
	if err == nil {
		// Found a user with matching mobile
		if mongoUser.UUID == user.UUID.String() {
			// Both mobile and UUID match, simple update is fine
			repo.mongoUsersColl.UpdateOne(context.TODO(), filter, bson.M{
				"name": user.Name,
				"family_id": user.FamilyID,
				"email": user.Email,
				"avatar": user.Avatar,
			})
			return
		}

		// UUID Mismatch, delete the all record with mobile or uuid
		repo.mongoUsersColl.DeleteOne(context.TODO(), filter)
	}
	// Delete target UUID user
	repo.mongoUsersColl.DeleteOne(context.TODO(), bson.M{"uuid": user.UUID})

	// At this point, no user with target mobile and UUID exists, create a new one
	mongoUser = models.UserMongo{
		UUID:          user.UUID.String(),
		FamilyID:      1,
		Name:          user.Name,
		Email:         user.Email,
		Mobile:        user.Mobile,
		Avatar:        user.Avatar,
		SchemaVersion: "1.0.0",
	}
	result, err := repo.mongoUsersColl.InsertOne(context.TODO(), mongoUser)
	if err != nil {
		utils.LogError("Insert mongo user error: ", err.Error(), result)
	}
}

func (repo *mongoRepo) DeleteUserInMongo(user *models.User) {
	filter := bson.M{"uuid": (*user).UUID}
	result, err := repo.mongoUsersColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		utils.LogError("Delete mongo user error: ", err.Error(), result)
	}
}
