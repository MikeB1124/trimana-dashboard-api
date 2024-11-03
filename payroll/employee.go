package payroll

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetEmployees() ([]Employee, error) {
	collection := getEmployeeCollection()
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var employees []Employee
	if err = cursor.All(context.Background(), &employees); err != nil {
		return nil, err
	}

	return employees, nil
}

func GetEmployeeByCardID(cardID string) (*Employee, error) {
	collection := getEmployeeCollection()
	var employee Employee
	err := collection.FindOne(context.Background(), bson.M{"cardID": cardID}).Decode(&employee)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("Employee not found")
	}
	if err != nil {
		return nil, err
	}
	return &employee, nil
}
