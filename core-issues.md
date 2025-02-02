# FORMAT

- **Type**: Fix / Enhance / Mod
- **Status**: To Do / WIP / For Review / Complete / Invoiced / Obsolete
- **Priority**: Low / High / Critical
- **Assigned to**: WW


# #2 - Add stock module

- **Type**: Fix / Enhance / Mod
- **Status**: To Do / WIP / For Review / Complete / Invoiced / Obsolete
- **Priority**: Low / High / Critical
- **Assigned to**: WW


# #1 - Migrate to PostgreSQL

- **Type**: Fix / Enhance / Mod
- **Status**: To Do / WIP / For Review / Complete / Invoiced / Obsolete
- **Priority**: Low / High / Critical
- **Assigned to**: WW

SQLite has been used thus far for OperationalCore and client OperationalPlatorms as it was proposed that a simple single file disk-based database has sufficient flexibility for our data needs. However, we have encountered 3 issues when using it:

1. There is no remote access as there is no client/server model. This was known and it was accepted that viewing a database remotely would be more difficult as there is an addtional step required to remotely backup the database and copy it to the local machine. However, this is more tedious than anticipated.
2. There is no native support for DECIMAL / NUMERIC data types and we are finding issues with floating point maths errors with British Rototherm in production.
3. SQLite only permits one writer at a time. With British Rototherm we have a dependecy on a Microsoft SQL Server database, specifically where we are running a number of SPs to modify data and require a long running transaction to be open in SQLite so that the stock transaction data between the two databases remains in sync. When SQLite goes to commit after success with MSQL Server, we are frequently facing a "database is locked" error due to only one writer being permitted at a time, causing the databases to become out of sync which is unacceptable. PostgreSQL has better support for multiple writers due to the client/server model.


- [ ] Create a connection on start up to postgres
- [ ] Restructure application to use
- [ ] Update casing of database columns
- [ ] Write a migration script
- [ ] Rename internal to pkg and move components and layouts inside

