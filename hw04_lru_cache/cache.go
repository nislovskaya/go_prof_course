package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, exists := c.items[key]; exists {
		item.Value.(*cacheItem).value = value
		c.queue.MoveToFront(item)
		return true
	}

	newCacheItem := &cacheItem{key: key, value: value}
	newListItem := c.queue.PushFront(newCacheItem)

	c.items[key] = newListItem

	if c.queue.Len() > c.capacity {
		lastItem := c.queue.Back()
		if lastItem != nil {
			c.queue.Remove(lastItem)
			delete(c.items, lastItem.Value.(*cacheItem).key)
		}
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
