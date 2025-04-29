# Coding Guidelines

-------------------------------------------------------------

## Database

### Use "Text" type instead of "varchar"
    Since varchar type does not provide noticeable performance benefits so we should use text type instead of it.

### Always Use "timestamp with time zone"
    Always Use "timestamp with time zone" data type for timestamp related column.

### Always Use "table_name_id" format for primary keys of table
    The primary key of a table should always have the following format where `_id` will be preceeded by the table name.

### Always Use "INTEGER GENERATED ALWAYS AS IDENTITY" type for primary keys
    The primary key of a table should always have the following type i.e. "INTEGER GENERATED ALWAYS AS IDENTITY".


-------------------------------------------------------------

## Standards

### Getting Structured data from request
    - Use appurl for parsing the query params from the url.

### Don't use html injection using javascript:
    - The primary source of html will be the backend instead of injecting it from javascript.

### Manage Data Primarily from backend
    -  Manage data from backend whenever possible.