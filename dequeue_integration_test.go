package blocking_dequeue

import (
	"sync"
	"testing"
	"time"
)

// Tests that when inserting items concurrently, each item is inserted once and only once.
func TestSyncedPushes(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	results := make([]int, 0, 10)

	dequeue := NewBlockingDequeue[int]()
	wg := sync.WaitGroup{}

	// Consume all values that are inserted into the dequeue concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(values); i++ {
			results = append(results, dequeue.PopBack())
		}
	}()

	// Insert the values concurrently
	for _, value := range values {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			dequeue.PushBack(v)
		}(value)
	}

	wg.Wait()

	// Make sure that all the number from values are in results
	for _, value := range values {
		found := false

		for _, result := range results {
			if result == value {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected %d to be in results", value)
		}
	}
}

// Test that when reading items concurrently, each item is read once and only once.
func TestSyncedPops(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	results := make([]int, 0, 10)

	dequeue := NewBlockingDequeue[int]()
	wg := sync.WaitGroup{}
	resultLock := sync.Mutex{}

	// Insert the values
	for _, value := range values {
		dequeue.PushBack(value)
	}

	// Consume all values that are popped from the dequeue concurrently
	for i := 0; i < len(values); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			value := dequeue.PopFront()

			// This only locks the results slice so PopFront is being tested correctly
			resultLock.Lock()
			defer resultLock.Unlock()

			results = append(results, value)
		}()
	}

	wg.Wait()

	// Make sure that all the number from values are in results
	for _, value := range values {
		found := false

		for _, result := range results {
			if result == value {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected %d to be in results", value)
		}
	}
}

// Test that when capacity is increased, blocking producers are freed up.
func TestCapacityChange(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	dequeue := NewBlockingDequeue[int]()
	dequeue.SetCapacity(5)
	wg := sync.WaitGroup{}

	// Insert the values concurrently
	for _, value := range values {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			dequeue.PushBack(v)
		}(value)
	}

	// Sleep to allow all the goroutines to execute
	time.Sleep(100 * time.Millisecond)

	// Make sure that only 5 items are in the dequeue
	if dequeue.Size() != 5 {
		t.Errorf("Expected dequeue to have 5 items, got %d", dequeue.Size())
	}

	// Update the capacity to 10
	err := dequeue.SetCapacity(0)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	wg.Wait()

	// Make sure that all the number from values are in the dequeue
	if dequeue.Size() != len(values) {
		t.Errorf("Expected dequeue to have 10 items, got %d", dequeue.Size())
	}
}
