package main

import (
	_ "ipack/com.ipack/init"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
