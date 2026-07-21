# plug-scaffold

[![CI](https://github.com/jedi-knights/plug-scaffold/actions/workflows/ci.yml/badge.svg)](https://github.com/jedi-knights/plug-scaffold/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/jedi-knights/plug-scaffold?include_prereleases&sort=semver)](https://github.com/jedi-knights/plug-scaffold/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

Neovim plugin project generator. Emits a fresh plugin skeleton — one directory tree per opinionated style — that passes [`plug-audit`](https://github.com/jedi-knights/plug-audit) cleanly out of the box.

> **Status:** in development — v0.1.0 not yet released.

## Positioning

`plug-scaffold` is to Neovim plugins what `cargo new` is to Rust crates: an opinionated `new <name>` command that produces a working project on the first commit. Three styles ship in the box, each mirroring a real-world author's conventions:

| `--style` | Convention source |
|---|---|
| `omar` | `M.setup(opts, deps)` with dependency injection, `detector.should_load()` gate, snacks-first pickers |
| `tj` | Metatable-lazy `init.lua`, hand-written vimdoc under `doc/`, plenary-busted tests |
| `prime` | Singleton `:setup()`, no default keymaps, plenary-busted, harpoon-shape CI |

Every generated tree is `plug-audit`-clean on the first run — the linter's ideal is the scaffolder's default, not the other way around.

## Install

### From source

Requires Go 1.23+.

```sh
go install github.com/jedi-knights/plug-scaffold/cmd/plug-scaffold@latest
```

### Pre-built binaries

Download the archive for your OS/arch from the [releases page](https://github.com/jedi-knights/plug-scaffold/releases), extract, and drop the binary somewhere on your `PATH`.

### Docker

```sh
docker run --rm -v "$PWD:/work" -w /work \
  ghcr.io/jedi-knights/plug-scaffold:latest new my-plugin.nvim --org=<you>
```

## Quickstart

```sh
plug-scaffold new my-plugin.nvim --org=<your-github-user-or-org>
cd my-plugin.nvim
plug-audit check .        # zero findings on the first run
git init && git add . && git commit -m "chore: initial commit"
```

Pick a style with `--style`:

```sh
plug-scaffold new my-plugin.nvim --org=<you> --style=tj
plug-scaffold new my-plugin.nvim --org=<you> --style=prime
plug-scaffold new my-plugin.nvim --org=<you> --style=omar   # default
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--org` | *(required)* | GitHub org or user that will own the repo |
| `--style` | `omar` | Template style: `tj` \| `prime` \| `omar` |
| `--author` | `git config user.name` | Author name stamped into LICENSE + README |
| `--dir` | current directory | Parent directory to scaffold into |

## What gets emitted

Every style emits the same top-level layout — the *contents* of the Lua files are where styles diverge. Common shape:

```
<plugin-name>/
├── LICENSE
├── README.md
├── .gitignore
├── Makefile
├── plugin/<module>.lua          # reload guard + named augroup
├── lua/<module>/
│   ├── init.lua                 # style-specific: DI / metatable-lazy / singleton class
│   ├── config.lua               # defaults + merge helper
│   └── health.lua               # :checkhealth entry point
├── doc/<module>.txt             # hand-written vimdoc (tj + prime; omar skips)
├── tests/<module>_spec.lua      # plenary-busted spec
└── scripts/minimal_init.lua     # test bootstrap
```

The `omar` style additionally emits `lua/<module>/detector.lua` for the `should_load()` gate; `tj` and `prime` additionally emit `doc/<module>.txt`.

## Ship criterion (v0.1.0)

```sh
plug-scaffold new my-thing.nvim --org=<you> --style=<any> && \
  cd my-thing.nvim && \
  plug-audit check .
# plug-audit: no findings.
```

## Development

```sh
go build ./...
go test ./...
golangci-lint run ./...
```

## License

MIT. See [LICENSE](./LICENSE).
