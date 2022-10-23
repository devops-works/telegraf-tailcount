package tailcount

import (
	"fmt"
	"io"
	"math"
	"os"
	"time"

	"github.com/nxadm/tail"
)

// Counter counts new lines every interval in a file
type Counter struct {
	file         string
	interval     int
	peakInterval int
	tags         string
	measurement  string
}

// NewCounter creates a new counter
func NewCounter(file string, o ...func(*Counter) error) (*Counter, error) {
	ctr := &Counter{file: file}

	// Apply functional options
	for _, f := range o {
		if err := f(ctr); err != nil {
			return nil, err
		}
	}

	if ctr.peakInterval > ctr.interval {
		return nil, fmt.Errorf("peak interval (%d) cannot be greater than interval (%d)", ctr.peakInterval, ctr.interval)
	}

	div := (float64)(ctr.interval) / (float64)(ctr.peakInterval)
	if div != math.Trunc(div) {
		return nil, fmt.Errorf("interval (%d) must be a multiple of peakinterval (%d)", ctr.interval, ctr.peakInterval)
	}

	return ctr, nil
}

// WithInterval sets the telegraf collection interval
func WithInterval(interval int) func(*Counter) error {
	return func(c *Counter) error {
		c.interval = interval
		return nil
	}
}

// WithPeakInterval sets the telegraf aggregation interval for peaks
func WithPeakInterval(interval int) func(*Counter) error {
	return func(c *Counter) error {
		c.peakInterval = interval
		return nil
	}
}

// WithMeasurement sets the measurement name
func WithMeasurement(measurement string) func(*Counter) error {
	return func(c *Counter) error {
		c.measurement = measurement
		return nil
	}
}

// WithTags sets the tags
func WithTags(tags string) func(*Counter) error {
	return func(c *Counter) error {
		c.tags = tags
		return nil
	}
}

// Run starts the counter
func (c *Counter) Run() {
	peakBuckets := make([]int, c.interval/c.peakInterval)
	peakStart := time.Now()

	// Create a tail
	t, err := tail.TailFile(c.file, tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
		Logger:    tail.DiscardingLogger,
		Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
		},
		Poll: true,
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	ticker := time.NewTicker((time.Duration)(c.interval) * time.Second)

	for {
		select {
		case <-t.Lines:
			index := (int)(time.Since(peakStart).Seconds() / (float64)(c.peakInterval))
			peakBuckets[int(index)]++
		case <-ticker.C:
			fmt.Printf("%s %s sum=%d,max=%d,min=%d\n", c.measurement, c.tags, sum(peakBuckets), max(peakBuckets), min(peakBuckets))
			peakStart = time.Now()
			peakBuckets = make([]int, c.interval/c.peakInterval)
		}
	}
}

func sum(s []int) int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func max(s []int) int {
	max := s[0]
	for _, v := range s[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func min(s []int) int {
	min := s[0]
	for _, v := range s[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// func dumpSlice(s []int) {
// 	for _, v := range s {
// 		fmt.Printf("%d ", v)
// 	}
// 	fmt.Println()
// }
