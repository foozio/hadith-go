# GitHub Actions Artifact Publishing Design

**Goal:** Provide a secure, auditable distribution path for the `hadith-api` binary by publishing it as a GitHub Actions artifact, while keeping the repository free of checked-in binaries.

**Architecture:** CI builds `hadith-api` on workflow runs and uploads the binary plus a checksum as a workflow artifact. The repository no longer stores the binary; the artifact is tied to a specific commit and build log. Artifacts are created on push events (and optionally tags) to avoid PR noise. This approach maintains provenance without introducing signing complexity.

**Tech Stack:** GitHub Actions, Go 1.21 (via `actions/setup-go`), `sha256sum`, `actions/upload-artifact`.

## Components

- **CI workflow update:** Add a build step that compiles `hadith-api` with explicit `GOOS`/`GOARCH` (default `linux/amd64`). Generate `SHA256SUMS` for the artifact(s). Upload both with stable names (e.g., `hadith-api-linux-amd64` and `SHA256SUMS`).
- **Repo hygiene:** Remove the checked-in `hadith-api` binary and add a `.gitignore` rule to prevent future accidental commits.
- **Docs:** Update `SECURITY.md` or `README.md` to describe how to retrieve artifacts and why binaries are not stored in git.

## Data Flow

1. CI runs on push.
2. Go toolchain builds `hadith-api` into a `dist/` directory.
3. `sha256sum` produces a checksum file.
4. `upload-artifact` publishes both files for download from the workflow run.

## Error Handling

- Build or checksum failures fail the CI job.
- Artifact upload failures fail the CI job (ensures no silent missing artifacts).

## Testing

- CI run validates `go build ./...` and `go test ./...` as before.
- Manual verification: download artifact from the workflow run and validate checksum.
