package structures

import (
	"context"
	"errors"
	"log"
	"time"
)

const (
	acceptableWastePercent = 25
	testSize               = 20
	checkRate              = time.Millisecond * 100
)

type TTLMonitor[T any] struct {
	store  *Store[T]
	cancel context.CancelFunc
}

func NewTTLMonitor[T any](store *Store[T]) *TTLMonitor[T] {
	return &TTLMonitor[T]{
		store: store,
	}
}

func (t *TTLMonitor[T]) StartMonitoring(ctx context.Context) error {
	if t.cancel != nil {
		return errors.New("monitoring has already been started")
	}

	c, cancel := context.WithCancel(ctx)
	t.cancel = cancel

	go t.activelyExpireCaches(c)

	return nil
}

func (t *TTLMonitor[T]) StopMonitoring() error {
	if t.cancel != nil {
		t.cancel()
		t.cancel = nil
		return nil
	}

	return errors.New("no running monitor")
}

func (t *TTLMonitor[T]) activelyExpireCaches(ctx context.Context) {
Loop:
	for {
		select {
		case <-ctx.Done():
			break Loop
		default:
			t.expireKeys()
		}

		time.Sleep(checkRate)
	}

	log.Println("no longer monitoring")
}

func (t *TTLMonitor[T]) expireKeys() {
	for {
		var values []keyVal[T]
		if testSize >= t.store.Length() {
			values = t.store.ToSlice()
		} else {
			values = t.store.GetRandomValues(testSize)
		}

		expired := 0
		for _, value := range values {
			if value.Expiration.Before(time.Now()) {
				t.store.Remove(value.Key)
				expired++
			}
		}

		if (float32(expired)/float32(testSize))*100 > acceptableWastePercent {
			continue
		} else {
			return
		}
	}
}
