-- name: GetArticle_ByID :many
SELECT *
FROM articles 
WHERE id = $1 
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

-- name: CreateManyArticles :exec
-- INSERT INTO articles (
-- 	src_id
-- 	, title
-- 	, content
-- ) VALUES $1
-- ;

-- name: DeleteArticle_ByID :exec
DELETE FROM articles
WHERE id = $1
;


