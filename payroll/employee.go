package payroll

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetEmployeeByID(id string) (*Employee, error) {
	collection := getEmployeeCollection()
	var employee Employee
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&employee)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("employee not found")
	}
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func GetEmployeeByCardID(cardID string) (*Employee, error) {
	collection := getEmployeeCollection()
	var employee Employee
	err := collection.FindOne(context.Background(), bson.M{"cardID": cardID}).Decode(&employee)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("employee not found")
	}
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func CreateEmployee(employee Employee) (*Employee, error) {
	collection := getEmployeeCollection()

	// Generate a new ObjectID and convert to string for _id
	if employee.ID == "" {
		employee.ID = primitive.NewObjectID().Hex()
	}

	_, err := collection.InsertOne(context.Background(), employee)
	if err != nil {
		return nil, err
	}

	return &employee, nil
}

func UpdateEmployee(id string, employee Employee) (*Employee, error) {
	collection := getEmployeeCollection()

	// Create update document with all fields
	update := bson.M{
		"$set": bson.M{
			"branchID":     employee.BranchID,
			"name":         employee.Name,
			"employeeID":   employee.EmployeeID,
			"hourlyRate":   employee.HourlyRate,
			"payFrequency": employee.PayFrequency,
			"earningsCode": employee.EarningsCode,
			"rateCode":     employee.RateCode,
			"cardID":       employee.CardID,
			"email":        employee.Email,
		},
	}

	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("employee not found")
	}

	// Get the updated employee
	return GetEmployeeByID(id)
}

func DeleteEmployee(id string) error {
	collection := getEmployeeCollection()

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}
