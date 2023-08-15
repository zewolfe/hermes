package orchestrator_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/zewolfe/hermes/internal/log"
	"github.com/zewolfe/hermes/pkg/orchestrator"
)

var noopLogger *log.Logger = log.NewNoopLogger()

type subs struct {
	chs []<-chan interface{}
	sync.RWMutex
}

func newSubs(limit int) *subs {
	return &subs{
		chs: make([]<-chan interface{}, limit),
	}
}

func TestOchestrator(t *testing.T) {
	id := "123456789"
	want := "Hello from the other side"
	o := orchestrator.New(
		orchestrator.WithLogger(noopLogger),
	)

	ch := o.Subscribe(id)

	go func() {
		o.Publish(id, want)
	}()

	got := <-ch

	switch got.(type) {
	case string:
		if got != want {
			t.Fatalf("Expected want %q to be the same as got %q, got=%q", want, got, got)
		}

	default:
		t.Fatalf("Expected got %v to be of type string but got type %T", got, got)
	}
}

func TestOchestratorWithLoad(t *testing.T) {
	limit := 100000
	msg := "Hello from the other side"
	err := make(chan error, limit)

	wg := sync.WaitGroup{}

	o := orchestrator.New(
		orchestrator.WithLogger(noopLogger),
	)

	subs := newSubs(limit)

	//subscribe for events
	for i := 0; i < limit; i++ {
		id := fmt.Sprint(i)

		ch := o.Subscribe(id)
		wg.Add(1)

		subs.Lock()
		subs.chs[i] = ch
		subs.Unlock()
	}

	go func() {
		for i := 0; i < limit; i++ {
			id := fmt.Sprint(i)

			want := fmt.Sprintf("%s: %s", msg, id)

			o.Publish(id, want)
		}
	}()

	go func() {
		//Output nil to the error channel in the end if no errors were encountered so we don't wait forever
		defer func() { err <- nil }()

		for i := 0; i < limit; i++ {
			subs.RLock()
			ch := subs.chs[i]
			subs.RUnlock()

			got := <-ch
			switch got.(type) {
			case string:
				want := fmt.Sprintf("%s: %s", msg, fmt.Sprint(i))

				if got != want {
					err <- errors.New(fmt.Sprintf("Expected want %q to be the same as got %q", want, got))
				}

				t.Logf("Published: %s", got)

			default:
				err <- errors.New(fmt.Sprintf("Expected got %v to be of type string but got type %T", got, got))
			}

			wg.Done()
		}
	}()

	select {
	case e := <-err:
		if e != nil {
			t.Fatalf(e.Error())
		}
	}

	wg.Wait()
}

func TestOrchestratorCleanUp(t *testing.T) {
	id := "369"
	want := "Parting is such sweet sorrow"
	err := make(chan error)
	wg := sync.WaitGroup{}

	o := orchestrator.New(
		orchestrator.WithLogger(noopLogger),
		orchestrator.WithEvictionInterval(time.Second*1),
		orchestrator.WithItemTTL(time.Second/2),
	)

	ch := o.Subscribe(id)
	wg.Add(1)

	go func() {
		o.Publish(id, want)
	}()

	go func() {
		got := <-ch

		switch got.(type) {
		case string:
			if got != want {
				err <- errors.New(fmt.Sprintf("Expected want %q to be the same as got %q, got=%q", want, got, got))
			}

		default:
			err <- errors.New(fmt.Sprintf("Expected got %v to be of type string but got type %T", got, got))
		}

		err <- nil
		wg.Done()
	}()

	select {
	case e := <-err:
		if e != nil {
			t.Fatal(e.Error())
		}
	}

	select {
	case <-time.After(time.Second * 2):
		_, open := <-ch
		if open {
			t.Fatalf("Expected ch for subscriber id %s to be closed after cleanup", id)
		}
	}

	wg.Wait()
}

/*
*
TODO:
Improve with variable subscription times
Variable publish times
Variable consumption through channels
*
*/
func FuzzTestOrchestrator(f *testing.F) {
	f.Add("123456", "It’s not a bug – it’s an undocumented feature.")
	f.Add("00055572", "Don’t worry if it doesn’t work right. If everything did, you’d be out of a job.")

	o := orchestrator.New(
		orchestrator.WithLogger(noopLogger),
		orchestrator.WithEvictionInterval(time.Second*3),
	)

	f.Fuzz(func(t *testing.T, id string, want string) {
		ch := o.Subscribe(id)

		go func() {
			o.Publish(id, want)
		}()

		got := <-ch

		switch got.(type) {
		case string:
			if got != want {
				t.Fatalf("Expected want %q to be the same as got %q, got=%q", want, got, got)
			}

		default:
			t.Fatalf("Expected got %v to be of type string but got type %T", got, got)
		}
	})

}
