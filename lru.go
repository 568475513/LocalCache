package local_cache

import (
	"sync"
	"time"
)

// Cache is an LRU cache. It is not safe for concurrent access.
type LruNode struct {
	Prev    *LruNode
	Next    *LruNode
	Timeout int64
	Value   interface{}
	Key     string
}

type LruLocalCache struct {
	maxSize int
	size    int
	head    *LruNode
	tail    *LruNode
	cache   map[interface{}]*LruNode
	locker  *sync.RWMutex
}

func NewLruLocalCache(maxSize int) *LruLocalCache {
	return &LruLocalCache{
		maxSize: maxSize,
		size:    0,
		head:    nil,
		tail:    nil,
		cache:   make(map[interface{}]*LruNode),
		locker:  new(sync.RWMutex),
	}
}

// Put adds a value to the cache.expire 有效期,单位秒
func (c *LruLocalCache) Put(key string, val interface{}, expire time.Duration) error {

	// 确认容量是否超出,如果超出了则进行清理

	locker := c.locker

	locker.Lock()
	defer locker.Unlock()
	c.ifFullRemoveLast()

	ts := time.Now().Add(expire).Unix()

	node := &LruNode{
		Prev:    nil,
		Next:    nil,
		Timeout: ts,
		Value:   val,
		Key:     key,
	}
	c.addToHead(node)
	return nil
}

func (c *LruLocalCache) Get(key string) (interface{}, bool) {
	locker := c.locker

	locker.RLock()
	node, ok := c.cache[key]
	if !ok {
		locker.RUnlock()
		return nil, false
	}
	locker.RUnlock()

	locker.Lock()
	defer locker.Unlock()
	if node.Timeout < time.Now().Unix() {
		c.delete(node)
		return nil, false
	}
	c.moveToHead(node)
	return node.Value, true
}

func (c *LruLocalCache) delete(node *LruNode) {
	if node == nil {
		return
	}

	key := node.Key

	if c.head == node {
		c.head = node.Next
	}
	if c.tail == node {
		c.tail = node.Prev
	}

	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}

	node.Next, node.Prev = nil, nil

	delete(c.cache, key)
	c.size -= 1
}

func (c *LruLocalCache) moveToHead(node *LruNode) {
	if node == nil {
		return
	}
	// 已经在队列头了
	if c.head == node {
		return
	}
	// 在队列尾了
	if c.tail == node {
		c.tail = node.Prev
	}

	//
	prev := node.Prev
	if prev != nil {
		prev.Next = node.Next
	}

	head := c.head
	node.Prev, node.Next = nil, head
	head.Prev = node
	c.head = node
	return
}

func (c *LruLocalCache) addToHead(node *LruNode) {
	if node == nil {
		return
	}

	head := c.head

	if head == nil {
		c.head, c.tail = node, node
	} else {

		node.Next = head
		node.Prev = nil

		head.Prev = node

		c.head = node
	}
	c.cache[node.Key] = node
	c.size++
}

func (c *LruLocalCache) ifFullRemoveLast() {
	if c.size+1 > c.maxSize {
		c.delete(c.tail)
	}
}
