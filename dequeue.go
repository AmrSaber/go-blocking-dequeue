package blocking_dequeue

import (
	"container/list"
	"fmt"
	"sync"
)

// Blocking dequeue, implemented with a linked list.
// By default the dequeue has infinite capacity.
// The dequeue is thread safe. And must not be copied.
type BlockingDequeue[T any] struct {
	list      *list.List
	writeCond *sync.Cond
	capacity  int

	OnFull  func() // Optional callback function invoked when the dequeue is full
	OnEmpty func() // Optional callback function invoked when the dequeue is empty
}

// Creates a new blocking dequeue with infinite capacity.
// The dequeue MUST only be created using this method.
func NewBlockingDequeue[T any]() *BlockingDequeue[T] {
	d := new(BlockingDequeue[T])
	d.list = list.New()
	d.writeCond = sync.NewCond(&sync.Mutex{})
	return d
}

// =================================[Push/Pop/Peek]=================================

// Add an item into the front (top) of the dequeue. Blocks if dequeue is full.
func (d *BlockingDequeue[T]) PushFront(item T) {
	// If the dequeue is full, wait until an item is removed
	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()
	for d.isFull_unsafe() {
		d.writeCond.Wait()
	}

	d.list.PushFront(item)

	// Notify the consumer that an item has been added
	defer d.writeCond.Broadcast()

	// Call the OnFull callback if the dequeue is full
	if d.isFull_unsafe() && d.OnFull != nil {
		d.OnFull()
	}
}

// Add an item to the back (bottom) of the dequeue. Blocks if dequeue is full.
func (d *BlockingDequeue[T]) PushBack(item T) {
	// If the dequeue is full, wait until an item is removed
	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()
	for d.isFull_unsafe() {
		d.writeCond.Wait()
	}

	d.list.PushBack(item)

	// Notify the consumer that an item has been added
	defer d.writeCond.Broadcast()

	// Call the OnFull callback if the dequeue is full
	if d.isFull_unsafe() && d.OnFull != nil {
		d.OnFull()
	}
}

// Read the first item (on the top/front) of the dequeue and remove it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PopFront() T {
	// If the dequeue is empty, wait until an item is added
	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()
	for d.isEmpty_unsafe() {
		d.writeCond.Wait()
	}

	item := d.list.Remove(d.list.Front()).(T)

	// Notify the producer that an item has been removed
	defer d.writeCond.Broadcast()

	// Call the OnEmpty callback if the dequeue is empty
	if d.isEmpty_unsafe() && d.OnEmpty != nil {
		d.OnEmpty()
	}

	return item
}

// Read the last item (at the end/back) of the dequeue and remove it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PopBack() T {
	// If the dequeue is empty, wait until an item is added
	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()
	for d.isEmpty_unsafe() {
		d.writeCond.Wait()
	}

	item := d.list.Remove(d.list.Back()).(T)

	// Notify the producer that an item has been removed
	defer d.writeCond.Broadcast()

	// Call the OnEmpty callback if the dequeue is empty
	if d.isEmpty_unsafe() && d.OnEmpty != nil {
		d.OnEmpty()
	}

	return item
}

// Read the first item of the dequeue without removing it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PeekFront() T {
	// If the dequeue is empty, wait until an item is added
	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()
	for d.isEmpty_unsafe() {
		d.writeCond.Wait()
	}

	element := d.list.Front()
	return element.Value.(T)
}

// Read the last item of the dequeue without removing it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PeekBack() T {
	// If the dequeue is empty, wait until an item is added
	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()
	for d.isEmpty_unsafe() {
		d.writeCond.Wait()
	}

	element := d.list.Back()
	return element.Value.(T)
}

// ================================[Size/Capacity related]================================

// Set dequeue capacity, if capacity is 0, dequeue is infinite.
// Capacity must also be greater than the current dequeue size.
// If an invalid capacity is sent, an error is returned and the dequeue capacity is not changed.
func (d *BlockingDequeue[T]) SetCapacity(capacity int) error {
	if capacity < 0 {
		return fmt.Errorf("capacity must be >= 0")
	}

	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()

	if capacity > 0 && capacity < d.list.Len() {
		return fmt.Errorf("capacity (%d) must be >= the current size (%d), or 0 for infinite capacity", capacity, d.list.Len())
	}

	d.capacity = capacity

	// Notify any blocked producer now that the capacity has changed (potentially increased)
	d.writeCond.Broadcast()

	return nil
}

// Get current capacity of the dequeue.
func (d *BlockingDequeue[T]) Capacity() int {
	return d.capacity
}

// Return the number of elements in the dequeue.
func (d *BlockingDequeue[T]) Size() int {
	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()

	return d.list.Len()
}

// Return true if the dequeue is empty, without acquiring any locks.
func (d *BlockingDequeue[T]) isEmpty_unsafe() bool {
	return d.list.Len() == 0
}

// Return true if the dequeue is empty.
func (d *BlockingDequeue[T]) IsEmpty() bool {
	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()

	return d.isEmpty_unsafe()
}

// Return true if the dequeue is full, without acquiring any locks.
func (d *BlockingDequeue[T]) isFull_unsafe() bool {
	return d.capacity > 0 && d.list.Len() == d.capacity
}

// Return true if the dequeue is full.
// i.e. the dequeue has limited capacity and the current size is equal to that capacity.
func (d *BlockingDequeue[T]) IsFull() bool {
	d.writeCond.L.Lock()
	defer d.writeCond.L.Unlock()

	return d.isFull_unsafe()
}
