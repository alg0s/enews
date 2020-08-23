-- name: GetAnnotatedArticles_ByID :one 
SELECT *
FROM annotated_articles
WHERE article_id = $1
;


-- name: GetAnnotatedArticles :many
SELECT * 
FROM annotated_articles
;

-- name: CreateAnnotatedArticle :exec 
INSERT INTO annotated_articles (
    article_id
    , annotation
) VALUES (
    $1, $2 
)
;