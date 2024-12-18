package main

import "testing"

func TestMapReduce(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	stream := Stream(done, 1, 2, 3, 4, 5)
	stream = Map(done, func(v interface{}) interface{} {
		return v.(int) * 2
	}, stream)

	for v := range stream {
		t.Log(v)
	}
}

func TestReduce(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	stream := Stream(done, 1, 2, 3, 4, 5)
	stream = Reduce(done, func(a, b interface{}) interface{} {
		return a.(int) + b.(int)
	}, stream)

	for v := range stream {
		t.Log(v)
	}
}

func Map(done <-chan interface{}, f func(interface{}) interface{}, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			case out <- f(v):
			}
		}
	}()

	return out
}

func Reduce(done <-chan interface{}, f func(interface{}, interface{}) interface{}, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		var val interface{}
		ok := false

		for v := range in {
			select {
			case <-done:
				return
			default:
				if !ok {
					val = v
					ok = true
					continue
				}
				val = f(val, v)
			}
		}
		out <- val
	}()

	return out
}
