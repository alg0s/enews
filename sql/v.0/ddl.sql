-- author: steve.dang

-- = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = -

/* Create tables for news entities network */

CREATE SCHEMA IF NOT EXISTS enet;
DROP TABLE IF EXISTS enet.entities;
DROP TABLE IF EXISTS enet.articles_entities;
DROP TABLE IF EXISTS enet.temp_articles_entities;
DROP TABLE IF EXISTS enet.matrix_article_entity_metadata;

-- ENTITIES --

CREATE TABLE IF NOT EXISTS enet.entities (
    ent_id      SERIAL PRIMARY KEY, 
    entity      VARCHAR NOT NULl, 
    entity_type VARCHAR(20) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(entity, entity_type)
);

CREATE INDEX idx_enet_entities ON enet.entities(entity, entity_type);


-- ARTICLES_ENTITIES --

CREATE TABLE IF NOT EXISTS enet.articles_entities (
    added_id    SERIAL,
    src         VARCHAR NOT NULL,
    src_id      VARCHAR NOT NULL, 
    entity      VARCHAR NOT NULL, 
    ent_id      INT NOT NULL,
    entity_type VARCHAR(20) NOT NULL, 
    counts      SMALLINT NOT NULL, 
    published_at TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    unique(src, src_id, entity, ent_id)
);

-- TEMP_ARTICLES_ENTITIES --

CREATE TABLE IF NOT EXISTS enet.temp_articles_entities (
    added_id    SERIAL,
    src         VARCHAR NOT NULL,
    src_id      VARCHAR NOT NULL, 
    entity      VARCHAR NOT NULL, 
    entity_type VARCHAR(20) NOT NULL, 
    counts      SMALLINT NOT NULL, 
    published_at TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    unique(src, src_id, entity_type, entity)
);

-- NER_LOG

CREATE TABLE IF NOT EXISTS enet.ner_logs (
    added_id    serial, 
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    src_id      VARCHAR, 
    src         VARCHAR, 
    src_tb      VARCHAR, 
    ops_type    VARCHAR(20) NOT NULL
);

-- MATRIX_ARTICLE_ENTITY_METADATA

DROP TABLE IF EXISTS enet.matrix_article_entity_metadata;

CREATE TABLE IF NOT EXISTS enet.matrix_article_entity_metadata (
    src         VARCHAR NOT NULL, 
    src_id      VARCHAR NOT NULL, 
    info        JSON  NOT NULL, 
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(src, src_id)
);

CREATE INDEX idx_enet_matrix_article_entity_metadata 
    ON enet.matrix_article_entity_metadata(src, src_id)
;



