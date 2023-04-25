package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/model"
	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func StartGatewayServer(cfg model.AppGrpc) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := GatewayServerImpl{
		Token: cfg.Key,
	}

	grpcServer := grpc.NewServer()
	model.RegisterGatewayServer(grpcServer, &server)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
}

type GatewayServerImpl struct {
	Token string
}

func (s *GatewayServerImpl) Stream(stream model.Gateway_StreamServer) error {
	clientToken := util.GetSteamHeader(stream.Context(), "authorization")
	if clientToken != s.Token {
		errMsg := "invalid authorization"
		util.PrintErrorLog("connect stream", errMsg)

		return status.Errorf(codes.PermissionDenied,
			errMsg)
	}

	util.AddStreamNum()
	util.PrintSuccessLog("connect stream", "client connected")

	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case requestID := <-model.ReceiveReqIDChan():
				go func(requestID string) {
					reqChan, closeFunc := model.ReceiveReqDataMapChan(requestID)
					request := <-reqChan
					closeFunc()

					err := stream.Send(&model.Request{
						RequestId: request.RequestId,
						Method:    request.Method,
						Route:     request.Route,
						Body:      request.Body,
					})
					if err != nil {
						util.PrintErrorLog("send stream", err.Error())
					}
				}(requestID)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			response, err := stream.Recv()
			if err != nil {
				cancelFunc()
				util.PrintErrorLog("receive stream", err.Error())

				return
			}

			go func(response *model.Response) {
				if !model.IsReqResDataMapAvailable("res", response.RequestId) {
					util.PrintErrorLog("receive stream", fmt.Sprintf("receive channel %s is not available", response.RequestId))

					return
				}

				model.SendResDataMapChan(response.RequestId, model.ReqResData{
					RequestId: response.RequestId,
					Body:      response.Body,
				})
			}(response)
		}
	}()

	wg.Wait()

	util.RemoveStreamNum()
	util.PrintSuccessLog("connect stream", "client disconnected")

	return nil
}
