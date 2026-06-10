// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Original source: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/5cdce9f43de329beba101eb056a6d7abb938fb61/internal/conns/mutexkv.go

package api

import (
	"context"
	"sync"
)

// GlobalMutexKV is a global MutexKV for use within this plugin.
var GlobalMutexKV = newMutexKV()

// mutexKV is a simple key/value store for arbitrary mutexes. It can be used to
// serialize changes across arbitrary collaborators that share knowledge of the
// keys they must serialize on.
type mutexKV struct {
	lock  sync.Mutex
	store map[string]*sync.Mutex
}

// Lock acquires the mutex for the given key, honouring ctx cancellation.
// On success the caller is responsible for calling Unlock for the same key.
// If ctx is cancelled before the lock can be acquired, ctx.Err() is returned
// and the caller must NOT call Unlock — the acquisition is abandoned and the
// inner mutex is released as soon as it can be taken.
func (m *mutexKV) Lock(key string, ctx context.Context) error {
	mu := m.Get(key)

	acquired := make(chan struct{})
	abandoned := make(chan struct{})
	go func() {
		mu.Lock()
		select {
		case acquired <- struct{}{}:
		case <-abandoned:
			mu.Unlock()
		}
	}()

	select {
	case <-acquired:
		return nil
	case <-ctx.Done():
		close(abandoned)
		return ctx.Err()
	}
}

// Unlock releases the mutex for the given key. The caller must have
// previously held the lock (a successful Lock for the same key).
func (m *mutexKV) Unlock(key string, ctx context.Context) {
	m.Get(key).Unlock()
}

// Returns a mutex for the given key, no guarantee of its lock status
func (m *mutexKV) Get(key string) *sync.Mutex {
	m.lock.Lock()
	defer m.lock.Unlock()
	mutex, ok := m.store[key]
	if !ok {
		mutex = &sync.Mutex{}
		m.store[key] = mutex
	}
	return mutex
}

// Returns a properly initialized MutexKV
func newMutexKV() *mutexKV {
	return &mutexKV{
		store: make(map[string]*sync.Mutex),
	}
}
