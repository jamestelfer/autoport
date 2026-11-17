# autoport

A simple CLI tool that generates a port number given stable inputs, shipped as a single static executable.

## Go version

This project uses **Go 1.27**, which was released after the AI knowledge cutoff. Do not rely on training data for Go stdlib or dependency APIs — always fetch current documentation via Context7 before using an unfamiliar API.

All JSON encoding and decoding uses Go 1.27's `encoding/json/v2` (and
`encoding/json/jsontext` where needed). Importing the v1 `encoding/json` package
is categorically forbidden; the repository-wide `depguard` rule enforces this
for production code, tests, and tools.

## Build and test

```
just verify    # fmt + build + lint + test (run before committing)
just build     # produces dist/autoport
just test      # go test ./...
just fmt       # gofmt -w .
just lint      # golangci-lint run ./...
just xbuild    # goreleaser snapshot build across the full release matrix
```

Tool versions (Go, golangci-lint, goreleaser, just) are pinned in `mise.toml`; run `mise install` to match CI. CI reads the same versions via `mise current`.

## Unit tests

Use the standard library `testing` package only. No assertion or snapshot libraries.

## Cross-platform rules

Linux, macOS and Windows on `amd64` and `arm64` are equal release targets. CI tests on Linux only, and `just xbuild` cross-compiles the full matrix, so keep the code portable: use `path/filepath` for paths, keep filenames legal on Windows, build-tag anything POSIX-only, and keep `CGO_ENABLED=0`.

## Project layout

`main.go` in the repository root holds the whole program, built as the `autoport` binary. Keep it a single file with no dependencies outside the Go standard library.

## Commits and PR titles

Use Conventional Commits for all commit messages and PR titles. `pr-title.yml` enforces the format on PRs (the squash-merge commit is taken from the PR title). **Commit types carry version-bump semantics**: release-please is wired up, so `feat:` drives a minor bump and `fix:` a patch bump, and a `!`/`BREAKING CHANGE` drives a major.

| Type | When to use |
|---|---|
| `feat: <description>` | new user-visible feature (minor bump) |
| `fix: <description>` | bug fix (patch bump) |
| `chore:`, `docs:`, `refactor:`, `test:`, `ci:` | maintenance, no release |

## Key conventions

- Service configuration is environment variables, not flags.
- `http.Server` always gets an explicit `ReadHeaderTimeout` and `IdleTimeout`.

## Releases

release-please opens a release PR from Conventional Commits on `main`. Merging it creates a draft release and pushes the `v*` tag; `release.yml` then runs the shared chinmina GoReleaser workflow, which builds the matrix, attests provenance, publishes the draft, and commits the Homebrew cask to `jamestelfer/homebrew-tap` and the Nix recipe to `jamestelfer/nixpkgs`. The workflows mint keyless tokens through octo-sts, against the `release-please-autoport`, `release-autoport` and `release-tap` trust policies in `jamestelfer/.github`. Repository configuration follows `docs/reference/repository-configuration.md` in `chinmina/.github`.

## Dependencies

The program uses the Go standard library only. Do not add third-party modules.

Use Context7 for up-to-date GoReleaser documentation (`/goreleaser/goreleaser`); ldflags inject `main.version`, `main.commit` and `main.date`. The Go stdlib is not on Context7: use upstream docs directly.

### CLI conventions

Arguments are handled directly from `os.Args`. `--version` as the only argument prints the build metadata. With no arguments, the program prints usage to stderr and exits `1`. Otherwise it joins the arguments with `::`, hashes them with FNV-1a, and prints a port in the range 10000–65535 to stdout.
