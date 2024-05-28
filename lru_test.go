package local_cache

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

var lruCache = NewLruLocalCache(100000)

func BenchmarkLruAsynDel(b *testing.B) {

	//lc := NewLocalCache(100, 100, 10*time.Second)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lruCache.Put(fmt.Sprintf("alive_id_%d", i), i, NoExpiration)
		}(i)
	}

	//time.Sleep(time.Second)

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lruCache.Get(fmt.Sprintf("alive_id_%d", i))
		}(i)
	}
	wg.Wait()
}

func BenchmarkLruMemDel(b *testing.B) {

	printMem("begin")
	lc := NewLruLocalCache(10000)
	printMem("new mem")

	for {
		time.Sleep(time.Second * 5)
		for i := 0; i < 10000; i++ {
			go func(i int) {
				lc.Put(GenerateRandomString(10), i, 20*time.Second)
			}(i)
		}

		runtime.GC()
		lc.Get("xxx")
		printMem("after gc")
	}

}
