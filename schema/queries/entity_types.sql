-- name: GetEntityType_ByName :one
SELECT *
FROM entity_types 
WHERE name = $1
;

