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
	list        *list.List
	itemAdded   *sync.Cond
	itemRemoved *sync.Cond
	capacity    int

	OnFull  func() // Callback function invoked when the dequeue is full
	OnEmpty func() // Callback function invoked when the dequeue is empty
}

// Creates a new blocking dequeue with infinite capacity.
func NewBlockingDequeue[T any]() *BlockingDequeue[T] {
	d := new(BlockingDequeue[T])
	d.list = list.New()
	d.itemAdded = sync.NewCond(&sync.Mutex{})
	d.itemRemoved = sync.NewCond(&sync.Mutex{})
	return d
}

// =================================[Push/Pop/Peek]=================================

// Add an item into the front (top) of the dequeue. Blocks if dequeue is full.
func (d *BlockingDequeue[T]) PushFront(item T) {
	// If the dequeue is full, wait until an item is removed
	d.itemRemoved.L.Lock()
	defer d.itemRemoved.L.Unlock()
	for d.isFull_unsafe() {
		d.itemRemoved.Wait()
	}

	// Add the item to the dequeue
	d.itemAdded.L.Lock()
	defer d.itemAdded.L.Unlock()

	d.list.PushFront(item)

	// Call the OnFull callback if the dequeue is full
	if d.isFull_unsafe() && d.OnFull != nil {
		d.OnFull()
	}

	// Notify the consumer that an item has been added
	d.itemAdded.Signal()
}

// Add an item to the end of the dequeue. Blocks if dequeue is full.
func (d *BlockingDequeue[T]) PushBack(item T) {
	// If the dequeue is full, wait until an item is removed
	d.itemRemoved.L.Lock()
	defer d.itemRemoved.L.Unlock()
	for d.isFull_unsafe() {
		d.itemRemoved.Wait()
	}

	// Add the item to the dequeue
	d.itemAdded.L.Lock()
	defer d.itemAdded.L.Unlock()

	d.list.PushBack(item)

	// Call the OnFull callback if the dequeue is full
	if d.isFull_unsafe() && d.OnFull != nil {
		d.OnFull()
	}

	// Notify the consumer that an item has been added
	d.itemAdded.Signal()
}

// Read the first item (on the top/front) of the dequeue and remove it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PopFront() T {
	// If the dequeue is empty, wait until an item is added
	d.itemAdded.L.Lock()
	defer d.itemAdded.L.Unlock()
	for d.isEmpty_unsafe() {
		d.itemAdded.Wait()
	}

	// Remove the item from the dequeue
	d.itemRemoved.L.Lock()
	defer d.itemRemoved.L.Unlock()

	item := d.list.Remove(d.list.Front()).(T)

	// Notify the producer that an item has been removed
	defer d.itemRemoved.Signal()

	// Call the OnEmpty callback if the dequeue is empty
	if d.isEmpty_unsafe() && d.OnEmpty != nil {
		d.OnEmpty()
	}

	return item
}

// Read the last item (at the end/back) of the dequeue and remove it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PopBack() T {
	// If the dequeue is empty, wait until an item is added
	d.itemAdded.L.Lock()
	defer d.itemAdded.L.Unlock()
	for d.isEmpty_unsafe() {
		d.itemAdded.Wait()
	}

	// Remove the item from the dequeue
	d.itemRemoved.L.Lock()
	defer d.itemRemoved.L.Unlock()

	item := d.list.Remove(d.list.Back()).(T)

	// Notify the producer that an item has been removed
	defer d.itemRemoved.Signal()

	// Call the OnEmpty callback if the dequeue is empty
	if d.isEmpty_unsafe() && d.OnEmpty != nil {
		d.OnEmpty()
	}

	return item
}

// Read the first item of the dequeue without removing it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PeekFront() T {
	// If the dequeue is empty, wait until an item is added
	d.itemAdded.L.Lock()
	defer d.itemAdded.L.Unlock()
	for d.isEmpty_unsafe() {
		d.itemAdded.Wait()
	}

	element := d.list.Front()
	return element.Value.(T)
}

// Read the first item of the dequeue without removing it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PeekBack() T {
	// If the dequeue is empty, wait until an item is added
	d.itemAdded.L.Lock()
	defer d.itemAdded.L.Unlock()
	for d.isEmpty_unsafe() {
		d.itemAdded.Wait()
	}

	element := d.list.Back()
	return element.Value.(T)
}

// ================================[Size/Capacity related]================================

// Acquires all update locks and returns a function to release them.
func (d *BlockingDequeue[T]) acquireWriteLocks() func() {
	d.itemAdded.L.Lock()
	d.itemRemoved.L.Lock()

	return func() {
		d.itemAdded.L.Unlock()
		d.itemRemoved.L.Unlock()
	}
}

// Set dequeue capacity, if capacity is 0, dequeue is infinite.
// Capacity must also be greater than the current dequeue size.
// If an invalid capacity is set, an error is returned and the dequeue capacity is not changed.
func (d *BlockingDequeue[T]) SetCapacity(capacity int) error {
	if capacity < 0 {
		return fmt.Errorf("capacity must be >= 0")
	}

	release := d.acquireWriteLocks()
	defer release()

	if capacity > 0 && capacity < d.list.Len() {
		return fmt.Errorf("capacity (%d) must be >= the current size (%d), or 0 for infinite capacity", capacity, d.list.Len())
	}

	d.capacity = capacity

	// Notify any blocked producer now that the capacity has changed (potentially increased)
	d.itemRemoved.Broadcast()

	return nil
}

func (d *BlockingDequeue[T]) Capacity() int {
	return d.capacity
}

// Return the number of elements in the dequeue.
func (d *BlockingDequeue[T]) Size() int {
	release := d.acquireWriteLocks()
	defer release()

	return d.list.Len()
}

// Return true if the dequeue is empty without acquiring any locks.
func (d *BlockingDequeue[T]) isEmpty_unsafe() bool {
	return d.list.Len() == 0
}

// Return true if the dequeue is empty.
func (d *BlockingDequeue[T]) IsEmpty() bool {
	release := d.acquireWriteLocks()
	defer release()

	return d.isEmpty_unsafe()
}

// Return true if the dequeue is full without acquiring any locks.
func (d *BlockingDequeue[T]) isFull_unsafe() bool {
	return d.capacity > 0 && d.list.Len() == d.capacity
}

// Return true if the dequeue is full.
// i.e. the dequeue has limited capacity and the current size is equal to that capacity.
func (d *BlockingDequeue[T]) IsFull() bool {
	release := d.acquireWriteLocks()
	defer release()

	return d.isFull_unsafe()
}
