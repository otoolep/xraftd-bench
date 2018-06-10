package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var addr string
var numReqs int
var modPrint int

const name = `xraftd-bench`
const desc = `xraftd-bench is a simple load testing utility for Raft-based KV stores.`

func init() {
	flag.StringVar(&addr, "a", "localhost:4001", "Node address")
	flag.IntVar(&numReqs, "n", 1000, "Number of requests")
	flag.IntVar(&modPrint, "m", 0, "Print progress every m requests")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s\n\n", desc)
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n", name)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	tester := NewHTTPTester(addr)
	if err := tester.Prepare(map[string]string{"foo": "bar"}); err != nil {
		fmt.Println("failed to prepare test:", err.Error())
		os.Exit(1)
	}

	d, err := run(tester, numReqs)
	if err != nil {
		fmt.Println("failed to run test:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Total duration:", d)
	fmt.Printf("Requests/sec: %.2f\n", float64((numReqs))/d.Seconds())
}

// Tester is the interface test executors must implement.
type Tester interface {
	Prepare(kv map[string]string) error
	Once() (time.Duration, error)
}

func run(t Tester, n int) (time.Duration, error) {
	var dur time.Duration

	for i := 0; i < n; i++ {
		d, err := t.Once()
		if err != nil {
			return 0, err
		}
		dur += d

		if modPrint != 0 && i != 0 && i%modPrint == 0 {
			fmt.Printf("%d requests completed in %s\n", i, d)
		}
	}
	return dur, nil
}
