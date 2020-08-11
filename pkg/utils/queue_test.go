package utils

import (
	"testing"
)

type container struct {
	name string
}

func (c *container) getName() string {
	return c.name
}

func TestStandardQueue(t *testing.T) {
	q := StandardQueue{}

	if q.Peek() != nil {
		t.Errorf(`Expect nil`)
	}

	if q.IsEmpty() == false {
		t.Errorf(`Expect true`)
	}

	q.Enqueue("first")
	q.Enqueue("second")

	if q.Length() != 2 {
		t.Errorf(`Expect length 2`)
	}

	if q.IsEmpty() == true {
		t.Errorf(`Expect false`)
	}

	if q.Dequeue() != "first" {
		t.Errorf(`Expect false`)
	}

	if q.Length() != 1 {
		t.Errorf(`Expect length 1`)
	}
}

func QueueCanTakeStructShouldSucceed(t *testing.T) {
	q := StandardQueue{}

	q.Enqueue(&container{name: "hello world"})

	if q.Peek().(*container).getName() != "hello world" {
		t.Errorf(`Expected name == 'hello world'`)
	}
}
