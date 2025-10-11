package payroll

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTimeCardByDate(employeeID string, dateNow time.Time) (*TimeCard, error) {
	collection := getTimeCardCollection()
	var timeCard TimeCard
	err := collection.FindOne(context.Background(), bson.M{"employeeID": employeeID, "createdAt": dateNow}).Decode(&timeCard)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &timeCard, nil
}

func GetAllTimeCards() ([]TimeCard, error) {
	collection := getTimeCardCollection()
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var timeCards []TimeCard
	if err = cursor.All(context.Background(), &timeCards); err != nil {
		return nil, err
	}

	return timeCards, nil
}

func GetTimeCardByID(id string) (*TimeCard, error) {
	collection := getTimeCardCollection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid timecard ID")
	}

	var timeCard TimeCard
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&timeCard)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("timecard not found")
	}
	if err != nil {
		return nil, err
	}
	return &timeCard, nil
}

func GetTimeCardsByEmployeeID(employeeID string) ([]TimeCard, error) {
	collection := getTimeCardCollection()
	cursor, err := collection.Find(context.Background(), bson.M{"employeeID": employeeID})
	if err != nil {
		return nil, err
	}

	var timeCards []TimeCard
	if err = cursor.All(context.Background(), &timeCards); err != nil {
		return nil, err
	}

	return timeCards, nil
}

func CreateNewTimeCard(employeeID string, dateNow time.Time) error {
	newTimeCard := TimeCard{
		ID:         primitive.NewObjectID(),
		EmployeeID: employeeID,
		CreatedAt:  dateNow,
		TimeBlocks: []TimeBlock{
			{CheckIn: time.Now().In(TimeZone)},
		},
	}
	collection := getTimeCardCollection()
	_, err := collection.InsertOne(context.Background(), newTimeCard)
	if err != nil {
		return err
	}
	return nil
}

func CreateTimeCard(timeCard TimeCard) (*TimeCard, error) {
	collection := getTimeCardCollection()

	// Generate new ObjectID if not provided
	if timeCard.ID.IsZero() {
		timeCard.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(context.Background(), timeCard)
	if err != nil {
		return nil, err
	}

	return &timeCard, nil
}

func UpdateTimeCard(id string, timeCard TimeCard) (*TimeCard, error) {
	collection := getTimeCardCollection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid timecard ID")
	}

	update := bson.M{
		"$set": bson.M{
			"employeeID": timeCard.EmployeeID,
			"createdAt":  timeCard.CreatedAt,
			"timeBlocks": timeCard.TimeBlocks,
		},
	}

	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("timecard not found")
	}

	return GetTimeCardByID(id)
}

func DeleteTimeCard(id string) error {
	collection := getTimeCardCollection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid timecard ID")
	}

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("timecard not found")
	}

	return nil
}

func CreateCheckInBlock(timeCardID primitive.ObjectID) error {
	collection := getTimeCardCollection()
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": timeCardID}, bson.M{"$push": bson.M{"timeBlocks": bson.M{"checkIn": time.Now().In(TimeZone)}}})
	if err != nil {
		return err
	}
	return nil
}

func StampCheckOut(timeCardID primitive.ObjectID, indexOfLastBlock int) error {
	collection := getTimeCardCollection()
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": timeCardID}, bson.M{"$set": bson.M{fmt.Sprintf("timeBlocks.%d.checkOut", indexOfLastBlock): time.Now().In(TimeZone)}})
	if err != nil {
		return err
	}
	return nil
}

func StampTimeCard(timeCard *TimeCard) (string, error) {
	indexOfLastBlock := len(timeCard.TimeBlocks) - 1
	lastBlock := timeCard.TimeBlocks[indexOfLastBlock]
	if lastBlock.CheckOut != (time.Time{}) {
		err := CreateCheckInBlock(timeCard.ID)
		if err != nil {
			return "", err
		}
		return "Checked In", nil
	} else {
		err := StampCheckOut(timeCard.ID, indexOfLastBlock)
		if err != nil {
			return "", err
		}
		return "Checked Out", nil
	}
}

func GetTimeCardsByPayPeriod(employeeID string, startPayPeriod time.Time, endPayPeriod time.Time) ([]TimeCard, error) {
	collection := getTimeCardCollection()
	cursor, err := collection.Find(context.Background(), bson.M{"employeeID": employeeID, "createdAt": bson.M{"$gte": startPayPeriod, "$lte": endPayPeriod}})
	if err != nil {
		return nil, err
	}

	var timeCards []TimeCard
	if err = cursor.All(context.Background(), &timeCards); err != nil {
		return nil, err
	}

	return timeCards, nil
}
