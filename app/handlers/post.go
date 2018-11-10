package handlers

import (
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"tech-db-server/app/models/post"
	"tech-db-server/app/models/thread"
	"tech-db-server/app/singletoneLogger"
)

func parsePostId(ctx *fasthttp.RequestCtx) int64 {
	id, _ := strconv.ParseInt(ctx.UserValue("id").(string), 10, 64)
	return id
}

func CreatePostAtThread(ctx *fasthttp.RequestCtx) {
	slug, id := getThreadSlugOrId(ctx)
	var posts post.PostPointList
	err := posts.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	status, posts := post.CreatePosts(slug, id, posts)
	switch status {
	case post.StatusOK:
		response(ctx, posts, fasthttp.StatusCreated)
	case post.StatusNotExist:
		responseWithDefaultError(ctx, fasthttp.StatusNotFound)
	case post.StatusBadParent:
		responseWithDefaultError(ctx, fasthttp.StatusConflict)
	}
}

func GetThreadPosts(ctx *fasthttp.RequestCtx) {
	slug, id := getThreadSlugOrId(ctx)
	id = thread.GetThreadId(slug, id)
	if id == 0 {
		responseWithDefaultError(ctx, fasthttp.StatusNotFound)
		return
	}
	limit := ctx.QueryArgs().GetUintOrZero("limit")
	desc := ctx.QueryArgs().GetBool("desc")
	since := ctx.QueryArgs().GetUintOrZero("since")
	sort := string(ctx.QueryArgs().Peek("sort"))
	response(ctx, post.GetPosts(id, limit, since, sort, desc), fasthttp.StatusOK)
}

func UpdatePostDetails(ctx *fasthttp.RequestCtx) {
	p := &post.Post{}
	err := p.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	p.ID = parsePostId(ctx)
	if p.ID == 0 {
		responseWithDefaultError(ctx, fasthttp.StatusNotFound)
		return
	}
	status := p.Update()
	if status == post.StatusOK {
		response(ctx, p, fasthttp.StatusOK)
		return
	}
	responseWithDefaultError(ctx, fasthttp.StatusNotFound)
}

func GetPostDetails(ctx *fasthttp.RequestCtx) {
	id := parsePostId(ctx)
	related := ctx.QueryArgs().Peek("related")
	data := post.PostDetails(id, strings.Split(string(related), ","))
	if data == nil {
		responseWithDefaultError(ctx, fasthttp.StatusNotFound)
		return
	}
	response(ctx, data, fasthttp.StatusOK)
}
