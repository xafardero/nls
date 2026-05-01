# How to Release a New Version

This guide walks through releasing a new version of `nls` from start to finish. No prior release experience required.

---

## Prerequisites

- You have write access to the GitHub repository.
- Your local `main` branch is up to date with the remote.
- All the changes you want to ship are already merged into `main`.

---

## Step 1 — Pick a version number

`nls` uses [Semantic Versioning](https://semver.org/) (`MAJOR.MINOR.PATCH`):

| Change type | Which number to bump | Example |
|---|---|---|
| Breaking change | MAJOR | `1.0.0` → `2.0.0` |
| New feature, backwards-compatible | MINOR | `0.1.4` → `0.2.0` |
| Bug fix or small improvement | PATCH | `0.1.4` → `0.1.5` |

For example, if the current version is `0.1.4` and you are shipping a bug fix, the next version is `0.1.5`.

---

## Step 2 — Create and push a Git tag

The release pipeline is triggered by pushing a tag that starts with `v`. The tag name becomes the version shown by `nls --version`.

```bash
git checkout main
git pull
git tag v0.1.5
git push origin v0.1.5
```

That's all you need to do manually. GitHub Actions takes over from here.

---

## Step 3 — Watch the release pipeline

Go to the **Actions** tab on GitHub and open the **Release** workflow run that was just triggered. It will:

1. Build three binaries:
   - `nls-linux-amd64`
   - `nls-linux-arm64`
   - `nls-macos-arm64`
2. Create a GitHub Release named `Release v0.1.5` and attach the binaries.

The whole process takes about two minutes.

---

## Step 4 — Verify the release

Once the workflow finishes:

1. Go to the **Releases** page on GitHub and confirm the new release appears with the three binary files attached.
2. Download the binary for your platform, make it executable, and check the version:

```bash
chmod +x nls-linux-amd64
./nls-linux-amd64 --version
# Expected output: nls v0.1.5
```

> **Note:** If you build locally with `go build`, the binary will report `nls dev`. That is expected — the real version is only injected by the release pipeline via `-ldflags`.

---

## Troubleshooting

**The workflow did not trigger.**
Confirm the tag starts with `v` (e.g. `v0.1.5`, not `0.1.5`). Only tags matching `v*` trigger the release workflow (see `.github/workflows/release.yml`).

**CI is failing.**
Run `go test ./...` locally to confirm tests pass before pushing the tag.
