package infocrawler

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)


type AutoTask func() error

type AutoUpdater struct {
	interval    int
	intervalMin int
	started     bool
	done        chan bool
	task        AutoTask
}

type autoUpdaterStatus struct {
	Interval    int  `json:"interval"`
	IntervalMin int  `json:"interval_min"`
	Started     bool `json:"started"`
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
	ticker := time.NewTicker(time.Duration(u.interval) * time.Second)
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
		Started:     u.started,
	}
}
