# Release Process

This document outlines the standard release procedure for the AI Hypervisor Platform.

## Versioning

We strictly follow [Semantic Versioning 2.0.0](https://semver.org/).

- **Major** (`X.y.z`): Incompatible API or structural changes.
- **Minor** (`x.Y.z`): Backwards-compatible new features.
- **Patch** (`x.y.Z`): Backwards-compatible bug fixes.

## Creating a Release

1.  **Prepare the Release Branch**: Create a release branch off `main` (e.g., `release-v1.2.0`).
2.  **Update Changelog**: Update `CHANGELOG.md`, moving contents from `[Unreleased]` to the new version header with today's date.
3.  **Bump Version**: Ensure any hardcoded version references (if any) are updated.
4.  **Create a PR**: Submit a Pull Request from the release branch to `main`.
5.  **Merge & Tag**: Once the PR is approved and merged, tag the commit on `main`.
    ```bash
    # git tag -a v1.2.0 -m "Release v1.2.0"
    # git push origin v1.2.0
    ```
6.  **GitHub Release**: The CI/CD pipeline will automatically build binaries, publish Docker images, and create a Draft GitHub Release.
7.  **Finalize**: Review the auto-generated GitHub Release notes and publish.
