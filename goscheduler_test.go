package goscheduler

import (
	"fmt"
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
	texts := make(chan string)
	for _ , c := range cases{
		t.Run(c.name, func(t *testing.T) {
			sc.AddJob(func(c chan string, text string) { c <- text }, c.execTime, false, 5*time.Second, texts, t.Name())
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
			sc.AddJob(func(c chan string, text string) { c <- text }, c.execTime, false, 5*time.Second, texts, t.Name())
		})
	}

	for range cases{
		fmt.Println(<-texts)
	}

}


