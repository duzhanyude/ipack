package main

import (
	_ "capture/com.capture/init"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
