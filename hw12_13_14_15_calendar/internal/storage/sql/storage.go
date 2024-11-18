package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/configs"
	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/storage/domain"
)

var ErrEventNotFound = errors.New("event not found")

type Storage struct {
	db  *sqlx.DB
	dsn string
	ctx context.Context
}

var ErrConnectFailed = errors.New("failed to connect database")

const (
	tableName          = "event"
	tableColumnsRead   = "id,title,description,datetime,duration,remind_time,user_id"
	tableColumnsInsert = "title,description,datetime,duration,remind_time,user_id"
)

func New(cfg configs.DatabaseConfig) *Storage {
	return &Storage{
		dsn: cfg.String(),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", s.dsn)
	if err != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+": %w", err)
	}

	s.ctx = ctx

	err = db.PingContext(s.ctx)
	if err != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+": %w", err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) Create(event domain.Event) (string, error) {
	var id string
	err := s.db.QueryRowContext(
		s.ctx,
		fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES($1,$2,$3,$4,$5,$6) RETURNING id",
			tableName, tableColumnsInsert,
		),
		event.Title,
		event.Description,
		event.DateTime.Format(time.RFC3339),
		event.Duration.Seconds(),
		event.UserID,
		event.NotifyIn.Format(time.RFC3339),
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Storage) Update(id string, event domain.Event) error {
	result, err := s.db.ExecContext(
		s.ctx,
		fmt.Sprintf(
			"UPDATE %s SET title=$1, description=$2, datetime=$3, duration=$4, user_id=$5, notify_in=$6 WHERE id=$7",
			tableName,
		),
		event.Title,
		event.Description,
		event.DateTime.Format(time.RFC3339),
		event.Duration.Seconds(),
		event.UserID,
		event.NotifyIn.Format(time.RFC3339),
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (s *Storage) Delete(id string) error {
	result, err := s.db.ExecContext(s.ctx,
		fmt.Sprintf("DELETE FROM %s WHERE id=$1", tableName),
		id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (s *Storage) ListOfEventsForDay(date time.Time) ([]domain.Event, error) {
	var events []domain.Event

	err := s.db.SelectContext(s.ctx,
		&events,
		fmt.Sprintf(
			"SELECT %s FROM %s WHERE DATE(datetime) = $1",
			tableColumnsRead, tableName,
		),
		date.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) ListOfEventsForWeek(date time.Time) ([]domain.Event, error) {
	var events []domain.Event

	startOfWeek := date.AddDate(0, 0, -int(date.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	err := s.db.SelectContext(s.ctx,
		&events,
		fmt.Sprintf(
			"SELECT %s FROM %s WHERE datetime BETWEEN $1 AND $2",
			tableColumnsRead, tableName,
		),
		startOfWeek.Format(time.RFC3339),
		endOfWeek.Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) ListOfEventsForMonth(date time.Time) ([]domain.Event, error) {
	var events []domain.Event

	err := s.db.SelectContext(s.ctx,
		&events,
		fmt.Sprintf(
			"SELECT %s FROM %s WHERE EXTRACT(MONTH FROM datetime)= $1 AND EXTRACT(YEAR FROM datetime)= $2",
			tableColumnsRead, tableName,
		),
		date.Month(),
		date.Year(),
	)
	if err != nil {
		return nil, err
	}

	return events, nil
}
