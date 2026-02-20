# Distribution

## Current: source only

Users need Go installed and build from source:

```sh
go install github.com/CaliLuke/plane-mcp-server@latest
```

## Planned

### GitHub Releases + GoReleaser

Cross-compile for all platforms (macOS arm64/amd64, Linux arm64/amd64, Windows) and publish to GitHub Releases via CI.

- Add `.goreleaser.yaml` config
- Add a GitHub Actions workflow that triggers on tag push
- Users download the binary for their platform from the releases page

### Homebrew tap

Set up a tap repo (`homebrew-tap`) so macOS users can install with:

```sh
brew install CaliLuke/tap/plane-mcp-server
```

GoReleaser can automate tap formula updates on each release.
