package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	flagAddr       = flag.String("addr", "", "addr to scan")
	flagPortMin    = flag.Int("port-min", 1, "min port to scan")
	flagPortMax    = flag.Int("port-max", 65535, "max port to scan")
	flagTimeout    = flag.Int("timeout", 1000, "timeout in millseconds")
	flagVerbose    = flag.Bool("v", false, "output details")
	flagConcurrent = flag.Int("c", 100, "concurrent")
)

func fatal(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}

func verbose(s string) {
	if *flagVerbose {
		fmt.Println("[VERBOSE]", s)
	}
}

func main() {
	flag.Parse()
	addr := *flagAddr
	if addr == "" {
		fatal("invalid addr")
	}
	timeout := time.Millisecond * time.Duration(*flagTimeout)
	min, max := *flagPortMin, *flagPortMax
	if max < min {
		fatal("invalid port range")
	}

	var wg sync.WaitGroup
	wg.Add(max - min + 1)
	sem := make(chan struct{}, *flagConcurrent)
	var openPorts []int
	for port := min; port <= max; port++ {
		sem <- struct{}{}
		go func(port int) {
			defer wg.Done()
			defer func() { <-sem }()

			verbose(fmt.Sprintf("scan %s:%d", addr, port))
			err := dial(addr, port, timeout)
			if err != nil {
				verbose(fmt.Sprintf("scan %s:%d failed: %v", addr, port, err))
			} else {
				openPorts = append(openPorts, port)
			}
		}(port)
	}
	wg.Wait()

	for _, port := range openPorts {
		fmt.Printf("%6d is open\n", port)
	}
}

func dial(addr string, port int, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(addr, strconv.Itoa(port)), timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
