package test

import (
	_ "net/http/pprof"
	"sync"
	"testing"
	"time"
)

func TestGMP(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(i%5))
		}(i)
	}
	wg.Wait()
}
