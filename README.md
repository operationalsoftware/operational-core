# Operational Core

## Installation

### Install CompileDaemon
```bash
go install github.com/githubnemo/CompileDaemon@latest
```

### Install go-assets-builder
```bash
go install github.com/jessevdk/go-assets-builder@latest
```

-------------------------------------------------------------

## Migration Guidelines

- The `initialise` function for migration runs for new clients who does not have the `app_user` table in their database.
- The `migrate` function for migration contains new migrations that have not reached production yet, for local migrations one can comment the ones that have already been applied so that any error could be avoided.

-------------------------------------------------------------

## GIT Merge Guidelines

- Develop in `dev` branch and then merge into other branches once updates are finalized in the following order:
    `dev` -> `staging` -> `production`
- For merging `Operational Core` updates to `Operational Platforms` perform the following steps in `dev` branch of operational platforms and then merge to other branches in same order as in the above step:
    - git fetch upstream
    - git merge upstream/production