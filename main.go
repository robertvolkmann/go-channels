package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type message struct {
	id   int
	time time.Time
}

type set map[int]bool

func main() {
	n, err := parseArgs(os.Args)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\nExactly one positive numeric argument is expected.\n", err)
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	messages := make(chan message, n)
	done := make(chan struct{})

	startReceiver(wg, messages, done, n)

	for i := 1; i <= n; i++ {
		startSender(wg, messages, done, i)
	}
	wg.Wait()
}

func parseArgs(args []string) (int, error) {
	if len(args) != 2 {
		return 0, errors.New("missing argument")
	}

	n, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, err
	}
	if n <= 0 {
		return 0, errors.New("non-positive integer")
	}
	return n, nil
}

func startReceiver(wg *sync.WaitGroup, messages <-chan message, done chan<- struct{}, numSenders int) {
	wg.Add(1)
	go func() {
		defer close(done)
		defer wg.Done()

		received := make(set)

		for len(received) != numSenders {
			m := <-messages
			received[m.id] = true
			fmt.Printf("Received message from send %d at %s\n", m.id, m.time.Format(time.StampMilli))
		}
	}()
}

func startSender(wg *sync.WaitGroup, messages chan<- message, done <-chan struct{}, i int) {
	wg.Add(1)
	go func(id int, d time.Duration) {
		defer wg.Done()

		for {
			select {
			case <-done:
				return
			case <-time.After(d):
				messages <- message{id, time.Now()}
			}
		}
	}(i, randomDuration())
}

func randomDuration() time.Duration {
	sec := rand.Intn(9) + 1
	return time.Duration(sec) * time.Second
}
