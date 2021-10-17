package dbrepo

import (
	"context"
	"errors"
	"github.com/ekateryna-tln/booking/internal/models"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation insert a reservation into the database
func (m *postgresDBRepo) InsertReservation(r models.Reservation) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	stmt := `insert into reservations
			(id ,first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = m.DB.ExecContext(ctx, stmt,
		uuid,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// InsertRoomRestriction insert a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	stmt := `insert into room_restrictions
			(id, start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = m.DB.ExecContext(ctx, stmt,
		uuid,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID and false if no availability exists
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	query := `select count(id) from room_restrictions
				where room_id = $1 and $2 < end_date and $3 > start_date`
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `select r.id, r.room_name from rooms r
				where id not in (select room_id from room_restrictions
				where $1 < end_date and $2 > start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err = rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *postgresDBRepo) GetRoomByID(roomID string) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms
				where id = $1`
	row := m.DB.QueryRowContext(ctx, query, roomID)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreateAt,
		&room.UpdateAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}

// GetUserByID returns a user by id
func (m *postgresDBRepo) GetUserByID(uuid uuid.UUID) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
				from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, uuid)
	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreateAt,
		&u.UpdateAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUserByID update a user in the database
func (m *postgresDBRepo) UpdateUserByID(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set select 
				first_name = $1,
				last_name = $2,
				email = $3,
				access_level = $4, 
				updated_at = $5
			where id = $1`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (uuid.UUID, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var uuid uuid.UUID
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(
		&uuid,
		&hashedPassword,
	)
	if err != nil {
		return uuid, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return uuid, "", errors.New("incorrect password")
	} else {
		return uuid, "", err
	}
	return uuid, hashedPassword, err
}
