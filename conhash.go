package conhash

import (
	"fmt"
	"hash/fnv"
	"io"
	"sync"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	rbu "github.com/emirpasic/gods/utils"
)

type Conhash struct {
	tree         *rbt.Tree
	lock         sync.RWMutex
	nodes        map[string]bool
	virtualNodes int
}

func New(virtualNodes int) *Conhash {
	c := Conhash{
		tree:         rbt.NewWith(rbu.UInt32Comparator),
		lock:         sync.RWMutex{},
		nodes:        make(map[string]bool),
		virtualNodes: virtualNodes,
	}
	return &c
}

// calculate sum32
func (c *Conhash) sum32(s string) uint32 {
	h := fnv.New32()
	io.WriteString(h, s)
	return h.Sum32()
}

// virtual node name
func (c *Conhash) vname(s string, i int) string {
	return fmt.Sprintf("%s_v%d", s, i)
}

func (c *Conhash) AddNode(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, found := c.nodes[node]
	if found {
		return
	}
	c.nodes[node] = true
	c.tree.Put(c.sum32(node), node)
	for i := 0; i < c.virtualNodes; i++ {
		vnode := c.vname(node, i)
		c.nodes[node] = true
		c.tree.Put(c.sum32(vnode), node)
	}
}

func (c *Conhash) DelNode(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, found := c.nodes[node]
	if !found {
		return
	}
	delete(c.nodes, node)
	c.tree.Remove(c.sum32(node))
	for i := 0; i < c.virtualNodes; i++ {
		vnode := c.vname(node, i)
		delete(c.nodes, vnode)
		c.tree.Remove(c.sum32(vnode))
	}
}

func (c *Conhash) Get(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	h := c.sum32(key)
	nearest, _ := c.tree.Ceiling(h)
	if nearest == nil {
		nearest = c.tree.Left()
	}
	n := nearest.Value.(string)
	return n
}
