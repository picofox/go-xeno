package server

import (
	"sync/atomic"

	"github.com/bytedance/gopkg/lang/fastrand"
)

// LoadBalance sets the load balancing method.
type LoadBalance int

const (
	// Random requests that connections are randomly distributed.
	Random LoadBalance = iota
	// RoundRobin requests that connections are distributed to a Poll
	// in a round-robin fashion.
	RoundRobin
)

// ILoadBalance sets the load balancing method for []*polls
type ILoadBalance interface {
	LoadBalance() LoadBalance
	// Choose the most qualified Poll
	Pick() (poll IPoll)

	Rebalance(polls []IPoll)
}

func newLoadbalance(lb LoadBalance, polls []IPoll) ILoadBalance {
	switch lb {
	case Random:
		return newRandomLB(polls)
	case RoundRobin:
		return newRoundRobinLB(polls)
	}
	return newRoundRobinLB(polls)
}

func newRandomLB(polls []IPoll) ILoadBalance {
	return &randomLB{polls: polls, pollSize: len(polls)}
}

type randomLB struct {
	polls    []IPoll
	pollSize int
}

func (b *randomLB) LoadBalance() LoadBalance {
	return Random
}

func (b *randomLB) Pick() (poll IPoll) {
	idx := fastrand.Intn(b.pollSize)
	return b.polls[idx]
}

func (b *randomLB) Rebalance(polls []IPoll) {
	b.polls, b.pollSize = polls, len(polls)
}

func newRoundRobinLB(polls []IPoll) ILoadBalance {
	return &roundRobinLB{polls: polls, pollSize: len(polls)}
}

type roundRobinLB struct {
	polls    []IPoll
	accepted uintptr // accept counter
	pollSize int
}

func (b *roundRobinLB) LoadBalance() LoadBalance {
	return RoundRobin
}

func (b *roundRobinLB) Pick() (poll IPoll) {
	idx := int(atomic.AddUintptr(&b.accepted, 1)) % b.pollSize
	return b.polls[idx]
}

func (b *roundRobinLB) Rebalance(polls []IPoll) {
	b.polls, b.pollSize = polls, len(polls)
}
