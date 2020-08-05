# Enews

Enews is a tool to extract entities from documents, primary designed to use for news articles. 

## Features & Requirements 

- Ability to extract entities from an article 
- An event based system with tasks:
  - Extract entities
  - Save unique entities 
  - Connect articles and entities  
- Support scheduled extraction jobs 
- Flexible in input and output definition:
  - Input: articles, and users can define designed columns of content 
  - Ouput: allowed pre-existing database otherwise generate tables in default database
- Have a metadata system to record all activities, i.e log system 
- Packaged like Django, with editable settings 
- Support both Sqlite3 and PostgreSQL
- Have interface agreement as part of the setting
- Extensible for future modules with different languages such as English, France. 
- Written primarily in Go but the schema could be ported in Python or other languages
- Data schemas are language agnostic, meaning that it should be stored in YAML or XML 
- Using ORM to communicate with databases 
- Using queue with multiple workers structure 

## Installation

``` sh
go mod init
```

#### Customers


#### Events
Enews is a event based system. There are three types of events:

  1. Extract Entity 
  2. Filter Unique Entities 
  3. Add New Nodes To Enews Graph

## Usage

  1. Install enews 
   
  2. Configure settings 
     - Input database and source table 
     - Output database 

  3. Run 
  Enews will automatically run the following procedure:
     - Verify connection with Input database + verify input table 
     - Verify connection with Output database + verify output tables 
     - Establish TaskQueue
     - Controller will start picking up articles from the input table and drop into TaskQueue
     - Executor will pick up tasks from TaskQueue and distributes to available workers to process in parallel. The number of works could be set in Settings

  4. Audit 
  At the end of the job, Enews will audit its work by running Auditor. Auditor will go through the log table as well as output tables to double check the results. Audit outcome will be saved in table `audit_instances`

#### Scheduled Job
Enews provides a scheduling feature to automatically run extraction jobs on an hourly or daily basis.

#### Automatic Retries

#### Configuring Logging

## Dependencies
1. **sqlc**

## Development


## Test 




<!--
# vim: set tw=79:
-->