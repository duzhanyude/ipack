package main

import (
	_ "com.capture/init"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
