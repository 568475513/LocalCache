package local_cache

import (
	"fmt"
	gc "github.com/patrickmn/go-cache"
	"sync"
	"testing"
	"time"
)

var (
	wg sync.WaitGroup
	gC = gc.New(60*time.Second, 60*time.Hour)
)

func BenchmarkGoCacheSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gC.Set(GenerateRandomString(10), i, NoExpiration)
	}
}

func BenchmarkGoCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gC.Get(GenerateRandomString(10))
	}
}

func BenchmarkGoCacheAsynDel(b *testing.B) {

	//lc := NewLocalCache(100, 100, 10*time.Second)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			gC.Set(fmt.Sprintf("alive_id_%d", i), i, NoExpiration)
		}(i)
	}

	//time.Sleep(time.Second)

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			gC.Get(fmt.Sprintf("alive_id_%d", i))
		}(i)
	}
	wg.Wait()
}
