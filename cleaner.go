package cacher

import (
	"time"
)

// A function used by current Cacher instance to clean
// expired keys on regular basis.
func (c *Cacher[C, T]) cleaner() {
	for {
		currTime := time.Now().Unix()
		c.mutex.Lock()
		for key, val := range c.cacheMap {
			// Skip the current clean window if cacher is reset or deleted.
			if c.status == cacherReset || c.status == cacherDeleted {
				c.status = noop
				break
			}
			if val.expiry <= currTime {
				delete(c.cacheMap, key)
			}
		}
		c.mutex.Unlock()
		// cleanup expired keys every c.cleanInterval duration
		time.Sleep(c.cleanInterval)
	}
}
