package util

import (
	"log"
	"sort"
	"sync"
)

// closeEvents 关闭事件
var closeEvents map[int]func() (errs []error)

// closedErrorHandler 关闭错误处理事件
var closedErrorHandler func(err error)

var m sync.Mutex

func init() {
	closeEvents = make(map[int]func() (errs []error), 2)
	closedErrorHandler = func(err error) {
		log.Printf("Close Event Trigger Error:%s\n", err)
	}
}

// Inject 关闭时发生错误处理方法
func Inject(handler func(err error)) {
	closedErrorHandler = handler
}

// CloserAdd 添加关闭事件 level 等级越低越先关闭
func CloserAdd(level int, event func() (errs []error)) {
	m.Lock()
	closeEvents[level] = event
	m.Unlock()
}

// Close .
func Close() {
	levels := make([]int, 0, len(closeEvents))
	for k := range closeEvents {
		levels = append(levels, k)
	}
	sort.Ints(levels)

	for i := 0; i < len(levels); i++ {
		event := closeEvents[levels[i]]
		if errs := event(); len(errs) > 0 {
			for j := 0; j < len(errs); j++ {
				closedErrorHandler(errs[j])
			}
		}
	}
}
