package orchestrator

import (
	"fmt"
	"sync"
	"time"

	"github.com/zewolfe/hermes/internal/log"
)

type Orchestrator struct {
	queueMux         sync.RWMutex
	store            store
	log              *log.Logger
	evictionInterval time.Duration
}

func New() *Orchestrator {
	o := &Orchestrator{
		queueMux:         sync.RWMutex{},
		store:            newStore(),
		log:              log.NewStdoutLogger(),
		evictionInterval: time.Minute * 5, //TODO: pass this as an option
	}

	o.startGarbageCollector()

	return o
}

func (o *Orchestrator) Subscribe(id string) <-chan interface{} {
	o.queueMux.Lock()
	defer o.queueMux.Unlock()

	ch := make(chan interface{}, 1)

	o.log.Info("Subscriber added", "subscriber", id)
	o.store.Add(id, ch)

	return ch
}

func (o *Orchestrator) Publish(id string, msg interface{}) {
	go func() {
		o.queueMux.RLock()
		defer o.queueMux.RUnlock()

		sub, err := o.store.Get(id)
		if err != nil {
			//TODO: Handle messages that aren't subscribed for
			//TODO: handle the error
			fmt.Println(err)
		}

		sub.ch <- msg
	}()
}

func (o *Orchestrator) startGarbageCollector() {
	ticker := time.NewTicker(o.evictionInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				o.log.Info("Cleaning up store...")
				o.store.CleanUp()
			}
		}
	}()
}
