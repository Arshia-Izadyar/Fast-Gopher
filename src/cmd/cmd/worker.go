package cmd

import (
	"fmt"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
)

type Task func(*dto.WhiteListAddDTO) error

type Result struct {
	workerID int
	i        *dto.WhiteListAddDTO
	// data     string
	err error
}

type Worker struct {
	id         int
	taskQueue  <-chan *dto.WhiteListAddDTO
	resultChan chan<- Result
}

func (w *Worker) Start() {
	go func() {
		for i := range w.taskQueue {
			err := test(i)
			w.resultChan <- Result{workerID: w.id, i: i, err: err}
		}
	}()
}

type WorkerPool struct {
	taskQueue   chan *dto.WhiteListAddDTO
	resultChan  chan Result
	workerCount int
}

var W *WorkerPool

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		taskQueue:   make(chan *dto.WhiteListAddDTO),
		resultChan:  make(chan Result),
		workerCount: workerCount,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		worker := Worker{id: i, taskQueue: wp.taskQueue, resultChan: wp.resultChan}
		worker.Start()
	}
}

func (wp *WorkerPool) Submit(i *dto.WhiteListAddDTO) {
	wp.taskQueue <- i
}

func (wp *WorkerPool) GetResult() Result {
	return <-wp.resultChan
}

func test(req *dto.WhiteListAddDTO) error {
	db := postgres.GetDB()
	tx, err := db.Begin()
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.InternalError}
	}

	q := `
	DELETE FROM active_devices where user_id = $1 AND device_id = $2;
	`

	if _, err = tx.Exec(q, req.UserId, req.UserDeviceID); err != nil {
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: "deletion failed"}
	}

	tx.Commit()
	return nil
}

func InitWorker(i int) *WorkerPool {
	W = NewWorkerPool(3)
	return W
}

func (wp *WorkerPool) StartResultLogger() {
	go func() {
		for result := range wp.resultChan {
			if result.err != nil {
				fmt.Printf("Error from worker %d processing item %s: %v\n", result.workerID, result.i, result.err)
			}
		}
	}()
}
