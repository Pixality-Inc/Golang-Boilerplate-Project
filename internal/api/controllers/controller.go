package controllers

import (
	"github.com/pixality-inc/golang-core/logger"
)

type ControllerImpl struct {
	*BooksController
}

func NewController(
	booksController *BooksController,
) Controller {
	return &ControllerImpl{
		BooksController: booksController,
	}
}

func NewLoggableController(name string) *logger.LoggableImpl {
	return logger.NewLoggableImplWithServiceAndFields(
		"controller",
		logger.Fields{
			"name": name,
		},
	)
}
