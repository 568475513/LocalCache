package local_cache

import (
	"fmt"
	gc "github.com/patrickmn/go-cache"
	"testing"
	"time"
)

var gC = gc.New(60*time.Second, 60*time.Hour)

func BenchmarkGoCacheSet(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gC.Set(fmt.Sprintf("alive_id_xxx%v", i), i, NoExpiration)
	}
	b.StopTimer()
}

func BenchmarkGoCacheAsynDel(b *testing.B) {

	startT := time.Now()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			gC.Set(fmt.Sprintf("alive_id_xxx%v", i), i, NoExpiration)
		}(i)
	}
	for i := 0; i < b.N; i++ {
		gC.Get(fmt.Sprintf("alive_id_xxx%v", i))
	}
	b.Log(time.Since(startT))
}
