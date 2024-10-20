package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	tasksChan := make(chan Task)
	errorsChan := make(chan error)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go executeTasks(ctx, &wg, tasksChan, errorsChan)
	}

	go sendTasks(ctx, tasksChan, tasks)

	isError := make(chan bool, 1)
	go handleErrors(isError, errorsChan, m, cancel)

	wg.Wait()
	close(errorsChan)

	if <-isError {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func executeTasks(ctx context.Context, wg *sync.WaitGroup, tasksChan <-chan Task, errorsChan chan<- error) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasksChan:
			if !ok {
				return
			}
			err := task()
			if err != nil {
				errorsChan <- err
			}
		}
	}
}

func sendTasks(ctx context.Context, tasksChan chan<- Task, tasks []Task) {
	defer close(tasksChan)

	for _, task := range tasks {
		select {
		case <-ctx.Done():
			return
		case tasksChan <- task:
		}
	}
}

func handleErrors(isError chan<- bool, errorsChan <-chan error, m int, cancel context.CancelFunc) {
	defer close(isError)

	errorsCount := 0
	for range errorsChan {
		errorsCount++
		if errorsCount == m {
			isError <- true
			cancel()
		}
	}
}
