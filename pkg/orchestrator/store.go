package orchestrator

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// TODO: Think about using generics?
// TODO: error handling
type item struct {
	ch         chan interface{}
	expiration time.Time
}

type lockedMap struct {
	sync.RWMutex
	data map[string]item
	ttl  time.Duration
}

type store interface {
	Add(key string, value chan interface{})
	Get(key string) (item, error)
	Delete(key string) error
	CleanUp()
}

func newStore() store {
	return &lockedMap{
		data: make(map[string]item),
		ttl:  time.Minute * 5, //TODO: Pass this as an option
	}
}

func (l *lockedMap) Add(key string, value chan interface{}) {
	l.Lock()
	defer l.Unlock()

	now := time.Now()
	expDate := now.Add(l.ttl)

	l.data[key] = item{
		ch:         value,
		expiration: expDate,
	}
}

func (l *lockedMap) Get(key string) (item, error) {
	l.RLock()
	defer l.RUnlock()

	i, ok := l.data[key]

	if !ok {
		return item{}, errors.New(fmt.Sprintf("Item not found for key: %s", key))
	}

	hasExpired := time.Now().After(i.expiration)
	if hasExpired {
		return item{}, errors.New(fmt.Sprintf("Item has expired for key: %s", key))
	}

	return i, nil
}

func (l *lockedMap) Delete(key string) error {
	l.Lock()
	defer l.Unlock()

	_, ok := l.data[key]
	if !ok {
		return errors.New(fmt.Sprintf("Item not found for key: %s", key))
	}

	delete(l.data, key)

	return nil
}

// Deletes expired items
func (l *lockedMap) CleanUp() {
	l.Lock()
	defer l.Unlock()

	for key, item := range l.data {
		hasExpired := time.Now().After(item.expiration)

		if hasExpired {
			delete(l.data, key)
		}
	}
}
