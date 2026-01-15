# Spacelift User Guides Library

A library of user guides for Spacelift, stored as code and embedded into the Spacelift backend.

## Purpose

This repository serves as the source of truth for Spacelift user guides content. By storing guides as structured YAML files in version control, content can be managed independently from the backend code. The library embeds this content at compile time using Go's `embed.FS`, providing type safety and compile-time validation that ensures invalid guides are caught during build rather than at runtime. This approach enables content contributors to add or modify guides through standard pull request workflows while maintaining consistency and reliability.

## Directory Structure

The repository uses a hierarchical directory structure that mirrors the GraphQL schema (Group → Chapter → Guide):

```
guides/
├── {group-slug}/
│   ├── group.yaml              # Group metadata
│   ├── {chapter-slug}/
│   │   ├── chapter.yaml       # Chapter metadata
│   │   ├── {guide-slug}.yaml  # Individual guide
│   │   └── ...
│   └── ...
└── ...
```

### Example Structure

```
guides/
├── getting-started/
│   ├── group.yaml
│   ├── first-stack/
│   │   ├── chapter.yaml
│   │   ├── create-stack.yaml
│   │   ├── configure-vcs.yaml
│   │   └── add-policy.yaml
│   └── collaboration/
│       ├── chapter.yaml
│       ├── invite-users.yaml
│       └── manage-teams.yaml
└── advanced-workflows/
    ├── group.yaml
    ├── drift-detection/
    │   ├── chapter.yaml
    │   └── enable-drift.yaml
    └── custom-inputs/
        ├── chapter.yaml
        └── terraform-variables.yaml
```

## Slug Generation

Slugs are automatically derived from the directory structure:

- **Group Slug**: Directory name (e.g., `getting-started`)
- **Chapter Slug**: Directory name (e.g., `first-stack`)
- **Guide Slug**: Filename without extension (e.g., `create-stack`)

Slugs provide human-readable identifiers at each level of the hierarchy. The hierarchical relationships (Group → Chapter → Guide) are maintained through the nested data structure. This approach eliminates manual ID management and ensures uniqueness within each level.

## File Formats

### group.yaml

Defines metadata for a group of guides.

```yaml
name: "Getting Started"
description: "Learn the basics of Spacelift"
skillLevel: BEGINNER  # BEGINNER, ENABLER, COMMANDER, or GUARDIAN
ordering: 1
```

**Required Fields:**
- `name` (string): Display name of the group
- `description` (string): Brief description of the group
- `skillLevel` (string): One of `BEGINNER`, `ENABLER`, `COMMANDER`, or `GUARDIAN`
- `ordering` (int): Display order (lower numbers appear first)

### chapter.yaml

Defines metadata for a chapter within a group.

```yaml
name: "Your First Stack"
description: "Create and manage your first stack"
ordering: 1
```

**Required Fields:**
- `name` (string): Display name of the chapter
- `description` (string): Brief description of the chapter
- `ordering` (int): Display order within the group (lower numbers appear first)

### {guide-slug}.yaml

Defines an individual guide with metadata, steps, and completion information.

```yaml
ordering: 1
metadata:
  title: "Create Your First Stack"
  description: "Learn how to create a stack in Spacelift"
  labels: ["terraform", "basics"]
  difficulty: "easy"
  minutesToComplete: 10

steps:
  - order: 1
    title: "Navigate to Stacks"
    instruction: "Click on **Stacks** in the left sidebar to open the stacks page."
    hint: "If you don't see the sidebar, click the menu icon in the top-left corner."
    docs:
      - title: "What is a Stack?"
        url: "https://docs.spacelift.io/concepts/stack"

  - order: 2
    title: "Create New Stack"
    instruction: "Click the **Add Stack** button in the top-right corner."
    docs:
      - title: "Stack Creation Guide"
        url: "https://docs.spacelift.io/concepts/stack/creating-a-stack"

completion:
  successMessage: "Congratulations! You've created your first stack. Next, learn how to connect it to your VCS."
  recommendedGuideIds:
    - "getting-started/first-stack/configure-vcs"
    - "getting-started/first-stack/add-policy"
```

**Required Fields:**

- `ordering` (int): Display order within the chapter (lower numbers appear first)

**metadata:**
- `title` (string): Display title of the guide
- `description` (string): Brief description
- `labels` ([]string): Tags for categorization
- `difficulty` (string): Difficulty level (e.g., "easy", "medium", "hard")
- `minutesToComplete` (int): Estimated time to complete (must be >= 0)

**steps:**
- `order` (int): Step number (must be > 0, unique within guide)
- `title` (string): Step title
- `instruction` (string): What the user should do (supports Markdown)
- `hint` (string, optional): Additional help text
- `docs` ([]object, optional): Related documentation links
  - `title` (string): Link text
  - `url` (string): Documentation URL

**completion:**
- `successMessage` (string): Message shown when guide is completed
- `recommendedGuideIds` ([]string): IDs of guides to suggest next

## How to Add New Guides

### 1. Create or Navigate to a Group

If the group doesn't exist, create a new directory under `guides/`:

```bash
mkdir -p guides/my-group
```

Create `group.yaml` with the group metadata:

```yaml
name: "My Group"
description: "Description of this group"
skillLevel: BEGINNER
```

### 2. Create a Chapter

Create a chapter directory within the group:

```bash
mkdir -p guides/my-group/my-chapter
```

Create `chapter.yaml` with the chapter metadata:

```yaml
name: "My Chapter"
description: "Description of this chapter"
```

### 3. Create a Guide

Create a YAML file for your guide in the chapter directory:

```bash
touch guides/my-group/my-chapter/my-guide.yaml
```

Fill in the guide content following the format described above. The filename (without extension) will become the guide's slug.

### 4. Test Your Changes

Run the tests to validate your guide structure:

```bash
go test -v
```

The library will validate:
- YAML syntax
- Required fields presence
- Skill level values
- Step ordering and uniqueness
- URL formats in documentation links
- Referenced guide IDs (if applicable)

### 5. Submit a Pull Request

Once tests pass, commit your changes and create a pull request. The CI pipeline will run validation automatically.

## Validation and Testing

### Compile-Time Validation

The library validates all content at build time by calling `Guides()`. Invalid content causes a panic, ensuring problems are caught before deployment:

```go
func Guides() (*Library, error) {
    lib, err := parse(guidesFS)
    if err != nil {
        panic("userguides: " + err.Error())
    }
    return lib, nil
}
```

### Validation Rules

- **Groups**: Must have name and valid skill level
- **Chapters**: Must have name
- **Guides**: Must have title, at least one step, and valid step ordering
- **Steps**: Must have positive order, unique within guide, and non-empty title/instruction
- **References**: Guide IDs in `recommendedGuideIds` should reference existing guides

### Running Tests

```bash
# Run all tests
go test -v

# Run specific test
go test -v -run TestGuidesLoad

# Check test coverage
go test -v -cover
```

## Integration with Backend

The Spacelift backend imports this library as a Go module:

```go
import userguidelib "github.com/spacelift-io/spacelift-user-guides-library"

// Load guides
lib, err := userguidelib.Guides()
```

Content is synced to the database during migrations, similar to policy templates. See the [design document](https://www.notion.so/spacelift/2e7251e5616a80e1afb8c72453a86566) for full integration details.

## Development Workflow

1. **Make changes** to guides in this repository
2. **Test locally** with `go test -v`
3. **Submit PR** for review
4. **Merge** to main branch
5. **Backend update**: Update `go.mod` in backend to new library version
6. **Deploy**: Backend deployment syncs new content to database

## Schema Versioning

Guide content is stored as JSONB in the database. If the schema needs to evolve, use the `schema_version` approach:

```go
switch schema_version {
  case 1:
    json.Decode(content, &v1Guide)
  case 2:
    json.Decode(content, &v2Guide)
}
```

## Design Document

For architectural details and implementation phases, see the full design document:
[User Guides Library - Design Document](https://www.notion.so/spacelift/2e7251e5616a80e1afb8c72453a86566)

## License

Copyright © Spacelift, Inc.
