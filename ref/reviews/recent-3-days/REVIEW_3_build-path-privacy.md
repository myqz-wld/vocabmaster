---
review_id: REVIEW_3
reviewed_at: 2026-09-16
baseline_commit: b249be6b740ea1eb5e47faae2fcec6a8bd57291e
expired: false
---

# Build Path Privacy

## Scope

```review-scope
Makefile
README.md
```

## Findings

The existing distribution executable retained absolute compiler and source
paths containing a local account identifier. Ignoring `build/` and `dist/`
does not remove those paths from an archive shared separately.

## Fixes Landed

- Build the executable with Go's `-trimpath` flag.
- Document the distribution build and retained CLI, symlink, and metadata.
- Regenerate the local archive and quarantine the superseded archive outside
  the repository and distribution inputs.

## Validation

- `mise exec -- make test`: Go tests and isolated installer tests passed.
- `mise exec -- make pack`: archive and checksum generation passed.
- Inspect the actual archive executable for the detected account identifier
  and absolute home-directory prefixes; neither remained.
- Verify that the archive retains `vm`, the README, and build metadata.

## Residual Risk

Existing installations and previously distributed copies require separate
handling if they need the same build-path cleanup. This change does not
reinstall or remove any installed program.
