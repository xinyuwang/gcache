package gcache

import (
	"fmt"
	"testing"
	"time"
)

type valNode struct {
	key string
	val int
}

func newValNode(k string, v int) *valNode {
	return &valNode{
		key: k,
		val: v,
	}
}

func createSampleCache(v int) Cache[valNode] {

	c := New[valNode]()

	for i := 0; i < v; i++ {

		key := fmt.Sprintf("%d", i)
		c.Add(key, newValNode(key, i))
	}

	return c
}

func TestNew(t *testing.T) {

	num := 100
	c := createSampleCache(num)
	t.Logf("Create sample cache with num: %d", num)
	if c.Len() != num {
		t.Errorf("Create sample cache but len invalid: %d", c.Len())
	}

	if v := c.Get("5"); v != nil {
		t.Logf("c.Get(5) = [%v]", v.val)
	}

	if v := c.Get("105"); v != nil {
		t.Errorf("c.Get return not nil")
	}

	c.Set("5", newValNode("200", 200))
	if v := c.Get("5"); v != nil {
		t.Logf("c.Get(5) after set to 200: %d", v.val)
	}

}

func TestAdd(t *testing.T) {

	num := 1000000
	arr := make([]*valNode, num)
	for i := 0; i < num; i++ {

		key := fmt.Sprintf("%d", i)
		arr[i] = newValNode(key, i)
	}

	t.Logf("Finish create #%d valNode", num)

	timeDiff()
	c := New[valNode]()
	for _, v := range arr {
		if err := c.Add(v.key, v); err != nil {
			t.Logf("c.Add(%s, %v) error: %v", v.key, v, err)
		}
	}
	tc := timeDiff()

	t.Logf("Add #%d cost: %d ms", num, tc)
}

func TestGet(t *testing.T) {

	num := 1000000
	arr := make([]*valNode, num)
	arrKey := make([]string, num)
	for i := 0; i < num; i++ {
		arrKey[i] = fmt.Sprintf("%d", i)
		arr[i] = newValNode(arrKey[i], i)
	}

	t.Logf("Finish create #%d valNode", num)

	c := New[valNode]()
	for _, v := range arr {
		if err := c.Add(v.key, v); err != nil {
			t.Logf("c.Add(%s, %v) error: %v", v.key, v, err)
		}
	}

	timeDiff()
	for i := 0; i < num; i++ {
		c.Get(arrKey[i])
	}
	tc := timeDiff()

	t.Logf("Get #%d cost: %d ms", num, tc)

}

func TestSet(t *testing.T) {

	num := 1000000
	arr := make([]*valNode, num)
	for i := 0; i < num; i++ {

		key := fmt.Sprintf("%d", i)
		arr[i] = newValNode(key, i)
	}

	t.Logf("Finish create #%d valNode", num)

	timeDiff()
	c := New[valNode]()
	for _, v := range arr {
		c.Set(v.key, v)
	}
	tc := timeDiff()

	t.Logf("Set #%d cost: %d ms", num, tc)
}

func TestDel(t *testing.T) {
	num := 1000000
	arr := make([]*valNode, num)
	arrKey := make([]string, num)
	for i := 0; i < num; i++ {
		arrKey[i] = fmt.Sprintf("%d", i)
		arr[i] = newValNode(arrKey[i], i)
	}

	t.Logf("Finish create #%d valNode", num)

	c := New[valNode]()
	for _, v := range arr {
		if err := c.Add(v.key, v); err != nil {
			t.Logf("c.Add(%s, %v) error: %v", v.key, v, err)
		}
	}

	timeDiff()
	for i := 0; i < num; i++ {
		c.Del(arrKey[i])
	}
	tc := timeDiff()

	t.Logf("Del #%d cost: %d ms", num, tc)
}

var lastTimeUs int64

func timeDiff() int64 {
	ts := time.Now().UnixNano() / 1e6
	var ret int64

	if lastTimeUs == 0 {
		lastTimeUs = ts
		ret = ts

	} else {
		ret = ts - lastTimeUs
		lastTimeUs = ts
	}

	return ret
}
