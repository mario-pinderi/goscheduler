# goScheduler

goScheduler is a Golang package for scheduling tasks. This tasks can be run once in a selected time ,repeatedly or repeatedly with a lifetime. All tasks will be run concurrently. 
Task properties:
* Function - function that will be run
* Arguments - parameters that will be passed to the function
* Execution time - time when task starts executing
* RepeatTime - how often function will be run
* LifeTime - time that task will be available


        package main
        
        import (
            "fmt"
            "time"
            "github.com/mariopinderist/goscheduler"
        )
        
        func main() {
            sc := goscheduler.NewScheduler()
            //normal job
            sc.AddJob(func(tt string) { fmt.Println(tt) }, time.Now().Add(5*time.Second), false, time.Second, "5 sec")
        
            // job repeated
            sc.Job(func(tt string) { fmt.Println(tt) }, "repeat").Args("repeat").RepeatEvery(3 * time.Second)
        
            // job repeated
            sc.Job(func(tt string) { fmt.Println(tt) }, "lifetime").Args("lifetime").RepeatEvery(3 * time.Second).LifeTime(time.Now().Add(20 * time.Second))
        
            <-sc.Start()
        }
    
# TODO:

* Running x task at a time to avoid overwhelming
* Prioritize tasks 