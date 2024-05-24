package local_cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLruLocalCache(t *testing.T) {

	var wg sync.WaitGroup
	printlnMem("first record")
	lc := NewLruLocalCache(10000)
	printlnMem("second record")

	now := time.Now()

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			_ = lc.Put(fmt.Sprintf("alive_id_xxx%v", i), i, time.Second*5)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			lc.Get(fmt.Sprintf("alive_id_xxx%v", i))
		}
	}()
	wg.Wait()
	fmt.Println(fmt.Sprintf("func deal with time: %v", time.Since(now)))
	//lc.Set("b", 2)
	//lc.Set("c", 3)
	//lc.Set("你好", 4)
	//lc.Set("什么", 5)
	//printlnMem()
}
