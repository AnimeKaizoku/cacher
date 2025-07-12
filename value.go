package cacher

type value[T any] struct {
	// expiry int64
	val T
	evictibleValue
}

// revaluate, ttl
func (v *value[T]) get() (value T, expired bool) {
	if v.isExpired(false) {
		expired = true
		return
	}
	value = v.val
	return
}

func (v *value[T]) getWithoutExpiry() (value T) {
	value = v.val
	return
}
