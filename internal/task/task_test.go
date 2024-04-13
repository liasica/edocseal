// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-13, by liasica

package task

import (
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestHalt(t *testing.T) {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	l.Info("TTT")
	task1 := create(5)
	task1.name = "test"
	go task1.run()

	tw := &sync.WaitGroup{}
	tw.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			t.Log(i)
			w := task1.AddJob(func() error {
				time.Sleep(time.Second * 5)
				return nil
			})
			t.Logf("任务结果: %v", <-w)
			tw.Done()
		}()
	}

	time.Sleep(time.Second * 2)
	t.Logf("task1 任务数量: %d", task1.pool.Running())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go task1.HaltAndWait(wg)

	wg.Wait()
	t.Log("Halted")

	tw.Wait()
}
