package main

import (
	"fmt"
	"reflect"
	"testing"
)

func assertState(t *testing.T, rb *RingBuffer, length, capacity int) {
	t.Helper()
	if got := rb.Len(); got != length {
		t.Errorf("Len() = %d, want %d", got, length)
	}
	if got := rb.Cap(); got != capacity {
		t.Errorf("Cap() = %d, want %d", got, capacity)
	}
	if got := rb.IsEmpty(); got != (length == 0) {
		t.Errorf("IsEmpty() = %t, want %t", got, length == 0)
	}
	if got := rb.IsFull(); got != (length == capacity) {
		t.Errorf("IsFull() = %t, want %t", got, length == capacity)
	}
}

func snapshot(rb *RingBuffer) RingBuffer {
	copy := *rb
	copy.data = append([]int(nil), rb.data...)
	return copy
}

func TestNewRingBufferInvalidCapacity(t *testing.T) {
	for _, capacity := range []int{0, -1, -10} {
		t.Run(fmt.Sprint(capacity), func(t *testing.T) {
			defer func() {
				if got := recover(); got != "capacity must be greater than zero" {
					t.Errorf("panic = %v, want capacity must be greater than zero", got)
				}
			}()
			NewRingBuffer(capacity)
		})
	}
}

func TestNewRingBufferInitialState(t *testing.T) {
	for _, capacity := range []int{1, 2, 5, 16} {
		t.Run(fmt.Sprint(capacity), func(t *testing.T) {
			rb := NewRingBuffer(capacity)
			assertState(t, rb, 0, capacity)
			if len(rb.data) != capacity || rb.head != 0 || rb.tail != 0 {
				t.Errorf("unexpected initial buffer: %+v", rb)
			}
			if !reflect.DeepEqual(rb.data, make([]int, capacity)) {
				t.Errorf("initial data = %v, want zero-filled storage", rb.data)
			}
		})
	}
}

func TestRingBufferCapacityOne(t *testing.T) {
	rb := NewRingBuffer(1)
	for _, want := range []int{42, 0, -7, 19} {
		assertState(t, rb, 0, 1)
		if !rb.Insert(want) {
			t.Fatal("Insert rejected a value in an empty buffer")
		}
		assertState(t, rb, 1, 1)
		before := snapshot(rb)
		if rb.Insert(99) {
			t.Error("Insert accepted a value in a full buffer")
		}
		if !reflect.DeepEqual(*rb, before) {
			t.Fatal("rejected Insert changed the buffer")
		}
		if got, ok := rb.Read(); !ok || got != want {
			t.Fatalf("Read() = (%d, %t), want (%d, true)", got, ok, want)
		}
		assertState(t, rb, 0, 1)
		if rb.data[0] != 0 {
			t.Error("Read did not clear the consumed slot")
		}
	}
}

func TestRingBufferFullInsertUnchanged(t *testing.T) {
	for _, capacity := range []int{2, 5} {
		t.Run(fmt.Sprint(capacity), func(t *testing.T) {
			rb := NewRingBuffer(capacity)
			// Move both indices so the full state also covers wrapped storage.
			rb.Insert(100)
			rb.Read()
			for value := 1; value <= capacity; value++ {
				if !rb.Insert(value) {
					t.Fatal("Insert rejected a value before the buffer was full")
				}
				assertState(t, rb, value, capacity)
			}
			before := snapshot(rb)
			for range 3 {
				if rb.Insert(-99) {
					t.Error("Insert accepted a value in a full buffer")
				}
				if !reflect.DeepEqual(*rb, before) {
					t.Fatal("rejected Insert changed the buffer")
				}
			}
			for want := 1; want <= capacity; want++ {
				if got, ok := rb.Read(); !ok || got != want {
					t.Fatalf("Read() = (%d, %t), want (%d, true)", got, ok, want)
				}
				assertState(t, rb, capacity-want, capacity)
			}
		})
	}
}

func TestRingBufferEmptyReadUnchanged(t *testing.T) {
	for _, drained := range []bool{false, true} {
		t.Run(fmt.Sprintf("drained=%t", drained), func(t *testing.T) {
			rb := NewRingBuffer(3)
			if drained {
				rb.Insert(-1)
				rb.Insert(0)
				rb.Read()
				rb.Read()
			}
			before := snapshot(rb)
			for range 3 {
				if got, ok := rb.Read(); ok || got != 0 {
					t.Errorf("Read() = (%d, %t), want (0, false)", got, ok)
				}
				if !reflect.DeepEqual(*rb, before) {
					t.Fatal("empty Read changed the buffer")
				}
				assertState(t, rb, 0, 3)
			}
		})
	}
}

func TestRingBufferRepeatedWraparoundFIFO(t *testing.T) {
	for _, capacity := range []int{1, 2, 3, 5, 8} {
		t.Run(fmt.Sprint(capacity), func(t *testing.T) {
			rb := NewRingBuffer(capacity)
			var queue []int
			next := -50
			read := func() {
				t.Helper()
				oldHead := rb.head
				want := queue[0]
				if got, ok := rb.Read(); !ok || got != want {
					t.Fatalf("Read() = (%d, %t), want (%d, true)", got, ok, want)
				}
				queue = queue[1:]
				if rb.data[oldHead] != 0 {
					t.Errorf("consumed slot %d was not cleared", oldHead)
				}
				assertState(t, rb, len(queue), capacity)
			}
			for cycle := range 120 {
				for len(queue) < capacity {
					if !rb.Insert(next) {
						t.Fatal("Insert rejected a value before the buffer was full")
					}
					queue = append(queue, next)
					next++
					assertState(t, rb, len(queue), capacity)
				}
				// Alternate partial drains and complete drains across index wraps.
				for n := 0; n < 1+cycle%capacity; n++ {
					read()
				}
			}
			for len(queue) > 0 {
				read()
			}
			assertState(t, rb, 0, capacity)
		})
	}
}
