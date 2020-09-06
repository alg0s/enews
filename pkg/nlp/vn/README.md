# enews Extractor

## Introduction
Extractor extracts entities from text. Each entity will come with an entity type, such as PER-PERSON, LOC-LOCATION, or ORG-ORGANIZATION. 

Extractors have the following components:
- Controller 
- ArticleQueue
- Workers 
- Ingestor

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

## Error Handling

1. **Unique constraint violated**: when the db returns an error
2. **Server Not Responding**
3. **Server Overloaded**: require restarting the server, a signal needs to be send upstream to `main`, and all ongoing article processes have to be canceled and re-run once the server is ready again.
4.  

## Remaining problems
1. Stage 3A and 3B of the pipeline have yet run in concurrently due to sharing the same `in` channel >> DONE

2. Finish the pipeline that add entities into:
   - `article_entities`
   - `entity_types`
   - `unique_entities`

   This will happen at the end of all batches. 
   
3. Preprocessing a batch:
   - Remove <missing> content
   - Remove empty content
   - Removing newline in content 
   - Prepare a very long article since the server cannot serve too long >8192 characters

4. Error & Exception Handling:
   - Handle empty content 
   - Handle content that servers can't annotate
   - Handle crashing pipeline
   - Handle restarting/pausing pipeline when restarting the server

5. Storing metadata of batch jobs and runs 
   - Details of a run 
   - Details about articles being processed 
  
6. Unit tests for all functions
   - Test restarting server with feedback from the pipeline