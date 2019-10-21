package utils

import (
	"sync"
)

type Set struct {
	_map   map[interface{}]interface{}
	rwlock sync.RWMutex
}

func (s *Set) Include(item interface{}) (exists bool) {
	s.rwlock.RLock()
	_, ok := s._map[item]
	s.rwlock.RUnlock()

	return ok
}

func (s *Set) Add(item interface{}) (exists bool) {
	if s.Include(item) {
		return true
	}

	s.rwlock.Lock()
	s._map[item] = struct{}{}
	s.rwlock.Unlock()

	return false
}

func (s *Set) Remove(item interface{}) (exists bool) {
	if !s.Include(item) {
		return false
	}

	s.rwlock.Lock()
	delete(s._map, item)
	s.rwlock.Unlock()

	return true
}

func (s *Set) GetAll() (items []interface{}) {
	s.rwlock.RLock()
	for k, _ := range s._map {
		items = append(items, k)
	}

	s.rwlock.RUnlock()

	return
}

func (s *Set) GetAllInt() (items []int) {
	s.rwlock.RLock()
	for k, _ := range s._map {
		intk, ok := k.(int)
		if ok {
			items = append(items, intk)
		}
	}

	s.rwlock.RUnlock()

	return
}

func (s *Set) GetAllString() (items []string) {
	s.rwlock.RLock()
	for k, _ := range s._map {
		strk, ok := k.(string)
		if ok {
			items = append(items, strk)
		}
	}

	s.rwlock.RUnlock()

	return
}

func NewSet() (h *Set) {
	h = &Set{
		_map: make(map[interface{}]interface{}),
	}

	return
}
