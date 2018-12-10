package goscheduler

import (
	"fmt"
	"github.com/prometheus/common/log"
	"testing"
	"time"
)

func TestAddJob(t *testing.T) {
	cases := [...]struct{
		name string
		execTime time.Time
	}{
		{"5 second",time.Now().Add(5*time.Second)},
		{"10 second",time.Now().Add(10*time.Second)},
		{"15 second",time.Now().Add(15*time.Second)},
	}

	sc := NewScheduler()
	go func() {
		<-sc.Start()
	}()
	defer sc.close()
	texts := make(chan string)
	for _ , c := range cases{
		t.Run(c.name, func(t *testing.T) {
			sc.AddFunction(func(c chan string, text string) { c <- text }, c.name)
			err := sc.AddJobById(c.name, c.execTime, false, 5*time.Second, texts, t.Name())
			if err != nil {
				log.Fatal(err)
			}
		})
	}

	for range cases{
		fmt.Println(<-texts)
	}

}

func BenchmarkAddJob(t *testing.B) {
	cases := [...]struct{
		name string
		execTime time.Time
	}{
		{"5 second",time.Now().Add(5*time.Second)},
		{"10 second",time.Now().Add(10*time.Second)},
		{"15 second",time.Now().Add(15*time.Second)},
	}

	sc := NewScheduler()
	go func() {
		<-sc.Start()
	}()
	texts := make(chan string)
	for _ , c := range cases{
		t.Run(c.name, func(t *testing.B) {
			sc.AddFunction(func(c chan string, text string) { c <- text }, c.name)
			err := sc.AddJobById(c.name, c.execTime, false, 5*time.Second, texts, t.Name())
			if err != nil {
				log.Fatal(err)
			}
		})
	}

	for range cases{
		fmt.Println(<-texts)
	}

}


func TestCurrentJobs (t *testing.T) {
	sc := NewScheduler()

	_ = sc.GetJobsFromBolt()

	for k, v := range sc.jobs {
		fmt.Println(k, v)
	}
}


