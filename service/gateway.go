package service

import (
	"errors"
	"time"

	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/model"
	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/util"
)

type GatewayService struct {
	SendingRetry       int
	SendingInterval    time.Duration
	ConnectingRetry    int
	ConnectingInterval time.Duration
}

func (s GatewayService) SendGRPC(requestID string, data model.ReqResData) (model.ReqResData, error) {
	var closeChannelFunc func()
	defer func() {
		if closeChannelFunc != nil {
			closeChannelFunc()
		}
	}()

loop:
	for i := 0; i < s.SendingRetry; i++ {
		for j := 0; j < s.ConnectingRetry; j++ {
			if util.IsStreamConnect() {
				model.SendReqDataMapChan(requestID, data)

				model.SendReqIDChan(requestID)

				resChan, closeFunc := model.ReceiveResDataMapChan(requestID)
				closeChannelFunc = closeFunc

				select {
				case resData := <-resChan:
					return resData, nil
				case <-time.After(s.SendingInterval):
					continue loop
				}
			}

			<-time.After(s.ConnectingInterval)
		}

		return model.ReqResData{}, errors.New("local api is not available")
	}

	return model.ReqResData{}, errors.New("gateway timeout")
}
