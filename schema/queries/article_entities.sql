-- name: GetArticleEntities_ByArticleID :many
SELECT *
FROM article_entities 
WHERE article_id = $1 
;

-- name: CreateArticleEntities :exec
INSERT INTO article_entities (
    article_id
    , entity
    , entity_type 
    , counts 
) VALUES (
    $1, $2, $3, $4
)
;