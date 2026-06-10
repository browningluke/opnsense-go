package api

import (
	"context"
	"testing"
	"time"
)

func TestMutexKV_LockUncontendedDoesNotHang(t *testing.T) {
	m := newMutexKV()

	for i := 0; i < 50000; i++ {
		done := make(chan struct{})
		go func() {
			defer close(done)
			if err := m.Lock("k", context.Background()); err != nil {
				t.Errorf("Lock returned error: %v", err)
				return
			}
			m.Unlock("k", context.Background())
		}()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatalf("Lock hung on iteration %d", i)
		}
	}
}

func TestMutexKV_LockHonoursContextCancellation(t *testing.T) {
	m := newMutexKV()

	if err := m.Lock("k", context.Background()); err != nil {
		t.Fatalf("first Lock: %v", err)
	}
	defer m.Unlock("k", context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := m.Lock("k", ctx)
	if err == nil {
		m.Unlock("k", ctx)
		t.Fatal("expected Lock to return ctx.Err(), got nil")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("Lock took too long to honour cancellation: %v", elapsed)
	}
}
