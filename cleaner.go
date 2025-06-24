package cacher

import (
	"time"
)

type Cleanable interface {
	cleanExpired()
	getCleanInterval() time.Duration
}

type Cleaner struct {
	cachers       []Cleanable
	cleanInterval time.Duration
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		cachers: make([]Cleanable, 0),
	}
}

// A function to help find GCD of 2 time Durations
func gcd(a, b time.Duration) time.Duration {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func (cl *Cleaner) calculateIntervalGCD() {
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

func (cl *Cleaner) Register(c Cleanable) {
	cl.cachers = append(cl.cachers, c)
	cl.calculateIntervalGCD()
}

func (cl *Cleaner) Run() {
	go func() {
		for {
			for _, c := range cl.cachers {
				c.cleanExpired()
			}
			time.Sleep(cl.cleanInterval)
		}
	}()
}

// A function used by current Cacher instance to clean
// expired keys on regular basis.
func (c *Cacher[C, T]) cleaner() {
	for {
		c.cleanExpired()
		// cleanup expired keys every c.cleanInterval duration
		time.Sleep(c.cleanInterval)
	}
}
