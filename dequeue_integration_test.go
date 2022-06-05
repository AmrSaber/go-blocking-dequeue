package blocking_dequeue

import (
	"sync"
	"testing"
)

// Tests that when inserting items concurrently, each item is inserted once and only once.
func TestSyncedPushes(t *testing.T) {
	values := make([]int, 0)
	for i := 1; i <= 1000; i++ {
		values = append(values, i)
	}

	results := make([]int, 0, len(values))

	dequeue := NewBlockingDequeue(make([]int, 5))
	wg := sync.WaitGroup{}

	// Consume all values that are inserted into the dequeue concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 2*len(values); i++ {
			results = append(results, dequeue.PopBack())
		}
	}()

	// Insert the values concurrently
	for _, value := range values {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			dequeue.PushBack(v)
			dequeue.PushFront(v)
		}(value)
	}

	wg.Wait()

	// Make sure that all the number from values are in results
	times := make(map[int]int)
	for _, value := range results {
		times[value]++
	}

	for value, count := range times {
		if count != 2 {
			t.Errorf("Expected %d to be in results twice, got %d", value, count)
		}
	}
}

// Test that when reading items concurrently, each item is read once and only once.
func TestSyncedPops(t *testing.T) {
	values := []int{}
	for i := 1; i <= 1000; i++ {
		values = append(values, i)
	}
	results := make([]int, 0, len(values))

	dequeue := NewBlockingDequeue(make([]int, 2000))
	wg := sync.WaitGroup{}
	resultLock := sync.Mutex{}

	// Insert the values
	for _, value := range values {
		dequeue.PushBack(value)
		dequeue.PushBack(value)
	}

	// Consume all values that are popped from the dequeue concurrently
	for i := 0; i < len(values); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			v1 := dequeue.PopFront()
			v2 := dequeue.PopBack()

			// This only locks the results slice so PopFront and PopBack are being tested correctly
			resultLock.Lock()
			defer resultLock.Unlock()

			results = append(results, v1, v2)
		}()
	}

	wg.Wait()

	// Make sure that all the number from values are in results
	times := make(map[int]int)
	for _, value := range results {
		times[value]++
	}

	for value, count := range times {
		if count != 2 {
			t.Errorf("Expected %d to be in results twice, got %d", value, count)
		}
	}
}

// Test that when reading and writing items at the same time, no value is lost. And that a small buffer is sufficient.
func TestSyncedMixedWrites(t *testing.T) {
	values := []int{}
	for i := 1; i <= 1000; i++ {
		values = append(values, i)
	}
	results := make([]int, 0, len(values))
	resultLock := sync.Mutex{}

	dequeue := NewBlockingDequeue(make([]int, 10))
	wg := sync.WaitGroup{}

	// Concurrent producers
	for i := 0; i < len(values); i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			dequeue.PushBack(v)
			dequeue.PushFront(v)
		}(values[i])
	}

	// Concurrent consumers
	for i := 0; i < len(values); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			v1 := dequeue.PopFront()
			v2 := dequeue.PopBack()

			// This only locks the results slice so PopFront and PopBack are being tested correctly
			resultLock.Lock()
			defer resultLock.Unlock()

			results = append(results, v1, v2)
		}()
	}

	wg.Wait()

	// Make sure that all the number from values are in results
	times := make(map[int]int)
	for _, value := range results {
		times[value]++
	}

	for value, count := range times {
		if count != 2 {
			t.Errorf("Expected %d to be in results twice, got %d", value, count)
		}
	}
}
