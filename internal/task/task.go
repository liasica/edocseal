// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-11, by liasica

package task

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/liasica/edocseal/internal/g"
)

var instance *Task

type Task struct {
	dispatch chan *Job
}

type Job struct {
	worker JobWorker
	waiter chan error
}

type JobWorker func() error
type JobWaiter func() chan error

func NewTask() *Task {
	if instance == nil {
		instance = &Task{
			dispatch: make(chan *Job, g.GetTaskNum()),
		}
	}
	return instance
}

func (t *Task) Run() {
	for {
		select {
		case job := <-t.dispatch:
			job.waiter <- t.do(job)
		}
	}
}

func (t *Task) do(job *Job) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			zap.L().Info("签约任务崩溃", zap.Error(err))
		}
	}()
	err = job.worker()
	return
}

func (t *Task) AddJob(worker JobWorker) (waiter chan error) {
	waiter = make(chan error)
	t.dispatch <- &Job{
		worker: worker,
		waiter: waiter,
	}
	return
}
