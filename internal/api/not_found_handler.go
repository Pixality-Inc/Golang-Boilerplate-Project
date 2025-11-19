package api

import (
	"github.com/pixality-inc/golang-core/http"

	"github.com/valyala/fasthttp"
)

func NotFoundRequestHandler(originalHandler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		originalHandler(ctx)

		if ctx.Response.StatusCode() == fasthttp.StatusNotFound {
			ctx.Response.ResetBody()

			http.NotFound(ctx, nil)
		}
	}
}
