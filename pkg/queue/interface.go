package queue

// Interface defines required components of a queue
type Interface interface {
	Enqueue(interface{})
	Dequeue() interface{}
	IsEmpty() bool
	Length() int
	Peek() interface{}
	Lock()
	Unlock()
	IsLocked() bool
}
