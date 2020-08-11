// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error)
	DeleteArticle(ctx context.Context, id int32) error
	GetArticleByID(ctx context.Context, id int32) ([]Article, error)
	GetArticle_Limit(ctx context.Context, limit int32) ([]Article, error)
	GetRawArticle(ctx context.Context, addedID sql.NullString) ([]EnewsRawArticle, error)
	GetRawArticle_Limit(ctx context.Context, limit int32) ([]EnewsRawArticle, error)
	ListArticles(ctx context.Context) ([]Article, error)
	ListRawArticles(ctx context.Context) ([]EnewsRawArticle, error)
}

var _ Querier = (*Queries)(nil)
