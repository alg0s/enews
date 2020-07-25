package main

import (
	"database/sql"
	"time"
)

// RawArticle represents raw table of news article in db
type RawArticle struct {
	AddedID      int            `db:"added_id"`
	SrcID        sql.NullString `db:"src_id"`
	ArticleType  sql.NullString `db:"article_type"`
	URL          sql.NullString `db:"img_sr"`
	Title        sql.NullString `db:"title"`
	PublishTime  sql.NullTime   `db:"publish_time"`
	PublishDate  sql.NullTime   `db:"publish_date"`
	Category     sql.NullString `db:"category"`
	Author       sql.NullString `db:"author"`
	ContentRaw   sql.NullString `db:"content_raw"`
	ContentText  sql.NullString `db:"content_text"`
	TagsRaw      sql.NullString `db:"tags_raw"`
	TagsText     sql.NullString `db:"tags_text"`
	Summary      sql.NullString `db:"summary"`
	LikeCount    sql.NullInt32  `db:"like_count"`
	DislikeCount sql.NullInt32  `db:"dislike_count"`
	RatingCount  sql.NullInt32  `db:"rating_count"`
	ViralCount   sql.NullInt32  `db:"viral_count"`
	CommentCount sql.NullInt32  `db:"comment_count"`
	TopicID      sql.NullInt32  `db:"topic_id"`
	PostedAt     time.Time      `db:"posted_at"`
	CreatedAt    time.Time      `db:"created_at"`
}
