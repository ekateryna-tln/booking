package repository

import (
	"github.com/ekateryna-tln/booking/internal/models"
	"github.com/gofrs/uuid"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (string, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID string) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(roomID string) (models.Room, error)

	GetUserByID(uuid uuid.UUID) (models.User, error)
	UpdateUserByID(u models.User) error
	Authenticate(email, testPassword string) (string, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id string) (models.Reservation, error)
	UpdateReservationByID(r models.Reservation) error
	DeleteReservationByID(id string) error
	UpdateProcessedForReservation(id string, processed int) error
	AllRooms() ([]models.Room, error)
}
