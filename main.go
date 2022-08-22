package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second)
		results <- j
	}
}

func worker2(id int) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	job := make(chan int, 10)
	result := make(chan int, 10)
	for w := 1; w <= 9; w++ {
		go worker(w, job, result)
	}
	for j := 1; j <= 9; j++ {
		job <- j
	}
	close(job)
	for a := 1; a <= 9; a++ {
		<-result
	}

	//var wg sync.WaitGroup
	//for i := 1; i <= 9; i++ {
	//	wg.Add(1)
	//
	//	i := i
	//	go func() {
	//		defer wg.Done()
	//		worker2(i)
	//	}()
	//}
	//wg.Wait()
}
