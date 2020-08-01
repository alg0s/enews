SELECT * FROM enet.matrix_article_entity_metadata;
SELECT count(1) FROM enet.matrix_article_entity_metadata;


-- = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = - 

------ -*- CAUTION: TRUNCATING -*- ------

TRUNCATE TABLE enet.entities RESTART IDENTITY CASCADE;
TRUNCATE TABLE enet.articles_entities RESTART IDENTITY CASCADE;
TRUNCATE TABLE enet.temp_articles_entities RESTART IDENTITY CASCADE;

------ TESTING -------

SELECT * FROM enet.entities;
SELECT * FROM enet.articles_entities;
SELECT * FROM enet.temp_articles_entities ORDER BY added_id desc;
SELECT * FROM enet.matrix_article_entity_metadata
SELECT count(1) FROM enet.temp_articles_entities;
SELECT count(1) FROM enet.articles_entities;
SELECT count(1) FROM enet.entities;
SELECT count(1) FROM news.zing;
SELECT count(1) FROM enet.matrix_article_entity_metadata;
SELECT count(DISTINCT src_id) FROM news.zing;               
SELECT count(DISTINCT src_id) FROM enet.articles_entities;
SELECT count(DISTINCT src_id) FROM enet.temp_articles_entities;

/* 
 * Check unique entities in ent.entities 
 * */

SELECT 
    entity, entity_type, count(1) AS cnt
FROM 
    enet.entities
GROUP BY 
    entity, entity_type
HAVING 
    count(1) > 1
ORDER BY cnt DESC;


/* 
 * Select list of unprocessed articles 
 * */

SELECT 
    n.src_id
FROM 
    news.zing n LEFT JOIN enet.articles_entities e
    ON  n.src_id = e.src_id
WHERE 
    e.src_id IS NULL
;

/* Select total unprocessed articles */

SELECT 
    count(n.src_id) /1000
FROM 
    news.zing n LEFT JOIN enet.articles_entities e 
    ON n.src_id = e.src_id
WHERE 
    e.src_id IS NULL
;

/* 
 * Select new entities from temp
 * */
SELECT
    t.entity, 
    t.entity_type,
    count(1) AS CNT
FROM 
    enet.temp_articles_entities t
    LEFT JOIN enet.entities e
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
;

/* Gather ent_id for temp */
SELECT 
    tmp.src, tmp.src_id, 
    tmp.entity, ent.ent_id,
    tmp.entity_type, tmp.counts,
    tmp.published_at
FROM 
    enet.temp_articles_entities tmp
    LEFT JOIN enet.entities ent
        ON tmp.entity = ent.entity
           AND tmp.entity_type = ent.entity_type
;




/* Select articles for given entity, entity_type 
 * */
WITH articles AS (
SELECT 
    src_id
FROM    
    enet.articles_entities
WHERE 
    entity = 'vingroup'
    AND entity_type = 'B-ORG'
)
SELECT 
    e.src_id, 
    e.ent_id, 
    e.entity, 
    e.entity_type
FROM 
    enet.articles_entities e
    INNER JOIN articles 
    ON e.src_id = articles.src_id
;



/* Select article not in enet.matrix_article_entity_metadata
 * */

SELECT 
    e.src_id, 
    e.ent_id
FROM 
    enet.articles_entities e 
    LEFT JOIN enet.matrix_article_entity_metadata m
    ON e.src = m.src 
       AND e.src_id = m.src_id
WHERE 
    m.src IS NULL
    AND m.src_id IS NULL
    AND e.src = 'zing'
ORDER BY 
    e.src, e.src_id, e.ent_id
;

-- = = = = = = = = = = = = = = = = = = = = = = = = --


/* Select src contains ent_id 
 * */

SELECT 
    a.*
FROM 
    enet.articles_entities a 
    INNER JOIN enet.articles_entities b 
    ON a.src_id = b.src_id
WHERE b.ent_id = 730
;

/* Select src_id contains ent_id 
 * */
SELECT *
FROM enet.articles_entities
WHERE ent_id = 730
;


/* Select article pairs containing ent_id 
 * */
SELECT  m.src_id, m.info
FROM    enet.matrix_article_entity_metadata m
        INNER JOIN enet.articles_entities a
        ON m.src_id = a.src_id 
WHERE   ent_id = 730
;

SELECT * FROM enet.matrix_article_entity_metadata;

