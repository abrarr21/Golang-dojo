package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Image struct {
	URL    string `bson:"url" json:"url"`
	Name   string `bson:"name" json:"name"`
	FileID string `bson:"file_id" json:"file_id"`
}

type User struct {
	ID    bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email string        `bson:"email" json:"email"`
	Image Image         `bson:"image" json:"image"`
}
