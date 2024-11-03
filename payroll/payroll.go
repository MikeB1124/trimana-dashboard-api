package payroll

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MikeB1124/trimana-dashboard-api/configuration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient        *mongo.Client
	db                 = "trimanaDB"
	employeeCollection = "employees"
	timeCardCollection = "timeCards"
)

var TimeZone, _ = time.LoadLocation("America/Los_Angeles")

func init() {
	config := configuration.GetConfig()
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.du0vf.mongodb.net", config.MongoDB.Username, config.MongoDB.Password))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	MongoClient = client
}

func getEmployeeCollection() *mongo.Collection {
	return MongoClient.Database(db).Collection(employeeCollection)
}

func getTimeCardCollection() *mongo.Collection {
	return MongoClient.Database(db).Collection(timeCardCollection)
}
