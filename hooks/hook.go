package hooks

import (
	"unsafe"
)

type HookFunc func(args ...interface{})
type HookFuncList []HookFunc
type EventHookFuncList map[int]HookFuncList

type Hook struct {
	execFunc EventHookFuncList
}

func (h *Hook) AddHook(event int, hook HookFunc) {
	funcList, ok := h.execFunc[event]

	if !ok {
		h.execFunc[event] = HookFuncList{}
		funcList = h.execFunc[event]
	}

	funcAddress := GetFuncAddress(hook)
	funcList[funcAddress] = hook
}

func (h *Hook) AddHooks(event int, hooks HookFuncList) {
	_, ok := h.execFunc[event]

	if !ok {
		h.execFunc[event] = HookFuncList{}
	}

	for _, hook := range hooks {
		h.execFunc[event] = append(h.execFunc[event], hook)
	}
}

func (h *Hook) RemoveHookByEvent(event int) {
	_, ok := h.execFunc[event]

	if !ok {
		return
	}

	delete(h.execFunc, event)
}

func (h *Hook) ExecuteFunc(event int, args ...interface{}) {
	funcList, ok := h.execFunc[event]

	if !ok {
		return
	}

	for _, execFunc := range funcList {
		execFunc(args...)
	}
}

func NewHook() *Hook {
	return &Hook{
		execFunc: make(EventHookFuncList, 0),
	}
}

func GetFuncAddress(i interface{}) int {
	return *(*int)(unsafe.Pointer(&i))
}