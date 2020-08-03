-- name: GetArticle :many
SELECT *
FROM articles 
WHERE id = $1 
;

-- name: ListArticles :many
SELECT * 
FROM articles 
;

-- name: CreateArticle :one 
INSERT INTO articles (
  src_id
  , title
  , content
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteArticle :exec
DELETE FROM articles
WHERE id = $1;