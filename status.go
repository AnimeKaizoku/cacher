package cacher

type status int

const (
	noop status = iota
	cleanerBusy
	cacherReset
	cacherDeleted
)

// type mutex struct {
// 	*sync.RWMutex
// 	sig statys
// }

// func (m *mutex) iLock() signal {
// 	if m.sig != noop {
// 		m.Lock()

// 	}
// 	return 0
// }
