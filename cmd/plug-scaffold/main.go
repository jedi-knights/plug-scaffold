// plug-scaffold is a Neovim plugin project generator. It emits a fresh
// plugin skeleton — one directory tree per opinionated style — that passes
// plug-audit cleanly out of the box.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jedi-knights/plug-scaffold/cmd/plug-scaffold/commands"
)

var version = "dev" // overridden at build time with -ldflags

func main() {
	root := &cobra.Command{
		Use:   "plug-scaffold",
		Short: "Neovim plugin project generator",
		Long: `plug-scaffold emits a fresh Neovim plugin skeleton — one directory tree
per opinionated style — that passes plug-audit cleanly out of the box.`,
		SilenceUsage: true,
	}

	root.AddCommand(
		commands.NewNewCmd(),
		commands.NewVersionCmd(version),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
