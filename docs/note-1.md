### Why do we need a separate `count` field if we already have `head` and `tail` ?

Think about these two situations:
- Empty buffer: head = 0, tail = 0
- Full buffer: After five insertions, head = 0, tail = 0

> Both have identical pointer positions!

---

We use count to distinguish between an empty and a full ring buffer.

With capacity = 5:
- Empty: count == 0
- Full: count == capacity (5)
- Partially filled: 0 < count < capacity

Even when head == tail, we can determine the buffer's state using count.
