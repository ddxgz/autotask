package autotask_test

import (
	"time"

	"fmt"
	"testing"

	"github.com/ddxgz/autotask"
	"github.com/pkg/errors"
)

func autoTaskFunc() error {
	for i := 1; i <= 5; i++ {
		fmt.Printf("auto task running %v step at: %v\n", i, time.Now())
		time.Sleep(1 * time.Second)
	}
	return nil
}

func autoTaskErrFunc() error {
	for i := 1; i <= 2; i++ {
		fmt.Printf("auto task running %v step at: %v\n", i, time.Now())
		time.Sleep(1 * time.Millisecond)
	}
	return errors.New("auto task err func error")
}

// func TestStart(t *testing.T){
// // use a chan to receive signal from the running task
// }

func TestStartStop(t *testing.T) {

	var tasker = autotask.New(autotask.Options{
		Interval:    7,
		IntervalMin: 5,
		Task:        autoTaskFunc,
	})

	if tasker.Started() == true {
		t.Errorf("task started before really start!")
	}

	tasker.SetTimeUnit(time.Second)

	go tasker.Start()

	time.Sleep(1 * time.Millisecond)
	if tasker.Started() != true {
		t.Errorf("task not started after start!")
	}

	tasker.Stop()
	time.Sleep(1 * time.Millisecond)
	if tasker.Started() == true {
		t.Errorf("task not stopped after call Stop!")
	}

}

func TestErrInTask(t *testing.T) {

	var tasker = autotask.New(autotask.Options{
		Interval:    3,
		IntervalMin: 3,
		Task:        autoTaskErrFunc,
	})
	tasker.SetTimeUnit(time.Millisecond)

	go tasker.Start()

	time.Sleep(1 * time.Millisecond)
	if tasker.Started() != true {
		t.Errorf("task not started after start!")
	}

	time.Sleep(9 * time.Millisecond)
	if tasker.Started() == true {
		t.Errorf("task not stopped after error occurred!")
	}

}

func TestStatus(t *testing.T) {
	var tasker = autotask.New(autotask.Options{
		Interval:    7,
		IntervalMin: 5,
		Task:        autoTaskFunc,
	})

	s := tasker.Status()

	if s.Interval != 7 {
		t.Errorf("Interval expected:%v, got:%v", 7, s.Interval)
	}
	if s.IntervalMin != 5 {
		t.Errorf("IntervalMin expected:%v, got:%v", 5, s.IntervalMin)
	}
	if s.TimeUnit != time.Hour {
		t.Errorf("TimeUnit expected:%v, got:%v", time.Hour, s.TimeUnit)
	}
	if s.Started != false {
		t.Errorf("Started expected:%v, got:%v", false, s.Started)
	}

}
