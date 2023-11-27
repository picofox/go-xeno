//go:build !windows
// +build !windows

package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

func setNumLoops(numLoops int) error {
	return pollmanager.SetNumLoops(numLoops)
}

func setLoadBalance(lb LoadBalance) error {
	return pollmanager.SetLoadBalance(lb)
}

func setLoggerOutput(w io.Writer) {
	logger = log.New(w, "", log.LstdFlags)
}

// manage all pollers
var pollmanager *manager
var logger *log.Logger

func init() {
	var loops = runtime.GOMAXPROCS(0)/20 + 1
	pollmanager = &manager{}
	pollmanager.SetLoadBalance(RoundRobin)
	pollmanager.SetNumLoops(loops)

	setLoggerOutput(os.Stderr)
}

// LoadBalance is used to do load balancing among multiple pollers.
// a single poller may not be optimal if the number of cores is large (40C+).
type manager struct {
	NumLoops int
	balance  ILoadBalance // load balancing method
	polls    []IPoll      // all the polls
}

// SetNumLoops will return error when set numLoops < 1
func (ego *manager) SetNumLoops(numLoops int) error {
	if numLoops < 1 {
		return fmt.Errorf("set invalid numLoops[%d]", numLoops)
	}

	if numLoops < ego.NumLoops {
		// if less than, close the redundant pollers
		var polls = make([]IPoll, numLoops)
		for idx := 0; idx < ego.NumLoops; idx++ {
			if idx < numLoops {
				polls[idx] = ego.polls[idx]
			} else {
				if err := ego.polls[idx].Close(); err != nil {
					logger.Printf("NETPOLL: poller close failed: %v\n", err)
				}
			}
		}
		ego.NumLoops = numLoops
		ego.polls = polls
		ego.balance.Rebalance(ego.polls)
		return nil
	}

	ego.NumLoops = numLoops
	return ego.Run()
}

// SetLoadBalance set load balance.
func (ego *manager) SetLoadBalance(lb LoadBalance) error {
	if ego.balance != nil && ego.balance.LoadBalance() == lb {
		return nil
	}
	ego.balance = newLoadbalance(lb, ego.polls)
	return nil
}

// Close release all resources.
func (ego *manager) Close() error {
	for _, poll := range ego.polls {
		poll.Close()
	}
	ego.NumLoops = 0
	ego.balance = nil
	ego.polls = nil
	return nil
}

// Run all pollers.
func (ego *manager) Run() (err error) {
	defer func() {
		if err != nil {
			_ = ego.Close()
		}
	}()

	// new poll to fill delta.
	for idx := len(ego.polls); idx < ego.NumLoops; idx++ {
		var poll IPoll
		poll, err = openPoll()
		if err != nil {
			return
		}
		ego.polls = append(ego.polls, poll)
		go poll.Wait()
	}

	// LoadBalance must be set before calling Run, otherwise it will panic.
	ego.balance.Rebalance(ego.polls)
	return nil
}

// Reset pollers, this operation is very dangerous, please make sure to do this when calling !
func (ego *manager) Reset() error {
	for _, poll := range ego.polls {
		poll.Close()
	}
	ego.polls = nil
	return ego.Run()
}

// Pick will select the poller for use each time based on the LoadBalance.
func (ego *manager) Pick() IPoll {
	return ego.balance.Pick()
}
