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

### Generate SSL certificates
```bash
./gen-dev-certs.sh
```

---

### Run Development Server

To start the development server, follow the steps below:

1. **Make the startup script executable**  
   Run the following command to ensure the `start-dev.sh` script has execute permissions:

   ```
   chmod +x start-dev.sh

2. **Start the development server**  
   Execute the script to launch the development environment:

   ```
   ./start-dev.sh

---

## Migration Guidelines

- The `initialise` function for migration runs for each new client.
- The `migrate` function for migration contains new migrations that have not reached production yet, for local migrations one can comment the ones that have already been applied so that any error could be avoided.

---

## GIT Merge Guidelines

- Develop in `dev` branch and then merge into other branches once updates are finalized in the following order:
    `dev` -> `staging` -> `production`
- For merging `Operational Core` updates to `Operational Platforms` perform the following steps in `dev` branch of operational platforms and then merge to other branches in same order as in the above step:
    - git fetch upstream
    - git merge upstream/production