package local_cache

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

var lc = NewLocalCache(10000, 5*time.Second)

func TestNewLocalCache(t *testing.T) {

	//var wg sync.WaitGroup
	printlnMem("first record")
	printlnMem("second record")
	//lc.Set("b", 2, time.Second * 5)
	for i := 0; i < 1000000; i++ {
		lc.Set(fmt.Sprintf("alive_id_xxx%v", i), i, time.Second*10)
	}
	printlnMem("third record")
	lc.Get("alive_id_xxx1000")
	for {
		time.Sleep(time.Second * 5)
		printlnMem("after gc")
	}

	//now := time.Now()

	//wg.Add(2)
	//go func() {
	//	defer wg.Done()
	//	for i:=0; i < 1000000; i++{
	//		lc.Set(fmt.Sprintf("alive_id_xxx%v", i), i , time.Second * 10)
	//	}
	//}()
	//
	//go func() {
	//	defer wg.Done()
	//	for i:=0; i < 1000000; i++{
	//		lc.Get(fmt.Sprintf("alive_id_xxx%v", i))
	//	}
	//}()
	//wg.Wait()
	//fmt.Println(fmt.Sprintf("func deal with time: %v",  time.Since(now)))
	//lc.Set("b", 2)
	//lc.Set("c", 3)
	//lc.Set("你好", 4)
	//lc.Set("什么", 5)
	//printlnMem()
}

func printlnMem(str string) {

	runtime.GC() // 强制进行垃圾收集以更准确地获取内存使用情况

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Println(fmt.Sprintf("%s MemAlloc: %dKB, TotalAlloc:%dKB, HeapAlloc:%dKB, HeapInuse:%dKB, StackInuse:%dKB", str, mem.Alloc/1024, mem.TotalAlloc/1024, mem.HeapAlloc/1024, mem.HeapInuse/1024, mem.StackInuse/1024))
}
