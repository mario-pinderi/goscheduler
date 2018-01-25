package goscheduler

import (
	"reflect"
	"time"
	"errors"
	"github.com/google/uuid"
)

type Scheduler struct {
	jobs map[string]*Job
}
type Job struct {
	Function    interface{}
	Arguments   []reflect.Value
	ExecTime    time.Time
	Repetitive  bool
	RepeatTime  time.Duration
	HasLifetime bool
	Lifetime    time.Time
}

func NewScheduler() *Scheduler {
	sc := Scheduler{}
	sc.jobs = make(map[string]*Job)
	return &sc
}

func (scheduler *Scheduler) AddJob(function interface{}, execTime time.Time, repetitive bool, repeatTime time.Duration, args ...interface{}) error {
	if function != nil && !execTime.IsZero() {
		var paramsArray []reflect.Value
		for _, typ := range args {
			paramsArray = append(paramsArray, reflect.ValueOf(typ))
		}
		id := uuid.New()
		scheduler.jobs[id.String()] = &Job{Function: function, ExecTime: execTime, Arguments: paramsArray, Repetitive: repetitive, RepeatTime: repeatTime}
		return nil
	}
	return errors.New("function and execTime can not be nil")
}

func (scheduler *Scheduler) Job(function interface{}, jobId string) *Job {
	job := Job{Function: function}
	job.Arguments = []reflect.Value{}
	job.ExecTime = time.Now()
	scheduler.jobs[jobId] = &job
	return &job
}

func (job *Job) Args(args ...interface{}) *Job {
	var paramsArray []reflect.Value
	for _, typ := range args {
		paramsArray = append(paramsArray, reflect.ValueOf(typ))
	}
	job.Arguments = paramsArray
	return job
}

func (job *Job) ExecutionTime(execTime time.Time) *Job {
	job.ExecTime = execTime
	return job
}

func (job *Job) RepeatEvery(repeatTime time.Duration) *Job {
	job.Repetitive = true
	job.RepeatTime = repeatTime
	return job
}

func (job *Job) LifeTime(lifeTime time.Time) *Job {
	job.HasLifetime=true
	job.Lifetime = lifeTime
	return job
}

func (scheduler *Scheduler) CleanJobs() {
	scheduler.jobs = make(map[string]*Job)
}

func (scheduler *Scheduler) DeleteJob(jobId string) {
	delete(scheduler.jobs, jobId)
}

func (scheduler *Scheduler) runRemainJobs() {
	currTime := time.Now()
	for index, job := range scheduler.jobs {
		if currTime.After(job.ExecTime) {
			go reflect.ValueOf(job.Function).Call(job.Arguments)
			if job.Repetitive {
				scheduler.jobs[index].ExecTime = job.ExecTime.Add(job.RepeatTime)
				if job.HasLifetime && scheduler.jobs[index].ExecTime.After(job.Lifetime) {
					delete(scheduler.jobs, index)
				}
			} else {
				delete(scheduler.jobs, index)
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
