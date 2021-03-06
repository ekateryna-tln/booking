package dbrepo

import (
	"errors"
	"github.com/ekateryna-tln/booking/internal/models"
	"github.com/gofrs/uuid"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation insert a reservation into the database
func (m *testDBRepo) InsertReservation(r models.Reservation) (string, error) {
	if r.RoomID == "" {
		return "", errors.New("the room not found")
	}
	uuid, _ := uuid.NewV4()
	return uuid.String(), nil
}

// InsertRoomRestriction insert a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == "fail-room-restriction-uuid" {
		return errors.New("fail insert room restriction")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID and false if no availability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID string) (bool, error) {
	if roomID == "room_not_available" {
		return false, nil
	} else if roomID == "fail-search-availability" {
		return false, errors.New("")
	}
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomByID(roomID string) (models.Room, error) {
	var room models.Room
	if roomID == "" {
		return room, errors.New("the room not found")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(uuid uuid.UUID) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUserByID(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (string, string, error) {
	return "", "", nil
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

// GetReservationByID returns one reservation by id
func (m *testDBRepo) GetReservationByID(id string) (models.Reservation, error) {
	var res models.Reservation
	return res, nil
}

// UpdateReservationByID updates a reservation in the database
func (m *testDBRepo) UpdateReservationByID(r models.Reservation) error {
	return nil
}

// DeleteReservationByID deletes a reservation by id
func (m *testDBRepo) DeleteReservationByID(id string) error {
	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id string, processed int) error {
	return nil
}

// AllRooms
func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}
