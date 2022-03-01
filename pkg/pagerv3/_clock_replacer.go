package pagerv3

import (
	"fmt"
)

// FrameID is the type for frame id
type FrameID int

// ClockReplacer represents the clock replacer algorithm
type ClockReplacer struct {
	cList     *ringBuff
	clockHand int
}

// Victim removes the victim frame as defined by the replacement policy
func (c *ClockReplacer) Victim() *FrameID {
	if c.cList.size == 0 {
		return nil
	}
	var victimFrameID *FrameID
	currentNode := c.clockHand
	for {

		if currentNode.value.(bool) {
			currentNode.value = false
			c.clockHand = &currentNode.next
		} else {
			frameID := currentNode.key.(FrameID)
			victimFrameID = &frameID

			c.clockHand = &currentNode.next

			c.cList.remove(currentNode.key)
			return victimFrameID
		}
	}
}

// Unpin unpins a frame, indicating that it can now be victimized
func (c *ClockReplacer) Unpin(id FrameID) {
	if !c.cList.hasKey(id) {
		c.cList.insert(id)
		if c.cList.size == 1 {
			c.clockHand = &c.cList.head
		}
	}
}

// Pin pins a frame, indicating that it should not be victimized until it is unpinned
func (c *ClockReplacer) Pin(id FrameID) {
	node := c.cList.find(id)
	if node == nil {
		return
	}

	if (*c.clockHand) == node {
		c.clockHand = &(*c.clockHand).next
	}
	c.cList.remove(id)

}

// Size returns the size of the clock
func (c *ClockReplacer) Size() int {
	return c.cList.size
}

// NewClockReplacer instantiates a new clock replacer
func NewClockReplacer(poolSize int) *ClockReplacer {
	cList := newCircularList(poolSize)
	return &ClockReplacer{cList, &cList.head}
}

type ringBuff struct {
	buf      []FrameID
	beg, end int
	size     int
}

func newRingBuff(size int) *ringBuff {
	return &ringBuff{
		buf:  make([]FrameID, size),
		size: size,
	}
}

func (r *ringBuff) insert(id FrameID) {
	r.buf[r.end] = id
	r.end++
	r.end %= r.size
}

func (r *ringBuff) get() FrameID {
	v := r.buf[r.beg]
	r.beg++
	r.beg %= r.size
	return v
}

func (r *ringBuff) remove(id FrameID) {
	n := r.find(id)
	if n == -1 {
		return
	}
	if n == r.beg {
		r.buf[n] = 0
		r.beg++
		r.beg %= r.size
	}
	if n == r.end {
		r.buf[n] = 0
		r.end--
		r.end %= r.size
	}
}

// find does a linear search and attempts to return the index at which the
// provided FrameID is located, otherwise returning -1 if it is not found
func (r *ringBuff) find(id FrameID) int {
	for i := r.beg; i < r.end; i++ {
		if r.buf[i] == id {
			return int(r.buf[i])
		}
	}
	return -1
}

func (r *ringBuff) hasKey(id FrameID) bool {
	return r.find(id) != -1
}

func (r *ringBuff) String() string {
	return fmt.Sprintf("(size=%d, beg=%d, end=%d) %+v", r.size, r.beg, r.end, r.buf)
}
