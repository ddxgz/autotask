package autotask

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type AutoTask func() error

type AutoUpdater struct {
	interval    int
	intervalMin int
	timeUnit time.Duration
	started     bool
	done        chan bool
	task        AutoTask
}

type autoUpdaterStatus struct {
	Interval    int `json:"interval"`
	IntervalMin int `json:"interval_min"`
	TimeUnit    time.Duration `json:"time_unit"`
	Started     bool `json:"started"`
}

type Options struct {
	Interval    int
	IntervalMin int
	Task        AutoTask
}

func New(o Options) *AutoUpdater {
	return &AutoUpdater{
		interval:    o.Interval,
		intervalMin: o.IntervalMin,
		timeUnit: time.Hour,
		started:     false,
		done:        make(chan bool),
		task:        o.Task,
	}
}

func (u *AutoUpdater) SetInterval(interval int) error {
	if interval < u.intervalMin {
		return errors.New("Cannot set interval smaller than IntervalMin()")
	}
	u.interval = interval
	return nil
}

func (u *AutoUpdater) Start() {
	// u := FeedUpdater
	// ticker := time.NewTicker(time.Duration(u.interval) * time.Hour)
	ticker := time.NewTicker(time.Duration(u.interval) * u.timeUnit)
	// ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	// done := make(chan bool)
	// u.done <- false
	if u.started {
		return
	}
	u.started = true

	for {
		select {
		case <-u.done:
			fmt.Println("Feed Auto Updater stopped!")
			return
		case <-ticker.C:
			if err := u.task(); err != nil {
				fmt.Println("task stopped due to err, ", err)
				u.Stop()
			}
		}
	}
}

func (u *AutoUpdater) Stop() {
	// u := FeedUpdater
	u.started = false
	u.done <- true
}

func (u *AutoUpdater) Started() bool {
	return u.started
}

func (u *AutoUpdater) IntervalMin() int {
	return u.intervalMin
}

func (u *AutoUpdater) Status() autoUpdaterStatus {
	return autoUpdaterStatus{
		Interval:    u.interval,
		IntervalMin: u.intervalMin,
		TimeUnit: u.timeUnit,
		Started:     u.started,
	}
}
