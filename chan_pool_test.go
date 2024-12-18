package main

import (
	"fmt"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	msg := make(chan interface{})

	for i := 0; i < 10; i++ {
		go Consume(i, msg)
	}
	go Produce(msg)
	time.Sleep(10 * time.Second)
}

func Consume(id int, msg chan interface{}) {
	for val := range msg {
		fmt.Println("worker", id, "receive value:", val)
	}
}

func Produce(msg chan interface{}) {
	for {
		select {
		case <-time.After(2 * time.Second):
			val := time.Now().Unix()
			fmt.Println("on time, send value:", val)
			msg <- val
		}
	}
}
