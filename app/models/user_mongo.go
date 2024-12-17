package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserMongo struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UUID          string             `bson:"uuid"`
	FamilyID      int64              `bson:"family_id"`
	Name          string             `bson:"name"`
	Email         string             `bson:"email"`
	Mobile        string             `bson:"mobile"`
	Avatar        string             `bson:"avatar"`
	SchemaVersion string             `bson:"schema_version,omitempty"`
}