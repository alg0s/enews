-- name: GetUniqueEntities_ByName :many
SELECT * 
FROM unique_entities 
WHERE name = $1 
;

-- name: GetUniqueEntities_ByEntityType :many
SELECT * 
FROM unique_entities
WHERE entity_type_id = $1
;

-- name: GetUniqueEntities_ByName_EntityType :one 
SELECT *
FROM unique_entities
WHERE 
    name = $1 
    AND entity_type_id = $2 
;

-- name: CreateUniqueEntity :exec 
INSERT INTO unique_entities (
    name
    , entity_type_id
) VALUES (
    $1, $2 
)
RETURNING *
;