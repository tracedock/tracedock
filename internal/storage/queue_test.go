package storage

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	suite.Suite
	queue *Queue
}

func Test_Queue_Enqueue(t *testing.T) {
	var wg sync.WaitGroup

	t.Run("should correctly add items", func(t *testing.T) {
		var queue = NewQueue()

		assert.Equal(t, 0, len(queue.items))

		for i := 1; i < 100; i++ {
			queue.Enqueue(i)

			assert.Equal(t, i, len(queue.items))
		}
	})

	t.Run("should handle concurrency", func(t *testing.T) {
		var goroutines = 100
		var enqueues = 1000
		var queue = NewQueue()

		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()

				for j := 0; j < enqueues; j++ {
					queue.Enqueue(rand.Int())
				}
			}()
		}

		wg.Wait()

		assert.Equal(t, goroutines*enqueues, len(queue.items))
	})
}

func Test_Queue_Dequeue(t *testing.T) {
	var wg sync.WaitGroup

	t.Run("should return nil when queue is empty", func(t *testing.T) {
		var queue = NewQueue()

		assert.Equal(t, 0, len(queue.items))
		assert.Equal(t, nil, queue.Dequeue())
	})

	t.Run("should dequeue in fifo way", func(t *testing.T) {
		var items = 100
		var queue = NewQueue()

		for i := 1; i <= items; i++ {
			queue.Enqueue(i)
		}

		for i := items; i == 0; i-- {
			assert.Equal(t, i, queue.Dequeue())
			assert.Equal(t, i, len(queue.items))
		}
	})

	t.Run("should handle concurrency", func(t *testing.T) {
		var mu sync.Mutex

		var queue = NewQueue()
		var enqueues = 1000
		var dequeued = make(map[string]bool)

		for i := 0; i < enqueues; i++ {
			queue.Enqueue(fmt.Sprintf("%v", i))
		}

		wg.Add(enqueues)

		for i := 0; i < enqueues; i++ {
			go func() {
				defer wg.Done()

				for {
					var item = queue.Dequeue()

					if item == nil {
						return
					}

					mu.Lock()
					dequeued[fmt.Sprintf("%v", item)] = true
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		assert.Equal(t, enqueues, len(dequeued))

		for i := 0; i < enqueues; i++ {
			assert.True(t, dequeued[fmt.Sprintf("%v", i)], fmt.Sprint(i))
		}
	})
}

func Benchmark_Queue_Enqueue(b *testing.B) {
	var queue = NewQueue()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.Enqueue(rand.Int())
		}
	})
}

func Benchmark_Queue_Dequeue(b *testing.B) {
	var queue = NewQueue()

	for i := 0; i < b.N*10; i++ {
		queue.Enqueue(i)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.Dequeue()
		}
	})
}
