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

func main() {
	fmt.Println("fdsfds")
	job := make(chan int, 10)
	result := make(chan int, 10)
	for w := 1; w <= 10; w++ {
		go worker(w, job, result)
	}
	for j := 1; j <= 9; j++ {
		job <- j
	}
	close(job)
	for a := 1; a <= 9; a++ {
		<-result
	}
}
