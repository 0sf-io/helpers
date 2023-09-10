package threadpool

import (
	"testing"
	"time"
)

func TestNewThreadPool(t *testing.T) {
	pool := NewThreadPool(5)
	if len(pool.workers) != 5 {
		t.Errorf("Expected 5 workers, got %d", len(pool.workers))
	}

	if cap(pool.queue) != 5 {
		t.Errorf("Expected 0 capacity, got %d", cap(pool.queue))
	}

	if len(pool.queue) != 0 {
		t.Errorf("Expected 0 jobs in queue, got %d", len(pool.queue))
	}

	if cap(pool.readyQueue) != 5 {
		t.Errorf("Expected 0 capacity, got %d", cap(pool.readyQueue))
	}

	if len(pool.readyQueue) != 0 {
		t.Errorf("Expected 0 jobs in ready queue, got %d", len(pool.readyQueue))
	}
}

func TestThreadPool_Start(t *testing.T) {
	pool := NewThreadPool(5)
	pool.Start()
	for _, worker := range pool.workers {
		if !worker.isActive {
			t.Errorf("Expected worker to be running, got %t", worker.isActive)
		}
	}
}

func TestThreadPool_AddJob(t *testing.T) {
	pool := NewThreadPool(5)
	pool.AddJob(func() {})

	if len(pool.queue) != 1 {
		t.Errorf("Expected 1 job in queue, got %d", len(pool.queue))
	}
}

func TestThreadPool_AddJob_Exec(t *testing.T) {
	pool := NewThreadPool(5)
	pool.Start()
	ran := false
	pool.AddJob(func() { ran = true })
	time.Sleep(time.Millisecond)
	if len(pool.queue) != 0 {
		t.Errorf("Expected 0 jobs in queue, got %d", len(pool.queue))
	}
	if !ran {
		t.Errorf("Expected job to run")
	}
}

func TestThreadPool_Stop(t *testing.T) {
	pool := NewThreadPool(5)
	pool.Start()
	pool.Stop()
	for _, worker := range pool.workers {
		if worker.isActive {
			t.Errorf("Expected worker to not be running, got %t", worker.isActive)
		}
	}
}

func BenchmarkThreadPool_AddJob(b *testing.B) {
	pool := NewThreadPool(5)
	pool.Start()
	for i := 0; i < b.N; i++ {
		pool.AddJob(func() {})
	}
}
