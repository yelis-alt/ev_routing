package additional

import (
	"runtime"
	"sync"
)

// ParallelWorkers is the default concurrency for CPU-bound parallel work,
// sized to the number of logical CPUs available to the process.
func ParallelWorkers() int {
	if n := runtime.GOMAXPROCS(0); n > 0 {
		return n
	}

	return 1
}

// ParallelFor runs worker(i) for every i in [0,n) across a bounded pool of
// goroutines (sized by ParallelWorkers), blocking until all calls return.
func ParallelFor(n int, worker func(i int)) {
	if n <= 0 {
		return
	}

	workers := min(ParallelWorkers(), n)

	jobs := make(chan int)

	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for i := range jobs {
				worker(i)
			}
		})
	}

	for i := range n {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
}
