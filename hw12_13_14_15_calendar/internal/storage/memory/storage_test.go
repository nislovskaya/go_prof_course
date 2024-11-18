package memorystorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/storage/domain"
	memorystorage "github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var event = domain.Event{
	Title:       "some event",
	DateTime:    time.Now(),
	Description: "this is some event",
	Duration:    60 * time.Minute,
	UserID:      "1",
}

func TestStorageModify(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		memStorage := memorystorage.New()
		err := memStorage.Connect(context.Background())
		require.NoError(t, err)

		id, err := memStorage.Create(event)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		storedEvents, err := memStorage.ListOfEventsForDay(event.DateTime)
		require.NoError(t, err)
		require.Len(t, storedEvents, 1)
		require.Equal(t, event.Title, storedEvents[0].Title)
	})

	t.Run("update", func(t *testing.T) {
		newTitle := "new title"
		memStorage := memorystorage.New()
		err := memStorage.Connect(context.Background())
		require.NoError(t, err)

		id, err := memStorage.Create(event)
		require.NoError(t, err)

		eventToUpdate := domain.Event{
			ID:          id,
			Title:       newTitle,
			DateTime:    event.DateTime,
			Description: event.Description,
			UserID:      event.UserID,
			Duration:    event.Duration,
		}

		err = memStorage.Update(id, eventToUpdate)
		require.NoError(t, err)

		storedEvents, err := memStorage.ListOfEventsForDay(event.DateTime)
		require.NoError(t, err)
		require.Len(t, storedEvents, 1)
		require.Equal(t, newTitle, storedEvents[0].Title)
	})

	t.Run("update unknown", func(t *testing.T) {
		memStorage := memorystorage.New()
		err := memStorage.Connect(context.Background())
		require.NoError(t, err)

		randomID := uuid.New().String()
		err = memStorage.Update(randomID, event)

		require.ErrorIs(t, err, memorystorage.ErrEventNotFound)
	})

	t.Run("delete", func(t *testing.T) {
		memStorage := memorystorage.New()
		err := memStorage.Connect(context.Background())
		require.NoError(t, err)

		id, err := memStorage.Create(event)
		require.NoError(t, err)

		err = memStorage.Delete(id)
		require.NoError(t, err)

		eventsAfterDelete, err := memStorage.ListOfEventsForDay(event.DateTime)
		require.NoError(t, err)
		require.Len(t, eventsAfterDelete, 0)
	})
}

func TestStorageRead(t *testing.T) {
	t.Run("read all for day", func(t *testing.T) {
		n := 3
		memStorage := memorystorage.New()
		err := memStorage.Connect(context.Background())
		require.NoError(t, err)

		for i := 0; i < n; i++ {
			_, err = memStorage.Create(event)
			require.NoError(t, err)
		}

		eventsForDay, err := memStorage.ListOfEventsForDay(event.DateTime)
		require.NoError(t, err)
		require.Len(t, eventsForDay, n)
	})

	t.Run("read all for week", func(t *testing.T) {
		initialDate := time.Date(2023, 10, 16, 13, 10, 0, 0, time.UTC)
		memStorage := memorystorage.NewWithEvents(map[string]domain.Event{
			"1": {ID: "1", Title: "1", DateTime: initialDate.Add(-time.Hour * 24 * 2), UserID: "1"},
			"2": {ID: "2", Title: "2", DateTime: initialDate.Add(-time.Hour * 24 * 1), UserID: "1"},
			"3": {ID: "3", Title: "3", DateTime: initialDate.Add(0), UserID: "1"},
			"4": {ID: "4", Title: "4", DateTime: initialDate.Add(time.Hour * 24 * 1), UserID: "1"},
			"5": {ID: "5", Title: "5", DateTime: initialDate.Add(time.Hour * 24 * 7), UserID: "1"},
		})

		eventsForWeek, err := memStorage.ListOfEventsForWeek(initialDate)

		require.NoError(t, err)
		require.Len(t, eventsForWeek, 4) // Ожидаем 4 события на эту неделю
	})

	t.Run("read all for month", func(t *testing.T) {
		initialDate := time.Date(2023, 10, 16, 13, 10, 0, 0, time.UTC)
		memStorage := memorystorage.NewWithEvents(map[string]domain.Event{
			"1": {ID: "1", Title: "1", DateTime: initialDate.Add(0), UserID: "1"},
			"2": {ID: "2", Title: "2", DateTime: initialDate.Add(time.Hour * 24 * 30), UserID: "1"},
			"3": {ID: "3", Title: "3", DateTime: initialDate.Add(time.Hour * 24 * 10), UserID: "1"},
			"4": {ID: "4", Title: "4", DateTime: initialDate.Add(time.Hour * 24 * 2), UserID: "1"},
			"5": {ID: "5", Title: "5", DateTime: initialDate.Add(-time.Hour * 24 * 30), UserID: "1"},
		})

		eventsForMonth, err := memStorage.ListOfEventsForMonth(initialDate)

		require.NoError(t, err)
		require.Len(t, eventsForMonth, 3)
	})
}
