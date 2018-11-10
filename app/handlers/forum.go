package handlers

import (
	"github.com/valyala/fasthttp"
	"tech-db-server/app/models/forum"
)

func CreateForum(ctx *fasthttp.RequestCtx) {
	f := &forum.Forum{}
	f.UnmarshalJSON(ctx.PostBody())

	f, status := f.Create()
	switch status {
	case forum.StatusOk:
		response(ctx, f, fasthttp.StatusCreated)
	case forum.StatusConflict:
		response(ctx, f, fasthttp.StatusConflict)
	case forum.StatusSomethingNotExist:
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

func GetForumUsers(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	limit := ctx.QueryArgs().GetUintOrZero("limit")
	desc := ctx.QueryArgs().GetBool("desc")
	since := string(ctx.QueryArgs().Peek("since"))
	users, status := forum.GetUsers(slug, limit, since, desc)
	if status == forum.StatusOk {
		response(ctx, users, fasthttp.StatusOK)
		return
	}
	responseWithDefaultError(ctx, fasthttp.StatusNotFound)
}
