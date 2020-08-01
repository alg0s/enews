-- version: 1.0 
-- PostgreSQL 11

-- The following tables are built in for enews. 
-- For different source systems, there will be 
-- a need to create a data-loading module to 
-- ingest data into these tables, in order to 
-- be extracted. 

/** 
  ARTICLE & ENTITY TABLES 
*/


CREATE TABLE articles (
    id              GENERATED ALWAYS AS IDENTITY
    , src_id        INT NOT NULL    
    , title         TEXT 
    , content       TEXT
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE extracted_articles (
    id              GENERATED ALWAYS AS IDENTITY
    , article_id    INT NOT NULL 
    , entities      TEXT
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
    FOREIGN KEY (article_id) 
        REFERENCES articles(added_id)
        ON DELETE CASCADE
);

-- Temp_Extracted_Articles stores extracted articles and their entities temporarily, 
-- acting as a staging table during the extraction process
CREATE TABLE temp_extracted_articles (
    article_id      INT NOT NULL 
    , entities      TEXT 
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
)

CREATE TABLE article_entities (
    id              GENERATED ALWAYS AS IDENTITY
    , article_id    INT NOT NULL
    , entity        TEXT
    , entity_type   VARCHAR(500)
    , counts        SMALLINT NOT NULL
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
    FOREIGN KEY (article_id) 
        REFERENCES articles(added_id)
        ON DELETE CASCADE
);

CREATE TABLE unique_entities (
    id                  GENERATED ALWAYS AS IDENTITY 
    , name              TEXT 
    , entity_type_id    INT  
    , created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
    FOREIGN KEY (entity_type_id) 
        REFERENCES entity_types(id)
        ON DELETE CASCADE
);

CREATE TABLE entity_types (
    id              GENERATED ALWAYS AS IDENTITY 
    , name          VARCHAR(500) NOT NULL
    , description   TEXT  
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


/** 
  GRAPH TABLES
*/


/** 
  METADATA TABLES 
*/

