package books

import (
	"context"
	"fmt"

	"github.com/pixality-inc/golang-boilerplate-project/internal/dao"
	"github.com/pixality-inc/golang-core/logger"
	"github.com/pixality-inc/golang-core/postgres"
)

type Service interface {
	ListBooks(ctx context.Context) ([]dao.Book, error)
}

type Impl struct {
	log      logger.Loggable
	db       postgres.Database
	booksDao dao.BooksDao
}

func New(
	db postgres.Database,
	booksDao dao.BooksDao,
) *Impl {
	return &Impl{
		log:      logger.NewLoggableImplWithService("books"),
		db:       db,
		booksDao: booksDao,
	}
}

func (s *Impl) ListBooks(ctx context.Context) ([]dao.Book, error) {
	booksRows, err := s.booksDao.List(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get books from db: %w", err)
	}

	return booksRows, nil
}
