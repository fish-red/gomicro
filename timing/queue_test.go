package timing

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	items := map[string]uint32{
		"banana": 4,
		"apple":  2,
		"pear":   1,
		"grape":  3,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	q := make(Queue, len(items))

	i := 0
	for k, v := range items {
		q[i] = &Item{
			Description: k,
			Timestamp:   v,
		}
		i++
	}

	heap.Init(&q)
	// Take the items out; they arrive in decreasing priority order.
	for q.Len() > 0 {
		item := heap.Pop(&q).(*Item)
		fmt.Printf("after init: %+v\n", item)
	}

	// insert a new item and modify it's priority
	item := &Item{
		Description: "orange",
		Timestamp:   13,
	}
	heap.Push(&q, item)

	// Take the items out; they arrive in decreasing priority order.
	for q.Len() > 0 {
		item := heap.Pop(&q).(*Item)
		fmt.Printf("after push: %+v\n", item)
	}
}
