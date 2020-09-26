-- name: GetArticle_ByID :one
SELECT *
FROM articles 
WHERE id = $1 
;

-- name: GetArticle_ByListID :many
SELECT *
FROM articles 
WHERE id = ANY($1::int[])
;

-- name: GetArticles_Limit :many 
SELECT * 
FROM articles 
LIMIT $1
;

-- name: GetArticles :many
SELECT *
FROM articles 
;

-- name: CreateArticle :exec
INSERT INTO articles (
	src_id
  	, title
  	, content
) VALUES (
  	$1, $2, $3
)
RETURNING *
;

-- name: DeleteArticle_ByID :exec
DELETE FROM articles
WHERE id = $1
;


-- name: GetUnprocessedArticleID :many
SELECT 
  a.id 
FROM 
  articles a
    LEFT JOIN 
  article_entities ae 
    ON a.id = ae.article_id 
    LEFT JOIN 
  stage_extracted_entities se 
    ON a.id = se.article_id
WHERE 
  ae.article_id IS NULL
  AND se.article_id IS NULL
;