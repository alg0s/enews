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
    id              SERIAL PRIMARY KEY
    , src_id        VARCHAR(250) NOT NULL    
    , title         TEXT 
    , content       TEXT
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Extracted_Articles stores extracted entities, playing a role of archiving article IDs
-- CREATE TABLE extracted_articles (
--     id              SERIAL PRIMARY KEY
--     , article_id    INT NOT NULL
--     , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
--     , FOREIGN KEY (article_id) 
--         REFERENCES articles(id)
--         ON DELETE CASCADE
--     , UNIQUE(article_id)
-- );

-- Annotated_Articles stores annotated articles and individual entities
CREATE TABLE annotated_articles (
    article_id INT NOT NULL 
    , annotation    TEXT NOT NULL 
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , UNIQUE(article_id)
);

-- Stage_Extracted_Articles stores extracted articles and their entities temporarily, 
-- acting as a staging table during the extraction process
CREATE TABLE stage_extracted_entities (
    article_id      INT NOT NULL 
    , entity        TEXT
    , entity_type   VARCHAR(250)
    , counts        SMALLINT NOT NULL
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , UNIQUE(article_id, entity, entity_type)
);

-- Article_Entities persists the entities extracted from each article
CREATE TABLE article_entities (
    id              SERIAL
    , article_id    INT NOT NULL
    , entity        TEXT
    , entity_type   VARCHAR(250)
    , counts        SMALLINT NOT NULL
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , FOREIGN KEY (article_id) 
        REFERENCES articles(id)
        ON DELETE CASCADE
    , UNIQUE(article_id, entity, entity_type)
);

-- Entity_Types persists the types of entities that could be extracted by NLP services
CREATE TABLE entity_types (
    id              SERIAL PRIMARY KEY
    , name          VARCHAR(250) NOT NULL
    , description   TEXT  
    , language      VARCHAR(50) NOT NULL
    , created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , UNIQUE(name)
);

CREATE TABLE unique_entities (
    id                  SERIAL 
    , name              TEXT
    , entity_type_id    INT 
    , created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , FOREIGN KEY (entity_type_id) 
        REFERENCES entity_types(id)
        ON DELETE CASCADE
    , UNIQUE(name, entity_type_id)
);


/** 
  GRAPH TABLES
*/


/** 
  METADATA TABLES 
*/

/** 
  SOURCE TABLES 
*/

CREATE SCHEMA enews;

CREATE TABLE enews.raw_articles (
	added_id      varchar(400) NULL,
	src_id        varchar(400) NULL,
	article_type  varchar(400) NULL,
	img_url       text NULL,
	title         text NULL,
	publish_time  time NULL,
	publish_date  date NULL,
	category      text NULL,
	author        text NULL,
	content_raw   text NULL,
	content_text  text NULL,
	tags_raw      text NULL,
	tags_text     text NULL,
	summary       text NULL,
	like_count    varchar(200) NULL,
	dislike_count varchar(200) NULL,
	rating_count  varchar(200) NULL,
	viral_count   varchar(200) NULL,
	comment_count varchar(200) NULL,
	topic_id int4 NULL,
	posted_at     timestamp NULL,
	created_at    timestamp NULL
);
