package randgen

import (
	"math/rand"
)

// Generator for number generation in go is not concurrently safe (for new source).
// In order to generate random number concurrently, a number of workers to generate the random numbers are needed.
// each worker will generate a random number and push it into a buffered channel
// Reference: https://golang.org/pkg/math/rand/
type Generator struct {
	randnum chan int
	woker   []RandWorker
}

// RandWorker struct
type RandWorker struct {
	randchan chan int
	randgen  *rand.Rand
	seed     int64
	min      int
	max      int
	stopped  bool
}

// Work via goroutines
// the worker will run forever when program runs
// but will get destroyed if the program exit
func (rw *RandWorker) Work() {
	go func() {
		for {
			// the worker will not stopped immediately
			// because it will wait the buffered channel to have space
			// closing channel might help this case
			if rw.stopped {
				return
			}
			randnum := rw.randgen.Intn(rw.max - rw.min)
			randnum = rw.min + randnum
			rw.randchan <- randnum
		}
	}()
}

// New random number generator
func New(workernumber, min, max int, seed int64) *Generator {
	randnumchan := make(chan int, workernumber*10)
	gen := Generator{
		randnum: randnumchan,
	}

	for i := 0; i < workernumber; i++ {
		r := rand.New(rand.NewSource(seed))
		w := RandWorker{
			randchan: randnumchan,
			randgen:  r,
			seed:     seed,
			min:      min,
			max:      max,
		}
		w.Work()
		gen.woker = append(gen.woker, w)
	}
	return &gen
}

// Generate a new random number
func (gen *Generator) Generate() int {
	return <-gen.randnum
}

// Stop generator's worker
func (gen *Generator) Stop() {
	for i := range gen.woker {
		gen.woker[i].stopped = true
	}
	// close the generator channel
	close(gen.randnum)
}
