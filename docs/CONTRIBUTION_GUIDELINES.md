  # Git Flow & Contribution Standards

This document defines how we name branches, write commits, and create pull requests across all repositories — including client forks — to ensure consistency and to **avoid unintended issue closures** when syncing code between repos.

---

## Branch Naming Convention (Git Flow)

We follow the [Git Flow](https://nvie.com/posts/a-successful-git-branching-model/) branching model.
| Branch Type | Purpose | Naming Convention | Example |
|--------------|----------|-------------------|----------|
| **main** | Production-ready code | `main` | `main` |
| **develop** | Active development base | `develop` | `develop` |
| **feature/** | New feature or enhancement | `feature/<short-desc>` | `feature/user-auth` |
| **bugfix/** | Fixes for non-critical bugs | `bugfix/<short-desc>` | `bugfix/fix-null-pointer` |
| **hotfix/** | Urgent fix from `main` | `hotfix/<short-desc>` | `hotfix/fix-prod-crash` |
| **release/** | Pre-release stabilization | `release/<version>` | `release/1.4.0` | | **task/** | Maintenance or CI work | `task/<short-desc>` | `task/update-dockerfile` |

### Example

```bash
# Create a feature branch from develop
git checkout develop
git checkout -b feature/add-user-api

```

## Commit Message Standard

We use a conventional commit style to keep history structured and searchable.

---

### Structure
```<type>: <short summary> (ref core#<issue-id>)```

---

### Allowed `<type>` Values

| Type | Purpose |
|------|----------|
| `feature` | For new features |
| `bugfix` | For fixing bugs |
| `hotfix` | For urgent fixes on `main` |
| `task` | For chores, maintenance, or CI updates |
| `refactor` | For code improvements without changing behavior |
| `docs` | For documentation changes |
| `test` | For adding or improving tests |

---

### Examples

```bash
feature: add user authentication (ref core#231)
bugfix: fix null pointer in barcode parser (ref core#145)
task: update CI workflow to Node 20

```

### Notes
- Use (ref core#123) instead of Fixes #123 to prevent auto-closing issues in client repos.
- Keep the subject line under 72 characters.
- Write commit bodies if needed to explain why, not what.

---

## Pull Request (PR) Standard

When opening a PR:

🔖 **PR Title**  
Prefix the title with `wip:` if the PR is still in progress.

Format:  
`<type>: <short summary>`  
or  
`wip <type>: <short summary>`

**Examples:**
- `feature: add user authentication`
- `bugfix: handle empty response from API`
- `task: migrate deployment script to bash`
- `wip feature: implement new caching layer`

---

🧾 **PR Description Template**

### Summary
Brief summary of what this PR does.

### Linked Issue
Fixes operationalsoftware/operational-core#<issue-id>

### Testing
- [ ] Verified locally
- [ ] Unit tests added/updated (if applicable)
- [ ] Deployment tested in staging

### Notes
(Optional) Any additional notes or follow-up tasks.

---

**Important:**  
- Only include `Fixes operationalsoftware/operational-core#123` in **main repo PRs**, not in commits.  
- This ensures that client repositories syncing from `operationalsoftware/operational-core` do **not auto-close their own issues**.

