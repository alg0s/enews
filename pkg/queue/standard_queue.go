package queue

/*
	ref:
	https://github.com/enriquebris/goconcurrentqueue/blob/master/fifo_queue.go
*/

import (
	"container/list"
)

type StandardQueue struct {
	elements *list.List
}

// New initializes and returns a new StandardQueue
func NewQueue() *StandardQueue {
	var l = list.New()
	var q = StandardQueue{elements: l}
	return &q
}

func (q *StandardQueue) Enqueue(element interface{}) interface{} {
	var e = q.elements.PushBack(element)
	return e.Value
}

// Dequeue removes and returns the first element in the queue  with FIFO
func (q *StandardQueue) Dequeue() interface{} {
	if q.IsEmpty() == true {
		return nil
	}
	var e = q.elements.Front()
	var v = e.Value
	q.elements.Remove(e)
	return v
}

func (q *StandardQueue) IsEmpty() bool {
	return q.elements.Len() == 0
}

func (q *StandardQueue) Length() int {
	return q.elements.Len()
}

func (q *StandardQueue) First() interface{} {
	if q.elements.Len() > 0 {
		var e = q.elements.Front()
		return e.Value
	}
	return nil
}

func (q *StandardQueue) Last() interface{} {
	if q.elements.Len() > 0 {
		var e = q.elements.Back()
		return e.Value
	}
	return nil
}
