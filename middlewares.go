package main

import (
	"github.com/savsgio/atreugo/v9"
)

func authMiddleware(ctx *atreugo.RequestCtx) error {

	if string(ctx.Path()) == "/login" {
		return ctx.Next()
	}

	jwtCookie := ctx.Request.Header.Cookie("atreugo_jwt")

	if len(jwtCookie) == 0 {
		// return ctx.ErrorResponse(errors.New("login required"), fasthttp.StatusForbidden)
		return ctx.RedirectResponse("/login", ctx.Response.StatusCode())
	}

	token, user, err := validateToken(string(jwtCookie))
	if err != nil {
		return err
	}

	if !token.Valid {
		// return ctx.ErrorResponse(errors.New("your session is expired, login again please"), fasthttp.StatusForbidden)
		return ctx.RedirectResponse("/login", ctx.Response.StatusCode())
	} else {
		ctx.SetUserValue("user", user)
	}

	return ctx.Next()
}
