# Ring Buffer

A small Go learning project that implements a fixed-capacity ring buffer for integers.

## Current features

- Create a buffer with a positive capacity.
- Insert and consume integers in FIFO order in constant time.
- Reuse freed slots through index wraparound.
- Inspect buffer length, capacity, and full or empty state.
- Reject insertions when the buffer is full, without overwriting existing values.
- Track the number of stored values separately from the head and tail positions.

The code is in `package main`, with a runnable demonstration rather than an importable library. It is not safe for concurrent use without external synchronization. Create buffers with `NewRingBuffer`; the zero value is not ready for use.

## Requirements

The Go version declared in `go.mod` is **1.27.1**. Use a toolchain that supports this version.

## Run the example

```sh
go run .
```

The example fills a buffer with capacity 5, consumes two values, and reuses their slots:

```text
Full: &{data:[10 20 30 40 50] capacity:5 head:0 tail:0 count:5}
Read: 10 true
Read: 20 true
Final: &{data:[60 70 30 40 50] capacity:5 head:2 tail:2 count:5}
```

The final FIFO order is `30, 40, 50, 60, 70`, even though the backing slice starts with `60, 70`.

## Usage

Within this package:

```go
rb := NewRingBuffer(2)

rb.Insert(10) // true
rb.Insert(20) // true
rb.Insert(30) // false: the buffer is full

value, ok := rb.Read() // 10, true: removes the oldest value
_ = value
_ = ok

rb.Insert(30) // true: reuses the freed slot
rb.Len()     // 2
rb.Cap()     // 2
rb.IsEmpty() // false
rb.IsFull()  // true
```

`NewRingBuffer` panics if the capacity is zero or negative. `Insert` returns `true` when it stores a value and `false` when the buffer is full. `Read` removes the oldest value and returns `(value, true)`, or `(0, false)` when empty. Zero is a valid stored value, so always check the boolean result.

## How it works

The buffer stores values in a preallocated slice. Each successful insertion advances `tail` using modulo arithmetic and increments `count`:

```go
rb.tail = (rb.tail + 1) % rb.capacity
```

A successful `Read` clears the consumed slot, advances `head` with the same modulo rule, and decrements `count`. A separate `count` distinguishes a full buffer from an empty buffer when `head` and `tail` are equal.

See [the ring buffer guide](docs/ring-buffer.md) for internals, a worked wraparound example, and an eBPF comparison. See [the note on the count field](docs/note-1.md) for full and empty states.

## Checks

```sh
go test ./...
go test -race ./...
go vet ./...
```

Tests cover invalid capacities, FIFO order, full and empty operations, valid zero values, and repeated wraparound.

## Project layout

- `main.go`: buffer implementation and demonstration.
- `main_test.go`: buffer behavior and boundary tests.
- `go.mod`: module path and Go version.
- `docs/ring-buffer.md`: data structure guide and eBPF context, adapted from the supplied note.
- `docs/note-1.md`: explanation of full and empty buffer states.
