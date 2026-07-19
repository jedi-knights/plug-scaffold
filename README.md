# plug-scaffold

Neovim plugin project generator. Emits a fresh plugin skeleton — one directory tree per opinionated style — that passes [`plug-audit`](https://github.com/jedi-knights/plug-audit) cleanly out of the box.

> **Status:** in development — v0.1.0 not yet released.

## Positioning

`plug-scaffold` is to Neovim plugins what `cargo new` is to Rust crates: an opinionated `new <name>` command that produces a working project on the first commit. Three styles ship in the box, each mirroring a real-world author's conventions:

| `--style` | Convention source |
|---|---|
| `tj` | Metatable-lazy `init.lua`, extension-protocol table shape, plenary-busted tests |
| `prime` | Singleton `:setup()`, no default keymaps, harpoon-shape CI |
| `omar` | `M.setup(opts, deps)` with dependency injection, `detector.should_load()` gate, snacks-first pickers |

## Ship criterion (v0.1.0)

```sh
plug-scaffold new my-thing --style=omar && cd my-thing && plug-audit check .
```

The generated tree must pass `plug-audit check` cleanly. `plug-audit`'s ideal is the scaffolder's default, not the other way around.

## License

MIT. See [LICENSE](./LICENSE).
