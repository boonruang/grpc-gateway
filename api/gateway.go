package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/model"
	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/service"
	"gitlab.com/ths-develops-team/iot/vms/aestiq/gateway/util"
)

type GatewayAPI struct {
	Service service.GatewayService
}

func (a GatewayAPI) Gateway(ctx *gin.Context) {
	requestID := requestid.Get(ctx)

	var jsonMap map[string]interface{}
	jsonData := ""
	rawData, _ := ctx.GetRawData()
	if len(rawData) > 0 {
		err := json.Unmarshal(rawData, &jsonMap)
		if err != nil {
			util.PrintRequestLog(ctx.Request.Method+":"+ctx.Request.URL.Path, requestID, jsonMap)

			resp := model.APIResponse(ctx, http.StatusBadRequest, model.API_RESPONSE_MESSAGE_400, nil)
			util.PrintErrorResponseLog(ctx.Request.Method+":"+ctx.Request.URL.Path, requestID, resp, err.Error())

			return
		}

		jsonBytes, err := json.Marshal(jsonMap)
		if err != nil {
			util.PrintRequestLog(ctx.Request.Method+":"+ctx.Request.URL.Path, requestID, jsonMap)

			resp := model.APIResponse(ctx, http.StatusBadRequest, model.API_RESPONSE_MESSAGE_400, nil)
			util.PrintErrorResponseLog(ctx.Request.Method+":"+ctx.Request.URL.Path, requestID, resp, err.Error())

			return
		}

		jsonData = string(jsonBytes)
	}

	util.PrintRequestLog(ctx.Request.Method+":"+ctx.Request.URL.Path, requestID, jsonMap)

	resData, err := a.Service.SendGRPC(requestID, model.ReqResData{
		RequestId: requestID,
		Method:    ctx.Request.Method,
		Route:     ctx.Request.URL.Path,
		Body:      jsonData,
	})

	if err != nil {
		resp := model.APIResponse(ctx, http.StatusGatewayTimeout, model.API_RESPONSE_MESSAGE_504, nil)
		util.PrintErrorResponseLog(ctx.Request.Method+":"+ctx.Request.URL.Path, requestID, resp, err.Error())

		return
	}

	var response model.APIResponsePayload
	err = json.Unmarshal([]byte(resData.Body), &response)
	if err != nil {
		resp := model.APIResponse(ctx, http.StatusBadGateway, model.API_RESPONSE_MESSAGE_502, nil)
		util.PrintErrorResponseLog(ctx.Request.Method+":"+ctx.Request.URL.Path, requestID, resp, err.Error())
	} else {
		resp := model.APIResponse(ctx, response.Code, response.Message, response.Data)
		util.PrintSuccessResponseLog(ctx.Request.Method+":"+ctx.Request.URL.Path, requestID, resp)
	}
}
