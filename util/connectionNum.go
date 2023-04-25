package util

import "sync"

var (
	connectionNum int16
	mu            sync.Mutex
)

func AddStreamNum() {
	mu.Lock()
	defer mu.Unlock()

	connectionNum += 1
}

func RemoveStreamNum() {
	mu.Lock()
	defer mu.Unlock()

	connectionNum -= 1
}

func IsStreamConnect() bool {
	mu.Lock()
	defer mu.Unlock()

	if connectionNum > 0 {
		return true
	}

	return false
}
