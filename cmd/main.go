package main

import (
	"fmt"
	"time"

	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/api"
	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/grpc"
	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/model"
	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/service"
	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "timestamp"

	cfg, err := model.ReadConfig()
	if err != nil {
		panic(err)
	}

	closeReqIDChan := model.InitReqIDChan()
	defer closeReqIDChan()

	closeReqResDataMapChan := model.InitReqResDataMap()
	defer closeReqResDataMapChan()

	grpc.StartGatewayServer(cfg.AppGrpc)

	gatewayService := service.GatewayService{
		SendingRetry:       cfg.SendingRetry,
		SendingInterval:    cfg.SendingInterval,
		ConnectingRetry:    cfg.ConnectingRetry,
		ConnectingInterval: cfg.ConnectingInterval,
	}

	gatewayAPI := api.GatewayAPI{
		Service: gatewayService,
	}

	r := gin.New()

	// add middleware here
	r.Use(
		gin.Recovery(),
		cors.Default(),
		requestid.New(),
		gin.BasicAuth(gin.Accounts{
			cfg.AppRest.Username: cfg.AppRest.Password,
		}),
	)
	r.NoRoute(gatewayAPI.Gateway)
	r.NoMethod(gatewayAPI.Gateway)

	restAddr := fmt.Sprintf("%s:%s", cfg.AppRest.Host, cfg.AppRest.Port)
	if err := r.Run(restAddr); err != nil {
		util.PrintErrorLog("main", err.Error())
	}
}
