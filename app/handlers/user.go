package handlers

import (
	"github.com/valyala/fasthttp"
	"tech-db-server/app/models/user"
	"tech-db-server/app/singletoneLogger"
)

func CreateUser(ctx *fasthttp.RequestCtx) {
	u := user.User{}
	err := u.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	u.Nickname = ctx.UserValue("nickname").(string)
	createdUser, existedUsers := u.Create()
	if createdUser != nil {
		response(ctx, createdUser, fasthttp.StatusCreated)
		return
	}
	response(ctx, existedUsers, fasthttp.StatusConflict)
}

func GetUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname")
	u := user.Get(nickname.(string))
	if u != nil {
		response(ctx, u, fasthttp.StatusOK)
		return
	}
	responseWithDefaultError(ctx, fasthttp.StatusNotFound)
}

func UpdateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)
	userUpdate := user.UserUpdate{}
	err := userUpdate.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		singletoneLogger.LogErrorWithStack(err)
	}
	u, status := user.Update(nickname, &userUpdate)
	if status == user.StatusOk {
		response(ctx, u, fasthttp.StatusOK)
		return
	}
	if status == user.StatusConflict {
		responseWithDefaultError(ctx, fasthttp.StatusConflict)
		return
	}
	responseWithDefaultError(ctx, fasthttp.StatusNotFound)
}
