// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package task

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

var (
	sign     *Task
	document *Task
)

type Task struct {
	name     string      // 任务名称
	dispatch chan *Job   // 分发器
	halt     atomic.Bool // 停止新任务进入

	pool *ants.PoolWithFunc
}

type Job struct {
	worker func() error
	waiter chan error
}

func SignTask() *Task {
	return sign
}

func DocumentTask() *Task {
	return document
}

func CreateTasks(signNumber, documentNumber int) {
	sign = create(signNumber)
	sign.name = "sign"

	document = create(documentNumber)
	document.name = "document"

	go sign.run()
	go document.run()
}

func create(num int) (t *Task) {
	pool, err := ants.NewPoolWithFunc(num, func(data interface{}) {
		t.do(data.(*Job))
	})
	if err != nil {
		zap.L().Fatal("创建任务池失败", zap.Error(err))
	}
	t = &Task{
		dispatch: make(chan *Job, num),
		pool:     pool,
	}
	return
}

// 启动任务队列
func (t *Task) run() {
	defer t.pool.Release()

	for {
		select {
		case job := <-t.dispatch:
			err := t.pool.Invoke(job)
			if err != nil {
				zap.L().Error("任务执行失败", zap.Error(err))
			}
		}
	}
}

// 执行任务
func (t *Task) do(job *Job) {
	start := time.Now()

	var err error
	defer func() {
		// 防止崩溃
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			zap.L().Error("任务崩溃", zap.String("task", t.name), zap.Error(err))
		}
		// 返回结果
		zap.L().Info("执行完成",
			zap.String("task", t.name),
			zap.Duration("cost", time.Since(start)),
			zap.Error(err),
			zap.Int("running", t.pool.Running()),
		)
		job.waiter <- err
	}()

	if t.halt.Load() {
		zap.L().Info("停止任务新增，直接返维护中", zap.String("task", t.name))
		err = errors.New("维护中")
		return
	}

	// 执行任务
	zap.L().Info("新增任务", zap.String("task", t.name), zap.Int("running", t.pool.Running()))
	err = job.worker()
	return
}

func (t *Task) AddJob(worker func() error) (waiter chan error) {
	waiter = make(chan error)
	t.dispatch <- &Job{
		worker: worker,
		waiter: waiter,
	}
	return
}

// Halt 停止任务队列
func (t *Task) Halt() {
	t.halt.Store(true)
}

// Wait 等待任务队列完成
func (t *Task) Wait() {
	for {
		isDone := t.pool.Running() == 0
		if isDone {
			return
		}
		fmt.Printf("%s 任务队列中还有 %d 个任务在执行, %d 个任务等待中\n", t.name, t.pool.Running(), t.pool.Waiting())
		time.Sleep(500 * time.Millisecond)
	}
}

// HaltAndWait 停止任务队列并等待执行中的队列完成
func (t *Task) HaltAndWait(wg *sync.WaitGroup) {
	t.Halt()
	t.Wait()
	wg.Done()
}
