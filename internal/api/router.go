package api

import (
	"github.com/pixality-inc/golang-core/http"
	"github.com/pixality-inc/golang-core/http/docs"
)

func NewRouter(
	docsEnabled bool,
	docsHandler *docs.Handler,
) http.Router {
	router := http.NewRouter()

	if docsEnabled {
		router.GET("/docs", docsHandler.Handle)
		router.GET("/docs/{file}", docsHandler.Handle)
	}

	return router
}
