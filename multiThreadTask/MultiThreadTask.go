package multiThreadTask

import (
	"context"
	"errors"
	"github.com/aadog/dict-go"
	"github.com/aadog/tasks-go/board"
	"github.com/gammazero/workerpool"
	"time"
)

type Result struct {
	StartTime time.Time
	RunTime   time.Duration
}

type MultiThreadTask struct {
	ctx             context.Context
	data            *dict.DictList
	thread          int
	processBeforeFn func(idx int, d *dict.Dict)
	processAfterFn  func(idx int, d *dict.Dict)
	errorFn         func(idx int, err error, d *dict.Dict)
	cboard          *board.CountBoard

	startTime time.Time
	runTime   time.Duration
}

func (m *MultiThreadTask) SetProcessBeforeFn(processFn func(idx int, d *dict.Dict)) *MultiThreadTask {
	m.processBeforeFn = processFn
	return m
}
func (m *MultiThreadTask) SetProcessAfterFunc(processFn func(idx int, d *dict.Dict)) *MultiThreadTask {
	m.processAfterFn = processFn
	return m
}
func (m *MultiThreadTask) SetErrorFunc(errorFn func(idx int, err error, d *dict.Dict)) *MultiThreadTask {
	m.errorFn = errorFn
	return m
}
func (m *MultiThreadTask) Board() *board.CountBoard {
	return m.cboard
}
func (m *MultiThreadTask) SetBoard(b *board.CountBoard) *MultiThreadTask {
	m.cboard = b
	return m
}

func (m *MultiThreadTask) Run(fn func(idx int, d *dict.Dict) error) *Result {
	taskres := &Result{}
	count := m.data.Len()
	taskres.StartTime = time.Now()
	defer func() {
		taskres.RunTime = time.Since(taskres.StartTime)
	}()
	pool := workerpool.New(m.thread)
	for i := 0; i < int(count); i++ {
		idx := i
		it, ok := m.data.Get(idx)
		if ok {
			pool.Submit(func() {
				if m.processBeforeFn != nil {
					m.processBeforeFn(idx, it)
				}
				defer func() {
					if m.processAfterFn != nil {
						m.processAfterFn(idx, it)
					}
				}()
				err := func() error {
					select {
					case <-m.ctx.Done():
						return errors.New("任务被取消")
					default:
						err := fn(idx, it)
						if err != nil {
							return err
						}
						m.cboard.AddSuccess(1)
					}
					return nil
				}()
				if err != nil {
					if m.errorFn != nil {
						m.errorFn(idx, err, it)
					}
					m.cboard.AddError(1)
				}
			})
		}
	}
	pool.StopWait()
	return taskres
}

func New(ctx context.Context, data *dict.DictList, thread int) *MultiThreadTask {
	m := &MultiThreadTask{}
	m.ctx = ctx
	m.data = data
	m.thread = thread
	return m
}
