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
            // If you wish to limit the number of tasks running at a time, enter the number as an argument while calling NewScheduler
            // The default value is max 100 tasks at a time ex. sc := goscheduler.NewScheduler(20)
           
            sc.AddFunction(func(tt string) { fmt.Println("first") }, "fnc1")
           	_ = sc.AddJobById("fnc1", time.Now().Add(5*time.Second), false, time.Second,"test1")
           
           	// job repeated
           	sc.AddFunction(func(tt string) { fmt.Println(tt) }, "fnc2")
           	_ = sc.AddJobById("fnc2", time.Now().Add(10*time.Second), false, time.Second,"test2")
           
           
           	// job repeated
           	sc.AddFunction(func(tt string) { fmt.Println(tt) }, "fnc3")
           	_ = sc.AddJobById("fnc3", time.Now().Add(15*time.Second), false, time.Second,"test3")
           
           	<-sc.Start()
        }
    
# TODO:

* Prioritize tasks 