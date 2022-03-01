package pagerv3

import (
	"fmt"
)

const (
	opGet = iota
	opAdd
)

type cQueue struct {
	buf []int // buf is the circular array buffer
	ptr int   // ptr is the current operating slot
	opr int   // opr stores the last operation
	cnt int   // cnt is number of entries
}

func NewCQueue(size int) *cQueue {
	return newCQueue(size)
}

func (c *cQueue) Insert(n int) {
	c.add(n)
}

func (c *cQueue) Remove(n int) {
	c.remove(n)
	c.decrptr()
	for c.peekRight() != -1 {
		c.swapRight()
	}
	c.incrptr()
}

func newCQueue(size int) *cQueue {
	c := &cQueue{
		buf: make([]int, size),
		ptr: 0,
		opr: opAdd,
		cnt: 0,
	}
	for i := range c.buf {
		c.buf[i] = -1
	}
	return c
}

func (c *cQueue) has(n int) bool {
	if c.cnt == 0 {
		return false
	}
	if c.cnt == 1 {
		return c.buf[c.ptr] == n
	}
	var was int
	was = c.ptr
	for i := 0; i < c.cnt; i++ {
		if c.buf[c.ptr] == n {
			c.ptr = was
			return true
		}
		c.next()
	}
	return false
}

func (c *cQueue) add(n int) {
	c.buf[c.ptr] = n
	c.incrptr()
	c.opr = opAdd
	if c.cnt < len(c.buf) {
		c.cnt++
	}
}

func (c *cQueue) swapLeft() {
	var cn, on int
	cn = c.buf[c.ptr]
	c.decrptr()
	on = c.buf[c.ptr]
	c.buf[c.ptr] = cn
	c.incrptr()
	c.buf[c.ptr] = on
}

func (c *cQueue) swapRight() {
	var cn, on int
	cn = c.buf[c.ptr]
	c.incrptr()
	on = c.buf[c.ptr]
	c.buf[c.ptr] = cn
	c.decrptr()
	c.buf[c.ptr] = on
}

func (c *cQueue) peekRight() int {
	c.incrptr()
	v := c.buf[c.ptr]
	c.decrptr()
	return v
}

func (c *cQueue) peekLeft() int {
	c.decrptr()
	v := c.buf[c.ptr]
	c.incrptr()
	return v
}

func (c *cQueue) get() int {
	v := c.buf[c.ptr]
	c.buf[c.ptr] = -1
	c.incrptr()
	c.opr = opGet
	c.cnt--
	return v
}

func (c *cQueue) remove(n int) {
	if c.cnt == 0 {
		return
	}
	if c.cnt == 1 {
		if c.buf[c.ptr] == n {
			_ = c.get()
			return
		}
	}
	// var cur int
	// cur = c.ptr
	for i := 0; i < c.cnt; i++ {
		if c.buf[c.ptr] == n {
			_ = c.get()
			// c.ptr = cur
			return
		}
		c.next()
	}
	return
}

func (c *cQueue) next() {
	if c.cnt < 2 {
		return
	}
	c.incrptr()
	for c.buf[c.ptr] == -1 {
		c.incrptr()
	}
}

func (c *cQueue) prev() {
	if c.cnt < 2 {
		return
	}
	c.decrptr()
	for c.buf[c.ptr] == -1 {
		c.decrptr()
	}
}

func (c *cQueue) nextEmpty() {
	if c.cnt == 0 {
		return
	}
	c.incrptr()
	for c.buf[c.ptr] != -1 {
		c.incrptr()
	}
	c.opr = opAdd
}

func (c *cQueue) incrptr() {
	c.ptr++
	c.ptr %= len(c.buf)
}

func (c *cQueue) decrptr() {
	if c.ptr == 0 {
		c.ptr = len(c.buf) - 1
	} else {
		c.ptr--
	}
	c.ptr %= len(c.buf)
}

type printer struct {
	mark int
	c    *cQueue
}

func NewCQPrinter(c *cQueue) *printer {
	return newPrinter(c)
}

func newPrinter(c *cQueue) *printer {
	return &printer{
		mark: 0,
		c:    c,
	}
}

func (p *printer) Print() {
	var aa, bb string
	for i := 0; i < len(p.c.buf); i++ {
		aa += fmt.Sprintf(" %d ", p.c.buf[i])
		if p.c.ptr == i {
			bb += fmt.Sprint(" ^ ")
		} else {
			bb += fmt.Sprint("   ")
		}
	}
	fmt.Printf("%.3d: [%s] (cnt=%d, ptr=%d, opr=%d)\n%.3d:  %s \n", p.mark, aa, p.c.cnt, p.c.ptr, p.c.opr, p.mark, bb)
	p.mark++
}
