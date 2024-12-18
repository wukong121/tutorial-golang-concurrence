package main

import "testing"

func TestStream(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	stream := Stream(done, 1, 2, 3, 4, 5)

	for v := range stream {
		t.Log(v)
	}
}

func TestTakeN(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	stream := Stream(done, 1, 2, 3, 4, 5)
	stream = TakeN(done, 3, stream)

	for v := range stream {
		t.Log(v)
	}
}

func Stream(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	c := make(chan interface{})

	go func() {
		defer close(c)
		for _, v := range values {
			select {
			case <-done:
				return
			case c <- v:
			}
		}
	}()

	return c
}

func TakeN(done <-chan interface{}, n int, c <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case out <- <-c:
			}
		}
	}()

	return out
}
