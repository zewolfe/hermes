package orchestrator

import (
	"fmt"
	"time"
)

type Orchestrator struct {
	store store
	options
}

func New(opts ...Options) *Orchestrator {
	options := applyOptions(opts...)

	o := &Orchestrator{
		store:   newStore(options.ttl),
		options: options,
	}

	o.startGarbageCollector()

	return o
}

func (o *Orchestrator) Subscribe(id string) <-chan interface{} {
	ch := make(chan interface{}, 1)

	o.log.Info("Subscriber added", "subscriber", id)
	o.store.Add(id, ch)

	return ch
}

func (o *Orchestrator) Publish(id string, msg interface{}) {
	go func() {
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
