package models

import "time"

// User is the user model
type User struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreateAt    time.Time
	UpdateAt    time.Time
}

// Room is the room model
type Room struct {
	ID       string
	RoomName string
	CreateAt time.Time
	UpdateAt time.Time
}

// Restriction is the restriction model
type Restriction struct {
	ID              string
	RestrictionName string
	CreateAt        time.Time
	UpdateAt        time.Time
}

// Reservation is the reservation model
type Reservation struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    string
	CreateAt  time.Time
	UpdateAt  time.Time
	Room      Room
}

// RoomRestriction is the room restriction model
type RoomRestriction struct {
	ID            string
	StartDate     time.Time
	EndDate       time.Time
	RoomID        string
	ReservationID string
	RestrictionID string
	CreateAt      time.Time
	UpdateAt      time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
