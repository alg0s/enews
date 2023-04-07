-- name: GetUniqueEntities_ByName :many
SELECT * 
FROM unique_entities 
WHERE name = $1 
;

-- name: GetUniqueEntities_ByType :many
SELECT * 
FROM unique_entities
WHERE entity_type = $1
;

-- name: GetUniqueEntities_ByName_Type :one 
SELECT *
FROM unique_entities
WHERE 
    name = $1 
    AND entity_type = $2 
;

-- name: CreateUniqueEntity :exec 
INSERT INTO unique_entities (
    name, entity_type
) VALUES (
    $1, $2 
)
RETURNING *
;

-- name: InsertNewEntitiesFromStagedEntities :exec 
INSERT INTO article_entities 
    (article_id, entity, entity_type, counts)
    SELECT 
        se.article_id, 
        se.entity, 
        se.entity_type, 
        se.counts
    FROM 
        stage_extracted_entities se 
            LEFT JOIN 
        article_entities ae 
            ON  se.article_id = ae.article_id
    WHERE   
        ae.article_id IS NULL
;

