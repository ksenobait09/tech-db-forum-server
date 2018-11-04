package handlers

import (
	"github.com/valyala/fasthttp"
	"tech-db-server/app/models/post"
	"tech-db-server/app/singletoneLogger"
)

func CreatePostAtThread(ctx *fasthttp.RequestCtx) {
	slug, id := getSlugOrId(ctx)
	var posts post.PostPointList
	err := posts.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	status, posts := post.CreatePosts(slug, id, posts)
	switch status {
	case post.StatusOK:
		response(ctx, posts, fasthttp.StatusCreated)
	case post.StatusNoThreadOrAuthor:
		responseWithDefaultError(ctx, fasthttp.StatusNotFound)
	case post.StatusNoParent:
		responseWithDefaultError(ctx, fasthttp.StatusConflict)
	}
}
