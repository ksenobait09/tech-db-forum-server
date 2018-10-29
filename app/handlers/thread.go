package handlers

import (
	"github.com/go-openapi/strfmt"
	"github.com/valyala/fasthttp"
	"strconv"
	"tech-db-server/app/models/thread"
	"tech-db-server/app/singletoneLogger"
)

func getSlugOrId(ctx *fasthttp.RequestCtx) (string, int) {
	slug := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slug, 10, 32)
	return slug, int(id)
}

func CreateThreadAtForum(ctx *fasthttp.RequestCtx) {
	t := &thread.Thread{}
	t.Forum = ctx.UserValue("slug").(string)
	err := t.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	status := t.Create()
	switch status {
	case thread.StatusOk:
		response(ctx, t, fasthttp.StatusCreated)
	case thread.StatusConflict:
		response(ctx, t, fasthttp.StatusConflict)
	case thread.StatusUserOrForumNotExist:
		responseWithDefaultError(ctx, fasthttp.StatusNotFound)
	}
}

func GetThreadDetails(ctx *fasthttp.RequestCtx) {
	t := &thread.Thread{}
	slug, id := getSlugOrId(ctx)
	status := t.Get(slug, id)
	if status == thread.StatusOk {
		response(ctx, t, fasthttp.StatusOK)
		return
	}
	responseWithDefaultError(ctx, fasthttp.StatusNotFound)
}

func UpdateThreadDetails(ctx *fasthttp.RequestCtx) {
	update := &thread.ThreadUpdate{}
	err := update.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	slug, id := getSlugOrId(ctx)
	t := update.Update(slug, id)
	if t != nil {
		response(ctx, t, fasthttp.StatusOK)
		return
	}
	responseWithDefaultError(ctx, fasthttp.StatusNotFound)
}

func GetForumThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	limit := ctx.QueryArgs().GetUintOrZero("limit")
	desc := ctx.QueryArgs().GetBool("desc")
	var since *strfmt.DateTime
	rawSince := ctx.QueryArgs().Peek("since")
	if rawSince != nil {
		datetime := strfmt.NewDateTime()
		err := datetime.UnmarshalText(rawSince)
		since = &datetime
		if err != nil {
			singletoneLogger.LogErrorWithStack(err)
		}
	}
	threads, status := thread.GetByForumSlug(slug, limit, since, desc)
	if status == thread.StatusUserOrForumNotExist {
		responseWithDefaultError(ctx, fasthttp.StatusNotFound)
		return
	}
	response(ctx, threads, fasthttp.StatusOK)

}
