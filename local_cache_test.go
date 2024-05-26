package local_cache

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestName(t *testing.T) {
}

func BenchmarkLocalCache_Set(b *testing.B) {
	lc := NewLocalCache(10000, 1000, 60*time.Second)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		lc.Set(fmt.Sprintf("alive_id_xxx%v", i), i, NoExpiration)
	}
	b.StopTimer()
}

func BenchmarkLocalCache_LruSet(b *testing.B) {
	lc := NewLocalCache(1, 15, 5*time.Second)
	for i := 0; i < 15; i++ {
		lc.Set(fmt.Sprintf("%v", i), i, NoExpiration)
	}
	for i := 0; i < 10; i++ {
		lc.Set(fmt.Sprintf("%v", i), i, NoExpiration)
	}
	lc.bucketsDta[0].ListLRUCache()
}

func BenchmarkLocalCache_Get(b *testing.B) {
	lc := NewLocalCache(10000, 1000, 60*time.Second)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		lc.Get(fmt.Sprintf("alive_id_xxx%v", i))
	}
	b.StopTimer()
}

func BenchmarkAsynMemDel(b *testing.B) {

	printMem("begin")
	lc := NewLocalCache(100, 100, 10*time.Second)
	printMem("new mem")

	for {
		time.Sleep(time.Second * 5)
		for i := 0; i < 10000; i++ {
			go func(i int) {
				lc.Set(GenerateRandomString(10), i, 20*time.Second)
			}(i)
		}

		runtime.GC()
		lc.Get("xxx")
		printMem("after gc")
	}

}

func BenchmarkAsynDel(b *testing.B) {

	lc := NewLocalCache(10000, 1000, 60*time.Second)
	startT := time.Now()
	for i := 0; i < b.N; i++ {
		go func() {
			lc.Set(fmt.Sprintf("alive_id_xxx%v", i), i, NoExpiration)
		}()
	}
	for i := 0; i < b.N; i++ {
		lc.Get(fmt.Sprintf("alive_id_xxx%v", i))
	}

	b.Log(time.Since(startT))
}

func BenchmarkLruLocalCache(b *testing.B) {

	lc := NewLocalCache(100, 10, 5*time.Second)
	for i := 0; i < 2000; i++ {
		lc.Set(fmt.Sprintf("alive_id_xxx%v", i), i, NoExpiration)
	}
	fmt.Println(lc)
}

func printMem(str string) {

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Println(fmt.Sprintf("\n %s MemAlloc: %dKB, TotalAlloc:%dKB, HeapAlloc:%dKB, HeapInuse:%dKB, StackInuse:%dKB", str, mem.Alloc/1024, mem.TotalAlloc/1024, mem.HeapAlloc/1024, mem.HeapInuse/1024, mem.StackInuse/1024))
}

func (b *bucket) ListLRUCache() {
	node := b.head.nextI
	for node != nil {
		fmt.Println(fmt.Sprintf("key: %s, value: %d", node.key, node.value))
		node = node.nextI
	}
	fmt.Println("xxxxxxxxxxxxxxxxx\\n")
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	// 使用加密安全的随机数生成器初始化种子
	rand.Seed(time.Now().UnixNano())
	// 使用加密安全的随机数填充切片
	rand.Read(b)
	// 将字节切片转换为base64字符串
	return base64.StdEncoding.EncodeToString(b)
}
