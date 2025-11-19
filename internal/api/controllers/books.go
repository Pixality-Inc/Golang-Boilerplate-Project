package controllers

import (
	"context"
	"fmt"

	"github.com/pixality-inc/golang-boilerplate-project/internal/books"
	"github.com/pixality-inc/golang-boilerplate-project/internal/dao"
	"github.com/pixality-inc/golang-boilerplate-project/internal/protocol"
	"github.com/pixality-inc/golang-core/logger"
)

type BooksController struct {
	log          logger.Loggable
	booksService books.Service
}

func NewBooksController(
	booksService books.Service,
) *BooksController {
	return &BooksController{
		log:          NewLoggableController("books"),
		booksService: booksService,
	}
}

func (c *BooksController) BooksGet(ctx context.Context) (*protocol.BooksResponse, error) {
	booksRows, err := c.booksService.ListBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list books: %w", err)
	}

	responseBooks := make([]*protocol.Book, 0, len(booksRows))

	for _, row := range booksRows {
		responseBooks = append(responseBooks, renderBook(row))
	}
	
	response := &protocol.BooksResponse{
		Total: int32(len(booksRows)),
		Books: responseBooks,
	}

	return response, nil
}

func renderBook(book dao.Book) *protocol.Book {
	return &protocol.Book{
		Id:    book.GetId().String(),
		Title: book.GetTitle(),
	}
}
