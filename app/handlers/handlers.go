package handlers

import (
	"github.com/valyala/fasthttp"
	"fmt"
)

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

func Hello (ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}
func Hello2 (ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s, %s!\n", ctx.UserValue("name"),  ctx.UserValue("name2"))
}

