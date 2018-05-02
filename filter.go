package main

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	minimumTime time.Time
)

func init() {
	minimumTime, _ = time.Parse(time.RFC3339, "2017-11-02T00:00:00.000Z")
}

func main() {
	var pattern = "./testdata/*.csv.gz"
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	total := 0
	matched := 0

	startTime := time.Now()
	if os.Getenv("GOPAR") == "" { // process one file at a time...
		fmt.Println("Strategy: One file at a time ...")
		for _, fn := range files {
			fmt.Printf("Processing: %s\n", fn)
			t, m, err := ProcessFile(minimumTime, fn)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				total += t
				matched += m
			}
		}
	} else { // do things in parallel
		var mutex sync.Mutex
		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		fmt.Printf("Strategy: Parallel, %d Workers ...\n", numCPU)

		fileCh := make(chan string, numCPU*2)
		go func() {
			for _, fn := range files {
				fileCh <- fn
			}
			fileCh <- "EOF"
		}()

		// start up numCPU number of working goroutines
		for i := 0; i < numCPU; i++ {
			wg.Add(1)
			go func(id int) {
				for fn := range fileCh {
					// there's no more files to process, doesn't matter
					// which worker recieves this first
					if fn == "EOF" {
						close(fileCh) // causes all workers routines to finish
						break
					}

					fmt.Printf("Processing(%d): %s\n", id, fn)
					t, m, err := ProcessFile(minimumTime, fn)
					if err != nil {
						fmt.Printf("Error: %s\n", err.Error())
					} else { // prevent races, need a mutex around writing these values
						mutex.Lock()
						total += t
						matched += m
						mutex.Unlock()
					}
				}

				wg.Done()
			}(i)
		}

		wg.Wait()
	}

	totalTime := time.Now().Sub(startTime)
	fmt.Printf("Total: %d, Matched: %d, Ratio: %0.2f%%\n", total, matched, float64(matched)/float64(total)*100)
	fmt.Printf("Time: %v\n", totalTime)
}

// ProcessFile processes the .csv.gz files as a stream of bytes counting all records that
// meet the minimum date
func ProcessFile(min time.Time, filename string) (total int, matched int, err error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return 0, 0, errors.Wrap(err, "Failed to open file")
	}

	gzreader, err := gzip.NewReader(file)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Failed gzip.NewReader")
	}

	csvreader := csv.NewReader(gzreader)

	for {
		record, err := csvreader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Error CSV Read error: %s\n", err.Error())
			break
		}

		ts, err := time.Parse(time.RFC3339, record[3])
		if err != nil {
			fmt.Printf("Error time.Parse %s\n", err.Error())
			continue
		}

		if ts.After(min) {
			matched += 1
		}

		total += 1
	}

	return total, matched, nil
}
