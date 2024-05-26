package local_cache

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func BenchmarkLruAsynDel(b *testing.B) {
	lruC := NewLruLocalCache(10000)

	startT := time.Now()
	for i := 0; i < b.N; i++ {
		go func() {
			lruC.Put(fmt.Sprintf("alive_id_xxx%v", i), i, NoExpiration)
		}()
	}
	for i := 0; i < b.N; i++ {
		lruC.Get(fmt.Sprintf("alive_id_xxx%v", i))
	}
	b.Log(time.Since(startT))
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
