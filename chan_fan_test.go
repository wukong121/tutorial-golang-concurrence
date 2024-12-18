package main

import (
	"reflect"
	"testing"
)

func TestFanIn(t *testing.T) {
	in := make(chan interface{})
	chans := FanOut(in, 10)

	go func() {
		defer close(in)
		for i := 0; i < 100; i++ {
			in <- i
		}
	}()

	for val := range FanIn(chans...) {
		t.Log(val)
	}
}

func FanIn(cs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		var cases []reflect.SelectCase
		for _, c := range cs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok {
				// remove the closed channel
				cases = append(cases[:i], cases[i+1:]...) // delete element i
				continue
			}
			out <- v.Interface()
		}
	}()

	return out
}

func FanOut(in <-chan interface{}, n int) []<-chan interface{} {
	var chans []<-chan interface{}

	for i := 0; i < n; i++ {
		ch := make(chan interface{})
		chans = append(chans, ch)
		go func() {
			defer close(ch)
			for val := range in {
				ch <- val
			}
		}()
	}

	return chans
}
