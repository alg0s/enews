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


-- name: InsertNewArticleEntitiesFromStagedEntities :exec
INSERT INTO unique_entities 
    ("name", entity_type)
    SELECT 
        se.entity
        , se.entity_type
    FROM 
        stage_extracted_entities se
            LEFT JOIN 
        unique_entities ue
            ON  se.entity = ue.name
                AND se.entity_type = ue.entity_type
    WHERE        
        ue.name IS NULL 
        AND ue.entity_type IS NULL 
    GROUP BY 
        se.entity
        , se.entity_type
;