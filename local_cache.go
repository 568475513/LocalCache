package local_cache

// 采用hash分片的链表桶的数据结构 分片桶上锁提高读写效率 初始化长度越大 锁颗粒度越细 内存占用也越大 算是一种以空间换时间逻辑
// 淘汰策略（懒加载+异步线程扫描+lru） gc回收 参考redis过期key淘汰策略
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
	// NoExpiration 默认不过期
	NoExpiration time.Duration = -1
)

type LC struct {
	*localCache
}

type localCache struct {
	bucketsDta []*bucket
	bucketLen  uint32
	itemLen    uint16
	seed       uint32
	stopCh     chan struct{}
}

type bucket struct {
	rwMu       sync.RWMutex
	items      map[string]*Item
	itemLen    uint16
	head, tail *Item //桶链头尾指针地址
}

type Item struct {
	key         string
	value       interface{}
	Expiration  int64
	preI, nextI *Item
}

func NewLocalCache(bucketLen, itemLen int, cleanTime time.Duration) *LC {

	lc := &localCache{
		bucketLen:  uint32(bucketLen), // 定义桶长度
		itemLen:    uint16(itemLen),   // 定义链表长度
		bucketsDta: make([]*bucket, bucketLen),
		seed:       newSeed(),
		stopCh:     make(chan struct{}),
	}
	LocalCache := &LC{lc}

	for i := 0; i < bucketLen; i++ {
		lc.bucketsDta[i] = lc.newBucket()
	}
	if cleanTime > 0 {
		go lc.clearBucket(cleanTime)
		runtime.SetFinalizer(LocalCache, exitTicker)
	}

	return LocalCache
}

// newBucket 初始化一个桶链表
func (l *localCache) newBucket() *bucket {

	b := &bucket{
		items:   make(map[string]*Item, l.itemLen),
		head:    new(Item),
		tail:    new(Item),
		itemLen: l.itemLen,
	}
	b.head.nextI = b.tail
	b.tail.preI = b.head
	return b
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
	return l.bucketsDta[djb33(l.seed, key)%l.bucketLen]
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

// 这个淘汰策略 遍历整个数据结构 频繁触发会引起性能抖动；这里尝试过new map的方式对gc更友好 但需要对整个数据结构上锁 实属没必要了
func (l *localCache) delBucketsData() {
	for _, v := range l.bucketsDta {
		v.rwMu.Lock()
		if len(v.items) > 0 { //旧桶key过期替换逻辑
			for key, val := range v.items {
				if val.isExpire() {
					delete(v.items, key)
					v.removeFromBucket(val)
				}
			}
		}
		v.rwMu.Unlock()
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

	b.moveToHead(item) // 移动至链表头部

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
	item, ok := b.items[key]
	if ok { //更新key
		item.value, item.Expiration = value, t
		b.rwMu.Unlock()
		b.moveToHead(item)
	} else { //新增key
		if len(b.items) >= int(b.itemLen) { //提前判断map容量 防止map扩容
			tailPre := b.tail.preI
			b.removeFromBucket(tailPre)
			delete(b.items, tailPre.key)
		}
		item = &Item{key: key, value: value, Expiration: t}
		b.items[key] = item
		b.moveToBucketHead(item)
		b.rwMu.Unlock()
	}

}

func (l *localCache) Del(key string) {
	l.getBucket(key).del(key)
}

func (b *bucket) del(key string) {

	b.rwMu.Lock()
	defer b.rwMu.Unlock()
	b.removeFromBucket(b.items[key])
	delete(b.items, key)
}

func (i *Item) isExpire() bool {
	if i.Expiration == int64(NoExpiration) {
		return false
	}
	return i.Expiration < time.Now().Unix()
}

func (b *bucket) moveToHead(item *Item) {
	b.rwMu.Lock()
	b.removeFromBucket(item)
	b.moveToBucketHead(item)
	b.rwMu.Unlock()
}

// removeFromBucket 从桶链表移除节点
func (b *bucket) removeFromBucket(item *Item) {
	item.preI.nextI = item.nextI
	item.nextI.preI = item.preI
}

// moveToBucketHead 将节点移至桶链表头部
func (b *bucket) moveToBucketHead(item *Item) {
	item.preI = b.head
	item.nextI = b.head.nextI
	b.head.nextI.preI = item
	b.head.nextI = item
}
