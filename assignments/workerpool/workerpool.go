package main

import (
	"fmt"
	"sync"
	"time"
)

// The 'jobs' channel is read-only (<-chan) and 'results' is write-only (chan<-).
// This ensures type safety and prevents accidental writes to the jobs channel or reads from the results channel.
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement waitgroup counter when worker exits

	// loop blocks until a job is available or channel closed
	for j := range jobs {
		fmt.Printf("Worker %d starting job %d\n", id, j)

		// Simulate a resource-intensive task
		time.Sleep(time.Second * 2)

		output := j * j // simulate task

		fmt.Printf("Worker %d finished job %d (Result: %d)\n", id, j, output)
		results <- output
	}
}

func main() {
	const numJobs = 100
	const numWorkers = 5

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1) // Increment waitgroup counter for each worker
		go worker(w, jobs, results, &wg)
	}

	// send tasks, producer q 100 jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Close the jobs channel to indicate no more jobs will be sent

	wg.Wait() // Wait for all workers to finish

	close(results) // Close the results channel to indicate no more results will be sent

	for r := range results {
		fmt.Println("Result:", r)
	}
}
