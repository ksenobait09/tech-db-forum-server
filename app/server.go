package app

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strconv"
	"tech-db-server/app/handlers"
)

/*
TODO: ХРАНИТЬ В ПАМЯТИ ГЛОБАЛЬНУЮ СТАТИСТИКУ
*/

// чертов роутер не позволяет зарегать вместе api/forum/create и /api/forum/:slug/create
func dirtyHack(router *fasthttprouter.Router) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		if path == "/api/forum/create" {
			handlers.CreateForum(ctx)
			return
		}
		router.Handler(ctx)
	}
}

func ListenAndServe(port int) error {
	router := fasthttprouter.New()
	// router.POST("/api/forum/create", handlers.CreateForum) см dirtyHack()
	router.POST("/api/forum/:slug/create", handlers.CreateThreadAtForum)
	router.GET("/api/forum/:slug/details", handlers.GetForumDetails)
	router.GET("/api/forum/:slug/threads", handlers.GetForumThreads)
	//router.GET("/api/forum/:slug/users", handlers.GetForumUsers)
	//router.GET("/api/post/:id/details", handlers.GetPostDetails)
	//router.POST("/api/post/:id/details", handlers.UpdatePostDetails)
	//router.POST("/api/service/clear", handlers.ServiceClear)
	//router.GET("/api/service/status", handlers.GetServiceStatus)
	//router.POST("/api/thread/:slug_or_id/create", handlers.CreatePostAtThread)
	router.GET("/api/thread/:slug_or_id/details", handlers.GetThreadDetails)
	router.POST("/api/thread/:slug_or_id/details", handlers.UpdateThreadDetails)
	//router.GET("/api/thread/:slug_or_id/posts", handlers.GetThreadPosts)
	//router.POST("/api/thread/:slug_or_id/vote", handlers.VoteThread)
	router.POST("/api/user/:nickname/create", handlers.CreateUser)
	router.GET("/api/user/:nickname/profile", handlers.GetUser)
	router.POST("/api/user/:nickname/profile", handlers.UpdateUser)

	stringPort := ":" + strconv.Itoa(port)
	return fasthttp.ListenAndServe(stringPort, dirtyHack(router))
}
