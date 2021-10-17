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
	Authenticate(email, testPassword string) (uuid.UUID, string, error)
}
