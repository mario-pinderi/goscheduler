package goscheduler

import (
	"reflect"
	"time"
	"errors"
)

type Scheduler struct {
	jobs []Job
}
type Job struct {
	Function   interface{}
	Arguments  []reflect.Value
	ExecTime   time.Time
	Repetitive bool
	RepeatTime time.Duration
}

func NewScheduler() *Scheduler {
	sc := Scheduler{}
	return &sc
}

func (scheduler *Scheduler) AddJob(function interface{}, execTime time.Time, repetitive bool, repeatTime time.Duration, args ...interface{}) error {
	if function != nil && !execTime.IsZero() {
		var paramsArray []reflect.Value
		for _, typ := range args {
			paramsArray = append(paramsArray, reflect.ValueOf(typ))
		}
		scheduler.jobs = append(scheduler.jobs, Job{Function: function, ExecTime: execTime, Arguments: paramsArray, Repetitive: repetitive, RepeatTime: repeatTime})
		return nil
	}
	return errors.New("function and execTime can not be nil")
}

func (scheduler *Scheduler) CleanJobs() {
	scheduler.jobs = []Job{}
}

func (scheduler *Scheduler) runRemainJobs() {
	currTime := time.Now()
	for index, job := range scheduler.jobs {
		if currTime.After(job.ExecTime) {
			go reflect.ValueOf(job.Function).Call(job.Arguments)
			if job.Repetitive {
				scheduler.jobs[index].ExecTime = scheduler.jobs[index].ExecTime.Add(job.RepeatTime)
			} else {
				scheduler.jobs = append(scheduler.jobs[:index], scheduler.jobs[index+1:]...)
			}
		}
	}
}

func (scheduler *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				scheduler.runRemainJobs()
			case <-stopped:
				return
			}
		}
	}()
	return stopped
}
