package utils

import (
	//   "encoding/json"
	"sync"
	//   "time"
)

var ()

// 实现了一个固定大小、线程安全的环状slice。
// 容量满了以后，添加新元素就会替代最旧的元素，保证长度不变且元素有序。

type RingSlice struct {
	offset int           // 目前的游标
	s      []interface{} // 存储元素

	lock sync.RWMutex
}

func (s *RingSlice) Add(element interface{}) {
	if element == nil {
		return
	}

	s.lock.Lock()
	if s.offset == cap(s.s) {
		s.offset = 0
	}

	s.s[s.offset] = element
	s.offset += 1

	s.lock.Unlock()
}

func (s *RingSlice) Clean() {
	s.lock.Lock()

	s.offset = 0
	for i, _ := range s.s {
		s.s[i] = nil
	}

	s.lock.Unlock()
}

func (s *RingSlice) GetAll() (slice []interface{}) {
	s.lock.RLock()

	capacity := cap(s.s)

	slice = make([]interface{}, 0, capacity)
	slice = append(slice, s.s[capacity:]...)
	slice = append(slice, s.s[:capacity]...)
	slice = slice[:s.Len()]

	s.lock.RUnlock()

	return
}

func (s *RingSlice) Len() (l int) {
	s.lock.RLock()
	if s.offset == cap(s.s) || s.s[s.offset] != nil {
		l = cap(s.s)
	} else {
		l = s.offset
	}
	s.lock.RUnlock()

	return l
}

func (s *RingSlice) Cap() (c int) {
	s.lock.RLock()
	c = cap(s.s)
	s.lock.RUnlock()

	return
}

func (s *RingSlice) Offset() (o int) {
	s.lock.RLock()
	o = s.offset
	s.lock.RUnlock()

	return
}

func NewRingSlice(capacity int) (s *RingSlice) {
	s = &RingSlice{
		s: make([]interface{}, capacity),
	}

	return
}
