package pattern

import "sync"

type SingletonInitFunc func() (interface{}, bool)

type Singleton interface {
	Get() (interface{}, bool)
}

type SingletonSettings struct {
	OnInit SingletonInitFunc
	DisableMutex bool
}

func NewSingleton(settings SingletonSettings) Singleton {
	return &singletonImpl{
		init: settings.OnInit,
		hasDisableMutex: settings.DisableMutex,
	}
}

type singletonImpl struct {
	mutex sync.Mutex

	hasDisableMutex bool
	instance interface{}
	init SingletonInitFunc
	hasInitialized bool
}

func (s *singletonImpl) Lock() {
	if s.hasDisableMutex {
		return
	}

	s.mutex.Lock()
}

func (s *singletonImpl) Unlock() {
	if s.hasDisableMutex {
		return
	}

	s.mutex.Unlock()
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