# autoport

A simple CLI tool that generates a port number given stable inputs.

## Installation

Pre-built binaries for Linux, macOS, and Windows (amd64 and arm64) are published to
[GitHub Releases](https://github.com/jamestelfer/autoport/releases). Every archive
carries a build provenance attestation.

Verify a downloaded archive with the [GitHub CLI](https://cli.github.com/):

```sh
gh attestation verify autoport_linux_amd64.tar.gz --repo jamestelfer/autoport \
  --signer-workflow chinmina/.github/.github/workflows/goreleaser-release.yml
```

Install with any of:

```sh
mise use -g github:jamestelfer/autoport
brew install jamestelfer/tap/autoport
nix profile install github:jamestelfer/nixpkgs#autoport
go install github.com/jamestelfer/autoport@latest
```

## Development

Run `mise install` for the pinned toolchain, then `just verify` before committing.
See `AGENTS.md` for conventions.

## Licence

Apache 2.0. See `LICENSE`.
