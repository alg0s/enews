package queue

import "testing"

func TestStandardQueue(t *testing.T) {
	var q = NewQueue()

	// 1. test empty queue
	if q.First() != nil {
		t.Errorf(`Expect nil`)
	}

	if q.IsEmpty() == false {
		t.Errorf(`Expect true`)
	}

	if q.Dequeue() != nil {
		t.Errorf(`Expect nil`)
	}

	// 2. test non-empty queue
	q.Enqueue(1)

	if q.First() == nil {
		t.Errorf(`Expect 1`)
	}

	if q.Last() != 1 {
		t.Error(`Expect 1`)
	}

	if q.Length() != 1 {
		t.Errorf(`Expect Length = 1`)
	}

	if q.IsEmpty() == true {
		t.Errorf(`Expect false`)
	}

	var d = q.Dequeue()

	if d != 1 {
		t.Errorf(`Expect 1`)
	}

	if q.IsEmpty() != true {
		t.Errorf(`Expect true`)
	}
}
