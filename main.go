package main

import "fmt"

type RingBuffer struct {
	data     []int
	capacity int
	head     int
	tail     int
	count    int
}

func NewRingBuffer(capacity int) *RingBuffer {
	if capacity <= 0 {
		panic("capacity must be grater than zero")
	}

	return &RingBuffer{
		data:     make([]int, capacity),
		capacity: capacity,
	}
}

func (rb *RingBuffer) Insert(value int) bool {
	if rb.count == rb.capacity {
		return false
	}

	rb.data[rb.tail] = value

	rb.tail = (rb.tail + 1) % rb.capacity

	rb.count++

	return true
}

func main() {
	rb := NewRingBuffer(5)

	for i := 1; i <= 6; i++ {
		value := i * 10

		success := rb.Insert(value)

		fmt.Printf(
			"Insert(%d): %t\n",
			value,
			success,
		)
	}

	fmt.Printf("\nFinal buffer: %+v\n", rb)
}
