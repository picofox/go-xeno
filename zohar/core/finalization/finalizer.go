package finalization

import (
	"sync"
)

var globalFinalizer Finalizer
var gfOnce sync.Once

func GetGlobalFinalizer() *Finalizer {
	gfOnce.Do(func() {
		globalFinalizer = Finalizer{
			items: make([]FinalizerItem, 0),
		}
	})
	return &globalFinalizer
}

type FinalizerItem struct {
	Name    string
	Subject any
	Handler func(any)
}

func NeoFinalizer() *Finalizer {
	return &Finalizer{
		items: make([]FinalizerItem, 0),
	}
}

type Finalizer struct {
	items []FinalizerItem
}

func (ego *Finalizer) Register(name string, sub any, m func(any)) {
	item := FinalizerItem{
		Name:    name,
		Subject: sub,
		Handler: m,
	}
	ego.items = append(ego.items, item)
}

func (ego *Finalizer) Finalize() {
	for i := len(ego.items) - 1; i >= 0; i-- {
		ego.items[i].Handler(ego.items[i].Subject)
	}
}
