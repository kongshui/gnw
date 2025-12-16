package message

import (
	"context"
	"sync"
)

type (
	sseCancelMap struct {
		sseMap map[string]context.CancelFunc
		lock   sync.RWMutex
	}
)

// getSseCancelMap 获取sseCancelMap
func getSseCancelMap() *sseCancelMap {
	return &sseCancelMap{
		sseMap: make(map[string]context.CancelFunc),
	}
}

// Add 添加
func (s *sseCancelMap) Add(key string, value context.CancelFunc) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.sseMap[key]; ok {
		s.sseMap[key]() // 调用cancel函数
	}
	s.sseMap[key] = value
	return true
}

// len 获取长度
func (s *sseCancelMap) Len() int {
	return len(s.sseMap)
}

// Get 获取cancel
func (s *sseCancelMap) Get(key string) context.CancelFunc {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if _, ok := s.sseMap[key]; !ok {
		return nil
	}
	return s.sseMap[key]
}

// Delete 删除
func (s *sseCancelMap) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// s.sseMap[key]() // 调用cancel函数
	if cancel, ok := s.sseMap[key]; ok {
		cancel()
	}
	delete(s.sseMap, key)
}
