package server

import (
	"eino-script/engine"
	"sync"
)

// 定义 LRU 缓存结构体
type LRUCache struct {
	capacity int
	cache    map[string]*CacheNode
	head     *CacheNode // 链表头节点（最近使用的）
	tail     *CacheNode // 链表尾节点（最久未使用的）
	mu       sync.Mutex // 互斥锁
}

// 定义链表节点
type CacheNode struct {
	key  string
	val  *engine.Engine
	prev *CacheNode
	next *CacheNode
}

// 初始化 LRU 缓存
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*CacheNode),
	}
}

// 添加或更新记录
func (lru *LRUCache) AddOrUpdate(val *engine.Engine) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	key := val.Id()

	// 如果记录已存在，更新值并移动到头部
	if node, exists := lru.cache[key]; exists {
		node.val = val
		lru.moveToHead(node)
		return
	}

	// 如果记录不存在，创建新节点并添加到头部
	newNode := &CacheNode{key: key, val: val}
	lru.cache[key] = newNode
	lru.addToHead(newNode)

	// 如果记录数量超过容量，移除尾部节点
	if len(lru.cache) > lru.capacity {
		lru.removeTail()
	}
}

// 获取记录
func (lru *LRUCache) Get(key string) (*engine.Engine, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	node, exists := lru.cache[key]
	if !exists {
		return nil, false
	}

	// 访问记录后，将其移动到头部
	lru.moveToHead(node)
	return node.val, true
}

// 将节点移动到头部
func (lru *LRUCache) moveToHead(node *CacheNode) {
	lru.removeNode(node)
	lru.addToHead(node)
}

// 将节点添加到头部
func (lru *LRUCache) addToHead(node *CacheNode) {
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
func (lru *LRUCache) removeNode(node *CacheNode) {
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
func (lru *LRUCache) removeTail() {
	if lru.tail == nil {
		return
	}

	delete(lru.cache, lru.tail.key)
	lru.removeNode(lru.tail)
}
