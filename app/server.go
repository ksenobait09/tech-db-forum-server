package app

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strconv"
)



func ListenAndServe(port int) error {
	router := fasthttprouter.New()
	//router.POST("/forum/create", handlers.CreateForum)
	//router.POST("/forum/:slug/create", handlers.CreateThreadAtForum)
	//router.GET("/forum/:slug/details", handlers.GetForumDetails)
	//router.GET("/forum/:slug/threads", handlers.GetForumThreads)
	//router.GET("/forum/:slug/users", handlers.GetForumUsers)
	//router.GET("/post/:id/details", handlers.GetPostDetails)
	//router.POST("/post/:id/details", handlers.UpdatePostDetails)
	//router.POST("/service/clear", handlers.ServiceClear)
	//router.GET("/service/status", handlers.GetServiceStatus)
	//router.POST("/thread/:slug_or_id/create", handlers.CreatePostAtThread)
	//router.GET("/thread/:slug_or_id/details", handlers.GetThreadDetails)
	//router.POST("/thread/:slug_or_id/details", handlers.UpdateThreadDetails)
	//router.GET("/thread/:slug_or_id/posts", handlers.GetThreadPosts)
	//router.POST("/thread/:slug_or_id/vote", handlers.VoteThread)
	//router.POST("/user/:nickname/create", handlers.CreateUser)
	//router.GET("/user/:nickname/profile", handlers.GetUser)
	//router.POST("/user/:nickname/profile", handlers.UpdateUser)

	stringPort := ":" + strconv.Itoa(port)
	return fasthttp.ListenAndServe(stringPort, router.Handler)
}