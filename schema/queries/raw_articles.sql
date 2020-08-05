-- name: GetRawArticle :many
SELECT *
FROM enews.raw_articles 
WHERE added_id = $1 
;

-- name: GetRawArticle_Limit :many
SELECT * 
FROM enews.raw_articles
LIMIT $1
;

-- name: ListRawArticles :many
SELECT * 
FROM enews.raw_articles 
;


