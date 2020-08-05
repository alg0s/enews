package utils

import (
	"sync"
)

// QueueInterface is an interface that wraps the basic of a queue
// to be implemented by specific queue types.
// A queue follows the principle of FIFO: First In First Out
type QueueInterface interface {
	Enqueue(interface{})
	Dequeue() interface{}
	IsEmpty() bool
	Length() int
	Peek() interface{}
}

// StandardQueue is the basic queue without any fancy features
type StandardQueue struct {
	elements []interface{}
}

// NewStandardQueue initates a new empty queue
func NewStandardQueue() *StandardQueue {
	newList := make([]interface{}, 0, 32)
	newQ := StandardQueue{elements: newList}
	return &newQ
}

// Enqueue puts a new element into the queue
func (q *StandardQueue) Enqueue(newElement interface{}) {
	// add new element
	q.elements = append(q.elements, newElement)
}

// Dequeue removes the first element in the queue and return its address
func (q *StandardQueue) Dequeue() interface{} {
	if q.IsEmpty() {
		return nil
	}

	// get address of the current first in queue
	currentFirst := q.elements[0]

	// reduce size of the queue slice
	q.elements = q.elements[1:]

	return currentFirst
}

// IsEmpty returns true if the queue is empty, else false
func (q *StandardQueue) IsEmpty() bool {
	if q.Length() == 0 {
		return true
	}
	return false
}

// Length returns the total number of elements in the queue
func (q *StandardQueue) Length() int {
	return len(q.elements)
}

// Peek returns the first element in the queue
// without dequeuing it
func (q *StandardQueue) Peek() interface{} {
	if q.IsEmpty() {
		return nil
	}
	return q.elements[0]
}

// BlockingQueue implements a block when Dequeuing
type BlockingQueue struct {
	lock     *sync.Mutex
	notEmpty *sync.Cond
	store    *StandardQueue
}

// NewBlockingQueue constructs a Blocking Queue
func NewBlockingQueue() *BlockingQueue {
	q := BlockingQueue{}
	q.store = NewStandardQueue()
	q.lock = new(sync.Mutex)
	q.notEmpty = sync.NewCond(q.lock)
	return &q
}

// IsEmpty returns true if queue collection is empty
func (q *BlockingQueue) IsEmpty() bool {
	return q.store.IsEmpty()
}

// Length returns the length of the queue collection
func (q *BlockingQueue) Length() int {
	return q.store.Length()
}

// Enqueue adds an object to queue collection
func (q *BlockingQueue) Enqueue(in interface{}) {
	q.lock.Lock()
	q.store.Enqueue(in)
	q.lock.Unlock()

	// signal to listeners that there is a new item in the queue
	q.notEmpty.Signal()
}

// Dequeue removes an object from the queue colleciton
func (q *BlockingQueue) Dequeue() interface{} {
	q.lock.Lock()

	for q.IsEmpty() {
		// wait here until there is an object in the queue
		q.notEmpty.Wait()
	}

	result := q.store.Dequeue()
	q.lock.Unlock()

	return result
}

// Peek returns the first object in the queue collection
func (q *BlockingQueue) Peek() interface{} {
	q.lock.Lock()

	for q.IsEmpty() {
		// wait here until there is an object in the queue
		q.notEmpty.Wait()
	}

	result := q.store.Peek()
	q.lock.Unlock()

	return result
}
