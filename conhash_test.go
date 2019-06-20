package conhash

import (
	"fmt"
	"testing"
)

type Dist struct {
	m map[string][]string // node -> keys
	k map[string]string   // key -> node
	c map[string]int      // node -> key count
}

func distribution(t *testing.T, c *Conhash) Dist {
	dist := Dist{
		m: make(map[string][]string),
		k: make(map[string]string),
		c: make(map[string]int),
	}
	count := 10000
	for i := 0; i < count; i++ {
		k := fmt.Sprintf("%d", i)
		node := c.Get(k)
		dist.m[node] = append(dist.m[node], k)
		dist.k[k] = node
		dist.c[node]++
	}
	for _, c := range dist.c {
		count -= c
	}
	if count != 0 {
		t.Fatal()
	}
	return dist
}

func TestConhash(t *testing.T) {
	c := New(100)
	nodes := []string{"A", "B", "C"}
	for _, node := range nodes {
		c.AddNode(node)
	}
	d1 := distribution(t, c)
	fmt.Println(d1.c)

	c.AddNode("D")
	d2 := distribution(t, c)
	fmt.Println(d2.c)

	c.DelNode("D")
	d3 := distribution(t, c)
	fmt.Println(d3.c)

	c.DelNode("C")
	d4 := distribution(t, c)
	fmt.Println(d4.c)

	c.DelNode("B")
	d5 := distribution(t, c)
	fmt.Println(d5.c)
}
