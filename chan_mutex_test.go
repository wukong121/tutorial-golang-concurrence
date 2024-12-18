package main

import (
	"sync"
	"testing"
	"time"
)

type Locker struct {
	c chan struct{}
}

func NewLocker() *Locker {
	return &Locker{
		c: make(chan struct{}, 1),
	}
}

func (l *Locker) Lock() {
	l.c <- struct{}{}
}

func (l *Locker) Unlock() {
	<-l.c
}

func (l *Locker) LockWithTimeout(timeout time.Duration) bool {
	select {
	case l.c <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

func TestCMutex(t *testing.T) {
	locker := NewLocker()
	var wg sync.WaitGroup
	a := 100

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer locker.Unlock()
			locker.Lock()
			a++
			wg.Done()
		}()
	}

	wg.Wait()
	t.Log("outside a = ", a)
}

func TestLockWithTimeout(t *testing.T) {
	var wg sync.WaitGroup
	locker := NewLocker()
	locker.Lock()

	wg.Add(1)
	go func() {
		defer wg.Done()
		isLocked := locker.LockWithTimeout(1 * time.Second)
		if isLocked {
			t.Log("lock success")
		} else {
			t.Log("lock failed")
		}
	}()

	time.Sleep(2 * time.Second)
	locker.Unlock()
	wg.Wait()
}

func TestMutex(t *testing.T) {
	var mu sync.Mutex
	a := 100

	for i := 0; i < 100; i++ {
		go func() {
			defer mu.Unlock()
			mu.Lock()
			a++
		}()
	}

	time.Sleep(1 * time.Second)
	t.Log("outside a = ", a)
}
