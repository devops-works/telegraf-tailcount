package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/devops-works/telegraf-tailcount/internal/tailcount"
)

var (
	// Version of current binary
	Version string
	// BuildDate of current binary
	BuildDate string
)

type config struct {
	file         string
	interval     int
	peakInterval int
	tags         string
	measurement  string
}

func main() {
	c, err := getArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("error parsing arguments: %v\n", err)
		usage()
		os.Exit(1)
	}

	// create a tailcount counter
	counter, err := tailcount.NewCounter(c.file,
		tailcount.WithInterval(c.interval),
		tailcount.WithPeakInterval(c.peakInterval),
		tailcount.WithMeasurement(c.measurement),
		tailcount.WithTags(c.tags),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	counter.Run()
}

func getArgs(args []string) (*config, error) {
	var err error

	if len(args)%2 != 1 {
		return nil, fmt.Errorf("wrong number of arguments")
	}

	c := config{
		interval:     10,
		peakInterval: 1,
		measurement:  "tailcount",
		tags:         "",
		file:         "",
	}

	i := slicePos("-i", args)
	if i != -1 {
		c.interval, err = strconv.Atoi(args[i+1])
		if err != nil {
			return nil, fmt.Errorf("unable to parse interval: %w", err)
		}
	}

	p := slicePos("-p", args)
	if p != -1 {
		c.peakInterval, err = strconv.Atoi(args[p+1])
		if err != nil {
			return nil, fmt.Errorf("unable to parse interval: %w", err)
		}
	}

	m := slicePos("-m", args)
	if m != -1 {
		c.measurement = args[m+1]
	}

	t := slicePos("-t", args)
	if t != -1 {
		c.tags = args[t+1]
	}

	c.file = args[len(args)-1]
	if c.file == "" || len(args) <= 1 {
		return nil, fmt.Errorf("no file provided")
	}

	if c.tags != "" {
		matched, err := regexp.Match(`^([a-z0-9]+=[a-z0-9]+,?)*$`, []byte(c.tags))
		if err != nil {
			return nil, fmt.Errorf("error checking tags format: %w", err)
		}
		if !matched {
			return nil, fmt.Errorf("invalid tags format")
		}

		c.tags = "file=" + c.file + "," + c.tags
	} else {
		c.tags = "file=" + c.file
	}

	return &c, nil
}

func usage() {
	fmt.Println("Usage: telegraf-tailcount [options] file")
	fmt.Println("Options:")
	fmt.Println("  -i int     Interval in seconds (default 10)")
	fmt.Println("  -p int     Peak interval in seconds (default 1)")
	fmt.Println("  -m string  Measurement name (default 'tailcount')")
	fmt.Println("  -t string  Comma-separated k=p pairs for tags (default 'file=<file>')")
}

// slicePos returns the position of the string in the slice, or -1 if not found
func slicePos(a string, list []string) int {
	for i, b := range list {
		if b == a {
			return i
		}
	}
	return -1
}
