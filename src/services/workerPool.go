package services

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gammazero/deque"
)

const (
	idleTimeout = 2 * time.Second
)

var Wp *WorkerPool

type WorkerPool struct {
	maxWorkers   int
	taskQueue    chan func()
	workerQueue  chan func()
	stoppedChan  chan struct{}
	stopSignal   chan struct{}
	waitingQueue deque.Deque[func()]
	stopLock     sync.Mutex
	stopOnce     sync.Once
	stopped      bool
	waiting      int32
	wait         bool
}

func New(maxWorkers int) *WorkerPool {

	if maxWorkers < 1 {
		maxWorkers = 1
	}

	Wp = &WorkerPool{
		maxWorkers:  maxWorkers,
		taskQueue:   make(chan func()),
		workerQueue: make(chan func()),
		stopSignal:  make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}

	go Wp.dispatch()
	return Wp
}

func GetPool() *WorkerPool {
	return Wp
}

func (p *WorkerPool) stop(wait bool) {
	p.stopOnce.Do(func() {

		close(p.stopSignal)

		p.stopLock.Lock()

		p.stopped = true
		p.stopLock.Unlock()
		p.wait = wait

		close(p.taskQueue)
	})
	<-p.stoppedChan
}

func (p *WorkerPool) Size() int {
	return p.maxWorkers
}

func (p *WorkerPool) Stop() {
	p.stop(false)
}

func (p *WorkerPool) StopWait() {
	p.stop(true)
}

func (p *WorkerPool) Stopped() bool {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()
	return p.stopped
}

func (p *WorkerPool) Submit(task func()) {
	if task != nil {
		p.taskQueue <- task
	}
}

func (p *WorkerPool) SubmitWait(task func()) {
	if task == nil {
		return
	}
	doneChan := make(chan struct{})
	p.taskQueue <- func() {
		task()
		close(doneChan)
	}
	<-doneChan
}

func (p *WorkerPool) WaitingQueueSize() int {
	return int(atomic.LoadInt32(&p.waiting))
}

func (p *WorkerPool) Pause(ctx context.Context) {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()
	if p.stopped {
		return
	}
	ready := new(sync.WaitGroup)
	ready.Add(p.maxWorkers)
	for i := 0; i < p.maxWorkers; i++ {
		p.Submit(func() {
			ready.Done()
			select {
			case <-ctx.Done():
			case <-p.stopSignal:
			}
		})
	}

	ready.Wait()
}

func (p *WorkerPool) dispatch() {
	defer close(p.stoppedChan)
	timeout := time.NewTimer(idleTimeout)
	var workerCount int
	var idle bool
	var wg sync.WaitGroup

Loop:
	for {

		if p.waitingQueue.Len() != 0 {
			if !p.processWaitingQueue() {
				break Loop
			}
			continue
		}

		select {
		case task, ok := <-p.taskQueue:
			if !ok {
				break Loop
			}

			select {
			case p.workerQueue <- task:
			default:

				if workerCount < p.maxWorkers {
					wg.Add(1)
					go worker(task, p.workerQueue, &wg)
					workerCount++
				} else {

					p.waitingQueue.PushBack(task)
					atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len()))
				}
			}
			idle = false
		case <-timeout.C:

			if idle && workerCount > 0 {
				if p.killIdleWorker() {
					workerCount--
				}
			}
			idle = true
			timeout.Reset(idleTimeout)
		}
	}

	if p.wait {
		p.runQueuedTasks()
	}

	for workerCount > 0 {
		p.workerQueue <- nil
		workerCount--
	}
	wg.Wait()

	timeout.Stop()
}

func worker(task func(), workerQueue chan func(), wg *sync.WaitGroup) {
	for task != nil {
		task()
		task = <-workerQueue
	}
	wg.Done()
}

func (p *WorkerPool) processWaitingQueue() bool {
	select {
	case task, ok := <-p.taskQueue:
		if !ok {
			return false
		}
		p.waitingQueue.PushBack(task)
	case p.workerQueue <- p.waitingQueue.Front():

		p.waitingQueue.PopFront()
	}
	atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len()))
	return true
}

func (p *WorkerPool) killIdleWorker() bool {
	select {
	case p.workerQueue <- nil:

		return true
	default:

		return false
	}
}

func (p *WorkerPool) runQueuedTasks() {
	for p.waitingQueue.Len() != 0 {

		p.workerQueue <- p.waitingQueue.PopFront()
		atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len()))
	}
}
