package pager

// Cache is a page cache. It is responsible for caching pages of data from
// pages on disk to have more fine-grained control of disk IO operations.
type Cache struct {
	pageCache [cacheSize]byte
}

func NewCache() *Cache {
	return &Cache{}
}
