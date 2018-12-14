package goscheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/prometheus/common/log"
	"go.etcd.io/bbolt"
	"reflect"
	"time"
)

const BUCKET_NOT_EXIST  = "bucket does not exist"
const NO_JOBS_IN_DB = "There are no jobs in database"
//Scheduler structure.
type Scheduler struct {
	bucketName []byte
	jobs      map[string]*Job
	functions map[string]interface{}
	dbCon *bbolt.DB
	max chan struct{}
	MaxCon int
}

//Job Structure.
type Job struct {
	FunctionId  string
	Arguments   []interface{}
	ExecTime    time.Time
	Repetitive  bool
	RepeatTime  time.Duration
	HasLifetime bool
	Lifetime    time.Time
}

//Create Scheduler instance.
func NewScheduler(opts ...int) *Scheduler {
	sc := Scheduler{}
	sc.jobs = make(map[string]*Job)
	sc.functions = make(map[string]interface{})
	sc.bucketName = []byte("functions")
	err := sc.boltConnection()
	if len(opts) > 0 && len(opts) < 2 {
		sc.max = make(chan struct{}, opts[0])
	} else {
		sc.max = make(chan struct{}, 100)
	}

	if err != nil {
		log.Debug(err)
	}
	return &sc
}

//Add functions to scheduler.
func (scheduler *Scheduler) AddFunction(function interface{}, functionId string) {
	if _, exist := scheduler.functions[functionId]; function != nil && functionId != "" && !exist {
		scheduler.functions[functionId] = function
	}
}

//Add Job to scheduler.
func (scheduler *Scheduler) AddJobById(functionId string, execTime time.Time, repetitive bool, repeatTime time.Duration, args ...interface{}) error {
	if _, exist := scheduler.functions[functionId]; !execTime.IsZero() && exist {
		/*var paramsArray []reflect.Value
		for _, typ := range args {
			paramsArray = append(paramsArray, reflect.ValueOf(typ))
		}*/
		id := uuid.New()
		scheduler.jobs[id.String()] = &Job{FunctionId: functionId, ExecTime: execTime, Arguments: args, Repetitive: repetitive, RepeatTime: repeatTime}
		return nil
	}
	return errors.New("function and execTime can not be nil")
}

//Create Job instance.
func (scheduler *Scheduler) Job(functionId string, jobId string) *Job {
	job := Job{FunctionId: functionId}
	job.ExecTime = time.Now()
	scheduler.jobs[jobId] = &job
	return &job
}

//Run Jobs if running time has reached.
func (scheduler *Scheduler) runRemainJobs() int {
	currTime := time.Now()
	cnt := 0
	for index, job := range scheduler.jobs {
		var paramsArray []reflect.Value
		if currTime.After(job.ExecTime) {
			if _,exists := scheduler.functions[job.FunctionId]; exists {
				cnt++
				for _, typ := range scheduler.jobs[index].Arguments {
					paramsArray = append(paramsArray, reflect.ValueOf(typ))
				}
				scheduler.max <- struct{}{}
				go func() {
					reflect.ValueOf(scheduler.functions[job.FunctionId]).Call(paramsArray)
					<-scheduler.max
				}()
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

	return cnt
}

func (scheduler *Scheduler) close() {
	fmt.Println("closed")
	err := scheduler.addRemainingJobsToBolt()
	if err != nil {
		log.Debug(err)
	}
	err = scheduler.dbCon.Close()
	if err != nil {
		log.Debug(err)
	}
	return
}

func (scheduler *Scheduler) boltConnection() error {
	if scheduler.dbCon == nil {
		dbCon, err := bbolt.Open("jobs.db", 0666, nil )
		if err != nil {
			return err
		}
		scheduler.dbCon = dbCon
		return nil
	}
	return nil
}

//Start the Scheduler.
func (scheduler *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)
	err := scheduler.GetJobsFromBolt()
	if err != nil && err.Error() != BUCKET_NOT_EXIST {
		log.Error(err)
	}
	if err != nil && err.Error() == NO_JOBS_IN_DB {
		log.Error("No jobs in database")
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				nr := scheduler.runRemainJobs()
				fmt.Println(nr)
				if nr != 0 {
					err := scheduler.addRemainingJobsToBolt()
					if err != nil {
						log.Debug(err)
						return
					}
				}
			case <-stopped:
				return
			}
		}
	}()
	return stopped
}

func (scheduler Scheduler) GetJobsFromBolt() error {
	err := scheduler.dbCon.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(scheduler.bucketName)
		if bucket == nil {
			return errors.New(BUCKET_NOT_EXIST)
		}

		val := bucket.Get([]byte("remainingJobs"))
		if val == nil {
			return errors.New(NO_JOBS_IN_DB)
		}
		err := json.Unmarshal(val, &scheduler.jobs)
		for k,v := range scheduler.jobs {
			fmt.Println(k, v)
		}
		if err != nil {
			return err
		}
		return nil
	})

	return err
}


func (scheduler *Scheduler) addRemainingJobsToBolt() error {
	err := scheduler.dbCon.Update(func(tx *bbolt.Tx) error {
		jobBytes, err := json.Marshal(scheduler.jobs)
		if err != nil {
			return err
		}
		bucket, err := tx.CreateBucketIfNotExists(scheduler.bucketName)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte("remainingJobs"), jobBytes)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

//Add arguments to Job.
func (job *Job) Args(args ...interface{}) *Job {
	/*var paramsArray []reflect.Value
	for _, typ := range args {
		paramsArray = append(paramsArray, reflect.ValueOf(typ))
	}
	job.Arguments = paramsArray*/
	job.Arguments = args
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
	err := scheduler.addRemainingJobsToBolt()
	if err != nil {
		fmt.Println(err)
	}
}

//Delete a Job by id.
func (scheduler *Scheduler) DeleteJob(jobId string) {
	delete(scheduler.jobs, jobId)
	err := scheduler.addRemainingJobsToBolt()
	if err != nil {
		fmt.Println(err)
	}
}