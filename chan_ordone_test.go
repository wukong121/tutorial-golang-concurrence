package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestOrDone(t *testing.T) {
	startTime := time.Now()
	<-Or(
		Sig(4*time.Second),
		Sig(5*time.Second),
		Sig(6*time.Second),
	)
	fmt.Println("done after", time.Since(startTime))

	<-OrWithReflect(
		Sig(4*time.Second),
		Sig(5*time.Second),
		Sig(6*time.Second),
	)
	fmt.Println("done after", time.Since(startTime))
}

func Sig(d time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(d)
	}()
	return c
}

func OrWithReflect(cs ...<-chan interface{}) <-chan interface{} {
	switch len(cs) {
	case 0:
		return nil
	case 1:
		return cs[0]
	}
	orDone := make(chan interface{})

	go func() {
		defer close(orDone)
		var cases []reflect.SelectCase
		for _, c := range cs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		reflect.Select(cases)
	}()

	return orDone
}

func Or(cs ...<-chan interface{}) <-chan interface{} {
	switch len(cs) {
	case 0:
		return nil
	case 1:
		return cs[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		switch len(cs) {
		case 2:
			select {
			case <-cs[0]:
			case <-cs[1]:
			}
		default:
			m := len(cs) / 2
			select {
			case <-Or(cs[:m]...):
			case <-Or(cs[m:]...):
			}
		}
	}()

	return orDone
}
