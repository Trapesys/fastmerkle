package fastmerkle

import (
	"fmt"
)

const (
	defaultBufferSize = 100 // arbitrary value; @Yoshiki, what do you think?
)

// workerPool is the pool of worker threads
// that parse hashing jobs
type workerPool struct {
	numWorkers int // Number of workers in the pool

	jobsCh    chan *workerJob    // The channel to receive jobs on
	resultsCh chan *workerResult // The channel to relay results to
}

// newWorkerPool spawns a new worker pool
func newWorkerPool(numWorkers int) *workerPool {
	return &workerPool{
		numWorkers: numWorkers,
		jobsCh:     make(chan *workerJob, defaultBufferSize),
		resultsCh:  make(chan *workerResult, defaultBufferSize),
	}
}

// addJob adds a new job asynchronously to be processed by the worker pool
func (wp *workerPool) addJob(job *workerJob) {
	go func() {
		wp.jobsCh <- job
	}()
}

// getResult takes out a result from the worker pool [Blocking]
func (wp *workerPool) getResult() *workerResult {
	return <-wp.resultsCh
}

// close closes the worker pool and their corresponding
// channels
func (wp *workerPool) close() {
	close(wp.jobsCh)
	close(wp.resultsCh)
}

// startWorkerPool starts the worker routines
func (wp *workerPool) startWorkerPool() {
	// Start the workers
	for i := 0; i < wp.numWorkers; i++ {
		go parseJobs(wp.jobsCh, wp.resultsCh)
	}
}

// workerJob is a single hashing job performed
// by the worker thread
type workerJob struct {
	storeIndex int      // the final store index after hashing
	sourceData [][]byte // the reference to the two items that need to be hashed
}

// workerResult is the result of the worker thread's hashing job
type workerResult struct {
	storeIndex int    // the final store index after hashing
	hashData   []byte // the actual hash result data
	error      error  // any kind of error that occurred during hashing
}

// parseJobs is the main activity method for the
// worker threads, there new jobs are parsed and results sent out
func parseJobs(
	jobsCh <-chan *workerJob,
	resultsCh chan<- *workerResult,
) {
	for job := range jobsCh {
		// Grab an instance of the fast hasher
		hasher := acquireFastHasher()

		// Concatenate all items that need to be hashed together
		preparedArray := make([]byte, 0)
		for i := 0; i < len(job.sourceData); i++ {
			preparedArray = append(preparedArray, job.sourceData[i]...)
		}

		// Hash the items in the job
		var err error
		if writeErr := hasher.addToHash(preparedArray); writeErr != nil {
			err = fmt.Errorf(
				"unable to write hash, %w",
				writeErr,
			)
		}

		// Construct a hash result from the fast hasher
		result := &workerResult{
			storeIndex: job.storeIndex,
			hashData:   hasher.getHash(),
			error:      err,
		}

		// Release the hasher as it's no longer needed
		releaseFastHasher(hasher)

		// Report the result back
		resultsCh <- result
	}
}
