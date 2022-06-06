# Blocking Dequeue

![Tests status](https://github.com/AmrSaber/go-blocking-dequeue/actions/workflows/tests.yaml/badge.svg?branch=master)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/AmrSaber/go-blocking-dequeue?color=blue&display_name=tag&sort=semver)
![GitHub license](https://img.shields.io/github/license/AmrSaber/go-blocking-dequeue)

This package (repo) provides an implementation of a thread-safe, blocking, generic dequeue that can be used as FIFO or LIFO or a hybrid between the 2.

## Installation

To install this package, you will need to setup your go workspace first. Also this package requires go **v1.18** or later.

1. To install the package run the following command:

   ```bash
   go get github.com/AmrSaber/go-blocking-dequeue
   ```

2. To import the package:

   ```go
   import "github.com/AmrSaber/go-blocking-dequeue"
   ```

3. Use the package in code using `blocking_dequeue` module (see usage below).

## Usage

### Initialization

To create a new dequeue use `blocking_dequeue.NewBlockingDequeue` function as follows:

```go
// Integers dequeue
buffer := make([]int, 10)
integersDequeue := blocking_dequeue.NewBlockingDequeue(buffer)
integersDequeue.PushBack(10)

// Strings dequeue
buffer := make([]string, 10)
stringsDequeue := blocking_dequeue.NewBlockingDequeue(buffer)
stringsDequeue.PushBack("hello")

type User struct {
  Username string
  Age      int
}

// Dequeue of custom type
buffer := make([]User, 10)
usersDequeue := blocking_dequeue.NewBlockingDequeue(buffer)
usersDequeue.PushBack(User{ "Amr", 25 })

// Pointer dequeue
buffer := make([]*User, 10)
usersPtrDequeue := blocking_dequeue.NewBlockingDequeue(buffer)
usersPtrDequeue.PushBack(&User{ "Amr", 25 })
```

The dequeue is implemented using generics, so it can hold any datatype, just create a buffer with the desired datatype and pass it to the creation function.

### Capacity

The capacity of the dequeue is the length of the provided buffer.

### Usage as Queue

```go
buffer := make([]int, 10)
dq := blocking_dequeue.NewBlockingDequeue(buffer)

dq.PushBack(1) // Pushed to the end of the dequeue
dq.PushBack(2) // Pushed to the end of the dequeue
dq.PushBack(3) // Pushed to the end of the dequeue

dq.PopFront() // Pops from the top, returns 1
dq.PopFront() // Pops from the top, returns 2
dq.PopFront() // Pops from the top, returns 3
```

### Usage as Stack

```go
buffer := make([]int, 10)
dq := blocking_dequeue.NewBlockingDequeue(buffer)

dq.PushFront(1) // Pushed to the start of the dequeue
dq.PushFront(2) // Pushed to the start of the dequeue
dq.PushFront(3) // Pushed to the start of the dequeue

dq.PopFront() // Pops from the top, returns 3
dq.PopFront() // Pops from the top, returns 2
dq.PopFront() // Pops from the top, returns 1
```

## API Documentation

The package itself exposes 1 function `NewBlockingQueue` that is used to create a new dequeue and return a pointer to it.

The dequeue itself exposes the following methods:

- `PushFront`, `PushBack`
- `PopFront`, `PopBack`
- `PeekFront`, `PeekBack`
- `Size`, `IsEmpty`, `IsFull`

The detailed documentation can be found at the related [go packages page](https://pkg.go.dev/github.com/AmrSaber/go-blocking-dequeue#section-documentation).

## Limitations and Drawbacks

This dequeue is implemented using ring (or circular) buffer so all of the operations are done in O(1) time complexity.

However, due to the thread-safe nature, and all the lock/unlock/wait/signal logic, it's expected to be a bit slower than plain ring buffer. If you intend to use this dequeue in a single threaded context (where only a single goroutine will have access to it) it's advised to use plain circular buffer or the built-in `container/list` instead.

If you intend to use it as a limited capacity queue to communicate between goroutines, it would be better to use built-in channels with buffer, so instead of

```go
buffer := make([]int, 10)
dq := blocking_dequeue.NewBlockingDequeue(buffer)

// Push to queue
dq.PushBack(1)

// Pop from queue
dq.PopFront()
```

You better use

```go
ch := make(chan int, 10)

// Push to queue
ch <- 1

// Pop from queue
<- ch
```

That is unless you need access the other provided methods, such as `Peek` variations, `Size`, `IsFull`, and so on...

## Benchmarking

No benchmarking against plain ring buffer or the built-in `container/list` nor channels yet. But it's in the plan.

## Contribution

If you find a bug, you are welcome to [create a ticket on github](https://github.com/AmrSaber/go-blocking-dequeue/issues) or create a PR with the fix directly mentioning in the description of the PR what is the problem and how your PR fixes it.

If you want to suggest a feature (even if you have no idea how it will be implemented) feel free to [open a ticket with it](https://github.com/AmrSaber/go-blocking-dequeue/issues).
