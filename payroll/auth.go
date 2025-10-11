package payroll

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

var userCollection = "users"

func getUserCollection() *mongo.Collection {
	return MongoClient.Database(db).Collection(userCollection)
}

func AuthenticateUser(username, password string) (*User, error) {
	collection := getUserCollection()
	var user User
	err := collection.FindOne(context.Background(), bson.M{"username": username, "password": password}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
