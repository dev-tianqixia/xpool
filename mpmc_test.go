package xpool

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestBasicOperation(t *testing.T) {
	q := NewQueue()

	q.Push(1)
	q.Push(2)
	q.Push(3)
	if q.length != 3 {
		t.Fail()
	}

	if q.Pop() != 1 {
		t.Fail()
	}
	if q.length != 2 {
		t.Fail()
	}
	if q.Pop() != 2 {
		t.Fail()
	}
	if q.length != 1 {
		t.Fail()
	}
	if q.Pop() != 3 {
		t.Fail()
	}
	if q.length != 0 {
		t.Fail()
	}
	if q.Pop() != nil {
		t.Fail()
	}
	if q.length != 0 {
		t.Fail()
	}
}

func TestConcurrentOperation(t *testing.T) {
	var (
		q       = NewQueue()
		done    = make(chan struct{})
		counter uint32
	)
	for i := 0; i < 100; i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					if q.Pop() != nil {
						atomic.AddUint32(&counter, 1)
					}
				}
			}
		}()
	}

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			q.Push(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(done)

	if counter != 10000 {
		t.Fail()
	}
}
