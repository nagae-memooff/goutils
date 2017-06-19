package utils

import (
	"encoding/json"
	"sync"
	//   "time"
)

var ()

type HashTable struct {
	_map   map[string]interface{}
	rwlock sync.RWMutex
}

func NewHashTable() (h *HashTable) {
	h = &HashTable{
		_map: make(map[string]interface{}),
	}

	return
}

func (h *HashTable) Get(key string) (value string) {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()

	value, _ = h._map[key].(string)

	return
}

func (h *HashTable) Set(key string, value interface{}) *HashTable {
	h.rwlock.Lock()
	defer h.rwlock.Unlock()

	h._map[key] = value

	return h
}

func (h *HashTable) GetBool(key string) (value bool) {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()

	value, _ = h._map[key].(bool)

	return
}

func (h *HashTable) GetInt(key string) (value int) {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()

	switch n := h._map[key].(type) {
	case int:
		value = n
	case float64:
		value = int(n)
	}

	return
}

func (h *HashTable) GetFloat64(key string) (value float64) {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()

	switch n := h._map[key].(type) {
	case int:
		value = float64(n)
	case float64:
		value = n
	}

	return
}

func (h *HashTable) GetChild(key string) (value *HashTable) {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()

	value, _ = h._map[key].(*HashTable)

	return
}

func (h *HashTable) GetInterface(key string) (value interface{}) {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()

	value, _ = h._map[key]

	return
}

func (h *HashTable) LoadFromJson(json_bytes []byte) (err error) {
	h.rwlock.Lock()
	defer h.rwlock.Unlock()
	if h._map == nil {
		h._map = make(map[string]interface{})
	}

	err = json.Unmarshal(json_bytes, h._map)
	return
}

func (h *HashTable) ToJson() (json_bytes []byte, err error) {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()
	if h._map == nil {
		h._map = make(map[string]interface{})
	}

	json_bytes, err = json.Marshal(h._map)
	return
}

func (h *HashTable) GetMap() map[string]interface{} {
	return h._map
}

func (h *HashTable) SetMap(_map map[string]interface{}) {
	h._map = _map
}
