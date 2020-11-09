package xpool

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	next *node
	item interface{}
}

// a lock free
type queue struct {
	aux      *node // aux points to the tail of the queue
	sentinel *node // sentinel.next is the head of the queue
	length   uint64
}

func NewQueue() *queue {
	dummy := &node{}
	return &queue{
		aux:      dummy,
		sentinel: dummy,
		length:   0,
	}
}

func (q *queue) Push(item interface{}) {
	n := &node{
		item: item,
	}

	prevAux := (*node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.aux)), unsafe.Pointer(n)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prevAux.next)), unsafe.Pointer(n))
	atomic.AddUint64(&q.length, 1)
}

func (q *queue) Pop() interface{} {
	for {
		currentSentinel := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.sentinel))))
		head := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&currentSentinel.next))))
		if head == nil {
			return nil
		}
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.sentinel)), unsafe.Pointer(currentSentinel), unsafe.Pointer(head)) {
			atomic.AddUint64(&q.length, ^uint64(0))
			ret := head.item
			head.item = nil
			return ret
		}
	}
}

func (q *queue) Len() uint64 {
	return atomic.LoadUint64(&q.length)
}

func (q *queue) IsEmpty() bool {
	currentSentinel := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.sentinel))))
	head := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&currentSentinel.next))))
	return head == nil
}
