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
