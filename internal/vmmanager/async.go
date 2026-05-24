package vmmanager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type asyncTask struct {
	ctx  context.Context
	name string
	run  func(context.Context) error
}

type asyncRunner struct {
	logger *logrus.Logger
	tasks  chan asyncTask
	wg     sync.WaitGroup
}

func newAsyncRunner(workerCount int, logger *logrus.Logger) *asyncRunner {
	if workerCount <= 0 {
		return nil
	}

	runner := &asyncRunner{
		logger: logger,
		tasks:  make(chan asyncTask, workerCount*2),
	}

	for i := 0; i < workerCount; i++ {
		runner.wg.Add(1)
		go runner.worker(i)
	}

	return runner
}

func (r *asyncRunner) Submit(ctx context.Context, name string, run func(context.Context) error) error {
	if r == nil {
		return run(ctx)
	}

	task := asyncTask{ctx: ctx, name: name, run: run}
	select {
	case r.tasks <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("async queue is full")
	}
}

func (r *asyncRunner) Close() {
	if r == nil {
		return
	}
	close(r.tasks)
	r.wg.Wait()
}

func (r *asyncRunner) worker(id int) {
	defer r.wg.Done()
	for task := range r.tasks {
		start := time.Now()
		ctx := task.ctx
		if ctx == nil {
			ctx = context.Background()
		}
		err := task.run(ctx)
		if err != nil {
			if r.logger != nil {
				r.logger.WithError(err).WithField("task", task.name).Error("Async task failed")
			}
			continue
		}
		if r.logger != nil {
			r.logger.WithField("task", task.name).WithField("duration_ms", time.Since(start).Milliseconds()).Debug("Async task completed")
		}
	}
}
