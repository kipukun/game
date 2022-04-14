package engine

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Registry struct {
	mu   sync.RWMutex
	m    map[string]any
	path string
}

func newRegistry(path string) *Registry {
	r := new(Registry)
	r.m = make(map[string]any)
	r.path = path
	return r
}

func (r *Registry) Save(s string, v any) {
	r.mu.Lock()
	r.m[s] = v
	r.mu.Unlock()
	r.mu.RLock()
	defer r.mu.RUnlock()

	bs, err := json.Marshal(r.m)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(r.path, bs, 0600)
	if err != nil {
		panic(err)
	}
}

func (r *Registry) Get(s string) any {
	r.mu.RLock()
	v := r.m[s]
	r.mu.RUnlock()
	return v
}

func (r *Registry) Load() {
	r.mu.RLock()
	bs, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			bs = []byte("{}")
		} else {
			panic(err)
		}
	}
	r.mu.RUnlock()
	m := make(map[string]any)
	json.Unmarshal(bs, &m)
	r.mu.Lock()
	r.m = m
	r.mu.Unlock()
}
