# Blocking Dequeue

![Tests status](https://github.com/AmrSaber/go-blocking-dequeue/actions/workflows/tests.yaml/badge.svg?branch=master)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/AmrSaber/go-blocking-dequeue?color=blue&display_name=tag&sort=semver)
![GitHub license](https://img.shields.io/github/license/AmrSaber/go-blocking-dequeue)

This package (repo) provides an implementation of a thread-safe, blocking, generic, infinite dequeue that can be used as FIFO or LIFO or a hybrid between the 2.

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
integersDequeue := blocking_dequeue.NewBlockingDequeue[int]()
integersDequeue.PushBack(10)

// Strings dequeue
stringsDequeue := blocking_dequeue.NewBlockingDequeue[string]()
stringsDequeue.PushBack("hello")

type User struct {
  Username string
  Age      int
}

// Dequeue of custom type
usersDequeue := blocking_dequeue.NewBlockingDequeue[User]()
usersDequeue.PushBack(User{ "Amr", 25 })

// Pointer dequeue
usersPtrDequeue := blocking_dequeue.NewBlockingDequeue[*User]()
usersPtrDequeue.PushBack(&User{ "Amr", 25 })
```

The dequeue is implemented using generics, so it can hold any datatype, just pass the data type you want between the square brackets `[ ]`.

### Capacity

You can set a capacity if you want (the default is unlimited capacity)

```go
dequeue.SetCapacity(100) // Limit the dequeue to only carry 100 elements
dequeue.SetCapacity(0) // Infinite capacity (the default)
```

Some notes about setting the capacity:

- Pushing to the dequeue will block if it has full capacity, the goroutine will be blocked until some values are popped from the dequeue.
- As in the above example, setting the capacity to 0 means having infinite capacity.
- If you attempt the set the capacity to be lower than the current size of the dequeue, the method `SetCapacity` will return an error, and the capacity won't be updated. The same thing will also happen if you try to set the capacity with a negative number.
- You can always set the capacity to 0 regardless of the current size.

### Usage as Queue

```go
dq := blocking_dequeue.NewBlockingDequeue[int]()

dq.PushBack(1) // Pushed to the end of the dequeue
dq.PushBack(2) // Pushed to the end of the dequeue
dq.PushBack(3) // Pushed to the end of the dequeue

dq.PopFront() // Pops from the top, returns 1
dq.PopFront() // Pops from the top, returns 2
dq.PopFront() // Pops from the top, returns 3
```

### Usage as Stack

```go
dq := blocking_dequeue.NewBlockingDequeue[int]()

dq.PushFront(1) // Pushed to the start of the dequeue
dq.PushFront(2) // Pushed to the start of the dequeue
dq.PushFront(3) // Pushed to the start of the dequeue

dq.PopFront() // Pops from the top, returns 3
dq.PopFront() // Pops from the top, returns 2
dq.PopFront() // Pops from the top, returns 1
```

### Listeners

You can attach listeners `onEmpty` and `onFull` to be invoked when the dequeue is empty or full respectively

```go
dq := blocking_dequeue.NewBlockingDequeue[int]()

dq.SetOnEmpty(func() {
  fmt.Println("Dequeue is now empty")
})

dq.PushBack(1)
dq.PopFront() // "Dequeue is now empty" is printed as onEmpty listener is invoked

// Dequeue can only be full if it has a positive capacity
dq.SetCapacity(2)

dq.SetOnFull(func() {
  fmt.Println("Dequeue is now full")
})

dq.PushBack(1)
dq.PushBack(2) // "Dequeue is now full" is printed as onFull listener is invoked
```

## API Documentation

The package itself exposes 1 function `NewBlockingQueue` that is used to create a new dequeue and return a pointer to it.

The dequeue itself exposes the following methods:

- `PushFront`, `PushBack`
- `PopFront`, `PopBack`
- `PeekFront`, `PeekBack`
- `Capacity`, `SetCapacity`, `IsFull`
- `Size`, `IsEmpty`

It also supports setting the listeners (see example above):

- `SetOnEmpty`
- `SetOnFull`

The detailed documentation can be found at the related [go packages page](https://pkg.go.dev/github.com/AmrSaber/go-blocking-dequeue).

## Limitations and Drawbacks

This dequeue is implemented using build-in `container/list` so all of the operations are done in O(1) time complexity.

However, due to the thread-safe nature, and all the lock/unlock/wait/signal/broadcast logic, it's expected to be a bit slower than the plain `container/list`. If you intend to use this dequeue in a single threaded context (where only a single goroutine will have access to it) it's advised to use the built-in `container/list` instead.

If you intend to use it as a limited capacity queue to communicate between goroutines, it would be better to use built-in channels with buffer, so instead of

```go
dq := blocking_dequeue.NewBlockingDequeue[int]()
dq.SetCapacity(10)

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

That is unless you need access the other provided methods and listeners, such as `Peek` variations, `Size`, `IsFull`, `onEmpty`, `onFull`, and so on...

## Benchmarking

No benchmarking against the built-in `container/list` nor channels yet. But it's in the plan.

## Contribution

If you find a bug, you are welcome to [create a ticket on github](https://github.com/AmrSaber/go-blocking-dequeue/issues) or create a PR with the fix directly mentioning in the description of the PR what is the problem and how your PR fixes it.

If you want to suggest a feature (even if you have no idea how it will be implemented) feel free to [open a ticket with it](https://github.com/AmrSaber/go-blocking-dequeue/issues).
