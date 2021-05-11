package pattern

import "sync"

type SingletonInitFunc func() (interface{}, bool)

type Singleton interface {
	Get() (interface{}, bool)
}

func NewSingleton(initFunc SingletonInitFunc) Singleton {
	return &singletonImpl{
		init: initFunc,
	}
}

type singletonImpl struct {
	sync.Mutex

	instance interface{}
	init SingletonInitFunc
	hasInitialized bool
}

func (s *singletonImpl) Get() (interface{}, bool) {
	if s.hasInitialized {
		return s.instance, true
	}

	s.Lock()
	defer s.Unlock()

	if s.hasInitialized {
		return s.instance, true
	}

	instance, ok := s.init()

	if !ok {
		return nil, false
	}

	s.instance = instance
	s.hasInitialized = true

	return s.instance, true
}