# bsync

**bsync** keeps your source-built binaries in sync with their upstream GitHub releases.

On older LTS systems, `apt` often lags years behind upstream. Tools like `fastfetch`, `ripgrep`, or `fd` may not exist in your package manager at all — so you build them from source. But then they never update. `bsync` solves that.

---

## How it works

You register each source-built package in a simple YAML config. `bsync` checks GitHub's release API for newer versions, rebuilds the ones that are outdated, and swaps the binary on your PATH — all without touching `apt` or your system package manager.

---

## Installation

### Prerequisites

- Go 1.25+
- `git`, `cmake`, or whatever build tools your packages require

### Build from source

```bash
git clone https://github.com/kaushtubhkanishk/bsync.git
cd bsync
go build -o bsync .
sudo mv bsync /usr/local/bin/
```

---

## Configuration

Two envs have to be set for bsync for work:

```bash
export SOURCE_BUILD_DIR='/path/to/source/repos/'
export SOURCE_MANIFEST_PATH='/path/to/config/'
```

Packages are defined in YAML files under SOURCE_MANIFEST.
