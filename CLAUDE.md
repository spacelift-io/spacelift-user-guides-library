# spacelift-user-guides-library

Go library (`github.com/spacelift-io/spacelift-user-guides-library`) that is the source of truth for Spacelift's in-product user guides. Embedded into the Spacelift backend at compile time via `embed.FS`.

## Architecture

Hierarchical structure mirroring the GraphQL schema: **Group → Chapter → Guide**

```
guides/
├── {group-slug}/
│   ├── group.yaml
│   └── {chapter-slug}/
│       ├── chapter.yaml
│       └── {guide-slug}.yaml
```

### group.yaml
```yaml
name: string          # required
description: string   # required
skillLevel: BEGINNER | ENABLER | COMMANDER | GUARDIAN  # required
ordering: int         # required
```

### chapter.yaml
```yaml
name: string          # required
description: string   # required
ordering: int         # required
variables:            # optional - template variables for guide steps
  - name: string
    description: string
    resourceType: stack | policy | aws_integration
```

### {guide-slug}.yaml
```yaml
slug: string          # required - the guide's identifier (NOT derived from filename)
ordering: int         # required
metadata:
  title: string       # required
  description: string
  labels: []string
  difficulty: easy | medium | hard
  minutesToComplete: int  # >= 0
  prerequisites: []string
steps:                # at least one required; must be sequentially ordered 1, 2, 3...
  - order: int        # required, > 0, unique, sequential
    title: string     # required
    instruction: string  # required, supports Markdown
    hint: string      # optional
    validationHint: string  # optional
    validation: string      # optional - Rego policy for step completion check
    docs:
      - title: string
        url: string   # must use http/https
completion:
  successMessage: string
  recommendedGuideIds: []string  # must reference existing guide slugs
```

## Key Implementation Notes

- Guide slugs come from the `slug:` field inside the YAML file, NOT the filename
- Variables in step instructions use `${variable_name}` syntax
- Validation policies are Rego (OPA) — they receive `input.stacks`, `input.runs`, `input.expectations` etc.
- `recommendedGuideIds` references are bare slugs (e.g. `"credentials-not-secrets"`), validated for existence at parse time
- `Guides()` panics on any validation error — invalid content breaks the build

## Validation Rules

- Required fields must be present
- Valid enums: skillLevel, difficulty
- Steps must be sequentially ordered starting at 1 (no gaps, no duplicates)
- Doc URLs must use http/https
- `recommendedGuideIds` must reference existing guide slugs
- No duplicate slugs within groups, chapters, or guides
- Labels must be non-empty strings
- `minutesToComplete` must be >= 0
- Chapter variables must have a valid `resourceType`

## Testing & CI

```bash
go test -v        # validates all content + runs tests
go test -v -cover
```

GitHub Actions:
- `validate.yml` — runs on every push/PR to main: tests, go.mod verify, build
- `release.yml` — handles releases

## Current Content

- `foundations/getting-started/` — 5 guides (the only fully populated chapter)
- `configuration-reuse/`, `delivery-at-scale/`, `operational-safety/` — stub groups (group.yaml only)

## Backend Integration

Backend imports this as a Go module and syncs content to DB during migrations. Workflow:
1. Change guides here
2. `go test -v` locally
3. PR → merge to main
4. Backend updates `go.mod` to new version
5. Backend deploy syncs to DB
