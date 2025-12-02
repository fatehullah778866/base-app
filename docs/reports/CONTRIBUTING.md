# Contributing to Base-App Reports

This guide explains how to create and maintain reports in the Base-App Reports Center.

## Quick Start

1. **Choose the correct directory** based on report type
2. **Use kebab-case naming** (lowercase with hyphens)
3. **Include required frontmatter** fields
4. **Follow the report template**
5. **Validate before committing** using `./scripts/validate-reports.sh`

## Directory Structure

```
docs/reports/
├── audits/security/          # Security audits
├── technical/               # Architecture, API specs
├── implementation/           # Status, progress
│   ├── stages/              # Development stages
│   └── milestones/         # Milestone reports
├── analysis/                # Performance, coverage
├── services/                # Service-specific
│   ├── auth/
│   ├── theme/
│   └── webhook/
└── planning/                # Roadmaps, planning
```

## Naming Conventions

### Rules

1. **Format**: kebab-case (lowercase with hyphens)
2. **Pattern**: `{category}-{descriptive-name}.md`
3. **No dates**: Dates go in frontmatter, not filename
4. **No underscores**: Use hyphens
5. **No spaces**: Use hyphens

### Examples

✅ **Correct:**
- `authentication-security-audit.md`
- `api-performance-analysis.md`
- `v1-implementation-complete.md`
- `services/auth/jwt-implementation-complete.md`

❌ **Incorrect:**
- `Authentication_Security_Audit.md` (wrong case, wrong separator)
- `apiPerformanceAnalysis.md` (camelCase)
- `api-performance-analysis-2025-01-15.md` (date in filename)
- `API Performance Analysis.md` (spaces, wrong case)

## Report Template

Every report must follow this structure:

```markdown
# Report Title

**Date:** January 2025  
**Status:** Complete | In Progress | Draft  
**Category:** Audit | Technical | Implementation | Analysis | Planning  
**Service:** auth | theme | webhook | all  
**Author:** [Optional]

## Summary

Brief summary of the report (2-3 sentences).

## Details

Detailed content of the report...

## Findings / Results

Key findings, results, or outcomes...

## Recommendations

Action items or next steps...

## Related Reports

- [Related Report Name](./path/to/report.md)
```

## Required Frontmatter Fields

### Date
- **Format**: Month YYYY (e.g., "January 2025")
- **Required**: Yes
- **Example**: `**Date:** November 2025`

### Status
- **Values**: Complete, In Progress, Draft
- **Required**: Yes
- **Example**: `**Status:** Complete`

### Category
- **Values**: Audit, Technical, Implementation, Analysis, Planning
- **Required**: Yes
- **Example**: `**Category:** Technical`

### Service
- **Values**: auth, theme, webhook, all
- **Required**: Yes (use "all" for general reports)
- **Example**: `**Service:** auth`

## Report Categories

### Audits (`audits/`)

Security audits, code quality reviews, compliance assessments.

**Location**: `docs/reports/audits/security/`

**Examples**:
- `security-audit-report.md`
- `api-security-review.md`
- `authentication-security-audit.md`

### Technical (`technical/`)

Architecture documentation, API specifications, database schemas.

**Location**: `docs/reports/technical/`

**Examples**:
- `base-app-technical-report.md`
- `api-specification.md`
- `database-schema-documentation.md`

### Implementation (`implementation/`)

Implementation status, stage completion, milestone tracking.

**Location**: `docs/reports/implementation/` or subdirectories

**Examples**:
- `implementation-summary.md`
- `stages/phase1-authentication-complete.md`
- `milestones/auth-service-complete.md`

### Analysis (`analysis/`)

Code analysis, performance metrics, optimization reports.

**Location**: `docs/reports/analysis/`

**Examples**:
- `code-coverage-analysis.md`
- `performance-analysis.md`
- `api-performance-report.md`

### Services (`services/`)

Service-specific reports for individual backend services.

**Location**: `docs/reports/services/{service-name}/`

**Examples**:
- `services/auth/authentication-implementation-complete.md`
- `services/theme/theme-service-complete.md`
- `services/webhook/webhook-reliability-report.md`

### Planning (`planning/`)

Roadmaps, backlogs, sprint planning, feature planning.

**Location**: `docs/reports/planning/`

**Examples**:
- `v2-roadmap.md`
- `sprint-planning.md`
- `feature-backlog.md`

## When to Create Reports

### After Major Milestones
- Create in `implementation/milestones/`
- Document completed features
- Include next steps

### After Security Reviews
- Create in `audits/security/`
- Document findings
- Include recommendations

### After Performance Testing
- Create in `analysis/`
- Include metrics
- Document optimizations

### After Service Completion
- Create in `services/{service}/`
- Document implementation details
- Include testing results

### After Architecture Changes
- Create in `technical/`
- Document design decisions
- Include diagrams if applicable

### During Planning
- Create in `planning/`
- Document roadmaps
- Include timelines

## Validation

Before committing reports, run the validation script:

```bash
./scripts/validate-reports.sh
```

The script checks:
- ✅ Directory structure
- ✅ Kebab-case naming
- ✅ Required frontmatter fields
- ✅ Broken links
- ✅ File locations

## Best Practices

### Content

1. **Be Specific**: Use clear, descriptive titles
2. **Be Complete**: Include all relevant information
3. **Be Current**: Update status and dates regularly
4. **Be Linked**: Reference related reports

### Structure

1. **Use Headers**: Organize with proper markdown headers
2. **Use Lists**: Break down complex information
3. **Use Code Blocks**: Include examples when relevant
4. **Use Links**: Link to related documentation

### Maintenance

1. **Update Status**: Keep status current
2. **Add Dates**: Update dates when modifying
3. **Link Reports**: Cross-reference related reports
4. **Archive Old**: Move outdated reports to `docs/archive/`

## Examples

### Example: Security Audit Report

**File**: `docs/reports/audits/security/api-security-review.md`

```markdown
# API Security Review

**Date:** November 2025  
**Status:** Complete  
**Category:** Audit  
**Service:** all  
**Author:** Security Team

## Summary

Comprehensive security review of Base-App API endpoints...

## Details

[Detailed content...]

## Findings

- Finding 1
- Finding 2

## Recommendations

- Recommendation 1
- Recommendation 2

## Related Reports

- [Initial Security Audit](./initial-security-audit.md)
- [Auth Service Report](../../services/auth/auth-service-complete.md)
```

### Example: Service Report

**File**: `docs/reports/services/auth/authentication-implementation-complete.md`

```markdown
# Authentication Service - Implementation Complete

**Date:** November 2025  
**Status:** Complete  
**Category:** Implementation  
**Service:** auth  
**Version:** 1.0

## Summary

The Authentication Service for Base-App v1.0 has been successfully implemented...

[Content...]
```

## Questions?

- Check [README.md](./README.md) for overview
- Review existing reports for examples
- Run validation script for guidance
- Follow the template structure

---

**Last Updated:** November 2025

