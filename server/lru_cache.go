package server

import (
	"sync"
)

type ElementInterface interface {
	Close()
}

// 定义 LRU 缓存结构体
type LRUCache[K comparable, V ElementInterface] struct {
	capacity int
	cache    map[K]*CacheNode[K, V]
	head     *CacheNode[K, V] // 链表头节点（最近使用的）
	tail     *CacheNode[K, V] // 链表尾节点（最久未使用的）
	mu       sync.Mutex       // 互斥锁
}

// 定义链表节点
type CacheNode[K comparable, V ElementInterface] struct {
	key  K
	val  V
	prev *CacheNode[K, V]
	next *CacheNode[K, V]
}

// 初始化 LRU 缓存
func NewLRUCache[K comparable, V ElementInterface](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*CacheNode[K, V]),
	}
}

// 添加或更新记录
func (lru *LRUCache[K, V]) AddOrUpdate(key K, val V) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	// 如果记录已存在，更新值并移动到头部
	if node, exists := lru.cache[key]; exists {
		node.val = val
		lru.moveToHead(node)
		return
	}

	// 如果记录不存在，创建新节点并添加到头部
	newNode := &CacheNode[K, V]{key: key, val: val}
	lru.cache[key] = newNode
	lru.addToHead(newNode)

	// 如果记录数量超过容量，移除尾部节点
	if len(lru.cache) > lru.capacity {
		node := lru.tail
		lru.removeTail()
		node.val.Close()
	}
}

// 获取记录
func (lru *LRUCache[K, V]) Get(key K) (V, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	node, exists := lru.cache[key]
	if !exists {
		var zeroValue V
		return zeroValue, false
	}

	// 访问记录后，将其移动到头部
	lru.moveToHead(node)
	return node.val, true
}

// 将节点移动到头部
func (lru *LRUCache[K, V]) moveToHead(node *CacheNode[K, V]) {
	lru.removeNode(node)
	lru.addToHead(node)
}

// 将节点添加到头部
func (lru *LRUCache[K, V]) addToHead(node *CacheNode[K, V]) {
	node.prev = nil
	node.next = lru.head

	if lru.head != nil {
		lru.head.prev = node
	}
	lru.head = node

	if lru.tail == nil {
		lru.tail = node
	}
}

// 移除节点
func (lru *LRUCache[K, V]) removeNode(node *CacheNode[K, V]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		lru.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		lru.tail = node.prev
	}
}

// 移除尾部节点
func (lru *LRUCache[K, V]) removeTail() {
	if lru.tail == nil {
		return
	}

	delete(lru.cache, lru.tail.key)
	lru.removeNode(lru.tail)
}
