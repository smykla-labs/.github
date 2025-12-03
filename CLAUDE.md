# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Purpose

Organization-wide defaults and synchronization for smykla-labs repositories. This is a special `.github` repository that:

1. **Community Health Files** - Provides default templates for all smykla-labs repos (CODE_OF_CONDUCT, CONTRIBUTING, SECURITY, issue/PR templates)
2. **Label Sync** - Automated synchronization of GitHub labels across all repositories using crazy-max/ghaction-github-labeler
3. **File Sync** - Automated synchronization of specified files across repositories using BetaHuhn/repo-file-sync-action

## Repository Structure

```text
.
├── .github/
│   ├── labels.yml              # Label definitions synced to all repos
│   ├── sync.yml                # File sync configuration
│   ├── workflows/
│   │   ├── sync-labels.yml     # Label sync workflow
│   │   └── sync-files.yml      # File sync workflow
│   ├── actions/                # Reusable composite actions
│   │   ├── generate-token/     # GitHub App token generation
│   │   └── get-org-repos/      # Fetch org repositories
│   ├── ISSUE_TEMPLATE/         # Default issue templates
│   └── PULL_REQUEST_TEMPLATE.md
├── templates/                  # Source files for file sync
│   ├── CODE_OF_CONDUCT.md
│   ├── CONTRIBUTING.md
│   ├── SECURITY.md
│   └── .github/
│       ├── ISSUE_TEMPLATE/
│       └── PULL_REQUEST_TEMPLATE.md
├── CODE_OF_CONDUCT.md          # Org-wide defaults (native GitHub)
├── CONTRIBUTING.md
└── SECURITY.md
```

## Key Concepts

### Community Health Files (Native GitHub)

Files in the root automatically apply to all smykla-labs repositories that don't have their own versions. This is a native GitHub feature for `.github` repositories.

### Synchronization System

The synchronization system uses a hybrid approach:

**Label Sync:**
- **Action**: Custom composite action (`.github/actions/sync-labels-to-repo`)
- **Workflow**: `.github/workflows/sync-labels.yml`
- **Config**: `.github/labels.yml`
- **Method**: Direct API updates via GitHub CLI

**File Sync:**
- **Action**: [BetaHuhn/repo-file-sync-action](https://github.com/BetaHuhn/repo-file-sync-action) v1.23.3 (346 ⭐)
- **Workflow**: `.github/workflows/sync-files.yml`
- **Config**: `.github/sync.yml` (expanded from templates)
- **Method**: Creates PRs with file changes

**Flow:**
- Each workflow triggers independently on relevant file changes
- Label sync runs for all repos in parallel
- File sync creates PRs with `org-sync` label
- Uses `chore/org-sync` branch for file changes

### Label Synchronization

- **Source**: `.github/labels.yml` contains all label definitions
- **Format**: YAML with `name`, `color` (with #), and optional `description`
- **Target**: All repositories in smykla-labs organization (except `.github` itself)
- **Method**: Custom composite action using GitHub CLI
- **Features**: Efficient map-based diff, parallel processing, dry-run support

### File Synchronization

- **Source**: `templates/` directory
- **Config**: `.github/sync.yml` (uses `%ALL_REPOS%` placeholder)
- **Target**: All repositories in smykla-labs organization (except `.github` itself)
- **Method**: BetaHuhn/repo-file-sync-action creates PRs
- **Features**: Auto PR creation, template support, per-file commits

### Authentication

All workflows use the **smyklot** GitHub App for authentication:
- `vars.SMYKLOT_APP_ID` - GitHub App ID (org-level variable)
- `secrets.SMYKLOT_PRIVATE_KEY` - GitHub App private key (org-level secret)

## Common Tasks

### Adding New Labels

1. Edit `.github/labels.yml`
2. Follow format:
   ```yaml
   - name: "label-name"          # Max 50 chars
     color: "#hex-color"         # 6-char hex with #
     description: "description"  # Max 100 chars (optional)
   ```
3. Commit and push to `main` - syncs automatically to all repos

### Adding New Sync Files

1. Add file to `templates/` directory (preserving the desired path structure)
2. Add file to `.github/sync.yml` in the `files` list:
   ```yaml
   - source: templates/YOUR_FILE.md
     dest: YOUR_FILE.md
   ```
3. Commit and push to `main` - syncs automatically to all repos

### Manual Sync

Trigger workflows manually via GitHub Actions:
1. Go to Actions tab in this repository
2. Select "Sync Labels" or "Sync Files"
3. Click "Run workflow"
4. Optionally enable "Dry run" to preview changes

### Workflow Behavior

#### Label Sync Workflow

- **Trigger**: Push to `main` when `labels.yml` changes, or manual dispatch
- **Flow**:
  1. Get list of all org repositories
  2. For each repo: Generate token and sync labels
- **Matrix**: Processes all repos in parallel (fail-fast: false)
- **Action**: Custom composite action handles sync logic efficiently
- **Dry run**: Available via workflow dispatch

#### File Sync Workflow

- **Trigger**: Push to `main` when `templates/**` or `sync.yml` changes, or manual dispatch
- **Flow**:
  1. Get list of all org repositories
  2. Expand `sync.yml` replacing `%ALL_REPOS%` placeholder
  3. Run BetaHuhn/repo-file-sync-action with expanded config
- **Action**: BetaHuhn/repo-file-sync-action handles PR creation and file sync
- **Branch**: Uses `chore/org-sync` prefix
- **Commits**: One commit per file with `chore(sync):` prefix
- **PR**: Labeled with `org-sync`, title: "chore(sync): sync organization files"
- **Dry run**: Available via workflow dispatch

## Files to Edit

### Labels

- **Definition**: `.github/labels.yml`
- **Workflow**: `.github/workflows/sync-labels.yml`

### Files

- **Source**: `templates/` directory
- **Config**: `.github/sync.yml`
- **Workflow**: `.github/workflows/sync-files.yml`

### Community Health Files

- **Default templates**: Root directory (CODE_OF_CONDUCT.md, CONTRIBUTING.md, SECURITY.md)
- **Issue templates**: `.github/ISSUE_TEMPLATE/`
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md`

## Label Categories

Labels are organized by prefix:
- `kind/*` - Issue/PR type (bug, enhancement, documentation, question, security)
- `area/*` - Affected component (ci, docs, api, testing, deps)
- `ci/*` - CI behavior control (skip-tests, skip-lint, skip-build, force-full)
- `release/*` - Release triggers (major, minor, patch, or auto-detect from commits)
- Automation: `org-sync` (sync PRs), `smyklot:*` (pending-ci with merge strategies)
- `triage/*` - Issue status (duplicate, wontfix, invalid, needs-info)
- `priority/*` - Priority levels (low, medium, high, critical)
- Community: `good first issue`, `help wanted`

## Architecture Decisions

This repository uses a hybrid approach for synchronization:

### Label Sync: Custom Composite Action
- **Location**: `.github/actions/sync-labels-to-repo`
- **Why**: Most label sync actions don't support multi-repo targeting
- **Implementation**: Simple bash script using GitHub CLI and jq
- **Benefits**: Full control, no external dependencies, easy to maintain

### File Sync: BetaHuhn/repo-file-sync-action
- **Stars**: 346
- **Version**: v1.23.3
- **Why**: Most popular file sync action, auto PR creation, template support
- **Repository**: https://github.com/BetaHuhn/repo-file-sync-action

## Important Notes

- This repository does NOT use standard `make` targets or testing frameworks
- All automation is via GitHub Actions workflows
- Changes to `labels.yml` or `templates/**` trigger automatic syncs
- Workflows exclude `.github` repository from sync targets
- Label sync happens directly via API (no PRs)
- File sync creates PRs for review
- Dry run mode available to preview changes without making them
- Community actions pinned to commit SHAs for security
- Custom label sync action is simple and maintainable (~100 lines)
- File sync handled by trusted community action (346 ⭐)
- ~50% code reduction from original orchestrator approach
