package mfe

import (
	"sync"
	"time"
)

// Job do function every duration in go routine or do function wait duration and then repeat
type Job struct {
	Mutex         sync.RWMutex
	Func          func()
	InGoroutine   bool
	do            bool
	IsMultiThread bool
	Duration      time.Duration
	Ch            chan bool
	imx           sync.RWMutex
	Name          string
}

// JobCreate create new job
func JobCreate(f func(), duration time.Duration) (j *Job) {
	j = &(Job{})

	j.Func = f
	j.Duration = duration
	j.Ch = make(chan bool, 1)

	return j
}

// Start job if it not started
func (j *Job) Start() {
	LogActionF(j.Name, "mfe.Job", "Start")
	j.Mutex.Lock()
	defer j.Mutex.Unlock()
	ch := j.Ch
	if !j.do {
		go func() {
			for j.do {
				LogActionF(j.Name, "mfe.Job", "Start.DoIteration")
				select {
				case <-ch:
					LogActionF(j.Name, "mfe.Job", "Start.Stopping")
					close(ch)
					return
				default:
					if j.InGoroutine {
						go j.Func()
					} else if j.IsMultiThread {
						j.Func()
					} else {
						j.imx.Lock()
						j.Func()
						j.imx.Unlock()
					}

					time.Sleep(j.Duration)
				}
			}
		}()
		j.do = true
	}
}

// Stop job
func (j *Job) Stop() {
	LogActionF(j.Name, "mfe.Job", "Stop")
	j.Mutex.Lock()
	defer j.Mutex.Unlock()
	if j.do {
		j.Ch <- true
		j.Ch = make(chan bool, 1)
		j.do = false
	}
}
