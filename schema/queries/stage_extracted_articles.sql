-- name: GetStageExtractedArticle_ByArticleID :many 
SELECT *
FROM stage_extracted_articles
WHERE article_id = $1 
;

-- name: CreateStageExtractedArticle :one 
INSERT INTO stage_extracted_articles (
	article_id
  , entity
  , entity_type
  , counts 
) VALUES (
  $1, $2, $3, $4
)
RETURNING *
;