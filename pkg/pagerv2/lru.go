package pagerv2

import (
	"fmt"
	"strings"
)

type entry struct {
	key   uint32
	val   int64
	dirty bool
	prev  *entry
	next  *entry
}

const defaultCacheSize = 8

type lru struct {
	c int               // c is capacity
	m map[uint32]*entry // m is a map of entries
	h *entry            // h is the head of the list
	t *entry            // t is the tail of the list
}

func (l *lru) print() {
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

func newLRU(c int) *lru {
	l := new(lru)
	l.init(c)
	return l
}

func (l *lru) init(c int) {
	if c < 2 {
		c = defaultCacheSize
	}
	l.c = c
	l.m = make(map[uint32]*entry)
	l.h = new(entry)
	l.t = new(entry)
	l.h.next = l.t
	l.t.prev = l.h
}

func (l *lru) evict() *entry {
	e := l.t.prev
	l.pop(e)
	delete(l.m, e.key)
	return e
}

func (l *lru) pop(e *entry) {
	e.prev.next = e.next
	e.next.prev = e.prev
}

func (l *lru) push(e *entry) {
	l.h.next.prev = e
	e.next = l.h.next
	e.prev = l.h
	l.h.next = e
}

func (l *lru) bump(e *entry) {
	l.pop(e)
	l.push(e)
}

func (l *lru) set(k uint32, v int64) {
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
		e.dirty = true
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
	e.dirty = true
	if l.h.next != e {
		l.bump(e)
	}
}

func (l *lru) get(k uint32) (int64, bool) {
	e := l.m[k]
	if e == nil {
		// we do not have the item in the
		// cache, so we simply return nil
		return -1, false
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

func (l *lru) free() (int64, bool) {
	e := l.m[0]
	if e == nil {
		// we do not have the item in the
		// cache, so we simply return nil
		return -1, false
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
