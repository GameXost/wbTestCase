package cache

import (
	"github.com/GameXost/wbTestCase/models"
	"log"
	"sync"
)

type Cache struct {
	mu       sync.Mutex
	data     map[string]*Node
	capacity uint64
	size     uint64
	head     *Node
	tail     *Node
}

type Node struct {
	order *models.Order
	prev  *Node
	next  *Node
	key   string
}

func NewCache(capacity uint64) *Cache {
	return &Cache{
		data:     make(map[string]*Node, capacity),
		capacity: capacity,
	}
}

func (c *Cache) Get(key string) (*models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	node, has := c.data[key]
	if !has {
		log.Println("cache miss")
		return nil, false
	}
	log.Println("cache hit")
	c.moveToTop(node)
	return node.order, true
}

func (c *Cache) Set(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	node, has := c.data[order.OrderUId]
	if has {
		node.order = order
		c.moveToTop(node)
		return
	}

	if c.size >= c.capacity {
		c.deleteBottom()
	}

	node = &Node{key: order.OrderUId, order: order}
	c.addToFront(node)
}

func (c *Cache) LoadFull(ids []*models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, v := range ids {
		if c.size >= c.capacity {
			break
		}
		node := &Node{order: v, key: v.OrderUId}
		c.addToFront(node)
	}
}

func (c *Cache) moveToTop(node *Node) {
	//если у нас элемент вверху - скип
	if node == c.head {
		return
	}
	//предыдущего элемента нет только у первого элемента
	if node.prev != nil {
		node.prev.next = node.next
	}
	// если элемент - хвост, то нужно сдвинуть указатель хвоста
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = c.tail.prev
	}
	//текущий элемент становится головой, т.е предыдущего нет, следующий - бывшая голова
	node.prev = nil
	node.next = c.head
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

}

func (c *Cache) deleteBottom() {
	if c.tail == nil {
		return
	}
	removed := c.tail
	if c.head == c.tail {
		c.tail = nil
		c.head = nil
	} else {
		c.tail = c.tail.prev
		c.tail.next = nil
	}
	delete(c.data, removed.key)
	c.size--

}

func (c *Cache) addToFront(node *Node) {
	node.next = c.head
	node.prev = nil
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = node
	}
	c.size++
	c.data[node.key] = node
}
