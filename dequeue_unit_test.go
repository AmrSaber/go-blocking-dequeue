package blocking_dequeue

import (
	"sync"
	"testing"
	"time"
)

// One rule for testing: Only what's been tested already can be used as a part of a later test

func TestSize(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()

	// Access the list directly, as PushFront() is not tested yet
	dequeue.list.PushFront(1)
	dequeue.list.PushFront(2)
	dequeue.list.PushFront(3)

	if dequeue.Size() != 3 {
		t.Errorf("Expected 3, got %d", dequeue.Size())
	}
}

func TestCapacity(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()

	// Test that Capacity() returns the actual capacity
	dequeue.capacity = 2
	if dequeue.Capacity() != 2 {
		t.Errorf("Expected 3, got %d", dequeue.Capacity())
	}

	// Test that SetCapacity() sets the capacity
	dequeue.SetCapacity(4)
	if dequeue.capacity != 4 {
		t.Errorf("Expected 4, got %d", dequeue.Capacity())
	}
}

func TestIsEmpty(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()

	if !dequeue.IsEmpty() {
		t.Errorf("Expected true, got %t", dequeue.IsEmpty())
	}

	dequeue.list.PushFront(1)
	if dequeue.IsEmpty() {
		t.Errorf("Expected false, got %t", dequeue.IsEmpty())
	}
}

func TestIsFull(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()

	dequeue.capacity = 2
	dequeue.list.PushFront(1)
	dequeue.list.PushFront(2)

	if !dequeue.IsFull() {
		t.Errorf("Expected true, got %t", dequeue.IsFull())
	}
}

func TestPushFront(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()
	dequeue.PushFront(1)
	dequeue.PushFront(2)
	dequeue.PushFront(3)

	if dequeue.list.Front().Value.(int) != 3 {
		t.Errorf("Expected 3, got %d", dequeue.list.Front().Value.(int))
	}

	if dequeue.list.Back().Value.(int) != 1 {
		t.Errorf("Expected 1, got %d", dequeue.list.Back().Value.(int))
	}
}

func TestBlockingPushFront(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()
	dequeue.SetCapacity(1)
	dequeue.PushFront(1)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		dequeue.PushFront(2)
	}()

	// Sleep to make sure that the above goroutine is started and blocked
	time.Sleep(100 * time.Millisecond)

	if dequeue.list.Front().Value.(int) != 1 {
		t.Errorf("Expected 1, got %d", dequeue.list.Front().Value.(int))
	}

	// Remove the element to empty the dequeue and unblock the goroutine
	dequeue.writeCond.L.Lock()
	dequeue.list.Remove(dequeue.list.Front())
	dequeue.writeCond.Signal()
	dequeue.writeCond.L.Unlock()
	wg.Wait()

	if dequeue.list.Front().Value.(int) != 2 {
		t.Errorf("Expected 2, got %d", dequeue.list.Front().Value.(int))
	}
}

func TestPushBack(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()
	dequeue.PushBack(1)
	dequeue.PushBack(2)
	dequeue.PushBack(3)

	if dequeue.list.Front().Value.(int) != 1 {
		t.Errorf("Expected 1, got %d", dequeue.list.Front().Value.(int))
	}

	if dequeue.list.Back().Value.(int) != 3 {
		t.Errorf("Expected 3, got %d", dequeue.list.Back().Value.(int))
	}
}

func TestBlockingPushBack(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()
	dequeue.SetCapacity(1)
	dequeue.PushBack(1)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		dequeue.PushBack(2)
	}()

	// Sleep to make sure that the above goroutine is started and blocked
	time.Sleep(100 * time.Millisecond)

	if dequeue.list.Back().Value.(int) != 1 {
		t.Errorf("Expected 1, got %d", dequeue.list.Back().Value.(int))
	}

	// Remove the element to empty the dequeue and unblock the goroutine
	dequeue.writeCond.L.Lock()
	dequeue.list.Remove(dequeue.list.Back())
	dequeue.writeCond.Signal()
	dequeue.writeCond.L.Unlock()
	wg.Wait()

	if dequeue.list.Back().Value.(int) != 2 {
		t.Errorf("Expected 2, got %d", dequeue.list.Back().Value.(int))
	}
}

func TestPopFront(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()
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
	dequeue := NewBlockingDequeue[int]()

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
	dequeue.writeCond.L.Lock()
	dequeue.list.PushFront(1)
	dequeue.writeCond.Signal()
	dequeue.writeCond.L.Unlock()
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}
}

func TestPopBack(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()
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
	dequeue := NewBlockingDequeue[int]()

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
	dequeue.writeCond.L.Lock()
	dequeue.list.PushBack(1)
	dequeue.writeCond.Signal()
	dequeue.writeCond.L.Unlock()
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}
}

func TestPeekFront(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()
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

	if dequeue.Size() != 3 {
		t.Errorf("Expected size 3, got %d", dequeue.Size())
	}
}

func TestBlockingPeekFront(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()

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
	dequeue.writeCond.L.Lock()
	dequeue.list.PushFront(1)
	dequeue.writeCond.Signal()
	dequeue.writeCond.L.Unlock()
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}

	if dequeue.Size() != 1 {
		t.Errorf("Expected size 1, got %d", dequeue.Size())
	}

	if dequeue.list.Front().Value.(int) != 1 {
		t.Errorf("Expected 1, got %d", dequeue.list.Front().Value.(int))
	}
}

func TestPeekBack(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()
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
	dequeue := NewBlockingDequeue[int]()

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
	dequeue.writeCond.L.Lock()
	dequeue.list.PushBack(1)
	dequeue.writeCond.Signal()
	dequeue.writeCond.L.Unlock()
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}

	if dequeue.Size() != 1 {
		t.Errorf("Expected size 1, got %d", dequeue.Size())
	}

	if dequeue.list.Back().Value.(int) != 1 {
		t.Errorf("Expected 1, got %d", dequeue.list.Back().Value.(int))
	}
}

func TestOnEmpty(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()

	called := 0
	dequeue.SetOnEmpty(func() {
		called++
	})

	dequeue.PushBack(1)
	dequeue.PopFront()

	if called != 1 {
		t.Errorf("Expected 1, got %d", called)
	}

	dequeue.PushBack(2)
	dequeue.PopFront()

	if called != 2 {
		t.Errorf("Expected 2, got %d", called)
	}
}

func TestOnFull(t *testing.T) {
	dequeue := NewBlockingDequeue[int]()
	dequeue.SetCapacity(3)

	called := 0
	dequeue.SetOnFull(func() {
		called++
	})

	dequeue.PushBack(1)
	dequeue.PushBack(2)
	dequeue.PushBack(3)

	if called != 1 {
		t.Errorf("Expected 1, got %d", called)
	}

	dequeue.PopFront()
	dequeue.PushFront(0)

	if called != 2 {
		t.Errorf("Expected 2, got %d", called)
	}
}
