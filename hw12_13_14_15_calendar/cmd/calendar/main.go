package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/cmd"
	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/app"
	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/logger"
	server "github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/server/http"
	"github.com/nislovskaya/go_prof_course/hw12_13_14_15_calendar/internal/storage"
)

const appName = "calendar"

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)
	defer cancel()

	config := cmd.GetConfig(cmd.ConfigFile)

	logg := logger.New(appName, config.Logger.Level, os.Stdout)

	eventStorage, err := storage.GetStorage(config)
	if err != nil {
		logg.Error(fmt.Sprintf("failed to get storage instance: %s", err))
	}
	err = eventStorage.Connect(ctx)
	if err != nil {
		logg.Error(fmt.Sprintf("failed to connect to storage: %s", err))
	}

	calendar := app.New(logg, eventStorage)

	srv := server.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err = srv.Stop(ctx); err != nil {
			logg.Error("failed to stop http srv: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err = srv.Start(ctx); err != nil {
		logg.Error("failed to start http srv: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
