package controls

import "sync"

type MVar struct {
	value interface{}
	sync.RWMutex
}

func (mv *MVar) Set(v interface{}) {
	mv.Lock()
	defer mv.Unlock()

	mv.value = v
}

func (mv *MVar) Get() interface{} {
	mv.RLock()
	defer mv.RUnlock()

	return mv.value
}
