# enews Extractor

## Workflow
Entity Extraction steps:
- Initiatize queue
- Initialize NLP servers
- Truncate table `stage_extracted_article_entities`
- Check table `articles`, terminates if empty
- Loader starts loading articles by batches, which is a parameter of Extractor and drop them into the queue
- Extractor workers start picking up articles from the queue and process them, inserting outputs into the table `stage_extracted_article_entities`
- Once the queue is empty, meaning workers have processed all articles in the batch, the Controller will run the SQL query to process the extracted articles:
    - insert article ID into table `extracted_articles`
    - insert new article and entities into table `article_entities`
    - insert new entity_type into table `entity_type`
      and new entities into table `unique_entities`
- Repeat for a new batch