package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/storage/domain"
)

var ErrEventNotFound = errors.New("event not found")

type Storage struct {
	mu   sync.RWMutex
	data map[string]domain.Event
}

func New() *Storage {
	return &Storage{
		data: make(map[string]domain.Event),
	}
}

func NewWithEvents(events map[string]domain.Event) *Storage {
	return &Storage{data: events}
}

func (s *Storage) Connect(_ context.Context) error {
	s.data = make(map[string]domain.Event)
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = nil
	return nil
}

func (s *Storage) Create(event domain.Event) (string, error) {
	event.ID = uuid.New().String()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[event.ID] = event

	return event.ID, nil
}

func (s *Storage) Update(id string, event domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[id]; !exists {
		return ErrEventNotFound
	}

	event.ID = id
	s.data[id] = event

	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[id]; !exists {
		return ErrEventNotFound
	}

	delete(s.data, id)

	return nil
}

func (s *Storage) ListOfEventsForDay(date time.Time) ([]domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []domain.Event
	for _, event := range s.data {
		if event.DateTime.Year() == date.Year() && event.DateTime.YearDay() == date.YearDay() {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) ListOfEventsForWeek(date time.Time) ([]domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startOfWeek := date.Truncate(24*time.Hour).AddDate(0, 0, -int(date.Weekday())-1)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	var events []domain.Event
	for _, event := range s.data {
		if event.DateTime.After(startOfWeek.Add(-time.Second)) && event.DateTime.Before(endOfWeek.Add(time.Second)) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) ListOfEventsForMonth(date time.Time) ([]domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []domain.Event
	for _, event := range s.data {
		if event.DateTime.Month() == date.Month() && event.DateTime.Year() == date.Year() {
			events = append(events, event)
		}
	}

	return events, nil
}
