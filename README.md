# Ring Buffer

A small Go learning project that implements a fixed-capacity ring buffer for integers.

## Current features

- Create a buffer with a positive capacity.
- Insert integers in constant time.
- Reject insertions when the buffer is full, without overwriting existing values.
- Track the number of stored values separately from the head and tail positions.

Removal and reading operations are not implemented yet. The code is in `package main`, with a runnable demonstration rather than an importable library.

## Requirements

The Go version declared in `go.mod` is **1.27.1**. Use a toolchain that supports this version.

## Run the example

```sh
go run .
```

The example creates a buffer with capacity 5 and tries to insert six values:

```text
Insert(10): true
Insert(20): true
Insert(30): true
Insert(40): true
Insert(50): true
Insert(60): false

Final buffer: &{data:[10 20 30 40 50] capacity:5 head:0 tail:0 count:5}
```

## Usage

Within this package:

```go
rb := NewRingBuffer(2)

rb.Insert(10) // true
rb.Insert(20) // true
rb.Insert(30) // false: the buffer is full
```

`NewRingBuffer` panics if the capacity is zero or negative. `Insert` returns `true` when it stores a value and `false` when the buffer is full.

## How it works

The buffer stores values in a preallocated slice. Each successful insertion advances `tail` using modulo arithmetic and increments `count`:

```go
rb.tail = (rb.tail + 1) % rb.capacity
```

The `head` field is reserved for future read or removal operations. A separate `count` distinguishes a full buffer from an empty buffer when `head` and `tail` are equal.

See [the note on the count field](docs/note-1.md) for more detail.

## Project layout

- `main.go`: buffer implementation and demonstration.
- `go.mod`: module path and Go version.
- `docs/note-1.md`: explanation of full and empty buffer states.
