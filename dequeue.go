package blocking_dequeue

import (
	"sync"
)

// Blocking dequeue, implemented with a circular buffer.
// The dequeue is thread safe. And must not be copied.
type BlockingDequeue[T any] struct {
	buffer []T

	lock              *sync.Mutex
	notEmpty, notFull *sync.Cond

	first, last int
	isEmpty     bool
}

// Creates a new blocking dequeue with infinite capacity.
// The dequeue MUST only be created using this method.
func NewBlockingDequeue[T any](buffer []T) *BlockingDequeue[T] {
	d := new(BlockingDequeue[T])

	d.buffer = buffer

	d.first = 0
	d.last = 0
	d.isEmpty = true

	d.lock = &sync.Mutex{}
	d.notEmpty = sync.NewCond(d.lock)
	d.notFull = sync.NewCond(d.lock)

	return d
}

// =================================[Buffer helpers]=================================
func (d BlockingDequeue[T]) nextIndex(i int) int {
	return (i + 1) % len(d.buffer)
}

func (d BlockingDequeue[T]) prevIndex(i int) int {
	return (i - 1 + len(d.buffer)) % len(d.buffer)
}

// =================================[Push/Pop/Peek]=================================

// Add an item into the front (top) of the dequeue. Blocks if dequeue is full.
func (d *BlockingDequeue[T]) PushFront(item T) {
	d.lock.Lock()
	defer d.lock.Unlock()
	defer d.notEmpty.Signal()

	// If the dequeue is full, wait until an item is removed
	for d.isFull_unsafe() {
		d.notFull.Wait()
	}

	if !d.isEmpty {
		d.first = d.prevIndex(d.first)
	}
	d.buffer[d.first] = item

	d.isEmpty = false
}

// Add an item to the back (bottom) of the dequeue. Blocks if dequeue is full.
func (d *BlockingDequeue[T]) PushBack(item T) {
	d.lock.Lock()
	defer d.lock.Unlock()
	defer d.notEmpty.Signal()

	// If the dequeue is full, wait until an item is removed
	for d.isFull_unsafe() {
		d.notFull.Wait()
	}

	if !d.isEmpty {
		d.last = d.nextIndex(d.last)
	}
	d.buffer[d.last] = item

	d.isEmpty = false
}

// Read the first item (on the top/front) of the dequeue and remove it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PopFront() T {
	d.lock.Lock()
	defer d.lock.Unlock()
	defer d.notFull.Signal()

	// If the dequeue is empty, wait until an item is added
	for d.isEmpty_unsafe() {
		d.notEmpty.Wait()
	}

	item := d.buffer[d.first]

	if d.first == d.last {
		d.isEmpty = true
	} else {
		d.first = d.nextIndex(d.first)
	}

	return item
}

// Read the last item (at the end/back) of the dequeue and remove it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PopBack() T {
	d.lock.Lock()
	defer d.lock.Unlock()
	defer d.notFull.Signal()

	// If the dequeue is empty, wait until an item is added
	for d.isEmpty_unsafe() {
		d.notEmpty.Wait()
	}

	item := d.buffer[d.last]

	if d.first == d.last {
		d.isEmpty = true
	} else {
		d.last = d.prevIndex(d.last)
	}

	return item
}

// Read the first item of the dequeue without removing it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PeekFront() T {
	d.lock.Lock()
	defer d.lock.Unlock()

	// If the dequeue is empty, wait until an item is added
	for d.isEmpty_unsafe() {
		d.notEmpty.Wait()
	}

	return d.buffer[d.first]
}

// Read the last item of the dequeue without removing it. Blocks if the dequeue is empty.
func (d *BlockingDequeue[T]) PeekBack() T {
	d.lock.Lock()
	defer d.lock.Unlock()

	// If the dequeue is empty, wait until an item is added
	for d.isEmpty_unsafe() {
		d.notEmpty.Wait()
	}

	return d.buffer[d.last]
}

// ================================[Size/Capacity related]================================
// Return the number of elements in the dequeue.
func (d *BlockingDequeue[T]) Size() int {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.isEmpty {
		return 0
	}
	if d.first <= d.last {
		return d.last - d.first + 1
	} else {
		return (len(d.buffer) - d.first) + (d.last + 1)
	}
}

// Return true if the dequeue is empty, without acquiring any locks.
// Dequeue is empty if the first and last indices are the same.
func (d *BlockingDequeue[T]) isEmpty_unsafe() bool {
	return d.isEmpty
}

// Return true if the dequeue is empty.
func (d *BlockingDequeue[T]) IsEmpty() bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.isEmpty_unsafe()
}

// Return true if the dequeue is full, without acquiring any locks.
// Dequeue is full if the next item to be added will be the first item in the dequeue.
func (d *BlockingDequeue[T]) isFull_unsafe() bool {
	return d.nextIndex(d.last) == d.first
}

// Return true if the dequeue is full.
// i.e. the dequeue has limited capacity and the current size is equal to that capacity.
func (d *BlockingDequeue[T]) IsFull() bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.isFull_unsafe()
}
