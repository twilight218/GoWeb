package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTaskControdl(t *testing.T) {
	taskNum := 5

	wg := sync.WaitGroup{}
	wg.Add(taskNum)

	for i := 0; i < taskNum; i++ {
		go func(i int) {
			fmt.Println("info", i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func Test(t *testing.T){
	test := make(chan int, 10)
   
	go func(info chan int){
		for {
			select {
			case val, ok := <- test:
				if !ok {
					t.Logf("Channel Closed!")
					return
				}
				t.Logf("data %d\n", val)
			}
		}
   }(test)
   
	go func(){
		test <- 1
		time.Sleep(1 * time.Second)
		test <- 2
   		close(test)
   	}()
   
	time.Sleep(5 * time.Second)
}

func TestA(t *testing.T) {
    test := make(chan int, 5)
    exit := make(chan struct{})

	go func(info chan int, exit chan struct{}) {
		for {
			select {
			case val := <-info:
				t.Logf("data %d\n", val)

			case <-exit:
				t.Logf("Task Exit!!\n")
				return
			}
		}
	}(test, exit)

    go func() {
        test <- 1
        time.Sleep(1 * time.Second)
        test <- 2
        close(exit)
    }()
	time.Sleep(5 *time.Second)
}

func TestTimeOut(t *testing.T){
	test := make(chan int, 5)
   
   	go func(info chan int) {
		for{
			select {
			case val := <- info:
				t.Logf("Data %d\n", val)
		
			case <- time.After(2 * time.Second):
				t.Logf("Time out!\n")
				return
			}
		}
	}(test)
   
	go func(){
		test <- 1
		time.Sleep(1 * time.Second) //>=2
		test <- 2
	}()
   
	time.Sleep(5 *time.Second)
}


func task (name string, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
   
	fmt.Printf("Task %s started\n", name)
   
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Task %s stopped\n", name)
			return
		default:	
			time.Sleep(500 * time.Millisecond)
		}
	}
}
   
func testFinish(t *testing.T) {
   
	ctxA, cancelA := context.WithCancel(context.Background())
   
   
	ctxB, cancelB := context.WithCancel(ctxA)
	ctxC, cancelC := context.WithCancel(ctxA)
	ctxD, _ := context.WithCancel(ctxA)
   
	ctxE, _ := context.WithCancel(ctxB)
	ctxF, _ := context.WithCancel(ctxB)
   
	ctxG, _ := context.WithCancel(ctxC)
   
	wg := sync.WaitGroup{}
   
	wg.Add(1)
	go task("A", ctxA, &wg)
   
	wg.Add(1)
	go task("B", ctxB, &wg)
   
	wg.Add(1)
	go task("C", ctxC, &wg)
   
	wg.Add(1)
	go task("D", ctxD, &wg)
   
	wg.Add(1)
	go task("E", ctxE, &wg)
   
	wg.Add(1)
	go task("F", ctxF, &wg)
   
	wg.Add(1)
	go task("G", ctxG, &wg)
   
	time.Sleep(2 * time.Second)
   
	cancelB()
	time.Sleep(1 * time.Second)
   
	cancelC()
	time.Sleep(1 * time.Second)
   
	cancelA()
	time.Sleep(1 * time.Second)
   
	wg.Wait()
	fmt.Println("All tasks stopped")
} 