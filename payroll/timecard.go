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
