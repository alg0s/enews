# Enews

Enews is a tool to extract entities from documents, primary designed to use for news articles. 

## Installation

``` sh
go mod init
```

### Customers


### Events


### Authentication with Connect

## Usage

### Automatic Retries

### Configuring Logging


## Development


## Test 

## Features + Requirements 

- Ability to extract entities from an article 
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




<!--
# vim: set tw=79:
-->