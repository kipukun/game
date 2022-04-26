package engine

import (
	"encoding/json"
	xerrors "errors"
	"os"
	"sync"

	"github.com/kipukun/game/engine/errors"
)

// Registry is a global key/value store that states can use to save values.
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

// Save registers and saves v associated to s.
func (r *Registry) Save(s string, v any) error {
	var op errors.Op = "Registry.Save"

	r.mu.Lock()
	r.m[s] = v
	r.mu.Unlock()
	r.mu.RLock()
	defer r.mu.RUnlock()

	bs, err := json.Marshal(r.m)
	if err != nil {
		return errors.Error(op, "error marshaling json", err)
	}
	err = os.WriteFile(r.path, bs, 0600)
	if err != nil {
		return errors.Error(op, "error writing file", err)
	}
	return nil
}

// Get returns the value associated with s, or dv if there was no value found.
func (r *Registry) Get(s string, dv any) any {
	r.mu.RLock()
	v := r.m[s]
	r.mu.RUnlock()
	if v == nil {
		return dv
	}
	return v
}

// Load reads the Registry from disk.
func (r *Registry) Load() error {
	var op errors.Op = "Registry.Load"

	r.mu.RLock()
	bs, err := os.ReadFile(r.path)
	if err != nil {
		if xerrors.Is(err, os.ErrNotExist) {
			bs = []byte("{}")
		} else {
			return errors.Error(op, "error reading file", err)
		}
	}
	r.mu.RUnlock()
	m := make(map[string]any)
	json.Unmarshal(bs, &m)
	r.mu.Lock()
	r.m = m
	r.mu.Unlock()
	return nil
}
