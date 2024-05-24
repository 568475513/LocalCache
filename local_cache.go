package local_cache

// 采用hash切片桶的数据结构 分片桶上锁提高读写效率 初始化长度越大 锁颗粒度越细 内存占用也越大 算是一种以空间换时间逻辑
// 淘汰策略-懒加载-异步线程 delete 下次gc回收 参考redis过期key淘汰策略
import (
	"crypto/rand"
	"math"
	"math/big"
	insecurerand "math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	NoExpiration time.Duration = -1
)

type LC struct {
	*localCache
}

type localCache struct {
	bucketsDta []*bucket
	bucketNum  uint32
	seed       uint32
	stopCh     chan struct{}
}

type bucket struct {
	rwMu  sync.RWMutex
	items map[string]*Item
}

type Item struct {
	value      interface{}
	Expiration int64
}

func NewLocalCache(bucketNum int, cleanTime time.Duration) *LC {

	lc := &localCache{
		bucketNum:  uint32(bucketNum),
		bucketsDta: make([]*bucket, bucketNum),
		seed:       newSeed(),
		stopCh:     make(chan struct{}),
	}
	LocalCache := &LC{lc}

	for i := 0; i < bucketNum; i++ {
		b := &bucket{
			items: map[string]*Item{},
		}
		lc.bucketsDta[i] = b
	}
	if cleanTime > 0 {
		go lc.clearBucket(cleanTime)
		runtime.SetFinalizer(LocalCache, exitTicker)
	}

	return LocalCache
}

func newSeed() uint32 {
	max := big.NewInt(0).SetUint64(uint64(math.MaxUint32))
	rnd, err := rand.Int(rand.Reader, max)
	var seed uint32
	if err != nil {
		seed = insecurerand.Uint32()
	} else {
		seed = uint32(rnd.Uint64())
	}
	return seed
}

func djb33(seed uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		d = 5381 + seed + l
		i = uint32(0)
	)
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint32(k[i])
			d = (d * 33) ^ uint32(k[i+1])
			d = (d * 33) ^ uint32(k[i+2])
			d = (d * 33) ^ uint32(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint32(k[i])
	case 3:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
	case 4:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
		d = (d * 33) ^ uint32(k[i+2])
	}
	return d ^ (d >> 16)
}

func (l *localCache) getBucket(key string) *bucket {
	//通过hash取模找到对应的桶数据
	return l.bucketsDta[djb33(l.seed, key)%l.bucketNum]
}

func (l *localCache) clearBucket(ct time.Duration) {

	timer := time.NewTicker(ct)
	for {
		select {
		case <-timer.C: //定时清理时间到
			l.delBucketsData()
		case <-l.stopCh: //资源gc被回收后，定时器要退出
			timer.Stop()
			return
		}
	}
}

// 这个淘汰策略 遍历整个数据结构 频繁触发会引起性能抖动；聪明的你一定知道还有更好的策略
func (l *localCache) delBucketsData() {
	for _, v := range l.bucketsDta {
		v.rwMu.RLock()
		if len(v.items) > 0 {
			newMap := make(map[string]*Item)
			for key, val := range v.items {
				if !val.isExpire() {
					newMap[key] = val
				}
			}
			v.items = newMap
		}
		v.rwMu.RUnlock()
	}
}

func exitTicker(l *LC) {
	l.stopCh <- struct{}{}
}

func (l *localCache) Get(key string) (any, bool) {
	return l.getBucket(key).get(key)
}

func (b *bucket) get(key string) (any, bool) {

	b.rwMu.RLock()
	item, found := b.items[key]
	if !found {
		b.rwMu.RUnlock()
		return nil, false
	}

	if item.Expiration != int64(NoExpiration) && item.Expiration < time.Now().Unix() {
		b.rwMu.RUnlock()
		//这里采用懒加载的方式删除
		b.del(key)
		return nil, false
	}
	b.rwMu.RUnlock()
	return item.value, true
}

func (l *localCache) Set(key string, value any, t time.Duration) {

	var et int64
	if t > NoExpiration {
		et = time.Now().Add(t).Unix()
	} else {
		et = int64(t)
	}
	l.getBucket(key).set(key, value, et)
}

func (b *bucket) set(key string, value any, t int64) {

	b.rwMu.Lock()
	defer b.rwMu.Unlock()
	b.items[key] = &Item{value: value, Expiration: t}
}

func (l *localCache) Del(key string) {
	l.getBucket(key).del(key)
}

func (b *bucket) del(key string) {

	b.rwMu.Lock()
	defer b.rwMu.Unlock()
	delete(b.items, key)
}

func (i *Item) isExpire() bool {
	if i.Expiration == int64(NoExpiration) {
		return false
	}
	return i.Expiration < time.Now().Unix()
}
