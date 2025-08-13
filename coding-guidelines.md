# Coding Guidelines

This document outlines the coding standards and best practices followed in this project. These guidelines ensure consistency, maintainability, and clarity across all components of the codebase.

---

## Database Guidelines

### Use `TEXT` Instead of `VARCHAR`
    Avoid using the `VARCHAR` data type unless a strict character limit is required. In most cases, `TEXT` offers greater flexibility with no significant performance difference, making it the preferred choice.

### Use `TIMESTAMP WITH TIME ZONE`
    Always use the `TIMESTAMP WITH TIME ZONE` data type for timestamp-related columns. This ensures proper handling of time across different time zones and environments.

### Use `table_name_id` for Primary Key columns
    The primary key of every table must follow the naming convention:
        `table_name_id`
    This format improves schema readability and prevents ambiguity when joining tables.


### Use `INTEGER GENERATED ALWAYS AS IDENTITY` for Primary Keys
    All primary key columns must be defined using the following SQL standard:
```sql
INTEGER GENERATED ALWAYS AS IDENTITY
```

---

## Development Standards

### Parsing Structured Data from Requests
    Use the appurl utility for parsing query parameters from the request URL. This ensures a consistent and reliable approach to handling structured request data across the application.

### Avoid HTML Injection via JavaScript
    Avoid injecting raw HTML using JavaScript. All dynamic HTML rendering should be handled by the backend to improve security, reduce complexity, and prevent XSS vulnerabilities.

### Prefer Backend-Driven Data Management
    Data should primarily be managed and rendered by the backend whenever possible. This improves maintainability, scalability, and separation of concerns between frontend and backend layers.

---

### Structuring Go Models
    * The data models used outside the repository should be using primitive Go data types while the ones used in repository should be using `pgtypes` for NOT NULL fields.
    * The model that will return response data should have primitive Go data types in models and *pointer types will be used for fields where null needs to be returned otherwise falsy value for that data type will be returned.
    * The ToDomain() method of the DB model should be used to map pgtype data to the primitive data which will be then be sent in the response.
