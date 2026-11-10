package main

import (
	"fmt"
)

// RingBuffer is a fixed-capacity FIFO buffer of integers.
// It is not safe for concurrent use. Create one with NewRingBuffer.
type RingBuffer struct {
	data     []int
	capacity int
	head     int
	tail     int
	count    int
}

// NewRingBuffer returns an empty buffer with the given capacity.
// It panics if capacity is zero or negative.
func NewRingBuffer(capacity int) *RingBuffer {
	if capacity <= 0 {
		panic("capacity must be greater than zero")
	}

	return &RingBuffer{
		data:     make([]int, capacity),
		capacity: capacity,
	}
}

// Insert appends value to the buffer and returns true.
// If the buffer is full, it returns false without changing the buffer.
func (rb *RingBuffer) Insert(value int) bool {
	if rb.count == rb.capacity {
		return false
	}

	rb.data[rb.tail] = value

	rb.tail = (rb.tail + 1) % rb.capacity

	rb.count++

	return true
}

// Read removes and returns the oldest value and true, clearing its slot.
// If the buffer is empty, it returns zero and false without changing the buffer.
func (rb *RingBuffer) Read() (int, bool) {
	// Step 1: Check whether buffer is empty
	if rb.count == 0 {
		return 0, false
	}

	value := rb.data[rb.head]

	rb.data[rb.head] = 0

	rb.head = (rb.head + 1) % rb.capacity

	rb.count--

	return value, true
}

// Len returns the number of values in the buffer.
func (rb *RingBuffer) Len() int {
	return rb.count
}

// Cap returns the maximum number of values the buffer can hold.
func (rb *RingBuffer) Cap() int {
	return rb.capacity
}

// IsEmpty reports whether the buffer contains no values.
func (rb *RingBuffer) IsEmpty() bool {
	return rb.count == 0
}

// IsFull reports whether the buffer has no free slots.
func (rb *RingBuffer) IsFull() bool {
	return rb.count == rb.capacity
}

func main() {
	rb := NewRingBuffer(5)

	rb.Insert(10)
	rb.Insert(20)
	rb.Insert(30)
	rb.Insert(40)
	rb.Insert(50)

	fmt.Printf("Full: %+v\n", rb)

	// Consume the two oldest elements
	value, ok := rb.Read()
	fmt.Println("Read:", value, ok)

	value, ok = rb.Read()
	fmt.Println("Read:", value, ok)

	// Reuse freed slots
	rb.Insert(60)
	rb.Insert(70)

	fmt.Printf("Final: %+v\n", rb)
}
