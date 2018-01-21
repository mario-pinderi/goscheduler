# goscheduler

go Scheduler is a Golang package for scheduling tasks. This tasks can be run once in a selected time or repeatedly. All tasks will be run concurrently. 

    package main
    
    import (
        "fmt"
        "time"
        "github.com/mariopinderist/goscheduler"
    )
    
    func main() {
        sc := goscheduler.NewScheduler()
        sc.AddJob(func(tt string) { fmt.Println(tt) }, time.Now().Add(5*time.Second), true, time.Second, "5 sec")
        sc.AddJob(func(tt string) { fmt.Println(tt) }, time.Now().Add(10*time.Second), false, time.Second, "10 sec")
        sc.AddJob(func(tt string) { fmt.Println(tt) }, time.Now().Add(15*time.Second), false, time.Second, "15 sec")
        sc.AddJob(func(tt string) { fmt.Println(tt) }, time.Now().Add(20*time.Second), false, time.Second, "20 sec")
        fakeJobs(sc)
        <-sc.Start()
    }
    
#TODO:

* Running x task at a time to avoid overwhelming
* Prioritize tasks 