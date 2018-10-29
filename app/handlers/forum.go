package handlers

import (
	"github.com/valyala/fasthttp"
	"tech-db-server/app/models/forum"
	"tech-db-server/app/singletoneLogger"
)

func CreateForum(ctx *fasthttp.RequestCtx) {
	f := &forum.Forum{}
	err := f.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	f, status := f.Create()
	switch status {
	case forum.StatusOk:
		response(ctx, f, fasthttp.StatusCreated)
	case forum.StatusConflict:
		response(ctx, f, fasthttp.StatusConflict)
	case forum.StatusUserNotExist:
		responseWithDefaultError(ctx, fasthttp.StatusNotFound)
	}
}

func GetForumDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	f := forum.Get(slug)
	if f != nil {
		response(ctx, f, fasthttp.StatusOK)
		return
	}
	responseWithDefaultError(ctx, fasthttp.StatusNotFound)
}

//func GetForumUsers(ctx *fasthttp.RequestCtx) {
//	slug := ctx.UserValue("slug").(string)
//	limit := ctx.QueryArgs().GetUintOrZero("limit")
//	desc := ctx.QueryArgs().GetBool("desc")
//	since := string(ctx.QueryArgs().Peek("since"))
//	users := forum.GetUsers(slug, limit, since, desc)
//}
