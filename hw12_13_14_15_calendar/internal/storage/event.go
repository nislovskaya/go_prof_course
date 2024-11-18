package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/configs"
	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/storage/domain"
	memorystorage "github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	Memory   string = "memory"
	Database string = "database"
)

type EventStore interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error

	Create(event domain.Event) (string, error)
	Update(id string, event domain.Event) error
	Delete(id string) error
	ListOfEventsForDay(date time.Time) ([]domain.Event, error)
	ListOfEventsForWeek(date time.Time) ([]domain.Event, error)
	ListOfEventsForMonth(date time.Time) ([]domain.Event, error)
}

func GetStorage(cfg *configs.Config) (EventStore, error) {
	switch cfg.Storage.Type {
	case Memory:
		return memorystorage.New(), nil
	case Database:
		return sqlstorage.New(cfg.Database), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Storage.Type)
	}
}
