---
name: evaluate-guide-group
description: Use when reviewing a Spacelift user guide group for content quality, validation coverage, completeness, or before making improvements to guides
---

# Evaluate Guide Group

## Overview

Systematic review of a guide group (Group → Chapters → Guides). Covers content completeness, Rego validation correctness, and cross-guide consistency. Produces structured findings and targeted follow-up questions to improve guides collaboratively.

## Process

### 1. Read Everything First

Before evaluating anything, read ALL files in the group:
- `group.yaml` — skill level, ordering
- Each `chapter.yaml` — name, ordering, variables defined
- Every guide YAML in full

### 2. Evaluate Each Guide

#### Content Completeness
- Each step instruction is specific and actionable (not vague like "configure it")
- Hints add context beyond restating the instruction
- Prerequisites accurately reflect what the user must have done
- Steps flow logically — no unexplained jumps or missing context
- `minutesToComplete` is realistic for the number and complexity of steps
- `successMessage` summarises what was actually learned, not just "good job"

#### Validation Coverage

| Step type | Needs validation? |
|-----------|-------------------|
| User takes an action in Spacelift that creates or changes verifiable state | YES |
| The climax / final action step of the guide | ALWAYS |
| User watches, reads, or observes a result | NO |
| Setup step whose state is verified by the immediately next step | Optional |

Flag: action steps with no validation. Flag: validations on pure observation steps.

#### Validation Correctness

For each Rego block:
- Does it actually verify what the instruction asked the user to do?
- Is it fragile? (e.g. `count(runs) >= N` assumes no extra user activity)
- Is there a corresponding `validationHint`?
- Could it pass *before* the user does the step? (false positive)
- Does it handle the case where required resources don't exist yet (undefined rule / undefined variable)?

#### Structural Integrity
- `prerequisiteGuideSlugs` chain is complete and reflects actual dependencies
- `recommendedGuideIds` points to a logical next guide (or is empty if the series ends)
- `ordering` is sequential with no gaps across guides in the chapter
- `difficulty` matches actual complexity
- `labels` are accurate

### 3. Cross-Guide Consistency

After evaluating individual guides:
- Every `${variable}` used in step instructions is defined in the chapter's `chapter.yaml`
- State assumed in later guides (e.g. "autodeploy is enabled", "S3 bucket exists") was actually guaranteed by an earlier guide's validation — not just mentioned in instructions
- `prerequisiteGuideSlugs` covers all implicit dependencies, not just the immediate one

### 4. Report Structure

1. **What's working well** — strengths first
2. **Issues by severity:**
   - **Critical** — missing validation on a key action step, broken cross-guide logic, false-positive validation
   - **Improvement** — fragile Rego, missing `validationHint`, unclear instruction, wrong difficulty
   - **Minor** — wording, hint quality, time estimate

For each issue: guide slug + step number + what's wrong + concrete suggested fix or question.

### 5. Ask Follow-Up Questions

Ask when:
- A validation seems missing but it's unclear whether the step is truly actionable vs. observational
- An instruction is ambiguous and the intended user action isn't clear from context
- Cross-guide state is assumed but unclear whether a prior guide guarantees it
- The purpose or learning goal of a step isn't evident from the content

Ask one focused question at a time. Don't ask about things determinable from the content.

## After Evaluation

If findings include Critical or Improvement issues, tell the user:
> "Use the `implement-guide-fixes` skill to work through these fixes."

## Common Issues to Watch For

- Final climax step of a guide has no validation (most common gap)
- `count(runs) >= N` Rego is fragile if users run extra triggered runs
- `validationHint` missing when `validation` is present
- Step says "verify X" but Rego checks something different
- `${variable}` used in instruction but not defined in `chapter.yaml`
- Later guide assumes state set up only in instructions (not validated) of an earlier guide
- `prerequisiteGuideSlugs` uses bare slug correctly but the slug doesn't match any guide's `slug:` field
