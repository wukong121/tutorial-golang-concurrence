package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("goroutine exit...")
				return
			default:
				fmt.Println("goroutine running...")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	wg.Wait()

}

func TestSelect2(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan struct{})

	go func() {
		defer func() { ch2 <- struct{}{} }()
		val := <-ch1
		fmt.Println("value is ", val)
		time.Sleep(5 * time.Second)
	}()

	for {
		select {
		case ch1 <- 1:
			fmt.Println("send success")
		case <-ch2:
			fmt.Println("business is over")
			return
		case <-time.After(2 * time.Second):
			fmt.Println("timeout")
			return
		}
	}
}
