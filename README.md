# OperationalCore

An open source platform for small and mid-sized manufacturers* to win at:

* 📈 Continuous improvement
* 💻 Digital transformation
* 🎯 Lean

OperationalCore is a full stack web-application which prioritises robust, manufacturing-specific building blocks over an opinionated finished product. These building blocks are presented to the developer as data-repositories and services which can be used to create custom user experiences, fast.

OperationalCore is designed to be forked and customised at the code-level rather than through extension/plugin develpment. This results in a unified user experience, a unified developer experience, and better maintainability.

\* _Also suitable for other goods-centric businesses such as e-commerce and wholesale distribution since the operational challenges faced by these companies are mostly a subset of manufacturing._

---
## Motivation

Until now, without the budget to develop 100% bespoke software, small and mid-size manufacturers have been constrained to off-the-shelf software and multi-tenant SaaS solutions. These systems force businesses to mould their processes to the software.

Customisation of OperationalCore enables innovation of processes and software in accordance with Lean Manufacturing and Continuous Improvement.

Only one significant player in the market offers an open source software solution for manufacturers - Odoo. However, Odoo is opinionated thus forcing process compliance. Customisation happens via app/plugin development, outside of the core codebase, meaning system behaviour becomes harder to reason about and extend as customisation increases.

OperationalCore aims to provide a different user, developer, and process-governance experience.

---

## Using

OperationalCore is undergoing rapid development due to the opportunity and demand in this space. If you wish to use this platform at this stage in its development, we suggest connecting with us at Operational Software on [hello@operationalsoftware.co](mailto:hello@operationalsoftware.co).

---

## Standards

As we figure out this problem space, we are commited to creating a set of standards which govern how we will develop and evolve the platform, and how we suggest developers fork and customise this platform. These standards are in use at Operational Software and will continue to be extended and improved. They are also intended to be understood by AI agents which are used to support development. See:

* [Commit Standard](./readme/COMMIT_STANDARD.md)
* More coming soon...

---

## Development

### Install `CompileDaemon`
```bash
go install github.com/githubnemo/CompileDaemon@latest
```

### Install `go-assets-builder`
```bash
go install github.com/jessevdk/go-assets-builder@latest
```

### Generate SSL certificates

Required for SSL/HTTPS in development (required for the [Web NFC API](https://w3c.github.io/web-nfc/) etc.)

```bash
./gen-dev-certs.sh
```

### Secrets and environment

Required environment variables (for either option):

- APP_ENV, GO_ENV, SITE_ADDRESS
- PG_USER, PG_PASSWORD, PG_HOST, PG_PORT, PG_DATABASE
- SWIFT_API_USER, SWIFT_API_KEY, SWIFT_AUTH_URL, SWIFT_TENANT_ID, SWIFT_CONTAINER
- SECURE_COOKIE_HASH_KEY, SECURE_COOKIE_BLOCK_KEY, AES_256_ENCRYPTION_KEY
- SYSTEM_USER_PASSWORD
- DUMP_PREFIX, ORBIT_BACKUP_CONTAINER (for backups)

### Run Development Server

To start the development server, follow the steps below:

1. **Make the startup script executable**  
   Run the following command to ensure the `start-dev.sh` script has execute permissions:

   ```
   chmod +x start-dev.sh
   ```

2. **Start the development server**  
   ```bash
   ./start-dev.sh
   ```
---
