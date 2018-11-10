package handlers

import (
	"github.com/go-openapi/strfmt"
	"github.com/valyala/fasthttp"
	"strconv"
	"tech-db-server/app/models/thread"
)

func getThreadSlugOrId(ctx *fasthttp.RequestCtx) (string, int) {
	slug := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.ParseInt(slug, 10, 32)
	return slug, int(id)
}

func CreateThreadAtForum(ctx *fasthttp.RequestCtx) {
	t := &thread.Thread{}
	t.Forum = ctx.UserValue("slug").(string)
	t.UnmarshalJSON(ctx.PostBody())
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
	slug, id := getThreadSlugOrId(ctx)
	status := t.Get(slug, id)
	if status == thread.StatusOk {
		response(ctx, t, fasthttp.StatusOK)
		return
	}
	responseWithDefaultError(ctx, fasthttp.StatusNotFound)
}

func UpdateThreadDetails(ctx *fasthttp.RequestCtx) {
	update := &thread.ThreadUpdate{}
	update.UnmarshalJSON(ctx.PostBody())
	slug, id := getThreadSlugOrId(ctx)
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
		datetime.UnmarshalText(rawSince)
		since = &datetime
	}
	threads, status := thread.GetByForumSlug(slug, limit, since, desc)
	if status == thread.StatusUserOrForumNotExist {
		responseWithDefaultError(ctx, fasthttp.StatusNotFound)
		return
	}
	response(ctx, threads, fasthttp.StatusOK)
}

func VoteThread(ctx *fasthttp.RequestCtx) {
	vote := &thread.Vote{}
	vote.UnmarshalJSON(ctx.PostBody())
	slug, id := getThreadSlugOrId(ctx)
	t := thread.VoteForThread(slug, id, vote)
	if t != nil {
		response(ctx, t, fasthttp.StatusOK)
		return
	}
	responseWithDefaultError(ctx, fasthttp.StatusNotFound)
}
