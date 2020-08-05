// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createArticleStmt, err = db.PrepareContext(ctx, createArticle); err != nil {
		return nil, fmt.Errorf("error preparing query CreateArticle: %w", err)
	}
	if q.deleteArticleStmt, err = db.PrepareContext(ctx, deleteArticle); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteArticle: %w", err)
	}
	if q.getArticleStmt, err = db.PrepareContext(ctx, getArticle); err != nil {
		return nil, fmt.Errorf("error preparing query GetArticle: %w", err)
	}
	if q.getRawArticleStmt, err = db.PrepareContext(ctx, getRawArticle); err != nil {
		return nil, fmt.Errorf("error preparing query GetRawArticle: %w", err)
	}
	if q.getRawArticle_LimitStmt, err = db.PrepareContext(ctx, getRawArticle_Limit); err != nil {
		return nil, fmt.Errorf("error preparing query GetRawArticle_Limit: %w", err)
	}
	if q.listArticlesStmt, err = db.PrepareContext(ctx, listArticles); err != nil {
		return nil, fmt.Errorf("error preparing query ListArticles: %w", err)
	}
	if q.listRawArticlesStmt, err = db.PrepareContext(ctx, listRawArticles); err != nil {
		return nil, fmt.Errorf("error preparing query ListRawArticles: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createArticleStmt != nil {
		if cerr := q.createArticleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createArticleStmt: %w", cerr)
		}
	}
	if q.deleteArticleStmt != nil {
		if cerr := q.deleteArticleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteArticleStmt: %w", cerr)
		}
	}
	if q.getArticleStmt != nil {
		if cerr := q.getArticleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getArticleStmt: %w", cerr)
		}
	}
	if q.getRawArticleStmt != nil {
		if cerr := q.getRawArticleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRawArticleStmt: %w", cerr)
		}
	}
	if q.getRawArticle_LimitStmt != nil {
		if cerr := q.getRawArticle_LimitStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRawArticle_LimitStmt: %w", cerr)
		}
	}
	if q.listArticlesStmt != nil {
		if cerr := q.listArticlesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listArticlesStmt: %w", cerr)
		}
	}
	if q.listRawArticlesStmt != nil {
		if cerr := q.listRawArticlesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listRawArticlesStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                      DBTX
	tx                      *sql.Tx
	createArticleStmt       *sql.Stmt
	deleteArticleStmt       *sql.Stmt
	getArticleStmt          *sql.Stmt
	getRawArticleStmt       *sql.Stmt
	getRawArticle_LimitStmt *sql.Stmt
	listArticlesStmt        *sql.Stmt
	listRawArticlesStmt     *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                      tx,
		tx:                      tx,
		createArticleStmt:       q.createArticleStmt,
		deleteArticleStmt:       q.deleteArticleStmt,
		getArticleStmt:          q.getArticleStmt,
		getRawArticleStmt:       q.getRawArticleStmt,
		getRawArticle_LimitStmt: q.getRawArticle_LimitStmt,
		listArticlesStmt:        q.listArticlesStmt,
		listRawArticlesStmt:     q.listRawArticlesStmt,
	}
}
