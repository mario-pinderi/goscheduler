package goscheduler

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"time"
)

//Scheduler structure.
type Scheduler struct {
	jobs map[string]*Job
}

//Job Structure.
type Job struct {
	Function    interface{}
	Arguments   []reflect.Value
	ExecTime    time.Time
	Repetitive  bool
	RepeatTime  time.Duration
	HasLifetime bool
	Lifetime    time.Time
}

//Create Scheduler instance.
func NewScheduler() *Scheduler {
	sc := Scheduler{}
	sc.jobs = make(map[string]*Job)
	return &sc
}

//Add Job to scheduler.
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

//Create Job instance.
func (scheduler *Scheduler) Job(function interface{}, jobId string) *Job {
	job := Job{Function: function}
	job.Arguments = []reflect.Value{}
	job.ExecTime = time.Now()
	scheduler.jobs[jobId] = &job
	return &job
}

//Add arguments to Job.
func (job *Job) Args(args ...interface{}) *Job {
	var paramsArray []reflect.Value
	for _, typ := range args {
		paramsArray = append(paramsArray, reflect.ValueOf(typ))
	}
	job.Arguments = paramsArray
	return job
}

//Set execution time of the Job.
func (job *Job) ExecutionTime(execTime time.Time) *Job {
	job.ExecTime = execTime
	return job
}

//Set reputation of the Job.
func (job *Job) RepeatEvery(repeatTime time.Duration) *Job {
	job.Repetitive = true
	job.RepeatTime = repeatTime
	return job
}

//Set time that the Job will be running.
func (job *Job) LifeTime(lifeTime time.Time) *Job {
	job.HasLifetime = true
	job.Lifetime = lifeTime
	return job
}

//Remove all jobs of the scheduler.
func (scheduler *Scheduler) CleanJobs() {
	scheduler.jobs = make(map[string]*Job)
}

//Delete a Job by id.
func (scheduler *Scheduler) DeleteJob(jobId string) {
	delete(scheduler.jobs, jobId)
}

//Run Jobs if running time has reached.
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

//Start the Scheduler.
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
