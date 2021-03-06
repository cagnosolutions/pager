package _pager

import (
	"fmt"
	"strings"
)

type keyType int
type valType int

type entry struct {
	key  keyType
	val  valType
	prev *entry
	next *entry
}

const defaultCacheSize = 8

type LRU struct {
	c int                // c is capacity
	m map[keyType]*entry // m is a map of entries
	h *entry             // h is the head of the list
	t *entry             // t is the tail of the list
	f *entry             // f is a free evicted entry
}

func (l *LRU) print() {
	var ss []string
	if h := l.h; h != nil {
		e := h.next
		for e != l.t {
			ss = append(ss, fmt.Sprintf("e{%d, %d}", e.key, e.val))
			e = e.next
		}
	}
	fmt.Println(strings.Join(ss, ","))
}

func NewLRU(c int) *LRU {
	if c < 2 {
		c = defaultCacheSize
	}
	l := &LRU{
		c: c,
		m: make(map[keyType]*entry),
		h: new(entry),
		t: new(entry),
		f: nil,
	}
	l.h.next = l.t
	l.t.prev = l.h
	return l
}

func (l *LRU) evict() *entry {
	e := l.t.prev
	l.pop(e)
	delete(l.m, e.key)
	return e
}

func (l *LRU) pop(e *entry) {
	e.prev.next = e.next
	e.next.prev = e.prev
}

func (l *LRU) push(e *entry) {
	l.h.next.prev = e
	e.next = l.h.next
	e.prev = l.h
	l.h.next = e
}

func (l *LRU) bump(e *entry) {
	l.pop(e)
	l.push(e)
}

func (l *LRU) Set(k keyType, v valType) {
	e := l.m[k]
	if e == nil {
		// if we are at capacity
		if len(l.m) == l.c {
			// we are at capacity, we must
			// evict an entry, and then we
			// can proceed with updating
			// the cache
			e = l.evict()
		} else {
			// we are not at capacity but
			// the entry is just empty so,
			// we need to make a new entry
			// then we can proceed with
			// updating the cache
			e = new(entry)
		}
		// we did what we needed to, so now
		// we simply update the cache entry
		e.key = k
		e.val = v
		l.push(e)
		l.m[k] = e
		// and then return
		return
	}
	// otherwise, the entry must exist and
	// being at capacity does not matter
	// because we are simply doing an update
	// of the cache entry in place
	e.val = v
	if l.h.next != e {
		l.bump(e)
	}
}

func (l *LRU) Get(k keyType) (valType, bool) {
	e := l.m[k]
	if e == nil {
		// we do not have the item in the
		// cache, so we simply return nil
		return *new(valType), false
	}
	// otherwise, we must have the value
	// in the cache, so let's return it
	// but not before bump it
	if l.h.next != e {
		l.bump(e)
	}
	// and finally we return the value
	return e.val, true
}
