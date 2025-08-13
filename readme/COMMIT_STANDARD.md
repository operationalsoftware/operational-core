# Git Commit Message Cheat Sheet

**Format:**
```
<type>(<scope>): <summary>

[optional body]

[optional footer]
```

---

## Types
| Type     | Purpose |
|----------|---------|
| **feat** | New feature |
| **fix**  | Bug fix |
| **refactor** | Code restructure (no behaviour change) |
| **perf** | Performance improvement |
| **style** | Formatting, no logic change |
| **test** | Add/modify tests |
| **docs** | Documentation change |
| **chore** | Maintenance (deps, configs) |
| **build** | Build system/dependency changes |
| **ci** | CI/CD configuration |

---

## Rules
- **Scope**: Optional, e.g. `auth`, `ui`, `api`
- **Summary**: <=50 chars, imperative mood, no full stop
- **Body**: Explain *why*, wrap at 72 chars
- **Footer**: Link issues (`Closes #123`)

---

## Examples
```
feat(api): add pagination to user list

Added `page` and `limit` query params to improve data loading
performance for large datasets.

Closes #87
```
