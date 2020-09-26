-- name: GetStageExtractedEntities_ByArticleID :many
SELECT 
  *
FROM 
  stage_extracted_entities
WHERE 
  article_id = $1 
;

-- name: CreateStageExtractedEntity :exec
INSERT INTO stage_extracted_entities (
	article_id
  , entity
  , entity_type
  , counts 
) VALUES (
  $1, $2, $3, $4
)
RETURNING *
;