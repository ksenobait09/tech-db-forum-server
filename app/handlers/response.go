package handlers

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	error2 "tech-db-server/app/models/error"
	"tech-db-server/app/singletoneLogger"
)

func generateMessageJSON(responseMessage json.Marshaler) []byte {
	responseMessageJSON, err := responseMessage.MarshalJSON()
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	return responseMessageJSON
}

func response(ctx *fasthttp.RequestCtx, responseStruct json.Marshaler, status int) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)
	ctx.Write(generateMessageJSON(responseStruct))
}

func responseWithDefaultError(ctx *fasthttp.RequestCtx, status int) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)
	ctx.Write(error2.DefaultMessage)
}
