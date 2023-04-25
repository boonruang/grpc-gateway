package model

import (
	"sync"
)

type ReqResData struct {
	RequestId string
	Method    string
	Route     string
	Body      string
}

var reqIDChan chan string
var reqResDataMap map[string]chan ReqResData
var mu sync.Mutex

// reqIDChan
func InitReqIDChan() func() {
	reqIDChan = make(chan string, 1024)

	return func() {
		close(reqIDChan)
	}
}

func SendReqIDChan(reqID string) {
	reqIDChan <- reqID
}

func ReceiveReqIDChan() chan string {
	return reqIDChan
}

// reqResDataMap
func InitReqResDataMap() func() {
	mu.Lock()
	reqResDataMap = make(map[string]chan ReqResData)
	mu.Unlock()

	return func() {
		mu.Lock()
		defer mu.Unlock()

		for key, channel := range reqResDataMap {
			close(channel)
			delete(reqResDataMap, key)
		}
	}
}

func GetReqResDataMap() map[string]chan ReqResData {
	mu.Lock()
	defer mu.Unlock()

	return reqResDataMap
}

func SendReqDataMapChan(reqID string, data ReqResData) {
	reqID = "req:" + reqID

	mu.Lock()

	if reqResDataMap[reqID] == nil {
		reqResDataMap[reqID] = make(chan ReqResData, 1)
	}
	channel := reqResDataMap[reqID]

	mu.Unlock()

	channel <- data
}

func ReceiveReqDataMapChan(reqID string) (chan ReqResData, func()) {
	reqID = "req:" + reqID

	mu.Lock()

	if reqResDataMap[reqID] == nil {
		reqResDataMap[reqID] = make(chan ReqResData, 1)
	}
	channel := reqResDataMap[reqID]

	mu.Unlock()

	return channel, func() {
		close(channel)

		mu.Lock()
		delete(reqResDataMap, reqID)
		mu.Unlock()
	}
}

func SendResDataMapChan(reqID string, data ReqResData) {
	reqID = "res:" + reqID

	mu.Lock()

	if reqResDataMap[reqID] == nil {
		reqResDataMap[reqID] = make(chan ReqResData, 1)
	}
	channel := reqResDataMap[reqID]

	mu.Unlock()

	channel <- data
}

func ReceiveResDataMapChan(reqID string) (chan ReqResData, func()) {
	reqID = "res:" + reqID

	mu.Lock()

	if reqResDataMap[reqID] == nil {
		reqResDataMap[reqID] = make(chan ReqResData, 1)
	}
	channel := reqResDataMap[reqID]

	mu.Unlock()

	return channel, func() {
		close(channel)

		mu.Lock()
		delete(reqResDataMap, reqID)
		mu.Unlock()
	}
}

func IsReqResDataMapAvailable(prefix string, reqID string) bool {
	reqID = prefix + ":" + reqID

	mu.Lock()
	defer mu.Unlock()

	if reqResDataMap[reqID] == nil {
		return false
	}

	return true
}
