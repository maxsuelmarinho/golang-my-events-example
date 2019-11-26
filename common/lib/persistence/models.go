package persistence

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	First    string        `json:"first"`
	Last     string        `json:"last"`
	Age      int           `json:"age"`
	Bookings []Booking     `json:"bookings"`
}

func (u *User) String() string {
	return fmt.Sprintf("id: %s, first_name: %s, last_name: %s, Age: %d, Bookings: %v", u.ID, u.First, u.Last, u.Age, u.Bookings)
}

type Booking struct {
	Date    int64  `json:"date"`
	EventID []byte `json:"eventId"`
	Seats   int    `json:"seats"`
}

type Event struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `dynamodbav:"EventName" json:"name"`
	Duration  int           `json:"duration"`
	StartDate int64         `json:"startDate"`
	EndDate   int64         `json:"endDate"`
	Location  Location      `json:"location"`
}

type Location struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `json:"name"`
	Address   string        `json:"address"`
	Country   string        `json:"country"`
	OpenTime  int           `json:"openTime"`
	CloseTime int           `json:"closeTime"`
	Halls     []Hall        `json:"halls"`
}

type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}
