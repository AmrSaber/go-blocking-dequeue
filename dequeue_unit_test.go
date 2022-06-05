package blocking_dequeue

import (
	"sync"
	"testing"
	"time"
)

// One rule for testing: Only what's been tested already can be used as a part of a later test

func TestSize(t *testing.T) {
	buffer := make([]int, 5)
	dequeue := NewBlockingDequeue(buffer)

	dequeue.isEmpty = false

	// * * * . . => size = 3
	dequeue.first = 0
	dequeue.last = 2

	if dequeue.Size() != 3 {
		t.Errorf("Expected 3, got %d", dequeue.Size())
	}

	// * . * * * => size = 4
	dequeue.first = 2
	dequeue.last = 0

	if dequeue.Size() != 4 {
		t.Errorf("Expected 4, got %d", dequeue.Size())
	}

	// . * * * . => size = 3
	dequeue.first = 1
	dequeue.last = 3

	if dequeue.Size() != 3 {
		t.Errorf("Expected 3, got %d", dequeue.Size())
	}

	// * . . . * => size = 2
	dequeue.first = 4
	dequeue.last = 0

	if dequeue.Size() != 2 {
		t.Errorf("Expected 2, got %d", dequeue.Size())
	}

}

func TestIsEmpty(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))

	if !dequeue.IsEmpty() {
		t.Errorf("Expected true, got %t", dequeue.IsEmpty())
	}
}

func TestPushFront(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))
	dequeue.PushFront(1)
	dequeue.PushFront(2)
	dequeue.PushFront(3)

	if dequeue.buffer[dequeue.first] != 3 {
		t.Errorf("Expected 3, got %d", dequeue.buffer[0])
	}

	if dequeue.buffer[dequeue.last] != 1 {
		t.Errorf("Expected 1, got %d", dequeue.buffer[0])
	}
}

func TestIsFull(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))

	for i := 0; i < 5; i++ {
		dequeue.PushBack(i)
	}

	if !dequeue.IsFull() {
		t.Errorf("Expected true, got %t", dequeue.IsFull())
	}
}

func TestBlockingPushFront(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 2))
	dequeue.PushFront(1)
	dequeue.PushFront(2)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		dequeue.PushFront(3)
	}()

	// Sleep to make sure that the above goroutine is started and blocked
	time.Sleep(100 * time.Millisecond)

	if dequeue.buffer[dequeue.first] != 2 {
		t.Errorf("Expected 2, got %d", dequeue.buffer[dequeue.first])
	}

	// Remove the element to empty the dequeue and unblock the goroutine
	dequeue.lock.Lock()
	dequeue.first = dequeue.nextIndex(dequeue.first)
	dequeue.notFull.Signal()
	dequeue.lock.Unlock()
	wg.Wait()

	if dequeue.buffer[dequeue.first] != 3 {
		t.Errorf("Expected 3, got %d", dequeue.buffer[dequeue.first])
	}
}

func TestPushBack(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))
	dequeue.PushBack(1)
	dequeue.PushBack(2)
	dequeue.PushBack(3)

	if dequeue.buffer[dequeue.first] != 1 {
		t.Errorf("Expected 1, got %d", dequeue.buffer[dequeue.first])
	}

	if dequeue.buffer[dequeue.last] != 3 {
		t.Errorf("Expected 3, got %d", dequeue.buffer[dequeue.last])
	}
}

func TestBlockingPushBack(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 2))
	dequeue.PushBack(1)
	dequeue.PushBack(2)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		dequeue.PushBack(3)
	}()

	// Sleep to make sure that the above goroutine is started and blocked
	time.Sleep(100 * time.Millisecond)

	if dequeue.buffer[dequeue.last] != 2 {
		t.Errorf("Expected 2, got %d", dequeue.buffer[dequeue.last])
	}

	// Remove the element to empty the dequeue and unblock the goroutine
	dequeue.lock.Lock()
	dequeue.last = dequeue.prevIndex(dequeue.last)
	dequeue.notFull.Signal()
	dequeue.lock.Unlock()
	wg.Wait()

	if dequeue.buffer[dequeue.last] != 3 {
		t.Errorf("Expected 3, got %d", dequeue.buffer[dequeue.last])
	}
}

func TestPopFront(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))
	dequeue.PushFront(1)
	dequeue.PushFront(2)
	dequeue.PushFront(3)

	value := dequeue.PopFront()
	if value != 3 {
		t.Errorf("Expected 3, got %d", value)
	}

	value = dequeue.PopFront()
	if value != 2 {
		t.Errorf("Expected 2, got %d", value)
	}

	value = dequeue.PopFront()
	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}
}

func TestBlockingPopFront(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))

	value := -1
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		value = dequeue.PopFront()
	}()

	// Sleep to make sure that the above goroutine is started and blocked
	time.Sleep(100 * time.Millisecond)

	if value != -1 {
		t.Errorf("Expected -1, got %d", value)
	}

	// Add element to empty the dequeue and unblock the goroutine
	dequeue.PushFront(1)
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}
}

func TestPopBack(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))
	dequeue.PushBack(1)
	dequeue.PushBack(2)
	dequeue.PushBack(3)

	value := dequeue.PopBack()
	if value != 3 {
		t.Errorf("Expected 3, got %d", value)
	}

	value = dequeue.PopBack()
	if value != 2 {
		t.Errorf("Expected 2, got %d", value)
	}

	value = dequeue.PopBack()
	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}
}

func TestBlockingPopBack(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))

	value := -1
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		value = dequeue.PopBack()
	}()

	// Sleep to make sure that the above goroutine is started and blocked
	time.Sleep(100 * time.Millisecond)

	if value != -1 {
		t.Errorf("Expected -1, got %d", value)
	}

	// Add element to empty the dequeue and unblock the goroutine
	dequeue.PushFront(1)
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}
}

func TestPeekFront(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))
	dequeue.PushFront(1)
	dequeue.PushFront(2)
	dequeue.PushFront(3)

	value := dequeue.PeekFront()
	if value != 3 {
		t.Errorf("Expected 3, got %d", value)
	}

	value = dequeue.PeekFront()
	if value != 3 {
		t.Errorf("Expected 3, got %d", value)
	}
}

func TestBlockingPeekFront(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))

	value := -1
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		value = dequeue.PeekFront()
	}()

	// Sleep to make sure that the above goroutine is started and blocked
	time.Sleep(100 * time.Millisecond)

	if value != -1 {
		t.Errorf("Expected -1, got %d", value)
	}

	// Add element to empty the dequeue and unblock the goroutine
	dequeue.PushFront(1)
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}

	if dequeue.Size() != 1 {
		t.Errorf("Expected size 1, got %d", dequeue.Size())
	}

	if dequeue.buffer[dequeue.first] != 1 {
		t.Errorf("Expected 1, got %d", dequeue.buffer[dequeue.first])
	}
}

func TestPeekBack(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))
	dequeue.PushBack(1)
	dequeue.PushBack(2)
	dequeue.PushBack(3)

	value := dequeue.PeekBack()
	if value != 3 {
		t.Errorf("Expected 3, got %d", value)
	}

	value = dequeue.PeekBack()
	if value != 3 {
		t.Errorf("Expected 3, got %d", value)
	}

	if dequeue.Size() != 3 {
		t.Errorf("Expected size 3, got %d", dequeue.Size())
	}
}

func TestBlockingPeekBack(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))

	value := -1
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		value = dequeue.PeekBack()
	}()

	// Sleep to make sure that the above goroutine is started and blocked
	time.Sleep(100 * time.Millisecond)

	if value != -1 {
		t.Errorf("Expected -1, got %d", value)
	}

	// Add element to empty the dequeue and unblock the goroutine
	dequeue.PushFront(1)
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}

	if dequeue.Size() != 1 {
		t.Errorf("Expected size 1, got %d", dequeue.Size())
	}

	if dequeue.buffer[dequeue.last] != 1 {
		t.Errorf("Expected 1, got %d", dequeue.buffer[dequeue.last])
	}
}

func TestIsEmptyAfterUpdates(t *testing.T) {
	dequeue := NewBlockingDequeue(make([]int, 5))

	dequeue.PushFront(1)

	if dequeue.IsEmpty() {
		t.Errorf("Expected false, got true")
	}

	dequeue.PopFront()

	if !dequeue.IsEmpty() {
		t.Errorf("Expected true, got false")
	}
}
