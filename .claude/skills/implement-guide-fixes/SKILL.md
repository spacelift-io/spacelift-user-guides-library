---
name: implement-guide-fixes
description: Use when implementing fixes to Spacelift user guides — adding missing validations, fixing fragile Rego, correcting content issues, or addressing findings from evaluate-guide-group
---

# Implement Guide Fixes

## Overview

Guided process for applying fixes to guide YAML files. Ensures Rego validations are correct, consistent with existing patterns, and verified against the test suite before finishing.

## Process

### 1. Understand the Fix Before Touching Files

For each finding:
- What state change is the step asking the user to make?
- What input fields would reflect that change? (see Input Schema below)
- Is this a Rego fix, a content fix, or both?

Read the guide file in full before editing — don't fix in isolation.

### 2. Fix Types

#### Missing validation on an action step
1. Identify what Spacelift state the step creates or changes
2. Write Rego using the patterns below
3. Add `validationHint` — a single sentence telling the user what to do before proceeding
4. Sanity-check: could this Rego pass *before* the user does the step? If yes, strengthen it.

#### Fragile `count(runs) >= N`
Replace with the `latest_tracked_run` pattern — it checks the most recent run's status rather than assuming a total count. See Rego Patterns below.

#### `count(unconfirmed_runs) == 0` false positive
Replace with `latest_tracked_run.status == "FINISHED"` — the count check passes if runs never existed or already failed.

#### Missing `validationHint`
Every step with a `validation` block must have a `validationHint`. Keep it to one sentence: what the user needs to do/see before clicking Next.

#### Content / wording fix
Edit the `instruction` or `hint` field directly. Ensure instructions remain specific and actionable, hints add context beyond restating the instruction.

### 3. Rego Patterns

All Rego blocks must start with `package spacelift`.

**Stack lookup** (reuse across steps in the same guide):
```rego
main_stack := stack if {
  some stack in input.stacks
  stack.name == input.expectations.main_stack_name
}
```

**Latest tracked run** (preferred over counting runs):
```rego
tracked_runs contains run if {
  some run in input.runs
  run.stack_id == main_stack.id
  run.type == "TRACKED"
}

latest_tracked_run := run if {
  some run in tracked_runs
  run.created_at == max([r.created_at | some r in tracked_runs])
}

valid if {
  latest_tracked_run.status == "FINISHED"  # or "FAILED", "UNCONFIRMED"
}
```

**Resource/integration exists:**
```rego
valid if {
  some integration in input.aws_integrations
  integration.name == input.expectations.aws_integration_name
}
```

**Attachment exists:**
```rego
valid if {
  some attachment in input.aws_attachments
  attachment.name == input.expectations.aws_integration_name
  attachment.attached_to == main_stack.id
}
```

**Stack dependency:**
```rego
valid if {
  dependency_stack.id in main_stack.dependencies
}
```

**Stack boolean field:**
```rego
valid if {
  main_stack.autodeploy == true
}
```

### 4. Input Schema Reference

| Field | Type | Description |
|-------|------|-------------|
| `input.stacks[]` | array | All stacks; fields: `id`, `name`, `autodeploy`, `dependencies[]` |
| `input.runs[]` | array | All runs; fields: `id`, `stack_id`, `type` (`TRACKED`/`PROPOSED`), `status`, `created_at` |
| `input.aws_integrations[]` | array | AWS integrations; fields: `id`, `name` |
| `input.aws_attachments[]` | array | Integration→stack attachments; fields: `name`, `attached_to` (stack id) |
| `input.policies[]` | array | Policies; fields: `id`, `name` |
| `input.policy_attachments[]` | array | Policy→stack attachments; fields: `name`, `stack_id` |
| `input.expectations` | object | Chapter variable values keyed by variable name (e.g. `input.expectations.main_stack_name`) |

Run statuses: `INITIALIZING`, `PLANNING`, `UNCONFIRMED`, `APPLYING`, `FINISHED`, `FAILED`

### 5. Verify After Every Change

```bash
go test -v ./...
```

All tests must pass before considering a fix complete. Pay attention to:
- `TestSchemaValidation_Guides` — catches YAML structural issues
- `TestGuidesLoad` — runs all Go-level validation including referential integrity
- `TestPrerequisiteGuideSlugsExist` — catches broken prerequisite slug references

### 6. Cross-Guide Consistency Check

After fixing any validation, ask:
- Does the new validation assume state that a prior guide actually guarantees (via its own validation)?
- Does this fix change what state is guaranteed for the *next* guide in the chain?
- If you added a new `${variable}` reference, is it declared in `chapter.yaml`?

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Adding validation but forgetting `validationHint` | Always add both together |
| `count(unconfirmed_runs) == 0` to verify confirmation | Use `latest_tracked_run.status == "FINISHED"` |
| Rego references `main_stack.id` before `main_stack` is defined | Define the stack rule first |
| `valid if { latest_tracked_run.status == ... }` with no runs → undefined, not false | This is correct — fails closed. Don't "fix" it. |
| Fixing a Rego issue in one guide without checking if the same pattern is in other guides | Grep for the pattern across all guides |
