/* - - - WORK FLOW FROM TEMP -> ENET.ARTICLES_ENTITIES - - - - */

/* 1. insert new entities to ent.entities 
 * NOTE: might remove CTE if too slow
 * */

WITH 
    new_ents AS (
        SELECT
            t.entity, 
            t.entity_type,
            count(1) AS CNT
        FROM 
            enet.temp_articles_entities t
                LEFT JOIN 
            enet.entities e
                ON  t.entity = e.entity
                    AND t.entity_type = e.entity_type
        WHERE 
            e.entity IS NULL 
            AND e.entity_type IS NULL
        GROUP BY 
            t.entity, 
            t.entity_type
        ORDER BY 
            CNT DESC
)
INSERT INTO enet.entities 
    (entity, entity_type)
    SELECT 
        new_ents.entity, 
        new_ents.entity_type
    FROM 
        new_ents
;

/* 2. insert all temp values to ent.articles_entities 
 * */
INSERT INTO enet.articles_entities
    (
        src
        , src_id
        , entity
        , ent_id
        , entity_type
        , counts
        , published_at
    )
    SELECT 
        tmp.src
        , tmp.src_id
        , tmp.entity
        , ent.ent_id
        , tmp.entity_type
        , tmp.counts
        , tmp.published_at
    FROM 
        enet.temp_articles_entities tmp
            LEFT JOIN 
        enet.entities ent
            ON  tmp.entity = ent.entity
                AND tmp.entity_type = ent.entity_type
;

/* 3. Truncate ent.temp_articles_entities 
 * reset seq. 
 * */

TRUNCATE TABLE enet.temp_articles_entities 
RESTART IDENTITY CASCADE
;

-- = = = = = = = = = = = = = = = = = = = = = = = = --