package handlers

import (
	"github.com/valyala/fasthttp"
	"tech-db-server/app/models/service"
)

func ServiceClear(ctx *fasthttp.RequestCtx) {
	service.ClearDatabase()
	responseWithDefaultError(ctx, fasthttp.StatusOK)
}

func GetServiceStatus(ctx *fasthttp.RequestCtx) {
	response(ctx, service.GetStatus(), fasthttp.StatusOK)
}