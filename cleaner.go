package cacher

import (
	"sync"
	"time"
)

type cleanable interface {
	cleanExpired()
	getCleanInterval() time.Duration
}

type cleaner struct {
	mu            sync.RWMutex
	cachers       []cleanable
	cleanInterval time.Duration
	once          sync.Once
}

func newCleaner() *cleaner {
	return &cleaner{
		cachers: make([]cleanable, 0),
	}
}

// A function to help find GCD of 2 time Durations
func gcd(a, b time.Duration) time.Duration {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func (cl *cleaner) calculateIntervalGCD() {
	if len(cl.cachers) == 0 {
		cl.cleanInterval = 0
		return
	}
	interval := cl.cachers[0].getCleanInterval()
	for _, c := range cl.cachers[1:] {
		interval = gcd(interval, c.getCleanInterval())
	}
	cl.cleanInterval = interval
}

func (cl *cleaner) Register(c cleanable) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	cl.cachers = append(cl.cachers, c)
	cl.calculateIntervalGCD()
	cl.once.Do(cl.Run)
}

func (cl *cleaner) Run() {
	go func() {
		for {
			cl.mu.RLock()
			cachers := append([]cleanable(nil), cl.cachers...)
			interval := cl.cleanInterval
			cl.mu.RUnlock()
			for _, c := range cachers {
				c.cleanExpired()
			}
			time.Sleep(interval)
		}
	}()
}

// A function used by current Cacher instance to clean
// expired keys on regular basis.
func (c *Cacher[C, T]) cleaner() {
	ticker := time.NewTicker(c.cleanInterval)
	defer ticker.Stop()
	for range ticker.C {
		c.cleanExpired()
	}
}
