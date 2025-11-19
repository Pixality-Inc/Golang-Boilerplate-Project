package dao

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pixality-inc/golang-boilerplate-project/internal/types"
	"github.com/pixality-inc/golang-core/postgres"
	"github.com/pixality-inc/golang-core/util"
)

//nolint:iface
type Book interface {
	BookGetter
}

type BooksDao interface {
	Create(ctx context.Context, queryRunner postgres.QueryRunner, a Book) error
	CreateBatch(ctx context.Context, queryRunner postgres.QueryRunner, arr []Book) error
	List(ctx context.Context, queryRunner postgres.QueryRunner) ([]Book, error)
	GetById(ctx context.Context, queryRunner postgres.QueryRunner, id types.BookId) (Book, error)
}

func (d *BooksDaoImpl) Create(ctx context.Context, queryRunner postgres.QueryRunner, a Book) error {
	query, err := postgres.BuildSimpleInsertQuery(d.baseInsertQuery(), d.insertColumns, a)
	if err != nil {
		return fmt.Errorf("sql failed to build insert for query %v: %w", query, err)
	}

	if err := postgres.ExecuteQuery(ctx, queryRunner, query); err != nil {
		return err
	}

	return nil
}

func (d *BooksDaoImpl) CreateBatch(ctx context.Context, queryRunner postgres.QueryRunner, arr []Book) error {
	query, err := postgres.BuildBulkInsertQuery(d.baseInsertQuery(), d.insertColumns, arr)
	if err != nil {
		return fmt.Errorf("sql failed to build insert for query %v: %w", query, err)
	}

	if err := postgres.ExecuteQuery(ctx, queryRunner, query); err != nil {
		return err
	}

	return nil
}

func (d *BooksDaoImpl) List(ctx context.Context, queryRunner postgres.QueryRunner) ([]Book, error) {
	query := d.baseSelectQuery()

	var rows []bookRow

	if err := postgres.ExecuteQueryRows(ctx, queryRunner, query, &rows); err != nil {
		return nil, err
	}

	return util.MapSimple(rows, convertBookRowToModel), nil
}

func (d *BooksDaoImpl) GetById(ctx context.Context, queryRunner postgres.QueryRunner, id types.BookId) (Book, error) {
	query := d.baseSelectQuery().
		Where(squirrel.Eq{"id": id}).
		Limit(1)

	var rows []bookRow

	if err := postgres.ExecuteQueryRows(ctx, queryRunner, query, &rows); err != nil {
		return nil, err
	}

	convertedRows := util.MapSimple(rows, convertBookRowToModel)

	if len(rows) > 0 {
		return convertedRows[0], nil
	} else {
		return nil, postgres.ErrNoRows
	}
}
