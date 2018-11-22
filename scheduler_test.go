package goscheduler

import (
	"testing"
	"fmt"
	"time"
)

func TestAddJob(t *testing.T) {
	sc := NewScheduler()
	sc.AddJob(func(tt string) { fmt.Println(tt) }, time.Now().Add(5*time.Second), false, time.Second, "5 sec")
	sc.AddJob(func(tt string) { fmt.Println(tt) }, time.Now().Add(10*time.Second), false, time.Second, "10 sec")
	sc.AddJob(func(tt string) { fmt.Println(tt) }, time.Now().Add(15*time.Second), false, time.Second, "15 sec")
	sc.AddJob(func(tt string) { fmt.Println(tt) }, time.Now().Add(20*time.Second), false, time.Second, "20 sec")
	<-sc.Start()
}

func TestJob(t *testing.T) {
	sc := NewScheduler()
	job := sc.Job(func(tt string) { fmt.Println(tt) },"test")
	job.Args("test").RepeatEvery(3*time.Second)
	go func() {
		time.Sleep(10*time.Second)
		sc.DeleteJob("test")
	}()
	<-sc.Start()
}

func TestLifeTime(t *testing.T) {
	sc := NewScheduler()
	job := sc.Job(func(tt string) { fmt.Println(tt) },"test")
	job.Args("test").RepeatEvery(3*time.Second).LifeTime(time.Now().Add(30*time.Second))

	<-sc.Start()
}