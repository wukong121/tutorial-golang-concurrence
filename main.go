package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func doClean(closed chan struct{}) {
	time.Sleep(2 * time.Second)
	close(closed)

}
func main() {
	closing := make(chan struct{})
	closed := make(chan struct{})

	go func() {
		for {
			select {
			case <-closing:
				return
			default:
				time.Sleep(60 * time.Second)
			}
		}
	}()

	notifyC := make(chan os.Signal)
	signal.Notify(notifyC, syscall.SIGINT, syscall.SIGTERM)
	<-notifyC

	close(closing)
	go doClean(closed)

	select {
	case <-closed:
	case <-time.After(10 * time.Second):
		fmt.Println("timeout")
	}

	fmt.Println("process is exited")

}
