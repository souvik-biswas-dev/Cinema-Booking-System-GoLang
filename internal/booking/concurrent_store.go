package booking

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ConcurentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurentStore() *ConcurentStore {
	return &ConcurentStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurentStore) Book(b Booking) (Booking, error) {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.bookings[b.SeatID]; exists {
		return Booking{}, ErrSeatAlreadyBooked
	}

	id := uuid.New().String()
	b.ID = id
	b.Status = "held"
	b.ExpiresAt = time.Now().Add(2 * time.Minute)
	s.bookings[b.SeatID] = b
	return b, nil
}

func (s *ConcurentStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock()

	var result []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieID {
			result = append(result, b)
		}
	}
	return result
}

func (s *ConcurentStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	s.Lock()
	defer s.Unlock()

	for _, b := range s.bookings {
		if b.ID == sessionID && b.UserID == userID {
			b.Status = "confirmed"
			s.bookings[b.SeatID] = b
			return b, nil
		}
	}
	return Booking{}, ErrSeatAlreadyBooked
}

func (s *ConcurentStore) Release(ctx context.Context, sessionID string, userID string) error {
	s.Lock()
	defer s.Unlock()

	for key, b := range s.bookings {
		if b.ID == sessionID && b.UserID == userID {
			delete(s.bookings, key)
			return nil
		}
	}
	return ErrSeatAlreadyBooked
}
