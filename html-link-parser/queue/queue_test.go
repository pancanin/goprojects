package queue_test

import (
	"testing"
	"example.com/link-spider/queue"
)

func TestEnqueueOneElement(t *testing.T) {
	q := make([]string, 0)
	q = queue.Enqueue[string](q, "first")

	if len(q) != 1 {
		t.Fatalf("expected queue size to be %v, but was %v", 1, len(q))
	}
}

func TestEnqueueTwoElements(t *testing.T) {
	q := make([]string, 0)
	q = queue.Enqueue[string](q, "first")
	q = queue.Enqueue[string](q, "second")

	if len(q) != 2 {
		t.Fatalf("expected queue size to be %v, but was %v", 2, len(q))
	}
}

func TestTopWithOneElement(t *testing.T) {
	q := make([]string, 0)
	firstVal := "first"
	q = queue.Enqueue[string](q, firstVal)

	firstInQ := queue.Top[string](q)

	if firstInQ != firstVal {
		t.Fatalf("expected %v, but got %v", firstVal, firstInQ)
	}
}

func TestTopWithTwoElements(t *testing.T) {
	q := make([]string, 0)
	firstVal := "first"
	q = queue.Enqueue[string](q, firstVal)
	q = queue.Enqueue[string](q, "second")

	firstInQ := queue.Top[string](q)

	if firstInQ != firstVal {
		t.Fatalf("expected %v, but got %v", firstVal, firstInQ)
	}
}

func TestPopQueueSize(t *testing.T) {
	q := make([]string, 0)
	q = queue.Enqueue[string](q, "first")
	q = queue.Enqueue[string](q, "second")

	q = queue.Pop[string](q)

	if len(q) != 1 {
		t.Fatalf("Queue size did not shrink after pop op. Expected %v, got %v", 1, len(q))
	}
}

func TestPopFirstInQueue(t *testing.T) {
	q := make([]string, 0)
	q = queue.Enqueue[string](q, "first")
	secondVal := "second"
	q = queue.Enqueue[string](q, secondVal)

	q = queue.Pop[string](q)

	// After the pop we expect the 'second' value to be first in line on the queue
	firstInQ := queue.Top[string](q)

	if firstInQ != secondVal {
		t.Fatalf("Expected first value in queue: %v, but got %v", secondVal, firstInQ)
	}
}

// func TestLoadQueue(t *testing.T) {
// 	values := make([]string)
// }